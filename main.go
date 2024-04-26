package main

import (
	"fmt"
	"os"
	"time"

	"github.com/denverdino/aliyungo/cdn"
)

func checkEnv(accessKeyId, accessKeySecret, cdnDomain, certPemPath, certKeyPath string) {
	var errList []string
	if accessKeyId == "" {
		errList = append(errList, "ALI_ACCESS_KEY_ID")
	}
	if accessKeySecret == "" {
		errList = append(errList, "ALI_ACCESS_KEY_SECRET")
	}
	if cdnDomain == "" {
		errList = append(errList, "CDN_DOMAIN")
	}
	if certPemPath == "" {
		errList = append(errList, "CERT_PEM_PATH")
	}
	if certKeyPath == "" {
		errList = append(errList, "CERT_KEY_PATH")
	}

	if len(errList) > 0 {
		fmt.Println("err: 检测到以下环境变量未配置")
		for _, value := range errList {
			fmt.Printf("%v \n", value)
		}
		os.Exit(-1)
	}

}
func main() {
	ACCESS_KEY_ID := os.Getenv("ALI_ACCESS_KEY_ID")
	ACCESS_KEY_SECRET := os.Getenv("ALI_ACCESS_KEY_SECRET")
	cdnDomain := os.Getenv("CDN_DOMAIN")
	certPemPath := os.Getenv("CERT_PEM_PATH")
	certKeyPath := os.Getenv("CERT_KEY_PATH")

	checkEnv(ACCESS_KEY_ID, ACCESS_KEY_SECRET, cdnDomain, certPemPath, certKeyPath)
	var cert []byte
	var key []byte
	var err error
	if cert, err = os.ReadFile(certPemPath); err != nil {
		panic(err)
	}
	if key, err = os.ReadFile(certKeyPath); err != nil {
		panic(err)
	}

	// 生成一个不重复的证书名称
	var certName = "cert" + time.Now().Format("20060102150405")

	// 记录日志
	fmt.Println("time: ", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("update cert domain: ", cdnDomain)
	fmt.Println("certName: ", certName)

	client := cdn.NewClient(ACCESS_KEY_ID, ACCESS_KEY_SECRET)
	res, err := client.SetDomainServerCertificate(cdn.CertificateRequest{
		DomainName:              cdnDomain,
		CertName:                certName,
		ServerCertificateStatus: "on",
		ServerCertificate:       string(cert),
		PrivateKey:              string(key),
	})

	fmt.Printf("res: %v, err: %v\n", res, err)
	/*
		打包Linux二进制
		git bash:
			$ export GOOS=linux
			$ export GOARCH=amd64
			$ go build -o aliyun-cdn-cert-bot main.go

		命令行设置临时环境变量如：
			Linux: $ export CDN_DOMAIN="example.com"
			Windows: $ set CDN_DOMAIN="example.com"
	*/
}
