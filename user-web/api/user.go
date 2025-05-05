package api

import (
	"fmt"
	"joyshop_api/user-web/global"
	"joyshop_api/user-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	// 将grpc错误转换为http错误状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			// grpc错误转换为http错误状态码
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
				return
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Message(),
				})
				return
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
				return

			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "未知错误" + e.Message(),
				})
				return
			}
		}

	}
}

func GetUserList(ctx *gin.Context) {
	zap.S().Debug("获取用户列表页数据")
	//ip := "127.0.0.1"
	//port := 50051
	//拨号连接用户grpc服务
	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) // ✅ 新方式)
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】", "msg", err.Error())
		return
	}
	// 生成grpc客户端
	userSrvClient := proto.NewUserClient(userConn)
	// 调用grpc服务
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("pageSize", "10")
	pageInt, err := strconv.Atoi(page)

	if err != nil {
		zap.S().Errorw("[GetUserList] page 转换失败", "msg", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}
	paseSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		zap.S().Errorw("[GetUserList] pageSize 转换失败", "msg", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	userResp, err := userSrvClient.GetUserList(ctx, &proto.PageInfo{
		Page:     uint32(pageInt),
		PageSize: uint32(paseSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 调用 【用户服务失败】", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, user := range userResp.Data {
		//result = append(result, user)
		data := make(map[string]interface{})
		data["id"] = user.Id
		data["nickName"] = user.Nickname
		data["username"] = user.Username
		data["mobile"] = user.Mobile
		data["email"] = user.Email
		data["avatar"] = user.Avatar
		data["gender"] = user.Gender
		data["birthday"] = user.Birthday
		data["role"] = user.Role
		result = append(result, data)

	}
	// 返回数据
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "success",
		"total": userResp.Total,
		"data":  result,
	})
}

func CreateUser(context *gin.Context) {

}
func UpdateUser(context *gin.Context) {

}
func DeleteUser(context *gin.Context) {

}

func PassWordLogin(ctx *gin.Context) {
	// 获取请求参数
	var request struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 连接用户服务
	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[PassWordLogin] 连接 【用户服务失败】", "msg", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务器内部错误",
		})
		return
	}
	defer userConn.Close()

	// 生成grpc客户端
	userSrvClient := proto.NewUserClient(userConn)

	// 1. 先通过手机号获取用户信息
	userInfo, err := userSrvClient.GetUserByMobile(ctx, &proto.MobileRequest{
		Mobile: request.Mobile,
	})
	if err != nil {
		zap.S().Errorw("[PassWordLogin] 获取用户信息失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 2. 验证密码
	checkResp, err := userSrvClient.CheckPassword(ctx, &proto.PasswordCheckInof{
		Password:        request.Password,
		EncryptPassword: userInfo.Password,
	})
	if err != nil {
		zap.S().Errorw("[PassWordLogin] 验证密码失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	if !checkResp.Success {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "密码错误",
		})
		return
	}

	// 3. 返回用户信息
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "登录成功",
		"data": gin.H{
			"id":       userInfo.Id,
			"nickName": userInfo.Nickname,
			"mobile":   userInfo.Mobile,
			"email":    userInfo.Email,
			"avatar":   userInfo.Avatar,
			"gender":   userInfo.Gender,
			"birthday": userInfo.Birthday,
			"role":     userInfo.Role,
		},
	})
}

func SendSms(context *gin.Context) {

}
