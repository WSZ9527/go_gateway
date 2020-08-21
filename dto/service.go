package dto

import (
	"github.com/WSZ9527/go_gateway/public"
	"github.com/gin-gonic/gin"
)

// ServiceListInput 登录接口接收参数
type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" validate:""`                      //关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" example:"1" validate:"required"`        //页数
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" validate:"required"` //每页条数
}

// BindValidParam 绑定登录接口参数并校验
func (param *ServiceListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

// ServiceListItemOutput 登录接口接收参数
type ServiceListItemOutput struct {
	ID          int64  `json:"id" form:"id" comment:"ID" example:"" validate:""`                       //ID
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名" example:"" validate:""`  //服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:""` //服务描述
	LoadType    string `json:"load_type" form:"load_type" comment:"负载类型" example:"" validate:""`       //负载类型
	ServiceAddr string `json:"service_addr" form:"service_addr" comment:"服务地址" example:"" validate:""` //服务地址
	QPS         int64  `json:"qps" form:"qps" comment:"QPS" example:"" validate:""`                    //QPS
	Qpd         int64  `json:"qpd" form:"qpd" comment:"日请求数" example:"" validate:""`                   //日请求数
	TotalNode   int    `json:"total_node" form:"total_node" comment:"节点总数" example:"" validate:""`     //节点总数
}

// ServiceListOutput 登录接口接收参数
type ServiceListOutput struct {
	Total int64                   `json:"total" form:"total" comment:"总数" example:"" validate:""`             //总数
	List  []ServiceListItemOutput `json:"page_no" form:"page_no" comment:"列表" example:"" validate:"required"` //列表
}
