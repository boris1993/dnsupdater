# dnsupdater

[![Build Status](https://travis-ci.org/boris1993/dnsupdater.svg?branch=master)](https://travis-ci.org/boris1993/dnsupdater)

Obtain your current external IP address and update to the specified DNS record on CloudFlare 

Primarily built for MIPS 74kc since my router has a MIPS 74kc CPU

## How-to

### Download pre-built binaries

+ Go to **Releases** and download the binary for your target

+ Extract the archive

+ Rename `config.yaml.template` to `config.yaml`

+ Replace your configuration in the `config.yaml`

+ Upload `dnsupdater` and `config.yaml` to where you want this app to run

+ Set up a cron job like

```cron
0 0,12 * * * /root/dnsupdater/dnsupdater > /var/log/update-dns.log 2>&1 &
```

### Build from source

+ Install Go

+ Get this repo

```bash
go get github.com/boris1993/dnsupdater
```

+ Build for MIPS 74kc

For Windows users:
 
```cmd
install.bat
```

For *nix users:

```bash
make mips-softfloat
```

Then you will find the executable file under `${GOPATH}/bin/dnsupdater` directory. 

+ Upload to your router

Just upload the executable to your router via FTP or SFTP.

And don't forget to give it execute permission.

+ Create a cron job

Sure you don't want this to be a disposable product right?

```crontab
0 0,12 * * * /root/dnsupdater/dnsupdater > /var/log/update-dns.log 2>&1 &
```

## Build for other platforms

You can check for all preset targets by

```cmd
install.bat /?
```

```bash
make help
```

Or you can also specify your own `GOARCH` and `GOOS` (and maybe `GOMIPS`) to build for your platform 
as long as Go provides support to it.  
