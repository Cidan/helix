# Games Documentation

## Get Games

This is an example of how to get games.

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.GetGames(context.Background(), &helix.GamesParams{
    Names: []string{"Sea of Thieves", "Fortnite"},
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Top Games

This is an example of how to get top games.

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.GetTopGames(context.Background(), &helix.TopGamesParams{
    First: 20,
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
