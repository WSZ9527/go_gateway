package dao

import (
	"github.com/WSZ9527/go_gateway/public"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// AccessControl 权限控制表
type AccessControl struct {
	ID                int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceID         int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	OpenAuth          int    `json:"open_auth" gorm:"column:open_auth" description:"是否开启权限 1=开启"`
	BlackList         string `json:"black_list" gorm:"column:black_list" description:"黑名单ip"`
	WhiteList         string `json:"white_list" gorm:"column:white_list" description:"白名单ip"`
	WhiteHostName     string `json:"white_host_name" gorm:"column:white_host_name" description:"白名单主机"`
	ClientipFlowLimit int64  `json:"clientip_flow_limit" gorm:"column:clientip_flow_limit" description:"客户端ip限流"`
	ServiceFlowLimit  int64  `json:"service_flow_limit" gorm:"column:service_flow_limit" description:"服务端ip限流"`
}

// TableName 映射当前表对应数据库中的表名
func (t *AccessControl) TableName() string {
	return "gateway_service_access_control"
}

// Find 根据已知条件查找一个admin对象查找
func (t *AccessControl) Find(c *gin.Context, tx *gorm.DB, search *AccessControl) (*AccessControl, error) {
	accessControl := &AccessControl{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(accessControl).Error
	if err != nil {
		return nil, err
	}

	return accessControl, nil
}

// Save 更新修改对象到数据库中
func (t *AccessControl) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
