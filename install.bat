@ECHO OFF

SETLOCAL

SET APP_NAME=dnsupdater
SET PACKAGE_NAME=github.com/boris1993/%APP_NAME%

IF EXIST install.bat GOTO check-option
ECHO install.bat must be run from its folder
GOTO end

:check-option
SET OPTION=%1

IF "%OPTION%" EQU "-h" GOTO help
IF "%OPTION%" EQU "--help" GOTO help
IF "%OPTION%" EQU "/?" GOTO help
GOTO :check-arch

:help
ECHO Usage: install.bat TARGET
ECHO TARGET could be windows-amd64, mips-softfloat, linux-amd64, darwin-amd64
ECHO Default target is mips-softfloat
GOTO exit

:check-arch
IF "%1"=="windows-amd64" GOTO build-windows-amd64
IF "%1"=="linux-amd64" GOTO build-linux-amd64
IF "%1"=="darwin-amd64" GOTO build-darwin-amd64
IF "%1"=="mips-softfloat" GOTO mips-softfloat
GOTO build-mips-softfloat

:build-windows-amd64
SET GOARCH=amd64
SET GOOS=windows
GOTO do-build-windows

:build-linux-amd64
SET GOARCH=amd64
SET GOOS=linux
GOTO do-build

:build-darwin-amd64
SET GOARCH=amd64
SET GOOS=darwin
GOTO do-build

:build-mips-softfloat
SET GOARCH=mips
SET GOOS=linux
SET GOMIPS=softfloat
GOTO do-build-mips

:do-build-windows
ECHO Building binary for Windows running under amd64

IF NOT EXIST bin\%APP_NAME%-%GOOS%-%GOARCH% (
    mkdir bin\%APP_NAME%-%GOOS%-%GOARCH%
)

go build -i -o bin\%APP_NAME%-%GOOS%-%GOARCH%\%APP_NAME%.exe

ECHO Copying template config file to target directory...
copy config.yaml.template bin\%APP_NAME%-%GOOS%-%GOARCH%

IF errorlevel 1 GOTO error
GOTO success

:do-build-mips
ECHO Building binary for Linux running under %GOARCH%-%GOMIPS%

IF NOT EXIST bin\%APP_NAME%-%GOOS%-%GOARCH%-%GOMIPS% (
    mkdir bin\%APP_NAME%-%GOOS%-%GOARCH%-%GOMIPS%
)

go build -i -o bin\%APP_NAME%-%GOOS%-%GOARCH%-%GOMIPS%\%APP_NAME%

ECHO Copying template config file to target directory...
copy config.yaml.template bin\%APP_NAME%-%GOOS%-%GOARCH%-%GOMIPS%

IF errorlevel 1 GOTO error
GOTO success

:do-build
ECHO Building binary for %GOOS% running under %GOARCH%

IF NOT EXIST bin\%APP_NAME%-%GOOS%-%GOARCH% (
    mkdir bin\%APP_NAME%-%GOOS%-%GOARCH%
)

go build -i -o bin\%APP_NAME%-%GOOS%-%GOARCH%\%APP_NAME%

ECHO Copying template config file to target directory...
copy config.yaml.template bin\%APP_NAME%-%GOOS%-%GOARCH%

IF errorlevel 1 GOTO error
GOTO success

:success
ECHO Build finished
GOTO exit

:error
ECHO Build failed

:exit
