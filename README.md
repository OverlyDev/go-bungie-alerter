# go-bungie-alerter

This is a simple project with the goal of monitoring various web sources and posting new items via a Discord webhook.

The repo is a mess, I know. I have no idea what I'm doing :)

Release binaries can be found to the right side, or [here](https://github.com/OverlyDev/go-bungie-alerter/releases/latest) is a direct link to the latest release.

## Usage
No matter which platform you run BungieAlerter on, it needs access to the variable `DISCORD_WEBHOOK`.

This variable holds the full discord webhook url needed to send alerts.

You can either:
1. Export while executing the binary:
    - Linux: `DISCORD_WEBHOOK="\<your webhook\>" ./BungieAlerter`
    - Windows (powershell): `$env:DISCORD_WEBHOOK="\<your webhook\>"; .\BungieAlerter-windows-amd64.exe; $env:DISCORD_WEBHOOK=$null`
2. Save it in a .env file beside BungieAlerter:
    - create .env file in the same directory as BungieAlerter
    - add DISCORD_WEBHOOK="\<your webhook\>" to it
    - Run the binary

## Future
There's a lot of refinement and features I'd like to add. Below should hopefully be an up-to-date listing of those.

- cli args
    - provide webhook
    - configuration options
    - flag to run without webhook notifications ("terminal mode"?)
- publish docker images (skeleton is already in place)
- test if minification breaks things (only briefly played around with this)
