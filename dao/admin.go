package dao

import (
	"errors"
	"time"

	"github.com/WSZ9527/go_gateway/dto"
	"github.com/WSZ9527/go_gateway/public"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// Admin 管理员信息输入对象
type Admin struct {
	ID        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName  string    `json:"user_name" gorm:"column:user_name" description:"管理员用户名"`
	Salt      string    `json:"salt" gorm:"column:salt" description:"盐"`
	Password  string    `json:"password" gorm:"column:password" description:"密码"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

// TableName 映射当前表对应数据库中的表名
func (t *Admin) TableName() string {
	return "gateway_admin"
}

// Find 根据已知条件查找一个admin对象查找
func (t *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	adminInfo := &Admin{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(adminInfo).Error
	if err != nil {
		return nil, err
	}

	return adminInfo, nil
}

// LoginCheck 检测登录信息是否正确
func (t *Admin) LoginCheck(c *gin.Context, tx *gorm.DB, param *dto.AdminLoginInput) (*Admin, error) {
	adminInfo, err := t.Find(c, tx, (&Admin{UserName: param.UserName, IsDelete: 0}))
	if err != nil {
		return nil, errors.New("用户信息不存在")
	}
	saltPassword := public.GenSaltPassword(adminInfo.Salt, param.Password)
	if adminInfo.Password != saltPassword {
		return nil, errors.New("密码错误，请重新输入")
	}
	return adminInfo, nil
}

// Save 更新修改对象到数据库中
func (t *Admin) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
