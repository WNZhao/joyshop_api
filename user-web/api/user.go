package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"joyshop_api/user-web/proto"
	"net/http"
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
	ip := "127.0.0.1"
	port := 50051
	//拨号连接用户grpc服务
	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials())) // ✅ 新方式)
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】", "msg", err.Error())
		return
	}
	// 生成grpc客户端
	userSrvClient := proto.NewUserClient(userConn)
	// 调用grpc服务
	userResp, err := userSrvClient.GetUserList(ctx, &proto.PageInfo{})
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
		"code": 0,
		"msg":  "success",
		"data": result,
	})
}

func CreateUser(context *gin.Context) {

}
func UpdateUser(context *gin.Context) {

}
func DeleteUser(context *gin.Context) {

}
