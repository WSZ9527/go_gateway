package controller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/WSZ9527/go_gateway/dao"
	"github.com/WSZ9527/go_gateway/dto"
	"github.com/WSZ9527/go_gateway/middleware"
	"github.com/WSZ9527/go_gateway/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AdminLoginController 管理员登录控制器
type AdminLoginController struct {
}

// AdminLoginRegister 注册管理员相关控制器接口到路由
func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin) //注册管理员登录接口
	group.GET("/logout", adminLogin.AdminLogout)
}

// AdminLogin godoc
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags 管理员接口
// @ID /admin/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/login [post]
func (adminlogin *AdminLoginController) AdminLogin(c *gin.Context) {
	// 1 绑定校验拿到登录参数
	params := &dto.AdminLoginInput{}
	params.BindValidParam(c)

	// 2 从数据库连接池中取数据库连接
	tx, err := lib.GetGormPool("default")
	if err != nil { //取连接出错
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 3 调用dao接口验证登录
	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 4 保存session对象到redis数据库
	sessionInfo := &dto.AdminSessionInfo{
		ID:        admin.ID,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessionStr, err := json.Marshal(sessionInfo) //序列化对象
	if err != nil {                              //序列化失败
		middleware.ResponseError(c, 2003, err)
		return
	}
	session := sessions.Default(c)
	session.Set(public.AdminSessionInfoKey, string(sessionStr))
	session.Save() //保存序列化后的结构到session中

	// 5 登录校验成功,返回结果
	out := &dto.AdminLoginOutput{Token: admin.UserName + time.Now().String()}
	middleware.ResponseSuccess(c, out)
}

// AdminLogout godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/logout [get]
func (adminlogin *AdminLoginController) AdminLogout(c *gin.Context) {
	// 1 拿到session
	session := sessions.Default(c)

	// 2 删除保存的session
	session.Delete(public.AdminSessionInfoKey)
	session.Save()

	middleware.ResponseSuccess(c, "")
}

// AdminController 管理员控制器
type AdminController struct {
}

// AdminRegister 注册管理员相关控制器接口到路由
func AdminRegister(group *gin.RouterGroup) {
	admin := &AdminController{}
	group.GET("/info", admin.AdminInfo)
	group.POST("/change_pwd", admin.ChangePwd)
}

// AdminInfo godoc
// @Summary 获取管理员信息
// @Description 获取管理员信息
// @Tags 管理员接口
// @ID /admin/info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/info [get]
func (admin *AdminController) AdminInfo(c *gin.Context) {
	// 1 拿到session
	session := sessions.Default(c)
	sessionStr := session.Get(public.AdminSessionInfoKey)

	// 2 反序列化 获取SessionInfo对象
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessionStr)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 3 返回请求响应结果
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}

	middleware.ResponseSuccess(c, out)
}

// ChangePwd godoc
// @Summary 修改管理员密码
// @Description 修改管理员密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (admin *AdminController) ChangePwd(c *gin.Context) {
	// 1 绑定校验拿到修改的密码
	params := &dto.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 2 拿到session,反序列化拿到session对象
	session := sessions.Default(c)
	sessionStr := session.Get(public.AdminSessionInfoKey)
	sessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessionStr)), sessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 4 从数据库连接池中取数据库连接
	tx, err := lib.GetGormPool("default")
	if err != nil { //取连接出错
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 5 拿到原来的用户信息
	adminInfo := &dao.Admin{}
	adminInfo, err = adminInfo.Find(c, tx, &dao.Admin{UserName: sessionInfo.UserName})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 6 生成新密码
	adminInfo.Password = public.GenSaltPassword(adminInfo.Salt, params.Password)

	// 7 调用dao更新到数据库
	if err = adminInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	// 8 删除session
	session.Delete(public.AdminSessionInfoKey)
	session.Save()

	middleware.ResponseSuccess(c, "")
}
