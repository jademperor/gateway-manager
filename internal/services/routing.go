package services

import (
	"context"
	// "encoding/json"
	"errors"
	
	"github.com/jademperor/common/etcdutils"
	"github.com/jademperor/common/configs"
	"github.com/jademperor/common/models"
	"github.com/jademperor/common/pkg/utils"
	"github.com/jademperor/gateway-manager/internal/logger"
	"go.etcd.io/etcd/client"
)

// AddRouting ...
func AddRouting(routing *models.Routing) (string, error) {
	routingID := utils.UUID()
	routing.Idx = routingID
	routingKey := utils.Fstring("%s%s", configs.RoutingsKey, routingID)

	data, _ := etcdutils.Encode(routing)
	if err := store.Set(routingKey, string(data), -1); err != nil {
		return "", err
	}
	return routingID, nil
}

// DelRouting ...
func DelRouting(routingID string) error {
	routingKey := utils.Fstring("%s%s", configs.RoutingsKey, routingID)
	return store.Delete(routingKey, false)
}

// UpdateRouting ...
func UpdateRouting(routing *models.Routing) error {
	routingKey := utils.Fstring("%s%s", configs.RoutingsKey, routing.Idx)

	data, _ := etcdutils.Encode(routing)
	if err := store.Set(routingKey, string(data), -1); err != nil {
		return err
	}
	return nil
}

// GetAllRoutings get all routing configs from store
func GetAllRoutings(limit, offset int) ([]*models.Routing, int, error) {
	routings := make([]*models.Routing, 0)
	total := 0
	resp, err := store.Kapi.Get(context.Background(), configs.RoutingsKey, nil)
	if err != nil {
		return nil, 0, err
	}
	if !resp.Node.Dir {
		return []*models.Routing{}, total, errors.New("not a directory")
	}
	total = len(resp.Node.Nodes)
	logger.Logger.Infof("GetAllRoutings(limit:%d, offset:%d)", limit, offset)

	// over limit
	if offset >= total {
		return routings, total, nil
	}

	var nodes client.Nodes
	if limit > total-offset {
		nodes = resp.Node.Nodes[offset:total]
	} else {
		nodes = resp.Node.Nodes[offset : offset+limit]
	}

	for _, node := range nodes {
		routing := new(models.Routing)
		if err := etcdutils.Decode(node.Value, routing); err != nil {
			logger.Logger.Errorf("GetAllRoutings got err: %v", err)
			continue
		}
		routings = append(routings, routing)
	}

	return routings, total, nil
}

// GetRoutingInfo ...
func GetRoutingInfo(routingID string) (*models.Routing, error) {
	routingKey := configs.RoutingsKey + routingID
	v, err := store.Get(routingKey)
	if err != nil {
		return nil, err
	}

	routing := new(models.Routing)
	if err = etcdutils.Decode(v, routing); err != nil {
		return nil, err
	}

	return routing, nil
}
