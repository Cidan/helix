# Predictions Documentation

## Get Predictions

This is an example of how to get predictions.

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.GetPredictions(context.Background(), &helix.PredictionsParams{
    BroadcasterID: "145328278",
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## Create Prediction

This is an example of how to create a prediction.

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.CreatePrediction(context.Background(), &helix.CreatePredictionParams{
    BroadcasterID: "145328278",
    Title: "Test",
    Outcomes: []helix.PredictionChoiceParam{
        helix.PredictionChoiceParam{ Title: "choice 1" },
        helix.PredictionChoiceParam{ Title: "choice 2" },
    },
    PredictionWindow: 300,
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```

## End Prediction

This is an example of how to end a prediction.

```go
client, err := helix.NewClient(context.Background()
    helix.WithClientID("your-client-id"),
)
if err != nil {
    // handle error
}

resp, err := client.EndPrediction(context.Background(), &helix.EndPredictionParams{
    BroadcasterID: "145328278",
    ID: "c36165d9-d5f5-4f81-ab56-17e7347110c8",
    Status: "RESOLVED",
    WinningOutcomeID: "d0c0194a-6016-4ca3-b8eb-0c61183758ab",
)
if err != nil {
    // handle error
}

fmt.Printf("%+v\n", resp)
```
