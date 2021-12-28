# Bits Documentation

## Get Bits Leaderboard

This is an example of how to get the last 20 top bits contributers over the past week.

```go
client, err := helix.NewClient(
    context.Background(),
    helix.WithClientID("your-client-id"),
    helix.WithUserAccessToken("your-user-access-token"),
)
if err != nil {
    // handle error
}

resp, err := client.GetBitsLeaderboard(context.Background(), &helix.BitsLeaderboardParams{
    Count:  20,
    Period: "week",
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
