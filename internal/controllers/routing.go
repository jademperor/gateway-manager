package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jademperor/common/models"
	"github.com/jademperor/common/pkg/code"
	"github.com/jademperor/common/pkg/ginutils"
	"github.com/jademperor/gateway-manager/internal/services"
)

type getAllRoutingsForm struct {
	Limit  int `form:"limit,default=10" binding:"gte=0"`
	Offset int `form:"offset,default=0" binding:"gte=0"`
}
type getAllRoutingsResp struct {
	code.CodeInfo
	Routings []*models.Routing `json:"routings"`
	Total    int               `json:"total"`
}

// GetAllRoutings get all Routing configs
func GetAllRoutings(c *gin.Context) {
	var (
		form = new(getAllRoutingsForm)
		resp = new(getAllRoutingsResp)
		err  error
	)

	if err = c.ShouldBind(form); err != nil {
		err = ginutils.HdlValidationErrors(err)
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	if resp.Routings, resp.Total, err = services.GetAllRoutings(form.Limit, form.Offset); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

type addRoutingForm struct {
	Prefix          string `json:"prefix" binding:"required"`
	ClusterID       string `json:"cluster_id" binding:"required"`
	NeedStripPrefix bool   `json:"need_strip_prefix" binding:"required"`
}

type addRoutingResp struct {
	code.CodeInfo
	RoutingID string `json:"routing_id,omitempty"`
}

// AddRouting add an Routing config
func AddRouting(c *gin.Context) {
	var (
		form = new(addRoutingForm)
		resp = new(addRoutingResp)
		err  error
	)

	if err = c.ShouldBindJSON(form); err != nil {
		err = ginutils.HdlValidationErrors(err)
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	routingCfg := &models.Routing{
		Prefix:          form.Prefix,
		ClusterID:       form.ClusterID,
		NeedStripPrefix: form.NeedStripPrefix,
	}

	if resp.RoutingID, err = services.AddRouting(routingCfg); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

// type delRoutingForm struct{}
type delRoutingResp struct {
	code.CodeInfo
}

// DelRouting del an Routing config
func DelRouting(c *gin.Context) {
	var (
		// form = new(delRoutingForm)
		resp = new(delRoutingResp)
	)

	routingID := c.Param("routingID")
	if err := services.DelRouting(routingID); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

type updateRoutingForm struct {
	Prefix          string `json:"prefix" binding:"required"`
	ClusterID       string `json:"cluster_id" binding:"required"`
	NeedStripPrefix bool   `json:"need_strip_prefix" binding:"required"`
}
type updateRoutingResp struct {
	code.CodeInfo
}

// UpdateRouting update an Routing config
func UpdateRouting(c *gin.Context) {
	var (
		form = new(updateRoutingForm)
		resp = new(updateRoutingResp)
		err  error
	)

	if err = c.ShouldBindJSON(form); err != nil {
		err = ginutils.HdlValidationErrors(err)
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	routingCfg := &models.Routing{
		Idx:             c.Param("routingID"),
		Prefix:          form.Prefix,
		ClusterID:       form.ClusterID,
		NeedStripPrefix: form.NeedStripPrefix,
	}

	if err = services.UpdateRouting(routingCfg); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

// type getRoutingInfoForm struct{}
type getRoutingInfoResp struct {
	code.CodeInfo
	Routing *models.Routing `json:"routing,omitempty"`
}

// GetRoutingInfo get Routing config
func GetRoutingInfo(c *gin.Context) {
	var (
		// form = new(getRoutingInfoForm)
		resp = new(getRoutingInfoResp)
		err  error
	)

	routingID := c.Param("routingID")
	if resp.Routing, err = services.GetRoutingInfo(routingID); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}
