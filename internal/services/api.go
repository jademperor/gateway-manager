package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jademperor/common/configs"
	"github.com/jademperor/common/models"
	"github.com/jademperor/common/pkg/utils"
	"github.com/jademperor/gateway-manager/internal/logger"
	"go.etcd.io/etcd/client"
)

// AddAPI ...
func AddAPI(api *models.API) (string, error) {
	apiID := utils.UUID()
	api.Idx = apiID
	apiKey := utils.Fstring("%s%s", configs.APIsKey, apiID)

	data, _ := json.Marshal(api)
	if err := store.Set(apiKey, string(data), -1); err != nil {
		return "", err
	}
	return apiID, nil
}

// DelAPI ...
func DelAPI(apiID string) error {
	apiKey := utils.Fstring("%s%s", configs.APIsKey, apiID)
	return store.Delete(apiKey, false)
}

// UpdateAPI ...
func UpdateAPI(api *models.API) error {
	apiKey := utils.Fstring("%s%s", configs.APIsKey, api.Idx)

	data, _ := json.Marshal(api)
	if err := store.Set(apiKey, string(data), -1); err != nil {
		return err
	}
	return nil
}

// GetAllAPIs get all api configs from store
func GetAllAPIs(limit, offset int) ([]*models.API, int, error) {
	apis := make([]*models.API, 0)
	total := 0
	resp, err := store.Kapi.Get(context.Background(), configs.APIsKey, nil)
	if err != nil {
		return nil, 0, err
	}
	if !resp.Node.Dir {
		return []*models.API{}, total, errors.New("not a directory")
	}
	total = len(resp.Node.Nodes)
	logger.Logger.Infof("GetAllAPIs(limit:%d, offset:%d)", limit, offset)

	// over limit
	if offset >= total {
		return apis, total, nil
	}

	var nodes client.Nodes
	if limit > total-offset {
		nodes = resp.Node.Nodes[offset:total]
	} else {
		nodes = resp.Node.Nodes[offset : offset+limit]
	}

	for _, node := range nodes {
		api := new(models.API)
		if err := json.Unmarshal([]byte(node.Value), api); err != nil {
			logger.Logger.Errorf("GetAllAPIs got err: %v", err)
			continue
		}
		apis = append(apis, api)
	}

	return apis, total, nil
}

// GetAPIInfo ...
func GetAPIInfo(apiID string) (*models.API, error) {
	apiKey := configs.APIsKey + apiID
	v, err := store.Get(apiKey)
	if err != nil {
		return nil, err
	}

	api := new(models.API)
	if err = json.Unmarshal([]byte(v), api); err != nil {
		return nil, err
	}

	return api, nil
}
