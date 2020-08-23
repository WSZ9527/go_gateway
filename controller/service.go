package controller

import (
	"fmt"

	"github.com/WSZ9527/go_gateway/dao"
	"github.com/WSZ9527/go_gateway/dto"
	"github.com/WSZ9527/go_gateway/middleware"
	"github.com/WSZ9527/go_gateway/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

// ServiceController 服务控制器
type ServiceController struct {
}

// ServiceRegister 路由注册
func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/list", service.ServiceList)
	group.GET("/delete", service.ServiceDelete)
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/list [get]
func (service *ServiceController) ServiceList(c *gin.Context) {
	// 1 获取参数
	params := &dto.ServiceListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	serviceInfo := &dao.Service{}
	// 2 从数据库连接池中取数据库连接
	tx, err := lib.GetGormPool("default")
	if err != nil { //取连接出错
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 3 dao分页查询
	list, total, err := serviceInfo.PageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 4 拼接dto返回对象
	outList := []dto.ServiceListItemOutput{}
	for _, listItem := range list {
		//拿到详情信息
		serviceDetail, err := listItem.PageListDetail(c, tx, &listItem)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}
		serviceAddr := "unknow"
		//组装http后缀接入
		clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")

		// 判断接入的方式
		// http 后缀接入 clusterIP+clusterPort+path
		// http 域名接入 domian
		// tcp/grpc接入 clusterIP+servicePort
		switch {
		case serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHTTPS == 0:
			{
				serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, serviceDetail.HTTPRule.Rule)
			}
		case serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHTTPS == 1:
			{
				serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, serviceDetail.HTTPRule.Rule)
			}
		case serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain:
			{
				serviceAddr = serviceDetail.HTTPRule.Rule
			}
		case serviceDetail.Info.LoadType == public.LoadTypeTCP:
			{
				serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port)
			}
		case serviceDetail.Info.LoadType == public.LoadTypeGRPC:
			{
				serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port)
			}
		}

		ipList := serviceDetail.LoadBalance.GetIPListByMode()
		outItem := &dto.ServiceListItemOutput{
			ID:          listItem.ID,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServeiceDesc,
			ServiceAddr: serviceAddr,
			QPS:         0,
			Qpd:         0,
			TotalNode:   len(ipList),
		}

		outList = append(outList, *outItem)
	}
	out := &dto.ServiceListOutput{
		Total: total,
		List:  outList,
	}
	middleware.ResponseSuccess(c, out)
}

// ServiceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/delete
// @Accept  json
// @Produce  json
// @Param id query string true "服务id"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/delete [get]
func (service *ServiceController) ServiceDelete(c *gin.Context) {
	// 1 获取参数
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	serviceInfo := &dao.Service{ID: params.ID}
	// 2 从数据库连接池中取数据库连接
	tx, err := lib.GetGormPool("default")
	if err != nil { //取连接出错
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 3 dao分页查询,读取服务基本信息
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 4 更改软删除属性,并保存
	serviceInfo.IsDelete = 1
	if err = serviceInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, "")
}
