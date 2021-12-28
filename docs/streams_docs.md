# Streams Documentation

## Get Streams

This is an example of how to get streams. Here we are requesting the first two streams from the English language.

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.GetStreams(context.Background(), &helix.StreamsParams{
    First: 10,
    Language: []string{"en"},
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Followed Streams

This is an example of how to get followed streams.

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.GetFollowedStream(context.Background(), &helix.FollowedStreamsParams{
    UserID: "123456",
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
