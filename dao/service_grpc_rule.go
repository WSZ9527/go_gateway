package dao

import (
	"github.com/WSZ9527/go_gateway/public"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// GrpcRule grpc规则表
type GrpcRule struct {
	ID             int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	Port           int    `json:"port" gorm:"column:port" description:"ip端口"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header头信息"`
}

// TableName 映射当前表对应数据库中的表名
func (t *GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}

// Find 根据已知条件查找一个admin对象查找
func (t *GrpcRule) Find(c *gin.Context, tx *gorm.DB, search *GrpcRule) (*GrpcRule, error) {
	grpcRule := &GrpcRule{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(grpcRule).Error
	if err != nil {
		return nil, err
	}

	return grpcRule, nil
}

// Save 更新修改对象到数据库中
func (t *GrpcRule) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
