# dnsupdater
[![Build Status](https://travis-ci.com/boris1993/dnsupdater.svg?branch=master)](https://travis-ci.com/boris1993/dnsupdater)
![GitHub](https://img.shields.io/github/license/boris1993/dnsupdater)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/boris1993/dnsupdater)](https://github.com/boris1993/dnsupdater/releases/latest)
![Total download](https://img.shields.io/github/downloads/boris1993/dnsupdater/total.svg)

[English](README.md)

本工具可以获取你当前的工网IP地址，并将其更新到指定的DNS记录。

建议在你的家庭服务器，或者在你的路由器中运行这个工具。

你**绝不应该**在代理服务器或VPN后面运行这个应用。这种场景是未被考虑且未经测试的。

## 如何使用

+ 前往 [下载页面](https://github.com/boris1993/dnsupdater/releases/latest) 为你的目标平台下载最新的版本。

+ 将下载到的压缩包解压。

+ 将`config.yaml.template`重命名为`config.yaml`。

+ 完成`config.yaml`中的配置。

+ 将`dnsupdater`和`config.yaml`上传到你想要运行这个应用的地方。注意，这两个文件必须位于相同目录下。

+ 设置一个定时任务(cron)，比如

```cron
0 0,12 * * * /home/yourname/dnsupdater/dnsupdater > /var/log/update-dns.log 2>&1 &
```

## 配置要点

+ CloudFlare配置中的`APIKey`必须是一个单独的API令牌。
你可以在 [这里](https://dash.cloudflare.com/profile/api-tokens) 通过套用`编辑区域 DNS`这个模版来生成。

+ 切勿修改阿里云配置中的`RegionID`。目前阿里云仅接受`cn-hangzhou`这一个值。

+ 如果你想要禁用更新IPv6记录这个功能，那么你可以删掉配置文件中的`IPv6AddrAPI`属性。
如果你的网络不支持IPv6，那么你也需要删掉这个配置。否则你将会在日志中看到一条警告信息。

## 为其他平台构建可执行文件

你可以通过运行 `scripts` 文件夹中的脚本文件来查看预设的目标平台。

对于Windows用户:

```cmd
build.bat /?
```

对于*NIX用户:

```bash
make help
```

如果你的目标平台不在预设的列表，那么你可以通过`go build`命令，
并手动指定环境变量`GOARCH`、`GOOS`（如果有必要的话，还可以指定`GOMIPS`）来为你的目标平台构建可执行文件。

## 许可协议

本程序依据[MIT](LICENSE)协议开源。
