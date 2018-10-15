# dnsupdater

Obtain your current external IP address and update to the specified DNS record on CloudFlare 

Primarily built for MIPS 74kc since my router has a MIPS 74kc CPU

## Usage

+ Install Go

+ Clone this repo

```bash
go get github.com/boris1993/dnsupdater
```

+ Customize your configuration

```bash
cd github.com/boris1993/dnsupdater/config
cp config.go.template config.go
vim config.go # Or use any text editor you like
```

+ Build for MIPS 74kc

For Windows users:
 
```cmd
install.bat
```

For *nix users:

```bash
make
```

Then you will find the executable file under `bin` directory. 

+ Upload to your router

Just upload the executable to your router via FTP or SFTP.

And don't forget to give it execute permission.

+ Create a cron job

Sure you don't want this to be a disposable product right?

```crontab
0 0,12 * * * /root/dnsupdater >> /var/log/update-dns.log
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
