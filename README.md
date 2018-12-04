# dnsupdater
|master| dev |
|:----:|:---:|
|[![Build Status](https://travis-ci.org/boris1993/dnsupdater.svg?branch=master)](https://travis-ci.org/boris1993/dnsupdater)|[![Build Status](https://travis-ci.org/boris1993/dnsupdater.svg?branch=dev)](https://travis-ci.org/boris1993/dnsupdater)|

Obtain your current external IP address and update to the specified DNS record on CloudFlare 

Primarily built for MIPS 74kc since my router has a MIPS 74kc CPU

## How-to

### Using pre-built binaries

+ Download the [latest release](https://github.com/boris1993/dnsupdater/releases/latest) for your target

+ Extract the archive

+ Rename `config.yaml.template` to `config.yaml`

+ Finish your configuration in the `config.yaml`

+ Upload `dnsupdater` and `config.yaml` to the device you want this app to run

+ Set up a cron job like

```cron
0 0,12 * * * /root/dnsupdater/dnsupdater > /var/log/update-dns.log 2>&1 &
```

### Build from source

+ Install Go >= 11.1

+ Get this repo

```bash
git clone https://github.com/boris1993/dnsupdater.git
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

Then you will find the executable file under the `bin/dnsupdater-linux-mips-softfloat` directory. 

+ Finish the configuration

Rename `config.yaml.template` to `config.yaml` and finish your configuration. 

+ Upload to your router

Upload `dnsupdater` and `config.yaml` to your router.

And don't forget to give it executable permission.

+ Create a cron job

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
