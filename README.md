# dnsupdater
[![Build Status](https://travis-ci.com/boris1993/dnsupdater.svg?branch=master)](https://travis-ci.com/boris1993/dnsupdater)
![GitHub](https://img.shields.io/github/license/boris1993/dnsupdater)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/boris1993/dnsupdater)](https://github.com/boris1993/dnsupdater/releases/latest)
![Total download](https://img.shields.io/github/downloads/boris1993/dnsupdater/total.svg)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=boris1993_dnsupdater&metric=alert_status)](https://sonarcloud.io/dashboard?id=boris1993_dnsupdater)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fboris1993%2Fdnsupdater.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fboris1993%2Fdnsupdater?ref=badge_shield)

This app allows you updating your DNS records with your current external IP address.

## How-to

+ Download the [latest release](https://github.com/boris1993/dnsupdater/releases/latest) for your target

+ Extract the archive.

+ Rename `config.yaml.template` to `config.yaml`.

+ Finish your configuration in the `config.yaml`

+ Upload `dnsupdater` and `config.yaml` to the device you want this app to run. These 2 files must be under the same directory.

+ Set up a cron job like

```cron
0 0,12 * * * /home/yourname/dnsupdater/dnsupdater > /var/log/update-dns.log 2>&1 &
```

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
