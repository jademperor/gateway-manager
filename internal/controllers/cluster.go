package controllers

import (
	"github.com/jademperor/common/models"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jademperor/common/pkg/code"
	"github.com/jademperor/common/pkg/ginutils"
	"github.com/jademperor/gateway-manager/internal/services"
)

type getAllClustersResp struct {
	code.CodeInfo
	Clusters []*services.Cluster `json:"clusters"`
}

// GetAllClusters load all clusters info from
func GetAllClusters(c *gin.Context) {
	var (
		resp = new(getAllClustersResp)
		err  error
	)

	if resp.Clusters, err = services.GetAllClusters(); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
	return
}

// JSON type request data
type addClusterJSON struct {
	Name      string                   `json:"name" binding:"required"`
	Instances []*models.ServerInstance `json:"instances" binding:"required"`
}
type addClusterResp struct {
	code.CodeInfo
	ClusterID string `json:"cluster_id,omitempty"`
}

// AddCluster add a cluster into store
func AddCluster(c *gin.Context) {
	var (
		jsForm = new(addClusterJSON)
		resp   = new(addClusterResp)
		err    error
	)

	if err = c.ShouldBindJSON(jsForm); err != nil {
		err = ginutils.HdlValidationErrors(err)
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	if resp.ClusterID, err = services.NewCluster(jsForm.Name, jsForm.Instances); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
	return
}

type delClusterResp struct {
	code.CodeInfo
}

// DelCluster del a cluster and all server instance
func DelCluster(c *gin.Context) {
	var (
		resp = new(delClusterResp)
	)

	clusterID := c.Param("clusterID")

	if err := services.DelCluster(clusterID); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
	return
}

type updateClusterInfoForm struct {
	Name string `form:"name" binding:"required"`
	// Instances []*models.ServerInstance `json:"instances" binding:"required"`
}

type updateClusterInfoResp struct {
	code.CodeInfo
}

// UpdateClusterInfo ... update cluster info
func UpdateClusterInfo(c *gin.Context) {
	var (
		form = new(updateClusterInfoForm)
		resp = new(updateClusterInfoResp)
	)

	if err := c.ShouldBind(form); err != nil {
		err = ginutils.HdlValidationErrors(err)
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	clusterID := c.Param("clusterID")

	if err := services.UpdateClusterInfo(clusterID, form.Name); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
	return
}

type getClusterInfoResp struct {
	code.CodeInfo
	Cluster *services.Cluster `json:"cluster,omitempty"`
}

// GetClusterInfo get single cluster info
func GetClusterInfo(c *gin.Context) {
	var (
		resp = new(getClusterInfoResp)
		err  error
	)

	clusterID := c.Param("clusterID")

	if resp.Cluster, err = services.GetClusterInfo(clusterID); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
	return
}

type addClusterInsForm struct {
	Name            string `form:"name" binding:"required"`
	Addr            string `form:"addr" binding:"required"`
	Weight          int    `form:"weight" binding:"required"`
	NeedCheckHealth bool   `form:"need_check_health"`
}
type addClusterInsResp struct {
	code.CodeInfo
	IntanceID string `json:"instance_id"`
}

// AddClusterInstance add a new instance into cluster
func AddClusterInstance(c *gin.Context) {
	var (
		form = new(addClusterInsForm)
		resp = new(addClusterInsResp)
		err  error
	)

	if err = c.ShouldBind(form); err != nil {
		err = ginutils.HdlValidationErrors(err)
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	clusterID := c.Param("clusterID")
	if resp.IntanceID, err = services.AddClusterInstance(clusterID, form.Name,
		form.Addr, form.Weight, form.NeedCheckHealth); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

// type delClusterInsForm struct{}
type delClusterInsResp struct {
	code.CodeInfo
}

// DelClusterInstance del a instance from the cluster
func DelClusterInstance(c *gin.Context) {
	var (
		resp = new(delClusterInsResp)
	)

	clusterID := c.Param("clusterID")
	instanceID := c.Param("instanceID")

	if err := services.DelClusterInstance(clusterID, instanceID); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

type updateClusterInsForm struct {
	Name            string `form:"name" binding:"required"`
	Addr            string `form:"addr" binding:"required"`
	Weight          int    `form:"weight" binding:"required"`
	NeedCheckHealth bool   `form:"need_check_health"`
}
type updateClusterInsResp struct {
	code.CodeInfo
}

// UpdateClusterInstance update a server intance in the cluster
func UpdateClusterInstance(c *gin.Context) {
	var (
		form = new(updateClusterInsForm)
		resp = new(updateClusterInsResp)
		err  error
	)

	if err = c.ShouldBind(form); err != nil {
		err = ginutils.HdlValidationErrors(err)
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	clusterID := c.Param("clusterID")
	instanceID := c.Param("instanceID")
	if err = services.UpdateClusterInstanceInfo(clusterID, instanceID,
		form.Name, form.Addr, form.Weight, form.NeedCheckHealth); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

// type getClusterInsForm struct{}
type getClusterInsResp struct {
	code.CodeInfo
	Instance *models.ServerInstance `json:"instance,omitempty"`
}

// GetClusterInstance get instance detail in the cluster
func GetClusterInstance(c *gin.Context) {
	var (
		resp = new(getClusterInsResp)
		err  error
	)

	clusterID := c.Param("clusterID")
	instanceID := c.Param("instanceID")
	if resp.Instance, err = services.GetClusterInstanceInfo(clusterID, instanceID); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}
