# Webhook Documentation

## Get Webhook Subscriptions

This is an example of how to get webhook subscriptions.

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
    helix.WithAppAccessToken("your-app-access-token"),
)
if err != nil {
    // handle error
}

resp, err := client.GetWebhookSubscriptions(context.Background(), &helix.WebhookSubscriptionsParams{
    First: 10,
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
