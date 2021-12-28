# Entitlement Grants Documentation

## Create Entitlements Upload URL

This is an example of how to create a new entitlements upload URL endpoint:

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
    helix.WithAppAccessToken("your-app-access-token"),
)
if err != nil {
    // handle error
}

manifestID := "your-manifest-id"
entitlementType := "bulk_drops_grant"

resp, err := client.CreateEntitlementsUploadURL(context.Background(), manifestID, entitlementType)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
