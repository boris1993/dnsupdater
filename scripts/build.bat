@ECHO OFF

SETLOCAL

SET ROOT_DIR=%~dp0\..
SET APP_NAME=dnsupdater
SET PACKAGE_NAME=github.com/boris1993/%APP_NAME%

IF EXIST build.bat GOTO check-option
ECHO build.bat must be run from its folder
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
ECHO Default target is linux-amd64
GOTO exit

:check-arch
IF "%1"=="windows-amd64" GOTO build-windows-amd64
IF "%1"=="linux-amd64" GOTO build-linux-amd64
IF "%1"=="darwin-amd64" GOTO build-darwin-amd64
IF "%1"=="mips-softfloat" GOTO build-mips-softfloat
GOTO build-linux-amd64

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

IF NOT EXIST %ROOT_DIR%\bin\%APP_NAME%-%GOOS%-%GOARCH% (
    mkdir %ROOT_DIR%\bin\%APP_NAME%-%GOOS%-%GOARCH%
)

go build -i -o %ROOT_DIR%\bin\%APP_NAME%-%GOOS%-%GOARCH%\%APP_NAME%.exe %ROOT_DIR%\cmd\dnsupdater\main.go

ECHO Copying template config file to target directory...
copy %ROOT_DIR%\configs\config.yaml.template %ROOT_DIR%\bin\%APP_NAME%-%GOOS%-%GOARCH% > nul 2>&1

IF errorlevel 1 GOTO error
GOTO success

:do-build-mips
ECHO Building binary for Linux running under %GOARCH%-%GOMIPS%

IF NOT EXIST %ROOT_DIR%\bin\%APP_NAME%-%GOOS%-%GOARCH%-%GOMIPS% (
    mkdir %ROOT_DIR%\bin\%APP_NAME%-%GOOS%-%GOARCH%-%GOMIPS%
)

go build -i -o %ROOT_DIR%\bin\%APP_NAME%-%GOOS%-%GOARCH%-%GOMIPS%\%APP_NAME% %ROOT_DIR%\cmd\dnsupdater\main.go

ECHO Copying template config file to target directory...
copy %ROOT_DIR%\configs\config.yaml.template %ROOT_DIR%\bin\%APP_NAME%-%GOOS%-%GOARCH%-%GOMIPS% > nul 2>&1

IF errorlevel 1 GOTO error
GOTO success

:do-build
ECHO Building binary for %GOOS% running under %GOARCH%

IF NOT EXIST %ROOT_DIR%\bin\%APP_NAME%-%GOOS%-%GOARCH% (
    mkdir %ROOT_DIR%\bin\%APP_NAME%-%GOOS%-%GOARCH%
)

go build -i -o %ROOT_DIR%\bin\%APP_NAME%-%GOOS%-%GOARCH%\%APP_NAME% %ROOT_DIR%\cmd\dnsupdater\main.go

ECHO Copying template config file to target directory...
copy %ROOT_DIR%\configs\config.yaml.template %ROOT_DIR%\bin\%APP_NAME%-%GOOS%-%GOARCH% > nul 2>&1

IF errorlevel 1 GOTO error
GOTO success

:success
ECHO Build finished
GOTO exit

:error
ECHO Build failed

:exit
