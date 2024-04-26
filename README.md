
使用非官方的aliyun go sdk 定时更新CDN证书。

在项目所在文件夹中安装依赖：
```go
mkdir aliyun-cdn-cert-bot && cd aliyun-cdn-cert-bot
go mod init sundayhk.com/aliyun-cdn-cert-bot
go get "github.com/denverdino/aliyungo/cdn"
```
编译 main.go 即可使用：
```go
go build -o aliyun-cdn-cert-bot main.go
```

运行之前请设置下列环境变量：   
- ACCESS_KEY_ID、ACCESS_KEY_SECRET 为阿里云有权限的 RAM 子账号信息  
- CDN_DOMAIN 为阿里云 CDN 域名（非源站域名）  
- CERT_PEM_PATH 为 PEM证书文件路径（注意要用 fullchain 的证书，否则可能有些客户端会报错）  
- CERT_KEY_PATH 为 KEY证书密钥路径  

acme 安装及泛域名申请
```sh
curl https://get.acme.sh | sh -s email=my@example.com
domain=sundayhk.com
mkdir -p /data/ssl/acme
acme.sh --issue --dns dns_dp -d $domain -d *.$domain --keylength ec-256
acme.sh --install-cert -d $domain --ecc --key-file "/data/ssl/acme/$domain.key" --fullchain-file "/data/ssl/acme/$domain.pem" --reloadcmd "systemctl reload nginx"
```

crontab 设置定时任务，每月执行一次
```sh
$ crontab -e
15 11 * */1 * /data/shell/aliyun_cdn_cert_bot.sh > /dev/null 2>&1
```

命令行设置临时环境变量如：
```sh
  Linux: $ export CDN_DOMAIN="example.com"
  Windows: $ set CDN_DOMAIN="example.com"
```
Windows打包Linux二进制
```sh
git bash:
  $ export GOOS=linux
  $ export GOARCH=amd64
  $ go build -o aliyun-cdn-cert-bot main.go
```
