package dto

import (
	"time"

	"github.com/WSZ9527/go_gateway/public"
	"github.com/gin-gonic/gin"
)

// AdminLoginInput 登录接口接收参数
type AdminLoginInput struct {
	UserName string `json:"username" form:"username" comment:"管理员用户名" example:"admin" validate:"required,valid_username"` //管理员用户名
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`                   //密码
}

// BindValidParam 绑定登录接口参数并校验
func (param *AdminLoginInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

// AdminLoginOutput 退出登录接口接收参数
type AdminLoginOutput struct {
	Token string `json:"token" form:"token" comment:"token" example:"token" validate:""` //token
}

// AdminInfoOutput 管理员信息接口输出参数
type AdminInfoOutput struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	LoginTime    time.Time `json:"login_time"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}

// ChangePwdInput 管理员改密接口接收参数
type ChangePwdInput struct {
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"` //密码
}

// BindValidParam 绑定改密接口参数并校验
func (param *ChangePwdInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

// AdminSessionInfo 管理员session参数
type AdminSessionInfo struct {
	ID        int       `json:"id"`
	UserName  string    `json:"user_name"`
	LoginTime time.Time `json:"login_time"`
}
