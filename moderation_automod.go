package helix

import "context"

type HeldMessageModerationResponse struct {
	ResponseCommon
}

type HeldMessageModerationParams struct {
	UserID string `query:"user_id"`
	MsgID  string `query:"msg_id"`
	Action string `query:"action"` // Must be "ALLOW" or "DENY".
}

// Required scope: moderator:manage:automod
func (c *Client) ModerateHeldMessage(ctx context.Context, params *HeldMessageModerationParams, opts ...Option) (*HeldMessageModerationResponse, error) {
	resp, err := c.postAsJSON(ctx, "/moderation/automod/message", nil, params, opts)
	if err != nil {
		return nil, err
	}

	moderation := &HeldMessageModerationResponse{}
	resp.HydrateResponseCommon(&moderation.ResponseCommon)

	return moderation, nil
}
