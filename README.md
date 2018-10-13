# dnsupdater

Obtain your current external IP address and update to the specified DNS record on CloudFlare 

Primarily built for MIPS 74kc

## Usage

+ Install Go

+ Clone this repo

```bash
git clone https://github.com/boris1993/dnsupdater.git
```

+ Customize your configuration

```bash
cd dnsupdater/src/config
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

+ Build for other platforms

You can check for all preset targets by

```cmd
install.bat /?
```

```bash
make help
```

Or you can also specify your own `GOARCH` and `GOOS` (and maybe `GOMIPS`) to build for your platform 
as long as Go provides support to it.  
