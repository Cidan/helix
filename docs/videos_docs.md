# Videos Documentation

## Get Videos

This is an example of how to get videos.

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.GetVideos(context.Background(), &helix.VideosParams{
    GameID: "21779",
    Period: "month",
    Type:   "highlight",
    Sort:   "views",
    First:  10,
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Delete Videos

This is an example of how to delete videos.

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.DeleteVideos(context.Background(), &helix.DeleteVideosParams{
    IDs: []string{"992599293"},
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
