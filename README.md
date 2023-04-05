# go-bungie-alerter

This is a simple project with the goal of monitoring various web sources and posting new items via a Discord webhook.

The repo is a mess, I know. I have no idea what I'm doing :)

Release binaries can be found to the right side, or [here](https://github.com/OverlyDev/go-bungie-alerter/releases/latest) is a direct link to the latest release.

The changlog can be found [here](https://github.com/OverlyDev/go-bungie-alerter/blob/master/CHANGELOG.md)

## Usage

### Webhook
No matter which platform you run BungieAlerter on, it needs access to the variable `DISCORD_WEBHOOK`.

This variable holds the full discord webhook url needed to send alerts.

You can either:
1. Provide the webhook via the flag:
    - `BungieAlerter --webhook/-w <your webhook>`
2. Export while executing the binary:
    - Linux: `DISCORD_WEBHOOK="<your webhook>" ./BungieAlerter`
    - Windows (powershell): `$env:DISCORD_WEBHOOK="<your webhook>"; .\BungieAlerter-windows-amd64.exe; $env:DISCORD_WEBHOOK=$null`
3. Save it in a .env file beside BungieAlerter:
    - create .env file in the same directory as BungieAlerter
    - add `DISCORD_WEBHOOK="<your webhook>"` to it
    - Run the binary

Note: Depending on how you're running the application (CLI or Docker), different methods of supplying `DISCORD_WEBHOOK` will be more suitable.

### CLI
There's now a basic CLI. Running the binary without any args will provide you with usage information. It won't actually start unless you give it the `go` arg.

A quick overview of the available options:

args:
- `go`      - starts BungieAlerter
- `update`  - checks for a new version of BungieAlerter
- `help`    - shows help menu

flags:
- `--webhook/-w`    - specify the webhook url to use
- `--silent/-s`     - run without firing the webhook
- `--debug/-d`      - logs additional information
- `logfile/-f`      - enables logging to file
- `--help/-h`       - shows help menu
- `--version/-v`    - shows the binary version information 

examples:
- Run with specified webhook: `BungieAlerter -w <your webhook> go`
- Run without webhook notifications: `BungieAlerter -s go`


### Docker
There's also docker images built with compressed binaries for ease of deployment.

These are literally just the binaries dropped into an alpine base image. You'll need to modify the entrypoint most likely, give it a webhook somehow, and then be off to the races

## Future
There's a lot of refinement and features I'd like to add over time. There's also the potential of making this more of a generic application, allowing users to add only the channels/accounts/websites that they personally want notifications for; who knows.

I'll try to keep the [issues](https://github.com/OverlyDev/go-bungie-alerter/issues) up to date with any upcoming features.

Additionally, if you encounter a bug of any kind, feel free to submit an issue and I'll see about getting it sorted.
