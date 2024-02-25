package ioc

import (
	"os"
	"webook/internal/service/sms"
	"webook/internal/service/sms/localsms"
	"webook/internal/service/sms/tencent"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func InitSMSService() sms.Service {
	return localsms.NewService()
	// return InitTencentSMSService()
}

func InitTencentSMSService() sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("找不到腾讯SMS的 secret id")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("找不到腾讯SMS的 secret key")
	}
	client, err := tencentSMS.NewClient(
		common.NewCredential(secretId, secretKey),
		"np-nanjing",
		profile.NewClientProfile())
	if err != nil {
		panic("没有找到环境变量 SMS_SECRET_KEY")
	}
	return tencent.NewService(client, "", "")
}
