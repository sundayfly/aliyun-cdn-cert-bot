#!/bin/bash
# https://github.com/denverdino/aliyungo
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin
export ALI_ACCESS_KEY_ID="xxxxxxx"
export ALI_ACCESS_KEY_SECRET="xxxxxxxxxxxxxxx" 

conf_file="/data/shell/conf/aliyun_cdn_domain.txt"
acme_path="/data/ssl/acme"

while read line;do
    domain=$line;
    if [[ -z "$domain" ]] || [[ $domain == \#* ]]; then
        continue;
    fi

    acme_domain="${domain#*.}" # 删除第一个.前字符
    acme_pem_path="$acme_path/$acme_domain.pem"
    acme_key_path="$acme_path/$acme_domain.key"

    if [[ ! -e "$acme_pem_path" ]] || [[ ! -e "$acme_key_path" ]]; then
        echo "文件路径不存在: $acme_pem_path" 
        echo "文件路径不存在: $acme_key_path"
        continue;
    fi

    # 判断证书日期是否相同
    cdn_domain_url="https://$domain"
    cdn_domain_expire_gmt_date=$(curl -Ivs https://${domain}"/ssl_check" --connect-timeout 7 -o /dev/null 2>&1 | grep "expire date" | awk -F ": " '{print $NF}')
    acme_pem_expire_gmt_date="$(openssl x509 -in $acme_pem_path -noout -enddate | awk -F '=' '{print $NF}')"
    cdn_domain_expire_date=$(date -d "$cdn_domain_expire_gmt_date - 8 hours" "+%Y%m%d%H%M%S")
    acme_pem_expire_date=$(date -d "$acme_pem_expire_gmt_date - 8 hours" "+%Y%m%d%H%M%S")

    if [[ $cdn_domain_expire_date == $acme_pem_expire_date ]]; then
        echo "$domain 证书过期日期相同,跳过"
        continue;
    fi

    # 使用AliyunGo SDK更新CDN证书
    export CDN_DOMAIN=$domain 
    export CERT_PEM_PATH=$acme_pem_path
    export CERT_KEY_PATH=$acme_key_path
    output=$(/usr/local/bin/aliyun-cdn-cert-bot)
    echo $output | grep 'err: <nil>' && status=1 || status=0
    
    if [[ $status -eq 1 ]]; then
        echo "update cdn domain cert success."
    else
        echo "update cdn domain cert failed."
    fi
done < "$conf_file"
