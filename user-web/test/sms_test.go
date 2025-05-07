/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-07 12:11:20
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-07 12:19:47
 * @FilePath: /joyshop_api/user-web/test/sms_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package test

import (
	"joyshop_api/user-web/api"
	"testing"
)

func TestSendSms(t *testing.T) {
	// 调用发送短信方法
	err := api.SendSms()
	if err != nil {
		t.Errorf("发送短信失败: %v", err)
		return
	}
	t.Log("发送短信成功")
}
