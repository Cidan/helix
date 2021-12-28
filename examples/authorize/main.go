package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/nicklaw5/helix/v2"
)

var (
	client *helix.Client
)

func httpServer() (<-chan string, func()) {
	nl, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal("unable to listen on port 8000", err)
	}

	out := make(chan string)
	srv := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")

			if code == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			out <- code
			w.Header().Add("content-type", "text/html")
			_, _ = w.Write([]byte("<h1>You may close this window and return to the terminal.</h1>"))
		}),
	}

	go srv.Serve(nl)

	return out, func() { srv.Close() }
}

func getAccessToken(ctx context.Context) (string, error) {
	cp, close := httpServer()
	defer close()

	url := client.GetAuthorizationURL(ctx, &helix.AuthorizationURLParams{
		ResponseType: "code",
		Scopes: []string{
			"user:read:email",
		},
	})

	fmt.Printf("Please open the following link in your browser: %s\n", url)

	select {
	case code := <-cp:
		fmt.Printf("code: %v\n", code)
		tok, err := client.RequestUserAccessToken(ctx, code)
		if err != nil {
			return "", fmt.Errorf("unable to get access token %w", err)
		}

		return tok.Data.AccessToken, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	clientId := flag.String("client-id", "", "The client id.")
	clientSecret := flag.String("client-secret", "", "The client secret.")

	flag.Parse()

	if *clientId == "" || *clientSecret == "" {
		log.Fatal("you must specify --client-id and --client-secret")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := helix.NewClient(
		ctx,
		helix.WithClientID(*clientId),
		helix.WithClientSecret(*clientSecret),
		helix.WithRedirectURI("http://localhost:8000/oauth/callback"),
	)
	if err != nil {
		log.Fatal("unable to build client", err)
	}
	client = c

	accessToken, err := getAccessToken(ctx)
	if err != nil {
		log.Fatal("unable to get access token", err)
	}

	_, u, err := c.ValidateToken(ctx, accessToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Logged in user ID: %s\n", u.Data.UserID)
}
