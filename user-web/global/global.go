/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-05 16:17:40
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-07 18:02:34
 * @FilePath: /joyshop_api/user-web/global/global.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package global

import (
	"joyshop_api/user-web/config"
	"joyshop_api/user-web/proto"

	ut "github.com/go-playground/universal-translator"
	"google.golang.org/grpc"
)

var (
	ServerConfig *config.ServeConfig
	Trans        ut.Translator
	UserClient   proto.UserClient
	UserConn     *grpc.ClientConn
)
