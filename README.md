# dnsupdater
[![Build](https://github.com/boris1993/dnsupdater/actions/workflows/build.yml/badge.svg)](https://github.com/boris1993/dnsupdater/actions/workflows/build.yml)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/boris1993/dnsupdater)](https://github.com/boris1993/dnsupdater/releases/latest)
![Total download](https://img.shields.io/github/downloads/boris1993/dnsupdater/total.svg)

[中文版](README_zh_cn.md)

This app allows you updating your DNS records with your current external IP address.

It is recommended to run this program in your home server, or in your router。

You should **NEVER** run this program behind a proxy or a VPN. 
Running it behind a proxy is an unconsidered and untested scenario.

## How-to

+ Download the [latest release](https://github.com/boris1993/dnsupdater/releases/latest) for your target

+ Extract the archive.

+ Rename `config.yaml.template` to `config.yaml`.

+ Finish your configuration in the `config.yaml`

+ Upload `dnsupdater` and `config.yaml` to the device you want this app to run. 
These 2 files must be under the same directory.

+ Set up a cron job like

```cron
0 0,12 * * * /home/yourname/dnsupdater/dnsupdater > /var/log/update-dns.log 2>&1 &
```

## Important notes in configuration

+ The `APIKey` for your CloudFlare records should be a dedicated API token. 
You can generate one [here](https://dash.cloudflare.com/profile/api-tokens) with template `Edit zone DNS`.

+ Do not modify the property `RegionID` for your Aliyun DNS records. `cn-hangzhou` is the only accepted value for now. 

+ About JSON path

Here's a list of operators used in JSON path:

| Operator                  | Description                                                     |
|:--------------------------|:----------------------------------------------------------------|
| `$`                       | The root element to query. This starts all path expressions.    |
| `@`                       | The current node being processed by a filter predicate.         |
| `*`                       | Wildcard. Available anywhere a name or numeric are required.    |
| `..`                      | Deep scan. Available anywhere a name is required.               |
| `.<name>`                 | Dot-notated child                                               |
| `['<name>' (, '<name>')]` | Bracket-notated child or children                               |
| `[<number> (, <number>)]` | Array index or indexes                                          |
| `[start:end]`             | Array slice operator                                            |
| `[?(<expression>)]`       | Filter expression. Expression must evaluate to a boolean value. |

So if you have a JSON like this:

```json

{
  "ip": "103.156.184.21",
  "tz": "Asia/Taipei"
}
```

You can use `$.ip` to obtain the value in the field `ip`.

## Build for other platforms

You can check for all preset targets by running the scripts in the `scripts` folder.

For Windows users:

```cmd
build.bat /?
```

For *NIX users:

```bash
make help
```

Or you can specify your own `GOARCH` and `GOOS` (and maybe `GOMIPS`) with `go build` command 
to build the executable for your platform as long as Go provides support to it.  

## License

Licensed under [MIT](LICENSE) license.
