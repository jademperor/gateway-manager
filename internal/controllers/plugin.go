package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jademperor/common/models"
	"github.com/jademperor/common/pkg/code"
	"github.com/jademperor/common/pkg/ginutils"
	"github.com/jademperor/gateway-manager/internal/services"
)

type getAllCacheRulesForm struct {
	Limit  int `form:"limit,default=10" binding:"gte=0"`
	Offset int `form:"offset,defualt=0" binding:"gte=0"`
}
type getAllCacheRulesResp struct {
	code.CodeInfo
	Rules []*models.NocacheCfg `json:"rules"`
	Total int                  `json:"total"`
}

// GetAllCacheRules ...
func GetAllCacheRules(c *gin.Context) {
	var (
		form = new(getAllCacheRulesForm)
		resp = new(getAllCacheRulesResp)
		err  error
	)

	if err = c.ShouldBind(form); err != nil {
		err = ginutils.HdlValidationErrors(err)
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	if resp.Rules, resp.Total, err = services.GetAllCacheRules(form.Limit, form.Offset); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

type addCacheRuleForm struct {
	Regexp  string `form:"regexp" binding:"required"`
	Enabled bool   `form:"enabled"`
}
type addCacheRuleResp struct {
	code.CodeInfo
	RuleID string `json:"ruleID,omitempty"`
}

// AddCacheRule ...
func AddCacheRule(c *gin.Context) {
	var (
		form = new(addCacheRuleForm)
		resp = new(addCacheRuleResp)
		err  error
	)

	if err = c.ShouldBind(form); err != nil {
		err = ginutils.HdlValidationErrors(err)
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	if resp.RuleID, err = services.AddCacheRule(form.Regexp, form.Enabled); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

// type delCacheRuleForm struct{}
type delCacheRuleResp struct {
	code.CodeInfo
}

// DelCacheRule ...
func DelCacheRule(c *gin.Context) {
	var (
		// form = new(delCacheRuleForm)
		resp = new(delCacheRuleResp)
		err  error
	)

	// if err = c.ShouldBind(form); err != nil {
	// 	err = ginutils.HdlValidationErrors(err)
	// 	code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
	// 	c.JSON(http.StatusOK, resp)
	// 	return
	// }

	ruleID := c.Param("ruleID")
	if err = services.DelCacheRule(ruleID); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

type updateCacheRuleForm struct {
	Regexp  string `form:"regexp" binding:"required"`
	Enabled bool   `form:"enabled"`
}
type updateCacheRuleResp struct {
	code.CodeInfo
}

// UpdateCacheRule ....
func UpdateCacheRule(c *gin.Context) {
	var (
		form = new(updateCacheRuleForm)
		resp = new(updateCacheRuleResp)
		err  error
	)

	if err = c.ShouldBind(form); err != nil {
		err = ginutils.HdlValidationErrors(err)
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	ruleID := c.Param("ruleID")
	if err = services.UpdateCacheRule(ruleID, form.Regexp, form.Enabled); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}

// type getCacheRuleForm struct{}
type getCacheRuleResp struct {
	code.CodeInfo
	Rule *models.NocacheCfg `json:"rule,omitempty"`
}

// GetCacheRule ...
func GetCacheRule(c *gin.Context) {
	var (
		// form = new(getCacheRuleForm)
		resp = new(getCacheRuleResp)
		err  error
	)

	// if err = c.ShouldBind(form); err != nil {
	// 	err = ginutils.HdlValidationErrors(err)
	// 	code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
	// 	c.JSON(http.StatusOK, resp)
	// 	return
	// }
	ruleID := c.Param("ruleID")
	if resp.Rule, err = services.GetCacheRule(ruleID); err != nil {
		code.FillCodeInfo(resp, code.NewCodeInfo(code.CodeSystemErr, err.Error()))
		c.JSON(http.StatusOK, resp)
		return
	}

	code.FillCodeInfo(resp, code.GetCodeInfo(code.CodeOk))
	c.JSON(http.StatusOK, resp)
}
