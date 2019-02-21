package controllers

import (
	"github.com/jademperor/gateway-manager/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jademperor/common/models"
	"github.com/jademperor/common/pkg/code"
	"github.com/jademperor/common/pkg/ginutils"
)

type getAllAPIsForm struct {
	Limit  int `form:"limit,default=10" binding:"gte=0"`
	Offset int `form:"offset,default=0" binding:"gte=0"`
}
type getAllAPIsResp struct {
	code.CodeInfo
	APIs  []*models.API `json:"apis"`
	Total int           `json:"total"`
}

// GetAllAPIs get all api configs
func GetAllAPIs(c *gin.Context) {
	var (
		form = new(getAllAPIsForm)
		resp = new(getAllAPIsResp)
		err  error
	)

	if err = c.ShouldBind(form); err != nil {
		err = ginutils.HdlValidationErrors(err)
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	if resp.APIs, resp.Total, err = services.GetAllAPIs(form.Limit, form.Offset); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

type addAPIForm struct {
	Path            string            `json:"path" binding:"required"`
	Method          string            `json:"method" binding:"required"`
	TargetClusterID string            `json:"target_cluster_id" binding:"required"`
	RewritePath     string            `json:"rewrite_path" binding:"required"`
	NeedCombine     bool              `json:"need_combine"`
	CombineReqCfgs  []*apiCombination `json:"api_combination" binding:"required"`
}

type apiCombination struct {
	Path            string `json:"path" binding:"required"`
	Field           string `json:"field" binding:"required"`
	Method          string `json:"method" binding:"required"`
	TargetClusterID string `json:"target_cluster_id" binding:"required"`
}

type addAPIResp struct {
	code.CodeInfo
	APIID string `json:"api_id,omitempty"`
}

// AddAPI add an api config
func AddAPI(c *gin.Context) {
	var (
		form = new(addAPIForm)
		resp = new(addAPIResp)
		err  error
	)

	if err = c.ShouldBindJSON(form); err != nil {
		err = ginutils.HdlValidationErrors(err)
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	combCfgs := make([]*models.APICombination, len(form.CombineReqCfgs))
	for idx, combCfg := range form.CombineReqCfgs {
		combCfgs[idx] = &models.APICombination{
			Path:            combCfg.Path,
			Field:           combCfg.Field,
			Method:          combCfg.Method,
			TargetClusterID: combCfg.TargetClusterID,
		}
	}

	apiCfg := &models.API{
		Path:            form.Path,
		Method:          form.Method,
		TargetClusterID: form.TargetClusterID,
		RewritePath:     form.RewritePath,
		NeedCombine:     form.NeedCombine,
		CombineReqCfgs:  combCfgs,
	}

	if resp.APIID, err = services.AddAPI(apiCfg); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

// type delAPIForm struct{}
type delAPIResp struct {
	code.CodeInfo
}

// DelAPI del an api config
func DelAPI(c *gin.Context) {
	var (
		// form = new(delAPIForm)
		resp = new(delAPIResp)
	)

	apiID := c.Param("apiID")
	if err := services.DelAPI(apiID); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

type updateAPIForm struct {
	Path            string            `json:"path" binding:"required"`
	Method          string            `json:"method" binding:"required"`
	TargetClusterID string            `json:"target_cluster_id" binding:"required"`
	RewritePath     string            `json:"rewrite_path" binding:"required"`
	NeedCombine     bool              `json:"need_combine"`
	CombineReqCfgs  []*apiCombination `json:"api_combination" binding:"required"`
}
type updateAPIResp struct {
	code.CodeInfo
}

// UpdateAPI update an api config
func UpdateAPI(c *gin.Context) {
	var (
		form = new(updateAPIForm)
		resp = new(updateAPIResp)
		err  error
	)

	if err = c.ShouldBindJSON(form); err != nil {
		err = ginutils.HdlValidationErrors(err)
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	combCfgs := make([]*models.APICombination, len(form.CombineReqCfgs))
	for idx, combCfg := range form.CombineReqCfgs {
		combCfgs[idx] = &models.APICombination{
			Path:            combCfg.Path,
			Field:           combCfg.Field,
			Method:          combCfg.Method,
			TargetClusterID: combCfg.TargetClusterID,
		}
	}

	apiCfg := &models.API{
		Idx:             c.Param("apiID"),
		Path:            form.Path,
		Method:          form.Method,
		TargetClusterID: form.TargetClusterID,
		RewritePath:     form.RewritePath,
		NeedCombine:     form.NeedCombine,
		CombineReqCfgs:  combCfgs,
	}

	if err = services.UpdateAPI(apiCfg); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

// type getAPIInfoForm struct{}
type getAPIInfoResp struct {
	code.CodeInfo
	API *models.API `json:"api"`
}

// GetAPIInfo get api config
func GetAPIInfo(c *gin.Context) {
	var (
		// form = new(getAPIInfoForm)
		resp = new(getAPIInfoResp)
		err  error
	)

	apiID := c.Param("apiID")
	if resp.API, err = services.GetAPIInfo(apiID); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}
