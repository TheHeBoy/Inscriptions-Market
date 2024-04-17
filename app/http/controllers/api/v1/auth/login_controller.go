package auth

import (
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/requests"
	"gohub/pkg/auth"
	"gohub/pkg/errorcode"
	"gohub/pkg/jwt"
	"gohub/pkg/logger"
	"gohub/pkg/response"

	"github.com/gin-gonic/gin"
)

// LoginController 用户控制器
type LoginController struct {
	v1.BaseAPIController
}

// LoginByPhone 手机登录
func (lc *LoginController) LoginByPhone(c *gin.Context) {

	// 1. 验证表单
	request := requests.LoginByPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPhone); !ok {
		return
	}

	// 2. 尝试登录
	user, err := auth.LoginByPhone(request.Phone)
	if err != nil {
		logger.Error(err)
		// 失败，显示错误提示
		response.Error(c, errorcode.AUTH_ACCOUNT_OR_PASSWD_NOT_EXIST)
	} else {
		// 登录成功
		token := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)

		response.SuccessData(c, gin.H{
			"token": token,
		})
	}
}

// LoginByPassword 多种方法登录，支持手机号、email 和用户名
func (lc *LoginController) LoginByPassword(c *gin.Context) {
	// 1. 验证表单
	request := requests.LoginByPasswordRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPassword); !ok {
		return
	}

	// 2. 尝试登录
	user, err := auth.Attempt(request.LoginID, request.Password)
	if err != nil {
		// 失败，显示错误提示
		response.Error(c, errorcode.AUTH_ACCOUNT_OR_PASSWD_NOT_EXIST)

	} else {
		token := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)
		response.SuccessData(c, gin.H{
			"token": token,
		})
	}
}

// RefreshToken 刷新 Access Token
func (lc *LoginController) RefreshToken(c *gin.Context) {

	token, err := jwt.NewJWT().RefreshToken(c)

	if err != nil {
		logger.Error(err)
		response.Error(c, errorcode.AUTH_TOKEN_REFRESH_FAIL)
	} else {
		response.SuccessData(c, gin.H{
			"token": token,
		})
	}
}
