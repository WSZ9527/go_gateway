package dao

import (
	"github.com/WSZ9527/go_gateway/public"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// HTTPRule 管理员信息输入对象
type HTTPRule struct {
	ID             int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	RuleType       int    `json:"rule_type" gorm:"column:rule_type" description:"匹配类型 0=url前缀url_prefix 1=域名domain "`
	Rule           string `json:"rule" gorm:"column:rule" description:"type=domain表示域名，type=url_prefix时表示url前缀"`
	NeedHTTPS      int    `json:"need_https" gorm:"column:need_https" description:"支持https 1=支持"`
	NeedStripURI   int    `json:"need_strip_uri" gorm:"column:need_strip_uri" description:"启用strip_uri 1=启用"`
	NeedWebsocket  int    `json:"need_websocket" gorm:"column:need_websocket" description:"是否支持websocket 1=支持"`
	URIRewrite     string `json:"url_rewrite" gorm:"column:url_rewrite" description:"url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔"`
}

// TableName 映射当前表对应数据库中的表名
func (t *HTTPRule) TableName() string {
	return "gateway_service_http_rule"
}

// Find 根据已知条件查找一个admin对象查找
func (t *HTTPRule) Find(c *gin.Context, tx *gorm.DB, search *HTTPRule) (*HTTPRule, error) {
	httpRule := &HTTPRule{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(httpRule).Error
	if err != nil {
		return nil, err
	}

	return httpRule, nil
}

// Save 更新修改对象到数据库中
func (t *HTTPRule) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
