package main

import (
	"flag"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jademperor/common/pkg/ginutils"
	"github.com/jademperor/common/pkg/utils"
	"github.com/jademperor/gateway-manager/internal/controllers"
	"github.com/jademperor/gateway-manager/internal/healthchecking"
	"github.com/jademperor/gateway-manager/internal/logger"
	"github.com/jademperor/gateway-manager/internal/services"
)

var (
	engine    *gin.Engine
	etcdAddrs utils.StringArray

	addr    = flag.String("addr", ":8999", "the addr http api server listen and serve on, default = 8999")
	debug   = flag.Bool("debug", false, "set debug mode on, default not open debug mode (false)")
	logpath = flag.String("logpath", "./logs", "the folder directory what log files would be stored at")
)

func prepare() {
	engine = gin.New()

	// install middlewares
	engine.Use(ginutils.CORS())
	engine.Use(ginutils.Recovery(logger.Logger.Out))
	engine.Use(ginutils.LogRequest(logger.Logger, false))

	// register http apis
	// engine.GET("/v1/local", controllers.LocalManageAPIS)

	engine.GET("/v1/clusters", controllers.GetAllClusters)
	engine.GET("/v1/cluster_ids", controllers.GetAllClustersIDs)
	engine.POST("/v1/cluster", controllers.AddCluster)
	engine.DELETE("/v1/clusters/:clusterID", controllers.DelCluster)
	engine.PUT("/v1/clusters/:clusterID", controllers.UpdateClusterInfo)
	engine.GET("/v1/clusters/:clusterID", controllers.GetClusterInfo)

	engine.POST("/v1/clusters/:clusterID/instance", controllers.AddClusterInstance)
	engine.DELETE("/v1/clusters/:clusterID/instance/:instanceID", controllers.DelClusterInstance)
	engine.PUT("/v1/clusters/:clusterID/instance/:instanceID", controllers.UpdateClusterInstance)
	engine.GET("/v1/clusters/:clusterID/instance/:instanceID", controllers.GetClusterInstance)

	engine.GET("/v1/apis", controllers.GetAllAPIs)
	engine.POST("/v1/apis/api", controllers.AddAPI)
	engine.DELETE("/v1/apis/:apiID", controllers.DelAPI)
	engine.PUT("/v1/apis/:apiID", controllers.UpdateAPI)
	engine.GET("/v1/apis/:apiID", controllers.GetAPIInfo)

	engine.GET("/v1/routings", controllers.GetAllRoutings)
	engine.POST("/v1/routings/routing", controllers.AddRouting)
	engine.DELETE("/v1/routings/:routingID", controllers.DelRouting)
	engine.PUT("/v1/routings/:routingID", controllers.UpdateRouting)
	engine.GET("/v1/routings/:routingID", controllers.GetRoutingInfo)

	// engine.GET("/v1/plugins", controllers.GetAllPlugins)
	// engine.PUT("/v1/plugins/:id/status", controllers.UpdatePluginsStatus)

	// engine.GET("/v1/plugins/cache/rules", controllers.GetAllCacheRules)
	// engine.POST("/v1/plugins/cache/rule", controllers.AddCacheRule)
	// engine.DELETE("/v1/plugins/cache/rules/:ruleID", controllers.DelCacheRule)
	// engine.PUT("/v1/plugins/cache/rules/:ruleID", controllers.UpdateCacheRule)
	// engine.GET("/v1/plugins/cache/rules/:ruleID", controllers.GetCacheRule)
}

func main() {
	flag.Var(&etcdAddrs, "etcd-addr", "set etcd endpoints to connect to etcd store")
	flag.Parse()
	if len(etcdAddrs) == 0 {
		log.Fatal("error: etcd-addr need one endpoint at least!")
	}

	// initilize work
	if err := logger.Init(*logpath); err != nil {
		log.Fatal(err)
	}
	if err := services.Init(etcdAddrs); err != nil {
		log.Fatal(err)
	}

	// close gin debug mode
	if !*debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// start health checker ....
	healthchecking.Init(etcdAddrs, 1*time.Second)

	// start the server
	prepare()
	logger.Logger.Infof("Listening and serving HTTP on %s", *addr)
	if err := engine.Run(*addr); err != nil {
		log.Fatal(err)
	}
}
