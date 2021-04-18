# VU Mod Manager CLI
<a href="https://github.com/BF3RM/vumm-cli/actions/workflows/release.yml">
  <img src="https://img.shields.io/github/workflow/status/BF3RM/vumm-cli/goreleaser" alt="build status">
</a>
<a href="https://github.com/BF3RM/vumm-cli/releases">
  <img src="https://img.shields.io/github/release/BF3RM/vumm-cli.svg" alt="releases">
</a>
<a href="https://sonarcloud.io/dashboard?id=BF3RM_vumm-cli">
  <img src="https://sonarcloud.io/api/project_badges/measure?project=BF3RM_vumm-cli&metric=alert_status" alt="quality gate">
</a>
<a href="https://sonarcloud.io/component_measures?id=BF3RM_vumm-cli&metric=Security">
  <img src="https://sonarcloud.io/api/project_badges/measure?project=BF3RM_vumm-cli&metric=security_rating" alt="security">
</a>
<a href="https://sonarcloud.io/component_measures?id=BF3RM_vumm-cli&metric=Maintainability">
  <img src="https://sonarcloud.io/api/project_badges/measure?project=BF3RM_vumm-cli&metric=sqale_rating" alt="maintainability">
</a>
Automatically install the latest versions of your favourite Venice Unleashed mods without having to worry about mod dependencies.

Venice Unleashed Mod Manager is a tool server owners can use to automatically install mods and their dependencies.
It also allows developers to easily distribute new versions and don't have to worry about shipping the proper dependencies.\
Furthermore the tool automatically checks compatibility between your installed mods and will warn you when something is off.

## Installation
Download the latest version of vumm  [here](https://github.com/BF3RM/vumm-cli/releases/latest).
Extract the executable file in your `Battlefield 3/Server/Admin` folder.

## Usage
VUMM has two purposes, to simplify the distribution of mods and to simplify the installation of mods. Below you can find different type of commands that you can use.

### Manage installed mods
Installing a mod is as simple as running:
```bash
vumm install <mod>
```
So the following command will install the latest version of [BlueprintManager](https://github.com/BF3RM/BlueprintManager):
```bash
vumm install blueprintmanager
```
You can also specify a specific release tag or [semver constraint](https://docs.npmjs.com/about-semantic-versioning):
```bash
vumm install blueprintmanager@^1.0.0  # install the latest minor version of 1.x
vumm install blueprintmanager@~1.0.0  # install the latest patch version of 1.0.x
vumm install blueprintmanager@dev     # install the latest dev version of blueprintmanager
```

### Authenticating
If you are a mod creator, or trying to install a private mod you have to be authenticated.
This can be done by using the login and register commands of vumm. There are two access levels, publish and readonly.
Publish means the current session is allowed to manage the user and it's mods.
Readonly means the current session is only allowed to read information from the registry.
The readonly type of authentication is mainly what you want on your server, whereas the publish access level is what you want if you are a mod creator.
```bash
vumm login [--type <publish|readonly>]
```

Registering has the exact same syntax. If you registered successfully it's not needed to login as the access token will already be stored.
```bash
vumm register [--type <publish|readonly>]
```

### Publish a mod
As a mod creator you can simply publish your mod by running the following command in your mod source folder
```bash
vumm publish [-t <tag>] [--private]
```
It will read the [mod.json](https://docs.veniceunleashed.net/modding/your-first-mod/#the-modjson-file) to resolve the version, and it's dependencies automatically.\
By default your mod will be published to the `latest` tag, when specifying the `-t` flag you can specify another tag like `dev` or `beta`.\
To use this command you have to be logged in with publish rights.

#### Granting access
If you decide to publish a private mod, or willing to add multiple contributors to your mod the `grant` and `revoke` commands can be used.
```bash
vumm grant <mod> <user> <publish|readonly>  # give someone publish or read rights over a mod
vumm revoke <mod> <user>                    # revoke all permissions of user over a mod
```
If your mod is private you can use the grant command to give people read rights over the mod like so:
```bash
vumm grant realitymod paulhobbel readonly
```

## Updating VUMM
VUMM has an automatic updater built in which checks on a weekly basis for updates. If an update is available it will notify you about it.
You can always manually check for updates by running:
```bash
vumm update
```
To see the currently installed version of VUMM you can run:
```bash
vumm --version
```

## License
The Venice Unleashed Mod Manager CLI is available under the MIT license. See the [LICENSE](./LICENSE) file for more info.