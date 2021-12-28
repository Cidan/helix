package helix

import "context"

// GetExtensionSecretResponse response structure received
// when generating or querying for generated secrets
type ExtensionSecretCreationResponse struct {
	Data ManyExtensionSecrets
	ResponseCommon
}

// GetExtensionSecretResponse response structure received
// when fetching secrets for an extension
type GetExtensionSecretResponse struct {
	Data ManyExtensionSecrets
	ResponseCommon
}

type SecretsInformation struct {
	Version int      `json:"format_version"`
	Secrets []Secret `json:"secrets"`
}

type ManyExtensionSecrets struct {
	SecretInfo []SecretsInformation `json:"data"`
}

// Secret information about a generated secret
type Secret struct {
	ActiveAt Time   `json:"active_at"`
	Content  string `json:"content"`
	Expires  Time   `json:"expires_at"`
}

type ExtensionSecretCreationParams struct {
	ActivationDelay int    `query:"delay,300"` // min 300
	ExtensionID     string `query:"extension_id"`
}

type GetExtensionSecretParams struct {
	ExtensionID string `query:"extension_id"`
}

func (c *Client) CreateExtensionSecret(ctx context.Context, params *ExtensionSecretCreationParams, opts ...Option) (*ExtensionSecretCreationResponse, error) {
	resp, err := c.post(ctx, "/extensions/jwt/secrets", &ManyExtensionSecrets{}, params, opts)
	if err != nil {
		return nil, err
	}

	events := &ExtensionSecretCreationResponse{}
	resp.HydrateResponseCommon(&events.ResponseCommon)
	events.Data.SecretInfo = resp.Data.(*ManyExtensionSecrets).SecretInfo

	return events, nil
}

func (c *Client) GetExtensionSecrets(ctx context.Context, params *GetExtensionSecretParams, opts ...Option) (*GetExtensionSecretResponse, error) {
	resp, err := c.postAsJSON(ctx, "/extensions/jwt/secrets", &ManyExtensionSecrets{}, params, opts)
	if err != nil {
		return nil, err
	}

	events := &GetExtensionSecretResponse{}
	resp.HydrateResponseCommon(&events.ResponseCommon)
	events.Data.SecretInfo = resp.Data.(*ManyExtensionSecrets).SecretInfo

	return events, nil
}
