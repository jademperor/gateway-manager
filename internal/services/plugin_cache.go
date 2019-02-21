package services

import (
	"github.com/jademperor/common/etcdutils"
	"context"
	// "encoding/json"
	"errors"

	"github.com/jademperor/common/configs"
	"github.com/jademperor/common/models"
	"github.com/jademperor/common/pkg/utils"
	"github.com/jademperor/gateway-manager/internal/logger"
	"go.etcd.io/etcd/client"
)

// GetAllCacheRules ...
func GetAllCacheRules(limit, offset int) ([]*models.NocacheCfg, int, error) {
	rules := make([]*models.NocacheCfg, 0)
	total := 0
	resp, err := store.Kapi.Get(context.Background(), configs.CacheKey, nil)
	if err != nil {
		return rules, 0, err
	}
	if !resp.Node.Dir {
		return []*models.NocacheCfg{}, total, errors.New("not a directory")
	}
	total = len(resp.Node.Nodes)
	logger.Logger.Infof("GetAllNocacheCfgs(limit:%d, offset:%d)", limit, offset)

	// over limit
	if offset >= total {
		return rules, total, nil
	}

	var nodes client.Nodes
	if limit > total-offset {
		nodes = resp.Node.Nodes[offset:total]
	} else {
		nodes = resp.Node.Nodes[offset : offset+limit]
	}

	for _, node := range nodes {
		rule := new(models.NocacheCfg)
		if err := etcdutils.Decode(node.Value, rule); err != nil {
			logger.Logger.Errorf("GetAllNocacheCfgs got err: %v", err)
			continue
		}
		rules = append(rules, rule)
	}

	return rules, total, nil
}

// AddCacheRule ...
func AddCacheRule(regexp string, enabled bool) (string, error) {
	rule := &models.NocacheCfg{
		Idx:     utils.UUID(),
		Regexp:  regexp,
		Enabled: enabled,
	}
	ruleKey := utils.Fstring("%s%s", configs.CacheKey, rule.Idx)

	data, _ := etcdutils.Encode(rule)
	if err := store.Set(ruleKey, string(data), -1); err != nil {
		return "", err
	}
	return rule.Idx, nil
}

// DelCacheRule ...
func DelCacheRule(ruleID string) error {
	ruleKey := utils.Fstring("%s%s", configs.CacheKey, ruleID)
	return store.Delete(ruleKey, false)
}

// UpdateCacheRule ...
func UpdateCacheRule(ruleID, regexp string, enabled bool) error {
	ruleKey := utils.Fstring("%s%s", configs.CacheKey, ruleID)

	rule := &models.NocacheCfg{
		Idx:     ruleID,
		Regexp:  regexp,
		Enabled: enabled,
	}
	data, _ := etcdutils.Encode(rule)
	if err := store.Set(ruleKey, string(data), -1); err != nil {
		return err
	}
	return nil
}

// GetCacheRule ...
func GetCacheRule(ruleID string) (*models.NocacheCfg, error) {
	ruleKey := utils.Fstring("%s%s", configs.CacheKey, ruleID)
	rule := new(models.NocacheCfg)
	v, err := store.Get(ruleKey)
	if err != nil {
		return nil, err
	}
	err = etcdutils.Decode((v), rule)
	return rule, err
}
