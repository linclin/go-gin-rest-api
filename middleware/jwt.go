package middleware

import (
	"fmt"
	"go-gin-rest-api/models"
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func InitAuth() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:            global.Conf.System.AppName, // jwt标识
		SigningAlgorithm: "RS512",
		PrivKeyFile:      "./conf/rsa/rsa-private.pem", // Private key
		PubKeyFile:       "./conf/rsa/rsa-public.pem",  // Public key
		IdentityKey:      "AppId",
		Timeout:          time.Hour * time.Duration(global.Conf.Jwt.Timeout),    // token过期时间
		MaxRefresh:       time.Hour * time.Duration(global.Conf.Jwt.MaxRefresh), // token更新时间
		PayloadFunc:      payloadFunc,                                           // 有效载荷处理
		IdentityHandler:  identityHandler,                                       // 解析Claims
		Authenticator:    auth,                                                  // 校验token的正确性, 处理登录逻辑
		Authorizator:     authorizator,                                          // 用户登录校验成功处理
		Unauthorized:     unauthorized,                                          // 用户登录校验失败处理
		LoginResponse:    loginResponse,                                         // 登录成功后的响应
		LogoutResponse:   logoutResponse,                                        // 登出后的响应
		TokenLookup:      "header: Authorization",                               // 自动在这几个地方寻找请求中的token header: Authorization, query: token, cookie: jwt
		TokenHeadName:    "Bearer",                                              // header名称
		TimeFunc:         time.Now,
	})
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		return jwt.MapClaims{
			jwt.IdentityKey: v["AppId"],
			"AppId":         v["AppId"],
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	// 此处返回值类型map[string]interface{}与payloadFunc和authorizator的data类型必须一致, 否则会导致授权失败还不容易找到原因
	return map[string]interface{}{
		"IdentityKey": claims[jwt.IdentityKey],
		"AppId":       claims["AppId"],
	}
}

// @Summary [系统内部]获取token
// @Id auth
// @Tags [系统内部]Token
// @version 1.0
// @Accept application/x-json-stream
// @Param	body body 	models.ReqToken	true "token请求"
// @Success 200 object models.Token 返回列表
// @Failure 400 object models.Resp 查询失败
// @Router /api/v1/base/auth [post]
func auth(c *gin.Context) (interface{}, error) {
	var req sys.SysSystem
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	var system sys.SysSystem
	query := global.Mysql.Where(req).First(&system)
	if query.Error != nil {
		return nil, fmt.Errorf("AppId:%s和AppSecret:%s不存在:%s", req.AppId, req.AppSecret, query.Error)
	}
	// 将AppId以json格式写入, payloadFunc/authorizator会使用到
	return map[string]interface{}{
		"AppId": req.AppId,
		"exp":   time.Now().Unix() + 7200,
	}, nil
}

func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(map[string]interface{}); ok {
		// 将用户保存到context, api调用时取数据方便
		c.Set("AppId", v["AppId"].(string))
		return true
	}
	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	global.Log.Error(fmt.Sprintf("JWT认证失败, 错误码%d, 错误信息%s", code, message))
	c.JSON(http.StatusUnauthorized, models.Resp{
		Code: http.StatusUnauthorized,
		Data: "认证失败",
		Msg:  message,
	})
	return
}

func loginResponse(c *gin.Context, code int, token string, expires time.Time) {
	c.JSON(http.StatusOK, models.Token{
		Token:   token,
		Expires: expires,
	})
}

func logoutResponse(c *gin.Context, code int) {
	c.JSON(http.StatusOK, models.Resp{
		Code: http.StatusOK,
		Data: models.OkMsg,
		Msg:  models.CustomError[models.Ok],
	})
}
