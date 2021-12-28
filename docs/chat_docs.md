# Chat Documentation

## Get Channel Chat Badges

This is an example of how to get channel chat badges

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.GetChannelChatBadges(context.Background(), &helix.GetChatBadgeParams{
    BroadcasterID: "145328278",
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Global Chat Badges

This is an example of how to get global chat badges

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.GetGlobalChatBadges(context.Background())
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Channel Emotes

This is an example of how to get channel emotes

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.GetChannelEmotes(context.Background(), &helix.GetChannelEmotesParams{
    BroadcasterID: "145328278",
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Global Emotes

This is an example of how to get global emotes

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.GetGlobalEmotes()
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Get Emote Sets

This is an example of how to get a set of emotes

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.GetEmoteSets(context.Background(), &helix.GetEmoteSetsParams{
    EmoteSetIDs: []string{"300678379"},
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
