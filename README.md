# AutoTomatoPay
TomatoPay自动提现

## 用法

安装golang
```
CentOS:
yum install go -y

Debian/Ubuntu:
apt-get install golang -y
```

下载源代码
```
curl -O https://raw.githubusercontent.com/TheCGDF/AutoTomatoPay/master/AutoTomatoPay.go
```

自行修改源代码中的以下内容
```
var email = ""         //TomatoPay邮箱
var password = ""      //TomatoPay密码
var threshold = 100.00 //提现阈值，单位：元
```

使用`crontab -e`添加定时提现任务

例：每日几时几分申请提现
```
分 时 * * * go run 刚刚下载的文件路径/AutoTomatoPay.go
```

>根据[~~可靠~~消息](https://t.me/fanqiepay/9539)，每天（下午？）10点处理提现，建议crontab设置时间略早于该时间

加入crontab前建议先手动执行一遍，错误信息会直接输出并且保存在`AutoTomatoPay.log`中

## 鸣谢

[regendsoh](https://github.com/regendsoh)的帮助与贡献
