package main

import (
	"PublicChainBrowser-Server/controllers"
	"PublicChainBrowser-Server/db/mongo"
	"PublicChainBrowser-Server/db/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
)

var Mgo controllers.Chain

func main() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		panic(err)
	}

	mgo, err := mongo.InitMongod()
	if err != nil {
		panic(err)
	}
	redisCli := new(redis.RedisCli)
	redisCli.Config = cfg.Section("redis")

	h := controllers.Chain{Mgo: mgo}
	Mgo = h

	controllers.InitChainInfo(h)

	router.GET("/PublicChainBrowser/chain/test", h.Test)
	router.GET("/PublicChainBrowser/chain/getTest", h.GetTest)
	router.GET("/PublicChainBrowser/chain/getChainInfo", h.GetChainInfo)
	router.GET("/PublicChainBrowser/chain/getChainCommittee", h.GetChainCommittee)
	router.GET("/PublicChainBrowser/chain/getMainPageInfo", h.GetMainPageInfo)
	router.GET("/PublicChainBrowser/chain/getChainInfoStruct", h.GetChainInfoStruct)
	router.GET("/PublicChainBrowser/chain/getChainStatByType", h.GetChainStatByType)
	router.GET("/PublicChainBrowser/chain/getBlockTxByFilter", h.GetBlockTxByFilter)
	router.POST("/PublicChainBrowser/chain/getAccountByAddress", h.GetAccountByAddress)
	router.POST("/PublicChainBrowser/chain/getChainStats", h.GetChainStats)
	router.POST("/PublicChainBrowser/chain/getMainChainStat", h.GetMainChainStat)
	router.POST("/PublicChainBrowser/chain/getBlockNewTx", h.GetBlockNewTx)
	router.POST("/PublicChainBrowser/chain/getBlockNewTxPage", h.GetBlockNewTxPage)
	router.POST("/PublicChainBrowser/chain/getTxByTxTypeAndChainId", h.GetTxByTxTypeAndChainId)
	router.POST("/PublicChainBrowser/chain/getTxByContractAndChainId", h.GetTxByContractAndChainId)
	router.POST("/PublicChainBrowser/chain/getBlockTxByAddress", h.GetBlockTxByAddress)
	router.POST("/PublicChainBrowser/chain/getBlockDataByPage", h.GetBlockDataByPage)
	router.POST("/PublicChainBrowser/chain/getBlockDataByEpoch", h.GetBlockDataByEpoch)

	router.POST("/PublicChainBrowser/chain/getBlockData", h.GetBlockData)

	router.POST("/PublicChainBrowser/chain/getBlockDataInfo", h.GetBlockDataInfo)
	router.GET("/PublicChainBrowser/file", h.GetFile)
	chainGroup := router.Group("/PublicChainBrowser/chain")
	{

		chainGroup.POST("/getTxTypeByHeight", h.GetTxTypeByHeight)
		chainGroup.POST("/getChildrenChainStatsById", h.GetChildrenChainStatsById)
		chainGroup.POST("/getTxByParentId", h.GetTxByParentId)
	}

	router.Static("/static/", "static/")
	port := cfg.Section("http_serv").Key("port").String()
	if err := router.Run(port); err != nil {
		panic(err)
	}
}
