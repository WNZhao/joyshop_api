package api

import (
	"context"
	"fmt"
	"joyshop_api/user-web/forms"
	"joyshop_api/user-web/global"
	"joyshop_api/user-web/middlewares"
	"joyshop_api/user-web/models"
	"joyshop_api/user-web/proto"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"joyshop_api/user-web/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/consul/api"
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

// GetUserSrvClient 获取用户服务客户端
func GetUserSrvClient() (proto.UserClient, *grpc.ClientConn, error) {
	// 从consul中获取服务信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	consulClient, err := api.NewClient(cfg)

	if err != nil {
		zap.S().Errorw("[GetUserSrvClient] 创建consul客户端失败", "msg", err.Error())
		return nil, nil, err
	}

	// 使用服务名称进行服务发现
	serviceName := global.ServerConfig.UserSrvInfo.Name
	services, _, err := consulClient.Health().Service(serviceName, "", true, nil)

	if err != nil {
		zap.S().Errorw("[GetUserSrvClient] 获取服务列表失败", "msg", err.Error())
		return nil, nil, err
	}

	if len(services) == 0 {
		zap.S().Errorw("[GetUserSrvClient] 未找到可用服务", "service", serviceName)
		return nil, nil, fmt.Errorf("未找到可用服务: %s", serviceName)
	}

	// 获取第一个健康的服务实例
	service := services[0].Service
	userSrvHost := service.Address
	userSrvPort := service.Port

	// 建立grpc连接
	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserSrvClient] 连接 【用户服务失败】", "msg", err.Error())
		return nil, nil, err
	}

	// 生成grpc客户端
	userSrvClient := proto.NewUserClient(userConn)
	return userSrvClient, userConn, nil
}

// GetUserList 获取用户列表
func GetUserList(ctx *gin.Context) {
	zap.S().Debug("获取用户列表页数据")
	claims, ok := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	if ok {
		zap.S().Infof("访问用户：%d, 访问用户昵称:%s", currentUser.ID, currentUser.NickName)
	}

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

	userResp, err := global.UserClient.GetUserList(ctx, &proto.PageInfo{
		Page:     uint32(pageInt),
		PageSize: uint32(paseSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 调用用户服务失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, user := range userResp.Data {
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

// 通过手机号查询
func PassWordLogin(ctx *gin.Context) {
	// 获取请求参数
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := ctx.ShouldBindJSON(&passwordLoginForm); err != nil {
		if utils.HandleValidatorError(ctx, err, "PassWordLogin") {
			return
		}
	}
	// 获取验证码
	if !storage.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "验证码错误",
		})
		return
	}
	// 获取用户服务客户端
	userSrvClient, userConn, err := GetUserSrvClient()
	if err != nil {
		zap.S().Errorw("[PassWordLogin] 连接 【用户服务失败】", "msg", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务器内部错误",
		})
		return
	}
	defer userConn.Close()

	// 1. 先通过手机号获取用户信息
	userInfo, err := userSrvClient.GetUserByMobile(ctx, &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})
	if err != nil {
		zap.S().Errorw("[PassWordLogin] 获取用户信息失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 2. 验证密码
	checkResp, err := userSrvClient.CheckPassword(ctx, &proto.PasswordCheckInof{
		Password:        passwordLoginForm.Password,
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

	j := middlewares.NewJWT()
	claim := models.CustomClaims{
		ID:          uint(userInfo.Id),
		NickName:    userInfo.Nickname,
		AuthorityId: uint(userInfo.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),             // 签名的生效时间
			ExpiresAt: time.Now().Unix() + 3600*24*7, // 过期时间 一周
			Issuer:    "joyshop",                     // 签名的发行者
		},
	}
	token, err := j.CreateToken(claim)
	if err != nil {
		zap.S().Errorw("[PassWordLogin] 生成token失败", "msg", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "服务器内部错误【生成token失败】",
		})
		return
	}
	// 3. 返回用户信息
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "登录成功",
		"data": gin.H{
			"token":    token,
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

// 用户注册
func Register(ctx *gin.Context) {
	registerForm := forms.RegisterForm{}
	if err := ctx.ShouldBindJSON(&registerForm); err != nil {
		if utils.HandleValidatorError(ctx, err, "Register") {
			return
		}
	}

	// 1. 验证短信验证码
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: "", // 如果有密码，在这里设置
		DB:       0,  // 使用默认DB
	})

	// 从Redis获取验证码
	code, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err != nil {
		if err == redis.Nil {
			zap.S().Errorf("验证码已过期: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "验证码已过期",
			})
			return
		}
		zap.S().Errorf("获取验证码失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "获取验证码失败",
		})
		return
	}

	// 验证码比对
	if code != registerForm.Code {
		zap.S().Errorf("验证码错误: 输入=%s, 正确=%s", registerForm.Code, code)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "验证码错误",
		})
		return
	}

	// 2. 生成默认值
	username := fmt.Sprintf("user_%s", registerForm.Mobile)    // 用户名默认为 user_手机号
	nickname := fmt.Sprintf("用户%s", registerForm.Mobile)       // 昵称默认为 用户手机号
	email := fmt.Sprintf("%s@sample.com", registerForm.Mobile) // 邮箱默认为 手机号@sample.com

	// 3. 调用用户服务创建用户
	userSrvClient, userConn, err := GetUserSrvClient()
	if err != nil {
		zap.S().Errorf("连接用户服务失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务器内部错误",
		})
		return
	}
	defer userConn.Close()

	userInfo, err := userSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Nickname: nickname,
		Password: registerForm.Password,
		Mobile:   registerForm.Mobile,
		Email:    email,
		Username: username,
	})

	if err != nil {
		zap.S().Errorf("创建用户失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "创建用户失败",
		})
		return
	}

	// 4. 生成token
	j := middlewares.NewJWT()
	claim := models.CustomClaims{
		ID:          uint(userInfo.Id),
		NickName:    userInfo.Nickname,
		AuthorityId: uint(userInfo.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),             // 签名的生效时间
			ExpiresAt: time.Now().Unix() + 3600*24*7, // 过期时间 一周
			Issuer:    "joyshop",                     // 签名的发行者
		},
	}
	token, err := j.CreateToken(claim)
	if err != nil {
		zap.S().Errorf("生成token失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "生成token失败",
		})
		return
	}

	// 5. 删除验证码
	rdb.Del(context.Background(), registerForm.Mobile)

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "注册成功",
		"data": gin.H{
			"id":       userInfo.Id,
			"nickname": userInfo.Nickname,
			"token":    token,
			"expired":  time.Now().Add(time.Duration(global.ServerConfig.JWTInfo.ExpireTime) * time.Hour).Unix(),
		},
	})
}
