package dao

import (
	"time"

	"github.com/WSZ9527/go_gateway/dto"
	"github.com/WSZ9527/go_gateway/public"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// Service 管理员信息输入对象
type Service struct {
	ID           int64     `json:"id" gorm:"primary_key" description:"自增主键"`
	LoadType     int       `json:"load_type" gorm:"column:load_type" description:"负载均衡类型"`
	ServiceName  string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServeiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	UpdatedAt    time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt    time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete     int8      `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

// ServiceDetail 详细的服务信息
type ServiceDetail struct {
	Info          *Service       `json:"info" description:"基本信息"`
	HTTPRule      *HTTPRule      `json:"http_rule" description:"http_rule"`
	TCPRule       *TCPRule       `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *GrpcRule      `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance" description:"load_balance"`
	AccessControl *AccessControl `json:"access_control" description:"access_control"`
}

// TableName 映射当前表对应数据库中的表名
func (t *Service) TableName() string {
	return "gateway_service_info"
}

// Find 根据已知条件查找一个Service对象查找
func (t *Service) Find(c *gin.Context, tx *gorm.DB, search *Service) (*Service, error) {
	serviceInfo := &Service{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(serviceInfo).Error
	if err != nil {
		return nil, err
	}

	return serviceInfo, nil
}

// Save 更新修改对象到数据库中
func (t *Service) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}

// PageList 服务对象分页查询
func (t *Service) PageList(c *gin.Context, tx *gorm.DB, param *dto.ServiceListInput) ([]Service, int64, error) {
	total := int64(0)
	list := []Service{}
	offset := (param.PageNo - 1) * param.PageSize

	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Where("is_delete=0")
	if param.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(param.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}

// PageListDetail 服务对象分页详情信息查询
func (t *Service) PageListDetail(c *gin.Context, tx *gorm.DB, search *Service) (*ServiceDetail, error) {
	// 构建各种表
	httpRule := &HTTPRule{ServiceID: search.ID}
	httpRule, err := httpRule.Find(c, tx, httpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	tcpRule := &TCPRule{ServiceID: search.ID}
	tcpRule, err = tcpRule.Find(c, tx, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	grpcRule := &GrpcRule{ServiceID: search.ID}
	grpcRule, err = grpcRule.Find(c, tx, grpcRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	accessControl := &AccessControl{ServiceID: search.ID}
	accessControl, err = accessControl.Find(c, tx, accessControl)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	loadBalance := &LoadBalance{ServiceID: search.ID}
	loadBalance, err = loadBalance.Find(c, tx, loadBalance)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	serviceDetail := &ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	return serviceDetail, nil
}
