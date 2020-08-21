package dao

import (
	"github.com/WSZ9527/go_gateway/public"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// TCPRule tcp规则表
type TCPRule struct {
	ID        int64 `json:"id" gorm:"primary_key"`
	ServiceID int64 `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	Port      int   `json:"port" gorm:"column:port" description:"端口	"`
}

// TableName 映射当前表对应数据库中的表名
func (t *TCPRule) TableName() string {
	return "gateway_service_tcp_rule"
}

// Find 根据已知条件查找一个admin对象查找
func (t *TCPRule) Find(c *gin.Context, tx *gorm.DB, search *TCPRule) (*TCPRule, error) {
	model := &TCPRule{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save 更新修改对象到数据库中
func (t *TCPRule) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}
