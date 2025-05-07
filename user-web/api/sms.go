/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-07 11:24:51
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-07 13:56:20
 * @FilePath: /joyshop_api/user-web/api/sms.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"joyshop_api/user-web/utils"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"joyshop_api/user-web/forms"
	"joyshop_api/user-web/global"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// Description:
// 使用凭据初始化账号Client
// @return Client
// @throws Exception
func CreateClient() (*dysmsapi20170525.Client, error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: tea.String(global.ServerConfig.AliyunSms.AccessKeyId),
		// 您的AccessKey Secret
		AccessKeySecret: tea.String(global.ServerConfig.AliyunSms.AccessSecret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	client := &dysmsapi20170525.Client{}
	client, err := dysmsapi20170525.NewClient(config)
	return client, err
}

// GenerateRandomCode 生成指定长度的随机数字验证码
// length: 验证码长度，默认为4位
func GenerateRandomCode(length ...int) string {
	// 设置默认长度为4
	codeLength := 4
	if len(length) > 0 && length[0] > 0 {
		codeLength = length[0]
	}

	// 使用新的随机数生成方式
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成随机数字
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = byte(r.Intn(10)) + '0'
	}

	return string(code)
}

// SendSms 发送短信
// phone: 接收短信的手机号，如果为空则使用配置文件中的默认值
func SendSms(phone ...string) error {
	client, err := CreateClient()
	if err != nil {
		return fmt.Errorf("创建短信客户端失败: %v", err)
	}

	// 生成6位随机验证码
	code := GenerateRandomCode(6)

	// 确定使用哪个手机号
	phoneNumber := global.ServerConfig.AliyunSms.PhoneNumbers
	if len(phone) > 0 && phone[0] != "" {
		phoneNumber = phone[0]
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(global.ServerConfig.AliyunSms.SignName),
		TemplateCode:  tea.String(global.ServerConfig.AliyunSms.TemplateCode),
		PhoneNumbers:  tea.String(phoneNumber),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%s\"}", code)),
	}
	runtime := &util.RuntimeOptions{}

	// 发送短信
	_, err = client.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil {
		var sdkErr *tea.SDKError
		if errors.As(err, &sdkErr) {
			zap.S().Errorf("短信发送失败: %v", tea.StringValue(sdkErr.Message))
			if sdkErr.Data != nil {
				var data interface{}
				if err := json.NewDecoder(strings.NewReader(tea.StringValue(sdkErr.Data))).Decode(&data); err == nil {
					if m, ok := data.(map[string]interface{}); ok {
						if recommend, ok := m["Recommend"]; ok {
							zap.S().Errorf("短信发送失败建议: %v", recommend)
						}
					}
				}
			}
			return fmt.Errorf("短信发送失败: %v", tea.StringValue(sdkErr.Message))
		}
		return fmt.Errorf("短信发送失败: %v", err)
	}

	// 短信发送成功，将验证码存入Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: "", // 如果有密码，在这里设置
		DB:       0,  // 使用默认DB
	})

	// 设置验证码，过期时间从配置文件读取
	if err := rdb.Set(context.Background(), phoneNumber, code, time.Duration(global.ServerConfig.RedisInfo.ExpireTime)*time.Minute).Err(); err != nil {
		zap.S().Errorf("Redis存储验证码失败: %v", err)
		return fmt.Errorf("验证码存储失败: %v", err)
	}

	zap.S().Infof("短信发送成功: phone=%s, code=%s", phoneNumber, code)
	return nil
}

// SendSmsHandler 发送短信的HTTP处理函数
func SendSmsHandler(c *gin.Context) {

	// 使用表单验证器验证手机号
	form := forms.SendSmsForm{}
	if err := c.ShouldBind(&form); err != nil {
		if utils.HandleValidatorError(c, err, "SendSmsForm") {
			return
		}
	}

	// 发送短信
	if err := SendSms(form.Mobile); err != nil {
		zap.S().Errorf("发送短信失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "发送成功",
	})
}
