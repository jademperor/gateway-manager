package services

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/jademperor/common/configs"
	"github.com/jademperor/common/models"
	"github.com/jademperor/common/pkg/utils"
	"github.com/jademperor/gateway-manager/internal/logger"
)

// Cluster service layer
type Cluster struct {
	Idx       string                   `json:"idx"`
	Name      string                   `json:"name"`
	Instances []*models.ServerInstance `json:"instances"`
}

// NewCluster generate a new cluster
func NewCluster(name string, srvInstances []*models.ServerInstance) (clusterID string, err error) {
	clusterID = utils.UUID()

	clsOpt := models.ClusterOption{
		Idx:  clusterID,
		Name: name,
	}
	clusterKey := utils.Fstring("%s%s", configs.ClustersKey, clusterID)
	clusterOptKey := utils.Fstring("%s/%s", clusterKey, configs.ClusterOptionsKey)
	data, err := json.Marshal(clsOpt)
	if err != nil {
		logger.Logger.Errorf("json.Marshal(clsOpt) got err: %v", err)
		return "", err
	}
	// save cluster option
	if err = store.Set(clusterOptKey, string(data), -1); err != nil {
		return "", err
	}

	// save instances
	for _, instance := range srvInstances {
		instanceID := utils.UUID()
		instanceKey := utils.Fstring("%s/%s", clusterKey, instanceID)
		instance.ClusterID = clusterID
		instance.Idx = instanceID
		data, _ := json.Marshal(instance)
		_ = store.Set(instanceKey, string(data), -1)
	}
	return
}

// DelCluster del a cluster using store
func DelCluster(clusterID string) error {
	clusterKey := utils.Fstring("%s%s", configs.ClustersKey, clusterID)
	return store.Delete(clusterKey, true)
}

// UpdateClusterInfo update the cluster info (ClusterOption)
func UpdateClusterInfo(clusterID, name string) error {
	clusterOptKey := utils.Fstring("%s%s/%s",
		configs.ClustersKey, clusterID, configs.ClusterOptionsKey)

	clsOpt := &models.ClusterOption{
		Idx:  clusterID,
		Name: name,
	}
	data, _ := json.Marshal(clsOpt)

	return store.Set(clusterOptKey, string(data), -1)
}

// GetAllClusters ...
func GetAllClusters() ([]*Cluster, error) {
	var (
		clusterCfgs = make([]*Cluster, 0)
	)
	resp, err := store.Kapi.Get(context.Background(), configs.ClustersKey, nil)
	if err != nil {
		return nil, err
	} else if !resp.Node.Dir {
		return clusterCfgs, nil
	}
	for _, clusterNode := range resp.Node.Nodes {
		clusterID := strings.Split(clusterNode.Key, "/")[2]
		clsOpt := new(models.ClusterOption)
		srvInses := make([]*models.ServerInstance, 0)
		if resp2, err := store.Kapi.Get(context.Background(), clusterNode.Key, nil); err == nil && resp2.Node.Dir {
			for _, srvInsNode := range resp2.Node.Nodes {
				// skip the option node
				if strings.Split(srvInsNode.Key, "/")[3] == configs.ClusterOptionsKey {
					if err := json.Unmarshal([]byte(srvInsNode.Value), clsOpt); err != nil {
						logger.Logger.Error(err)
					}
					continue
				}

				srvInsCfg := new(models.ServerInstance)
				if err := json.Unmarshal([]byte(srvInsNode.Value), srvInsCfg); err != nil {
					logger.Logger.Error(err)
					continue
				}
				srvInses = append(srvInses, srvInsCfg)
			}
			if len(srvInses) != 0 {
				clusterCfgs = append(clusterCfgs, &Cluster{
					Idx:       clusterID,
					Name:      clsOpt.Name,
					Instances: srvInses,
				})
			}
		}
	}

	return clusterCfgs, nil
}

// GetClusterInfo ...
func GetClusterInfo(clusterID string) (*Cluster, error) {
	clusterKey := utils.Fstring("%s%s", configs.ClustersKey, clusterID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := store.Kapi.Get(ctx, clusterKey, nil)
	if err != nil {
		return nil, err
	}

	clsOpt := new(models.ClusterOption)
	srvInses := make([]*models.ServerInstance, 0)
	// load server instance ...
	for _, srvInsNode := range resp.Node.Nodes {
		// skip the option node
		if strings.Split(srvInsNode.Key, "/")[3] == configs.ClusterOptionsKey {
			if err := json.Unmarshal([]byte(srvInsNode.Value), clsOpt); err != nil {
				logger.Logger.Error(err)
			}
			continue
		}

		srvInsCfg := new(models.ServerInstance)
		if err := json.Unmarshal([]byte(srvInsNode.Value), srvInsCfg); err != nil {
			logger.Logger.Error(err)
			continue
		}
		srvInses = append(srvInses, srvInsCfg)
	}

	return &Cluster{
		Idx:       clusterID,
		Name:      clsOpt.Name,
		Instances: srvInses,
	}, nil
}

// AddClusterInstance add a instance into the cluster
func AddClusterInstance(clusterID, name, addr string,
	weight int, need bool) (instanceID string, err error) {
	instanceID = utils.UUID()
	instanceKey := utils.Fstring("%s%s/%s", configs.ClustersKey, clusterID, instanceID)

	srvInstance := &models.ServerInstance{
		Idx:             instanceID,
		Name:            name,
		Addr:            addr,
		ClusterID:       clusterID,
		Weight:          weight,
		NeedCheckHealth: need,
	}
	data, _ := json.Marshal(srvInstance)

	err = store.Set(instanceKey, string(data), -1)
	return
}

// DelClusterInstance del a instance from a cluster instance sets
func DelClusterInstance(clusterID, instanceID string) error {
	instanceKey := utils.Fstring("%s%s/%s", configs.ClustersKey, clusterID, instanceID)
	return store.Delete(instanceKey, false)
}

// UpdateClusterInstanceInfo update a instance info in a cluster sets
func UpdateClusterInstanceInfo(clusterID, instanceID, name, addr string,
	weight int, need bool) error {
	instanceKey := utils.Fstring("%s%s/%s", configs.ClustersKey, clusterID, instanceID)
	srvInstance := &models.ServerInstance{
		Idx:             instanceID,
		Name:            name,
		Addr:            addr,
		ClusterID:       clusterID,
		Weight:          weight,
		NeedCheckHealth: need,
	}
	data, _ := json.Marshal(srvInstance)

	store.Set(instanceKey, string(data), -1)
	return nil
}

// GetClusterInstanceInfo load cluster instance from cluster
func GetClusterInstanceInfo(clusterID, instanceID string) (*models.ServerInstance, error) {
	instance := new(models.ServerInstance)
	instanceKey := utils.Fstring("%s%s/%s", configs.ClustersKey, clusterID, instanceID)
	v, err := store.Get(instanceKey)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(v), instance); err != nil {
		return nil, err
	}
	return instance, nil
}
