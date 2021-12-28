# Authorize Example

## Provisioning Twitch

1. First you must have two factor authentication enabled to register an application in the developer console.
1. Goto https://dev.twitch.tv/console/apps/create
1. Fill in a name of your choice.
1. Add a OAuth redirect URL to `http://localhost:8000/oauth/callback`.
1. Set the category as `Application Integration`
1. Click the create.
1. Copy the client ID and the client secret as you'll need them next.

