package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	"shop-api/user-web/forms"
	"shop-api/user-web/global"
	"shop-api/user-web/global/response"
	"shop-api/user-web/middlewares"
	"shop-api/user-web/models"
	"shop-api/user-web/proto"
	"strconv"
	"strings"
	"time"
)

func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// HandleGrpcErrorToHttp 转换grpc error错误码为HTTP状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": s.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": s.Message(),
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": s.Message(),
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其它错误",
				})
			}
			return
		}
	}
}

// 表单验证
func handleValidatorErr(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func GetUserList(c *gin.Context) {

	var (
		err         error
		userConn    *grpc.ClientConn
		options     []grpc.DialOption
		userListRsp *proto.UserListResponse
	)
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 拨号连接grpc服务器获取连接
	if userConn, err = grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), options...); err != nil {
		zap.S().Errorw("[GetUserList] 连接 用户服务失败",
			"msg", err.Error(),
		)
		HandleGrpcErrorToHttp(err, c)
		return
	}

	// 基于连接生成grpc client并调用API
	userSrvClient := proto.NewUserClient(userConn)

	// 分页
	pn, _ := strconv.Atoi(c.DefaultQuery("pn", "0"))
	pSize, _ := strconv.Atoi(c.DefaultQuery("psize", "10"))
	pageInfo := &proto.PageInfo{
		Pn:    uint32(pn),
		PSize: uint32(pSize),
	}

	if userListRsp, err = userSrvClient.GetUserList(context.Background(), pageInfo); err != nil {
		zap.S().Errorw("[GetUserList] 查询 用户列表失败",
			"msg", err.Error(),
		)
		HandleGrpcErrorToHttp(err, c)
		return
	}
	res := make([]interface{}, 0)
	for _, userInfoRsp := range userListRsp.Data {
		user := response.UserResponse{
			Id:       userInfoRsp.Id,
			NickName: userInfoRsp.NickName,
			Birthday: response.JsonTime(time.Unix(int64(userInfoRsp.Birthday), 0)),
			Gender:   userInfoRsp.Gender,
			Mobile:   userInfoRsp.Mobile,
		}

		res = append(res, user)
	}
	c.JSON(http.StatusOK, res)
}

// PassWordLogin 账号密码方式登陆
func PassWordLogin(c *gin.Context) {
	var (
		err              error
		token            string
		userConn         *grpc.ClientConn
		options          []grpc.DialOption
		userInfoRsp      *proto.UserInfoResponse
		checkPasswordRsp *proto.CheckResponse
	)

	// 表单验证
	passwordLoginForm := forms.PassWordLoginForm{}
	// 通过shouldBind自己判断请求类型 json or form
	if err = c.ShouldBind(&passwordLoginForm); err != nil {
		// 通过表单验证器验证参数有效性
		handleValidatorErr(c, err)
		return
	}
	// 验证验证码是否准确
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	// 基于连接生成grpc client并调用API
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 拨号连接grpc服务器获取连接
	if userConn, err = grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), options...); err != nil {
		zap.S().Errorw("[GetUserList] 连接 用户服务失败",
			"msg", err.Error(),
		)
		HandleGrpcErrorToHttp(err, c)
		return
	}

	userSrvClient := proto.NewUserClient(userConn)

	// 业务逻辑
	// 查看用户是否存在
	if userInfoRsp, err = userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: passwordLoginForm.Mobile}); err != nil {
		zap.S().Errorw("[GetUserList] 查询 根据mobile查询用户失败",
			"msg", err.Error(),
		)
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, gin.H{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"mobile": "登陆失败",
				})
			}
			return
		}
	} else {
		// 查询用户存在, 那么通过RPC进一步验证密码有效性
		if checkPasswordRsp, err = userSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: userInfoRsp.Password,
			// 尝试RPC调用接口验证密码有效性
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"password": "校验密码异常",
			})
			return
		} else {
			// RPC调用成功查看验证结果, 成功则登陆成功
			if checkPasswordRsp.Success {
				// 生成token返回给用户
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(userInfoRsp.Id),
					NickName:    userInfoRsp.NickName,
					AuthorityId: uint(userInfoRsp.Role),
					StandardClaims: jwt.StandardClaims{
						// 签名的生效时间
						NotBefore: time.Now().Unix(),
						// 过期时间 60s * 60m * 24 * 30 = 30天过期
						ExpiresAt: time.Now().Unix() + 60*60*24*30,
						// 颁发者
						Issuer: "user_srv",
					},
				}
				// 生成token失败返回
				if token, err = j.CreateToken(claims); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}
				// 返回用户信息 + token
				c.JSON(http.StatusOK, gin.H{
					"id":        userInfoRsp.Id,
					"name":      userInfoRsp.NickName,
					"token":     token,
					"expire_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
				return
			} else {
				// 否则密码不正确
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "密码错误",
				})
			}
		}
	}
}
