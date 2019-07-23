package controllers

import (
	m "PublicChainBrowser-Server/db/mongo"
	"PublicChainBrowser-Server/log"
	"PublicChainBrowser-Server/rpc"
	"container/list"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
	"web3.go/common/hexutil"
	"web3.go/encoding"
)

/* 	初始化链的节点链信息等 */
func InitChainInfo(c Chain) {
	rsp := new(Rsp)
	body, err := rpc.GetChainInfo("[]")
	if err != nil {
		rsp.Msg = "faild"
		rsp.Code = CODE_Error
		//g.JSON(http.StatusOK, rsp)
		return
	}
	var ChainNodeinfos []ChainNodeInfo
	err = json.Unmarshal(body, &ChainNodeinfos)
	if err != nil {
		rsp.Msg = "faild"
		rsp.Code = CODE_Error
		//g.JSON(http.StatusOK, rsp)
		return
	}
	for i := 0; i < len(ChainNodeinfos); i++ {
		result := findAndUpdateChinInfo(c, ChainNodeinfos[i])
		if !result {
			rsp.Msg = "faild"
			rsp.Code = CODE_Error
			//g.JSON(http.StatusOK, rsp)
			return
		}
	}
	rsp.Msg = "success"
	rsp.Code = CODE_Success
	//g.JSON(http.StatusOK, rsp)
	return
}

func findAndUpdateChinInfo(c Chain, info ChainNodeInfo) bool {
	filter := bson.D{{"chainid", info.ChainId}}

	optCount := new(options.CountOptions)
	optCount.Hint = Index_id

	count, err := c.Mgo.Collection(m.ChainInfo).CountDocuments(context.Background(), filter, optCount)
	if err != nil {
		return false
	}
	if count == 0 {
		c.Mgo.Collection(m.ChainInfo).InsertOne(context.Background(), info)
	}
	if count > 0 {
		SingleResult := c.Mgo.Collection(m.ChainInfo).FindOneAndReplace(context.Background(), filter, info)
		fmt.Println(SingleResult)
	}
	return true
}

/* 获取链信息 */
func (c *Chain) GetChainInfo(g *gin.Context) {
	rsp := new(Rsp)
	cur, err := c.Mgo.Collection(m.ChainInfo).Find(context.Background(), bson.D{}, nil)
	if err != nil {
		log.Error(err.Error())
		rsp.Msg = "faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	var infos = make([]ChainNodeInfo, 0)
	for cur.Next(context.Background()) {
		elem := new(ChainNodeInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		infos = append(infos, *elem)
	}
	rsp.Data = infos
	rsp.Msg = "success"
	rsp.Code = CODE_Success
	g.JSON(http.StatusOK, rsp)
	return
}

/* 获取链信息 父子结构 */
func (c *Chain) GetChainInfoStruct(g *gin.Context) {
	rsp := new(Rsp)
	cur, err := c.Mgo.Collection(m.ChainInfo).Find(context.Background(), bson.D{}, nil)
	if err != nil {
		rsp.Msg = "get data faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	var infos = make([]ChainNodeInfoStruct, 0)
	for cur.Next(context.Background()) {
		elem := new(ChainNodeInfoStruct)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		elem.Childrens = make([]ChainNodeInfoStruct, 0)
		infos = append(infos, *elem)
	}

	info := new(ChainNodeInfoStruct)
	for i := 0; i < len(infos); i++ {
		if infos[i].ChainId == 0 {
			info = &infos[i]
		}
	}
	info.Childrens = getChildren(info.ChainId, infos)
	for i := 0; i < len(info.Childrens); i++ {
		info.Childrens[i].Childrens = getChildren(info.Childrens[i].ChainId, infos)
	}

	rsp.Data = info
	rsp.Msg = "success"
	rsp.Code = CODE_Success
	g.JSON(http.StatusOK, rsp)
	return
}

func getChainInfoStruct(c *Chain) *ChainNodeInfoStruct {

	var infos = make([]ChainNodeInfoStruct, 0)
	info := new(ChainNodeInfoStruct)
	cur, err := c.Mgo.Collection(m.ChainInfo).Find(context.Background(), bson.D{}, nil)
	if err != nil {
		return nil
	}

	for cur.Next(context.Background()) {
		elem := new(ChainNodeInfoStruct)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		elem.Childrens = make([]ChainNodeInfoStruct, 0)
		infos = append(infos, *elem)
	}

	for i := 0; i < len(infos); i++ {
		if infos[i].ChainId == 0 {
			info = &infos[i]
		}
	}
	info.Childrens = getChildren(info.ChainId, infos)
	for i := 0; i < len(info.Childrens); i++ {
		info.Childrens[i].Childrens = getChildren(info.Childrens[i].ChainId, infos)
	}
	return info
}

func getChildren(parentId int, infos []ChainNodeInfoStruct) []ChainNodeInfoStruct {
	var newinfos = make([]ChainNodeInfoStruct, 0)
	for i := 0; i < len(infos); i++ {
		if infos[i].Parent == parentId {
			newinfos = append(newinfos, infos[i])
		}
	}
	return newinfos
}

func getAllChainInfo(c *Chain) []ChainNodeInfo {
	var infos = make([]ChainNodeInfo, 0)
	cur, err := c.Mgo.Collection(m.ChainInfo).Find(context.Background(), bson.D{}, nil)
	if err != nil {
		log.Error(err.Error())
		return infos
	}
	for cur.Next(context.Background()) {
		elem := new(ChainNodeInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		infos = append(infos, *elem)
	}
	return infos
}

func getAllChainInfoParent(c *Chain) []ChainInfoParents {
	var infos = make([]ChainInfoParents, 0)
	cur, err := c.Mgo.Collection(m.ChainInfo).Find(context.Background(), bson.D{}, nil)
	if err != nil {
		log.Error(err.Error())
		return infos
	}
	for cur.Next(context.Background()) {
		elem := new(ChainInfoParents)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		if elem.Parent == 0 && elem.ChainId != 0 { //子链
			infos = append(infos, *elem)
		}
	}

	for cur.Next(context.Background()) {
		elem := new(ChainInfoParents)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		if elem.Parent != 0 && elem.ChainId != 0 { //子链
			AddChildren(infos, *elem)
		}
	}

	return infos
}

func AddChildren(infos []ChainInfoParents, info ChainInfoParents) {
	for i := 0; i < len(infos); i++ {
		if info.Parent == infos[i].ChainId {
			infos[i].Childrens = append(infos[i].Childrens, info)
		}
	}
}

/* 获取主页信息 */
func (c *Chain) GetMainPageInfo(g *gin.Context) {
	rsp := new(Rsp)
	filter := bson.D{{}}
	var sumCount int64
	var tpsCount int64
	cur, err := c.Mgo.Collection(m.ChainStats).Find(context.Background(), filter)
	if err == nil {
		for cur.Next(context.Background()) {
			elem := new(ChainStat)
			if err := cur.Decode(elem); err != nil {
				log.Error(err.Error())
			}
			sumCount = sumCount + int64(elem.Txcount)
			tpsCount = tpsCount + int64(elem.Tps)
		}
	}

	var info = new(MainPageInfo)
	info.ThkPrice = 12.5
	info.MarketValue = 12.5 * 10000000000
	info.TxCount = sumCount
	info.Tps = tpsCount

	rsp.Data = info
	rsp.Msg = "success"
	rsp.Code = CODE_Success
	g.JSON(http.StatusOK, rsp)
	return
}

/* 获取已分片子链 */
func (c *Chain) GetChainStatByType(g *gin.Context) {
	rsp := new(Rsp)
	var all = make([]ChainNodeInfo, 0)
	all = getAllChainInfo(c)

	var resultMap = make(map[string]interface{})
	//var typeA=make([]ChainNodeInfo, 0)
	//var typeB=make([]ChainNodeInfo, 0)

	var chainIdsA []int
	var chainIdsB []int
	for i := 0; i < len(all); i++ {
		if all[i].Parent == 0 {
			if existChildren(all, all[i]) > 0 { //分片了
				//typeA = append(typeA, all[i])
				//resultMap["typeA"]=typeA
				chainIdsA = append(chainIdsA, all[i].ChainId)
			} else {
				//typeB = append(typeB, all[i]) //没有分片
				//resultMap["typeB"]=typeB
				chainIdsB = append(chainIdsB, all[i].ChainId)
			}
		}
	}

	//1. 已经分片

	var ChainStatsA = make([]ChainStat, 0)
	ChainStatsA = getChainStats(c, chainIdsA, 1)

	//2. 未分片
	var ChainStatsB = make([]ChainStat, 0)
	ChainStatsB = getChainStats(c, chainIdsB, 0)

	resultMap["chainStatA"] = ChainStatsA //已分片
	resultMap["chainStatB"] = ChainStatsB //未分片

	rsp.Msg = "success"
	rsp.Code = CODE_Success
	rsp.Data = resultMap
	g.JSON(http.StatusOK, rsp)
	return
}

func getChainStats(c *Chain, queryIds []int, flag int) []ChainStat {
	var ChainStats = make([]ChainStat, 0)
	onlineFilter := bson.M{"chainid": bson.M{"$in": queryIds}}
	cur, err := c.Mgo.Collection(m.ChainStats).Find(context.Background(), onlineFilter)
	if err != nil {
		return ChainStats
	}

	for cur.Next(context.Background()) {
		elem := new(ChainStat)
		if err := cur.Decode(elem); err != nil {
			log.Error(err.Error())
		}
		elem.ChildrenTxcount = int64(getTxByChainIdParent(c, elem.ChainId))
		ChainStats = append(ChainStats, *elem)
	}

	for i := 0; i < len(ChainStats); i++ {

		if ChainStats[i].ChainId == 0 { //子链
			//count, err := c.Mgo.Collection(m.ChainInfo).CountDocuments(context.Background(), bson.D{{"parent", result.ChainId}})
			//if err == nil {
			//	result.ChildrenCount = int(count) //子链个数
			//}
		} else { //分片

			count, err := c.Mgo.Collection(m.ChainInfo).CountDocuments(context.Background(), bson.D{{"parent", ChainStats[i].ChainId}})
			if err == nil {
				ChainStats[i].ChildrenCount = int(count) //分片个数

				if ChainStats[i].ChildrenCount > 0 {
					cur, err := c.Mgo.Collection(m.ChainInfo).Find(context.Background(), bson.D{{"parent", ChainStats[i].ChainId}})
					var queryId []int32
					if err != nil {
						log.Error(err.Error())
					}
					for cur.Next(context.Background()) {
						elme := new(ChainStat)
						err := cur.Decode(elme)
						if err != nil {
							log.Error(err.Error())
						}
						queryId = append(queryId, elme.ChainId)
					}
					var sumCount int64
					for i := 0; i < len(queryId); i++ {
						info := new(ChainStat)
						c.Mgo.Collection(m.ChainStats).FindOne(context.Background(), bson.D{{"chainid", queryId[i]}}).Decode(&info)
						sumCount = sumCount + int64(info.Txcount)

					}
					ChainStats[i].ChildrenTxcount = sumCount
				}
			}
		}

	}

	return ChainStats
}

func existChildren(all []ChainNodeInfo, info ChainNodeInfo) int {
	var result int
	for i := 0; i < len(all); i++ {
		if all[i].Parent == info.ChainId {
			result++
		}
	}
	return result
}

/* 获取委员会成员信息 */
func (c *Chain) GetChainCommittee(g *gin.Context) {

	rsp := new(Rsp)
	chainIdstr := g.Query("chainId")
	chainId, _ := strconv.Atoi(chainIdstr)

	epochstr := g.Query("epoch")
	epoch, _ := strconv.Atoi(epochstr)

	filter := bson.D{}
	if chainIdstr == "" {
		rsp.Msg = "chainId or epochstr is empty!"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	} else {
		filter = bson.D{{"chainid", chainId}}
	}

	if epochstr != "" {
		filter = bson.D{{"chainid", chainId}, {"epoch", epoch}}
	}

	opts := new(options.FindOptions)
	limit := int64(1)
	//skip := 0

	sortMap := make(map[string]interface{})
	sortMap["epoch"] = -1
	opts.Sort = sortMap

	opts.Limit = &limit
	//opts.Skip =0

	cur, err := c.Mgo.Collection(m.Committee).Find(context.Background(), filter, opts)
	if err != nil {
		rsp.Msg = "faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	var infos = make([]Committee, 0)
	for cur.Next(context.Background()) {
		elem := new(Committee)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		if elem.MemberDetails == nil {
			elem.MemberDetails = make([]string, 0)
		}

		infos = append(infos, *elem)
		break //取一条
	}
	rsp.Data = infos
	rsp.Msg = "success"
	rsp.Code = CODE_Success
	g.JSON(http.StatusOK, rsp)
	return
}

/* 根据子链ID 获得旗下分片的交易总量 */
func getTxByChainIdParent(c *Chain, ChainId int32) (txCount int64) {
	filter := bson.D{{"parent", ChainId}}
	cur, err := c.Mgo.Collection(m.ChainInfo).Find(context.Background(), filter)
	if err != nil {
		return 0
	}
	var queryId []int32
	for cur.Next(context.Background()) {
		elme := new(ChainStat)
		err := cur.Decode(elme)
		if err != nil {
			log.Error(err.Error())
		}
		queryId = append(queryId, elme.ChainId)
	}
	if len(queryId) > 0 {
		filter := bson.M{"chainid": bson.M{"$in": queryId}}
		cur, err := c.Mgo.Collection(m.ChainStats).Find(context.Background(), filter)
		if err == nil {
			var txCount int64
			for cur.Next(context.Background()) {
				elem := new(ChainStat)
				err := cur.Decode(elem)
				if err != nil {
					log.Error(err.Error())
				}
				txCount = txCount + int64(elem.Txcount)
			}
			return txCount
		}
	}
	return 0
}

func (c *Chain) GetBlockTxByFilter(g *gin.Context) {
	rsp := new(Rsp)
	id := g.Query("chainId")
	chainId, err := strconv.Atoi(id)
	if err != nil {
		rsp := new(Rsp)
		rsp.Code = CODE_Error
		rsp.Msg = "chainId is error!"
		g.JSON(http.StatusOK, rsp)
		return
	}
	page := 1
	opts := new(options.FindOptions)
	limit := int64(6)
	skip := int64((page - 1) * 6)

	sortMap := make(map[string]interface{})
	sortMap["height"] = -1
	opts.Sort = sortMap

	opts.Limit = &limit
	opts.Skip = &skip
	filter := bson.D{{"chainid", chainId}}
	cur, err := c.Mgo.Collection(m.BlockTxs).Find(context.Background(), filter, opts)
	if err != nil {
		log.Error(err.Error())
		return
	}
	txs := make([]TxInfo, 0)
	for cur.Next(context.Background()) {
		elem := new(TxInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}

		txs = append(txs, *elem)
	}

	rsp.Data = txs
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

//获取已分片子链 的分片链信息
func (c *Chain) GetChildrenChainStatsById(g *gin.Context) {
	rsp := new(Rsp)
	var info PostParamTx
	err := g.BindJSON(&info)
	if err != nil {
		rsp.Msg = "json faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	chainId, _ := strconv.Atoi(info.ChainId)

	if chainId == 0 {
		rsp.Msg = "chainId must not 0"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	filter := bson.D{{"parent", chainId}}
	cur, err := c.Mgo.Collection(m.ChainInfo).Find(context.Background(), filter)
	if err != nil {
		log.Error(err.Error())
	}
	var QueryIds []int32
	for cur.Next(context.Background()) {
		elem := new(ChainInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		QueryIds = append(QueryIds, elem.ChainId)
	}

	chainStatsFilter := bson.M{"chainid": bson.M{"$in": QueryIds}}
	curChainStats, err := c.Mgo.Collection(m.ChainStats).Find(context.Background(), chainStatsFilter)
	if err != nil {
		log.Error(err.Error())
	}
	var ChainStats = make([]ChainStat, 0)
	for curChainStats.Next(context.Background()) {
		elem := new(ChainStat)
		err := curChainStats.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		ChainStats = append(ChainStats, *elem)
	}

	//给web套层
	rsp.Data = ChainStats
	rsp.Msg = "success"
	rsp.Code = CODE_Success
	g.JSON(http.StatusOK, rsp)
}

//获取已分片子链下的最新交易信息
func (c *Chain) GetTxByParentId(g *gin.Context) {
	rsp := new(Rsp)
	var info PostParamTx
	err := g.BindJSON(&info)
	if err != nil {
		rsp.Msg = "json faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	chainId, _ := strconv.Atoi(info.ChainId)

	if chainId == 0 {
		rsp.Msg = "chainId must not 0"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	filter := bson.D{{"parent", chainId}}
	cur, err := c.Mgo.Collection(m.ChainInfo).Find(context.Background(), filter)
	if err != nil {
		log.Error(err.Error())
	}
	var QueryIds []int32
	for cur.Next(context.Background()) {
		elem := new(ChainInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		QueryIds = append(QueryIds, elem.ChainId)
	}

	page := info.Page
	pageSize := info.Pagesize

	/* 默认 */
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	data := new(AllBlockData)
	data.PageSize = pageSize
	data.PageNum = page

	opts := new(options.FindOptions)
	limit := int64(pageSize)
	skip := int64((page - 1) * pageSize)

	sortMap := make(map[string]interface{})
	sortMap["timestamp"] = -1
	opts.Sort = sortMap

	opts.Limit = &limit
	opts.Skip = &skip

	chainStatsFilter := bson.M{"chainid": bson.M{"$in": QueryIds}}
	curTxs, err := c.Mgo.Collection(m.BlockTxs).Find(context.Background(), chainStatsFilter, opts)

	var Txs = make([]TxInfo, 0)
	for curTxs.Next(context.Background()) {
		elem := new(TxInfo)
		err := curTxs.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		Txs = append(Txs, *elem)
	}

	rsp.Data = Txs
	rsp.Msg = "success"
	rsp.Code = CODE_Success
	g.JSON(http.StatusOK, rsp)
	return
}

/* 模拟交易类型  后面从同步里面数据里取*/
func getTxType(info TxInfo) string {
	var result string

	////模拟数据
	if strings.Contains(info.Hash, "huo3") {
		result = TxType3
		return result
	}
	if strings.Contains(info.Hash, "huo4") {
		result = TxType4
		return result
	}

	//合约发布1
	if info.To == "" {
		result = TxType1
		return result
	}
	//合约交易2
	if len(info.Input) > 2 {
		result = TxType2
		return result
	}

	//链内转账3
	//跨链转账4
	return result
}

// 浏览器controller
func (c *Chain) GetAllTxChainCount(g *gin.Context) {
	rsp := new(Rsp)
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	data := new(ChainInfoList)
	countFilter := bson.M{"chainid": bson.M{"$in": []int{1, 2}}}
	opts := new(options.FindOptions)
	sortMap := make(map[string]interface{})
	sortMap["chainid"] = 1
	opts.Sort = sortMap
	cur, err := c.Mgo.Collection(m.ChainStats).Find(context.Background(), countFilter, opts)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for cur.Next(context.Background()) {
		elem := new(ChainInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		data.ChainInfolist.DataList = append(data.ChainInfolist.DataList, *elem)
	}
	rsp.Data = data
	g.JSON(http.StatusOK, rsp)
	return
}

/* 获取ChainStats*/
func (c *Chain) GetChainStats(g *gin.Context) {
	rsp := new(Rsp)
	var info PostParamTx
	errjson := g.BindJSON(&info)
	if errjson != nil {
		rsp.Msg = "json faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	chainid, _ := strconv.Atoi(info.ChainId)
	filter := bson.D{{"chainid", chainid}}
	if info.ChainId == "" {
		filter = bson.D{}
	}

	cur, err := c.Mgo.Collection(m.ChainStats).Find(context.Background(), filter)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "faild"
		g.JSON(http.StatusOK, rsp)
		return
	}
	var chainStats = make([]ChainStat, 0)
	for cur.Next(context.Background()) {
		elem := new(ChainStat)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		chainStats = append(chainStats, *elem)
	}
	rsp.Data = chainStats
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

/* 根据 链ID 和交易类别 查询交易   */
func (c *Chain) GetTxByTxTypeAndChainId(g *gin.Context) {
	rsp := new(Rsp)
	var info PostParamTx
	err := g.BindJSON(&info)
	if err != nil {
		rsp.Msg = "json faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	optCount := new(options.CountOptions)

	page := info.Page
	pageSize := info.Pagesize

	if info.Page == 0 {
		page = 1
	}
	if info.Pagesize == 0 {
		pageSize = 10
	}

	chainidstr := info.ChainId
	chainid, err := strconv.Atoi(chainidstr)

	data := new(TxList)
	data.TransactionsList.PageNum = page
	data.TransactionsList.PageSize = pageSize

	filter := bson.D{}
	if info.TxType != "" {
		filter = bson.D{{"txtype", info.TxType}}
		optCount.Hint = Index_txtype
	}
	if info.ChainId != "" && info.TxType == "" { //不查
		filter = bson.D{{"chainid", chainid}}
		optCount.Hint = Index_chainid
	}
	if info.TxType != "" && info.ChainId != "" {
		filter = bson.D{{"chainid", chainid}, {"txtype", info.TxType}}
		optCount.Hint = Index_txtype_chainid
	}

	opts := new(options.FindOptions)
	limit := int64(pageSize)
	skip := int64((page - 1) * pageSize)

	sortMap := make(map[string]interface{})
	sortMap["timestamp"] = -1
	opts.Sort = sortMap

	opts.Limit = &limit
	opts.Skip = &skip
	Txs, err := c.Mgo.Collection(m.BlockTxs).Find(context.Background(), filter, opts)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "get data error"
		g.JSON(http.StatusOK, rsp)
		return
	}
	data.TransactionsList.DataList = make([]TxInfo, 0)
	for Txs.Next(context.Background()) {
		elem := new(TxInfo)
		err := Txs.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		data.TransactionsList.DataList = append(data.TransactionsList.DataList, *elem)
	}
	var sumCount int64
	if len(filter) == 0 {
		sumCount = getEstimatedDocumentCount(c, m.BlockTxs)
	} else {
		sCount, err := c.Mgo.Collection(m.BlockTxs).CountDocuments(context.Background(), filter, optCount)
		if err == nil {
			sumCount = sCount
		}
	}

	data.TransactionsList.Pages = int(math.Ceil(float64(sumCount) / float64(pageSize)))
	data.TransactionsList.Total = int(sumCount)
	rsp.Data = data
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

/* 根据 链ID 和合约hash 查询交易   */
func (c *Chain) GetTxByContractAndChainId(g *gin.Context) {
	rsp := new(Rsp)
	var info PostParamTx
	err := g.BindJSON(&info)
	if err != nil {
		rsp.Msg = "json faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	optCount := new(options.CountOptions)
	optCount.Hint = Index_id
	page := info.Page
	pageSize := info.Pagesize

	if info.Page == 0 {
		page = 1
	}
	if info.Pagesize == 0 {
		pageSize = 10
	}

	chainidstr := info.ChainId
	chainid, err := strconv.Atoi(chainidstr)

	data := new(TxList)
	data.TransactionsList.PageNum = page
	data.TransactionsList.PageSize = pageSize

	filter := bson.D{{"chainid", chainid}}
	if info.Contract != "" && info.ChainId != "" {
		filter = bson.D{{"chainid", chainid}, {"to", info.Contract}}
		optCount.Hint = Index_to
	}

	opts := new(options.FindOptions)
	limit := int64(pageSize)
	skip := int64((page - 1) * pageSize)

	sortMap := make(map[string]interface{})
	sortMap["timestamp"] = -1
	opts.Sort = sortMap

	opts.Limit = &limit
	opts.Skip = &skip
	Txs, err := c.Mgo.Collection(m.BlockTxs).Find(context.Background(), filter, opts)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "get data error"
		g.JSON(http.StatusOK, rsp)
		return
	}
	var flag int
	data.TransactionsList.DataList = make([]TxInfo, 0)
	for Txs.Next(context.Background()) {
		elem := new(TxInfo)
		err := Txs.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		data.TransactionsList.DataList = append(data.TransactionsList.DataList, *elem)
		flag++
	}
	var sumCount int64
	if len(filter) == 0 {
		sumCount = getEstimatedDocumentCount(c, m.BlockTxs)
	} else {
		sCount, err := c.Mgo.Collection(m.BlockTxs).CountDocuments(context.Background(), filter, optCount)
		if err == nil {
			sumCount = sCount
		}
	}

	data.TransactionsList.Pages = int(math.Ceil(float64(sumCount) / float64(pageSize)))
	data.TransactionsList.Total = int(sumCount)
	rsp.Data = data
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

/* 获取主链链详情 */
func (c *Chain) GetMainChainStat(g *gin.Context) {
	rsp := new(Rsp)
	var info PostParamTx
	errjson := g.BindJSON(&info)
	if errjson != nil {
		rsp.Msg = "json faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	if info.ChainId == "" {
		rsp.Msg = "chainId is  empty"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	chainidstr := info.ChainId
	chainid, _ := strconv.Atoi(chainidstr)
	result := new(ChainStat)

	filter := bson.D{{"chainid", chainid}}

	err := c.Mgo.Collection(m.ChainStats).FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "faild"
		g.JSON(http.StatusOK, rsp)
		return
	}
	if result.ChainId == 0 { //子链
		count, err := c.Mgo.Collection(m.ChainInfo).CountDocuments(context.Background(), bson.D{{"parent", result.ChainId}})
		if err == nil {
			result.ChildrenCount = int(count) //子链个数
		}
	} else { //分片
		count, err := c.Mgo.Collection(m.ChainInfo).CountDocuments(context.Background(), bson.D{{"parent", result.ChainId}})
		if err == nil {
			result.ChildrenCount = int(count) //分片个数
		}
	}

	rsp.Data = result
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

func (c *Chain) GetAllTxChain(g *gin.Context) {
	// 为将来增加用户准备
	address := g.Query("address")
	chainType := g.Query("chainType")
	rsp := new(Rsp)
	data := new(ChainList)
	var chainids []int
	// var countFilter bson.M
	chainidFilter := bson.M{"address": address}
	if chainType != "" {
		chainidFilter["chaintype"] = "admin"
	}
	cur, err := c.Mgo.Collection(m.UserChaininfo).Find(context.Background(), chainidFilter)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "get data error"
		g.JSON(http.StatusOK, rsp)
		return
	}
	for cur.Next(context.Background()) {
		elem := new(UserChainInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		chainids = append(chainids, int(elem.ChainId))
	}

	countFilter := bson.M{"chainid": bson.M{"$in": chainids}}
	count, _ := c.Mgo.Collection(m.ChainStats).CountDocuments(context.Background(), countFilter)
	//log.Info("count", count)
	opts := new(options.FindOptions)
	limit := int64(6)
	// skip := int64((page - 1) * 6)
	opts.Limit = &limit
	sortMap := make(map[string]interface{})
	sortMap["chainid"] = 1
	opts.Sort = sortMap
	// opts.Skip = &(skip)
	chaincur, err := c.Mgo.Collection(m.ChainStats).Find(context.Background(), countFilter)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "get data error"
		g.JSON(http.StatusOK, rsp)
		return
	}

	for chaincur.Next(context.Background()) {
		chainelem := new(ChainInfo)
		err := chaincur.Decode(chainelem)
		if err != nil {
			log.Error(err.Error())
		}
		if chainelem.ChainId == 0 {
			continue
		}
		data.Chainlist.DataList = append(data.Chainlist.DataList, *chainelem)
	}

	data.Chainlist.Total = int(count)
	data.Chainlist.PageNum = 1
	data.Chainlist.PageSize = 6
	data.Chainlist.Pages = (6 + int(count)) / 6
	rsp.Code = CODE_Success
	rsp.Data = data
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

/* 获取交易信息*/
func (c *Chain) GetBlockNewTx(g *gin.Context) {
	page, err := strconv.Atoi(g.PostForm("page"))

	if err != nil {
		log.Error(err.Error())
	}
	chainidstr := g.PostForm("chainId")
	chainid, err := strconv.Atoi(chainidstr)

	rsp := new(Rsp)
	// start := page - 1
	// end := start*5 + 5
	data := new(TxList)
	data.TransactionsList.PageNum = page

	// result := c.RedisCli.Client.LRange("blockTxs-"+chainid, int64(start), int64(end)).Val()
	filter := bson.D{{"chainid", chainid}}

	// for _, item := range result {
	// 	elem := new(TxInfo)
	// 	json.Unmarshal([]byte(item), elem)
	// 	data.TransactionsList.DataList = append(data.TransactionsList.DataList, *elem)
	// }
	opts := new(options.FindOptions)
	limit := int64(6)
	skip := int64((page - 1) * 6)

	sortMap := make(map[string]interface{})
	sortMap["height"] = -1
	opts.Sort = sortMap

	opts.Limit = &limit
	opts.Skip = &skip
	Txs, err := c.Mgo.Collection(m.BlockTxs).Find(context.Background(), filter, opts)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "get data error"
		g.JSON(http.StatusOK, rsp)
		return
	}
	for Txs.Next(context.Background()) {
		elem := new(TxInfo)
		err := Txs.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		data.TransactionsList.DataList = append(data.TransactionsList.DataList, *elem)
	}
	log.Infof("GetBlockNewTx count:%s", len(data.TransactionsList.DataList))
	rsp.Data = data
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

/* 获取交易信息  分页时间倒叙*/
func (c *Chain) GetBlockNewTxPage(g *gin.Context) {
	rsp := new(Rsp)
	var info PostParamTx
	err := g.BindJSON(&info)
	if err != nil {
		rsp.Msg = "json faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	optCount := new(options.CountOptions)
	optCount.Hint = Index_id //默认包含chainId

	page := info.Page
	pageSize := info.Pagesize

	chainidstr := info.ChainId
	chainid, err := strconv.Atoi(chainidstr)

	data := new(TxList)
	data.TransactionsList.PageNum = page
	data.TransactionsList.PageSize = pageSize

	filter := bson.D{{"chainid", chainid}}
	optCount.SetHint(Index_chainid)
	height, _ := strconv.Atoi(info.Height)

	if info.ChainId == "" {
		filter = bson.D{}
	}
	//直接hash
	if info.Hash != "" {
		filter = bson.D{{"hash", info.Hash}}
		optCount.Hint = Index_hash
	}

	//後加的test
	if info.Height != "" {
		filter = bson.D{{"height", height}}
		optCount.Hint = Index_height_chainid
	}

	//直接hash和ChainID
	if info.ChainId != "" && info.Hash != "" {
		filter = bson.D{{"chainid", chainid}, {"hash", info.Hash}}
		optCount.Hint = Index_hash_chainid
	}
	//直接Chainid和Height
	if info.ChainId != "" && info.Height != "" {
		filter = bson.D{{"chainid", chainid}, {"height", height}}
		optCount.Hint = Index_height_chainid
	}

	opts := new(options.FindOptions)
	limit := int64(pageSize)
	skip := int64((page - 1) * pageSize)
	sortMap := make(map[string]interface{})
	sortMap["timestamp"] = -1
	opts.Sort = sortMap
	opts.Limit = &limit
	opts.Skip = &skip

	Txs, err := c.Mgo.Collection(m.BlockTxs).Find(context.Background(), filter, opts)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "get data error"
		g.JSON(http.StatusOK, rsp)
		return
	}
	var flag int
	data.TransactionsList.DataList = make([]TxInfo, 0)
	for Txs.Next(context.Background()) {
		elem := new(TxInfo)
		err := Txs.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		//模拟关联交易
		elem.ConnectedTx.TxType = ""
		elem.ConnectedTx.Hash = ""
		if elem.To == QK {
			elem.Value = float64(getAmount(elem.Input))
		}
		data.TransactionsList.DataList = append(data.TransactionsList.DataList, *elem)
		flag++
	}

	var sumCount int64
	if len(filter) == 0 {
		sumCount = getEstimatedDocumentCount(c, m.BlockTxs) // c.Mgo.Collection(m.BlockTxs).EstimatedDocumentCount(context.Background(),options.EstimatedDocumentCount().SetMaxTime(90000000000000))
	} else {
		if len(filter) == 1 {
			if optCount.Hint == Index_chainid {
				sumCount = getTxCountByChainId(c, chainid)
			}
		} else {
			sCount, err := c.Mgo.Collection(m.BlockTxs).CountDocuments(context.Background(), filter, optCount)
			if err == nil {
				sumCount = sCount
			}
		}
	}

	data.TransactionsList.Pages = int(math.Ceil(float64(sumCount) / float64(pageSize)))
	data.TransactionsList.Total = int(sumCount)

	rsp.Data = data
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

/**
* get table doc all count
 */
func getEstimatedDocumentCount(c *Chain, tableName string) int64 {
	count, _ := c.Mgo.Collection(tableName).EstimatedDocumentCount(context.Background(), options.EstimatedDocumentCount().SetMaxTime(9000000000000))
	return count
}

//获取链同步回来的交易个数
func getTxCountByChainId(c *Chain, chainId int) int64 {
	filter := bson.D{{"chainid", chainId}}
	var info = new(HeightTx)
	err := c.Mgo.Collection(m.HeightTx).FindOne(context.Background(), filter).Decode(&info)
	if err == nil {
		return int64(info.TxCount)
	}
	return 0

}

/* 获取交易信息  分页时间倒叙*/
func (c *Chain) GetBlockTxByAddress(g *gin.Context) {
	rsp := new(Rsp)
	var info PostParamTx
	err := g.BindJSON(&info)
	if err != nil {
		rsp.Msg = "json faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	optCount := new(options.CountOptions)

	page := info.Page
	pageSize := info.Pagesize

	chainidstr := info.ChainId
	chainid, err := strconv.Atoi(chainidstr)

	data := new(TxList)
	data.TransactionsList.PageNum = page
	data.TransactionsList.PageSize = pageSize

	filter := bson.D{{"chainid", chainid}}

	if info.ChainId != "" && info.Address != "" {
		filter = bson.D{{"chainid", chainid}, {"$or", []interface{}{bson.D{{"from", info.Address}}, bson.D{{"to", info.Address}}}}}
	}

	opts := new(options.FindOptions)
	limit := int64(pageSize)
	skip := int64((page - 1) * pageSize)

	sortMap := make(map[string]interface{})
	sortMap["timestamp"] = -1
	opts.Sort = sortMap

	opts.Limit = &limit
	opts.Skip = &skip
	Txs, err := c.Mgo.Collection(m.BlockTxs).Find(context.Background(), filter, opts)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "get data error"
		g.JSON(http.StatusOK, rsp)
		return
	}
	data.TransactionsList.DataList = make([]TxInfo, 0)
	for Txs.Next(context.Background()) {
		elem := new(TxInfo)
		err := Txs.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		data.TransactionsList.DataList = append(data.TransactionsList.DataList, *elem)
	}

	filter = bson.D{{"from", info.Address}, {"chainid", chainid}}
	sumCount, err := c.Mgo.Collection(m.BlockTxs).CountDocuments(context.Background(), filter, optCount.SetHint(Index_from_chainid))

	filter = bson.D{{"to", info.Address}, {"chainid", chainid}}
	sumCount1, _ := c.Mgo.Collection(m.BlockTxs).CountDocuments(context.Background(), filter, optCount.SetHint(Index_to_chainid))

	sumCount = sumCount + sumCount1
	if err != nil {
		log.Error("get sumcount err:" + err.Error())
	}
	//sumCount:=10 //默认10条

	data.TransactionsList.Pages = int(math.Ceil(float64(sumCount) / float64(pageSize)))
	data.TransactionsList.Total = int(sumCount)

	rsp.Data = data
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

/* 获取交易信息  分页时间倒叙*/
func (c *Chain) GetBlockDataInfo(g *gin.Context) {
	rsp := new(Rsp)
	var info PostParamTx
	err := g.BindJSON(&info)
	if err != nil {
		rsp.Msg = "json faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	//允许height chainid hash 参数为空
	heightstr := info.Height
	height, _ := strconv.Atoi(heightstr)
	chainidstr := info.ChainId
	chainid, _ := strconv.Atoi(chainidstr)

	hash := info.Hash

	page := info.Page
	pageSize := info.Pagesize

	/* 默认 */
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	data := new(AllBlockData)
	data.PageSize = pageSize
	data.PageNum = page

	opts := new(options.FindOptions)
	limit := int64(pageSize)
	skip := int64((page - 1) * pageSize)

	sortMap := make(map[string]interface{})
	sortMap["timestamp"] = -1
	opts.Sort = sortMap

	opts.Limit = &limit
	opts.Skip = &skip

	filter := bson.D{{"chainid", chainid}}

	if chainidstr == "" {
		filter = bson.D{}
	}

	if heightstr != "" {
		filter = bson.D{{"height", height}}
	}
	if hash != "" && chainidstr == "" {
		filter = bson.D{{"hash", hash}}
	}
	if hash != "" && chainidstr != "" {
		filter = bson.D{{"hash", hash}, {"chainid", chainid}}
	}
	if heightstr != "" && chainidstr != "" {
		filter = bson.D{{"height", height}, {"chainid", chainid}}
	}
	if hash != "" && chainidstr != "" && heightstr != "" {
		filter = bson.D{{"height", height}, {"chainid", chainid}, {"hash", hash}}
	}

	cur, err := c.Mgo.Collection(m.BlockData).Find(context.Background(), filter, opts)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "faild"
		g.JSON(http.StatusOK, rsp)
		return
	}
	data.DataList = make([]BlockInfo, 0)
	for cur.Next(context.Background()) {
		elem := new(BlockInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		tm := time.Unix(int64(elem.Timestamp), 0)
		elem.TheDateString = tm.Format("2006-01-02 15:04:05")
		//elem.NextHeight=elem.Height+1
		data.DataList = append(data.DataList, *elem)
		break //重复数据时 跳出
	}

	rsp.Data = data.DataList
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

func (c *Chain) GetTxTypeByHeight(g *gin.Context) {
	rsp := new(Rsp)
	var info PostParamTx
	err := g.BindJSON(&info)
	if err != nil {
		rsp.Msg = "json faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	//允许height chainid hash 参数为空
	heightstr := info.Height
	height, _ := strconv.Atoi(heightstr)
	chainidstr := info.ChainId
	chainid, _ := strconv.Atoi(chainidstr)
	filter := bson.D{{"chainid", chainid}, {"height", height}}

	//c.Mgo.Collection(m.BlockTxs).Find(context.Background(), filter)

	cur, err := c.Mgo.Collection(m.BlockTxs).Find(context.Background(), filter)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "faild"
		g.JSON(http.StatusOK, rsp)
		return
	}
	var data = new(HeightTxInfo)
	infos := make([]TxInfo, 0)
	for cur.Next(context.Background()) {
		elem := new(TxInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		infos = append(infos, *elem)
	}
	data.TxtSumCount = len(infos)

	var type1 = new(TxTypeInfo)
	type1.TxType = TxType1
	var type2 = new(TxTypeInfo)
	type2.TxType = TxType2
	var type3 = new(TxTypeInfo)
	type3.TxType = TxType3
	var type4 = new(TxTypeInfo)
	type4.TxType = TxType4
	var type5 = new(TxTypeInfo)
	type5.TxType = TxType5
	var type6 = new(TxTypeInfo)
	type6.TxType = TxType6

	for i := 0; i < len(infos); i++ {
		result := infos[i].TxType //getTxType(infos[i])
		if result == TxType1 {
			type1.TxCount++
		}
		if result == TxType2 {
			type2.TxCount++
		}
		if result == TxType3 {
			type3.TxCount++
		}
		if result == TxType4 {
			type4.TxCount++
		}
		if result == TxType5 {
			type5.TxCount++
		}
		if result == TxType6 {
			type6.TxCount++
		}
	}
	data.TxByTypeInfo = append(data.TxByTypeInfo, *type1, *type2, *type3, *type4, *type5, *type6)
	rsp.Data = data
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

/* 块数据 分页 */
func (c *Chain) GetBlockDataByPage(g *gin.Context) {

	optCount := new(options.CountOptions)
	optCount.Hint = Index_id
	rsp := new(Rsp)
	var info PostParamTx
	err := g.BindJSON(&info)
	if err != nil {
		rsp.Msg = "json faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	//允许height chainid hash 参数为空
	heightstr := info.Height
	height, _ := strconv.Atoi(heightstr)
	chainidstr := info.ChainId
	chainid, _ := strconv.Atoi(chainidstr)

	hash := info.Hash

	page := info.Page
	pageSize := info.Pagesize

	/* 默认 */
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	data := new(AllBlockData)
	data.PageSize = pageSize
	data.PageNum = page

	opts := new(options.FindOptions)
	limit := int64(pageSize)
	skip := int64((page - 1) * pageSize)

	sortMap := make(map[string]interface{})
	sortMap["timestamp"] = -1
	opts.Sort = sortMap

	opts.Limit = &limit
	opts.Skip = &skip

	filter := bson.D{{"chainid", chainid}}

	if chainidstr == "" {
		filter = bson.D{}
	}

	if heightstr != "" {
		filter = bson.D{{"height", height}}
		optCount.Hint = Index_height_chainid
	}
	if hash != "" && chainidstr == "" {
		filter = bson.D{{"hash", hash}}
		optCount.Hint = Index_hash
	}
	if hash != "" && chainidstr != "" {
		filter = bson.D{{"hash", hash}, {"chainid", chainid}}
		optCount.Hint = Index_hash_chainid
	}
	if heightstr != "" && chainidstr != "" {
		filter = bson.D{{"height", height}, {"chainid", chainid}}
		optCount.Hint = Index_height_chainid
	}
	if hash != "" && chainidstr != "" && heightstr != "" {
		filter = bson.D{{"height", height}, {"chainid", chainid}, {"hash", hash}}
		optCount.Hint = Index_hash
	}

	cur, err := c.Mgo.Collection(m.BlockData).Find(context.Background(), filter, opts)
	//cur, err := c.Mgo.Collection("blockData_bak201907011700").Find(context.Background(), filter, opts)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "faild"
		g.JSON(http.StatusOK, rsp)
		return
	}
	data.DataList = make([]BlockInfo, 0)
	for cur.Next(context.Background()) {
		elem := new(BlockInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		tm := time.Unix(int64(elem.Timestamp), 0)
		elem.TheDateString = tm.Format("2006-01-02 15:04:05")
		elem.ChainIdCommitteeCount = getCommitteeCount(c, elem.Chainid)
		data.DataList = append(data.DataList, *elem)
	}

	var sumCount int64
	if len(filter) == 0 {
		sumCount = getEstimatedDocumentCount(c, m.BlockData)
	} else {
		sCount, _ := c.Mgo.Collection(m.BlockData).CountDocuments(context.Background(), filter, optCount)
		sumCount = sCount
	}

	data.Pages = int(math.Ceil(float64(sumCount) / float64(pageSize))) //页数向上取整
	data.Total = int(sumCount)
	rsp.Data = data
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

/* 根据参选轮次 块数据 分页 */
func (c *Chain) GetBlockDataByEpoch(g *gin.Context) {
	rsp := new(Rsp)
	var info PostParamTx
	err := g.BindJSON(&info)
	if err != nil {
		rsp.Msg = "json faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	//Epoch chainid hash 参数为空
	Epochstr := info.Epoch
	epoch, _ := strconv.Atoi(Epochstr)
	chainidstr := info.ChainId
	chainid, _ := strconv.Atoi(chainidstr)

	page := info.Page
	pageSize := info.Pagesize

	/* 默认 */
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	data := new(AllBlockData)
	data.PageSize = pageSize
	data.PageNum = page

	opts := new(options.FindOptions)
	limit := int64(pageSize)
	skip := int64((page - 1) * pageSize)

	sortMap := make(map[string]interface{})
	sortMap["height"] = -1
	opts.Sort = sortMap

	opts.Limit = &limit
	opts.Skip = &skip

	//filter := bson.D{{"chainid", chainid}}

	startheight := epoch * 300
	endheight := startheight + 300

	filter := bson.D{{"chainid", chainid}, {"height", bson.M{"$gte": startheight, "$lt": endheight}}}

	cur, err := c.Mgo.Collection(m.BlockData).Find(context.Background(), filter, opts)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "faild"
		g.JSON(http.StatusOK, rsp)
		return
	}
	data.DataList = make([]BlockInfo, 0)
	for cur.Next(context.Background()) {
		elem := new(BlockInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		tm := time.Unix(int64(elem.Timestamp), 0)
		elem.TheDateString = tm.Format("2006-01-02 15:04:05")
		elem.ChainIdCommitteeCount = getCommitteeCount(c, elem.Chainid)
		data.DataList = append(data.DataList, *elem)
	}

	optCount := new(options.CountOptions)
	optCount.Hint = Index_height_chainid

	sumCount, _ := c.Mgo.Collection(m.BlockData).CountDocuments(context.Background(), filter, optCount)
	data.Pages = int(math.Ceil(float64(sumCount) / float64(pageSize))) //页数向上取整
	data.Total = int(sumCount)

	rsp.Data = data
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

///* 获取关联交易 */
//func getConnectedTx(c Chain, txInfo TxInfo) ConnectedTx {
//
//	resultInfo := new(ConnectedTx)
//	resultInfo.ChainId = 0
//	resultInfo.TxType = TxType5
//	resultInfo.Hash = ""
//
//	/* 从库里查对应的关系才行  才行*/
//	if txInfo.TxType == TxType4 || txInfo.TxType == TxType5 || txInfo.TxType == TxType6 {
//		//这种
//		filter := bson.D{{"hashqk", txInfo.Hash}}
//		var info = new(CashCheckHash)
//		err := c.Mgo.Collection(m.CashCheckHash).FindOne(context.Background(), filter).Decode(&info)
//		if err != nil {
//			log.Error(err.Error())
//		}
//
//	}
//
//	return *resultInfo
//}

/* 块数据 分页 */
func (c *Chain) GetBlockData(g *gin.Context) {

	rsp := new(Rsp)
	page, _ := strconv.Atoi(g.PostForm("page"))
	pageSize, _ := strconv.Atoi(g.PostForm("pagesize"))

	/* 默认 */
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	data := new(AllBlockData)
	data.PageSize = pageSize
	data.PageNum = page

	opts := new(options.FindOptions)
	limit := int64(pageSize)
	skip := int64((page - 1) * pageSize)

	sortMap := make(map[string]interface{})
	sortMap["timestamp"] = -1
	opts.Sort = sortMap

	opts.Limit = &limit
	opts.Skip = &skip

	filter := bson.D{}

	//filter := bson.D{{"chainid", chainid}}
	//
	//if chainidstr == "" {
	//	filter = bson.D{}
	//}
	//
	//if heightstr != "" {
	//	filter = bson.D{{"height", height}}
	//}
	//if hash != "" && chainidstr == "" {
	//	filter = bson.D{{"hash", hash}}
	//}
	//if hash != "" && chainidstr != "" {
	//	filter = bson.D{{"hash", hash}, {"chainid", chainid}}
	//}
	//if heightstr != "" && chainidstr != "" {
	//	filter = bson.D{{"height", height}, {"chainid", chainid}}
	//}
	//if hash!=""&&chainidstr!=""&&heightstr!=""{
	//	filter = bson.D{{"height", height}, {"chainid", chainid}, {"hash", hash}}
	//}

	cur, err := c.Mgo.Collection(m.BlockData).Find(context.Background(), filter, opts)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "faild"
		g.JSON(http.StatusOK, rsp)
		return
	}
	data.DataList = make([]BlockInfo, 0)
	for cur.Next(context.Background()) {
		elem := new(BlockInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}
		tm := time.Unix(int64(elem.Timestamp), 0)
		elem.TheDateString = tm.Format("2006-01-02 15:04:05")
		elem.ChainIdCommitteeCount = getCommitteeCount(c, elem.Chainid)
		data.DataList = append(data.DataList, *elem)
	}

	rsp.Data = data.DataList
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}
func getCommitteeCount(c *Chain, chainId int) int {
	filter := bson.D{{"chainid", chainId}}
	var info = new(ChainStat)
	err := c.Mgo.Collection(m.ChainStats).FindOne(context.Background(), filter).Decode(&info)
	if err != nil {
		return 0
	} else {
		return len(info.Currentcomm)
	}
}

/* 获取块信息 */
func (c *Chain) GetNewBlock(g *gin.Context) {
	// chainid := g.PostForm("chainId")
	// rsp := new(Rsp)
	// data := new(BlockList)
	// result := c.RedisCli.Client.LRange("blockData-"+chainid, 0, 5).Val()
	// for _, item := range result {
	// 	elem := new(BlockInfo)
	// 	json.Unmarshal([]byte(item), elem)
	// 	data.Blocklist = append(data.Blocklist, *elem)
	// }
	// rsp.Code = CODE_Success
	// rsp.Msg = "success"
	// rsp.Data = data
	// g.JSON(http.StatusOK, rsp)

	// filter := bson.D{{"chainid", chainid}}

	// for _, item := range result {
	// 	elem := new(TxInfo)
	// 	json.Unmarshal([]byte(item), elem)
	// 	data.TransactionsList.DataList = append(data.TransactionsList.DataList, *elem)
	// }

	rsp := new(Rsp)
	data := new(BlockList)
	chainidstr := g.PostForm("chainId")
	chainid, err := strconv.Atoi(chainidstr)

	filter := bson.D{{"chainid", chainid}}

	opts := new(options.FindOptions)
	limit := int64(6)
	skip := int64((1 - 1) * 6) //
	opts.Limit = &limit
	opts.Skip = &skip

	sortMap := make(map[string]interface{})
	sortMap["height"] = -1
	opts.Sort = sortMap

	Txs, err := c.Mgo.Collection(m.BlockData).Find(context.Background(), filter, opts)
	if err != nil {
		rsp.Code = CODE_Error
		rsp.Msg = "get data error"
		g.JSON(http.StatusOK, rsp)
		return
	}
	for Txs.Next(context.Background()) {
		elem := new(BlockInfo)
		err := Txs.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}

		tm := time.Unix(int64(elem.Timestamp), 0)
		elem.TheDateString = tm.Format("2006-01-02 15:04:05")

		data.Blocklist = append(data.Blocklist, *elem)
	}
	//data.Blocklist=sortBlock(data.Blocklist)
	//log.Info("GetNewBlock count:", len(data.Blocklist))
	rsp.Data = data
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)

	return
}

/*

func (c *Chain) GetBlockTransactions(g *gin.Context) {
	heightstr := g.PostForm("height")
	height, _ := strconv.Atoi(heightstr)

	chainidstr := g.PostForm("chainId")
	chainid, _ := strconv.Atoi(chainidstr)
	pagestr := g.PostForm("page")
	if pagestr == "" {
		pagestr = "1"
	}
	page, _ := strconv.Atoi(pagestr)
	sizestr := g.PostForm("size")
	size10, _ := strconv.Atoi(sizestr)
	size64 := int64(size10)
	skip := (page - 1) * size10
	skip64 := int64(skip)
	rsp := new(Rsp)

	data := new(BlockTxs)
	filter := bson.D{{"height", height}, {"chainid", chainid}}
	opts := new(options.FindOptions)
	opts.Limit = &size64
	opts.Skip = &(skip64)
	// sortMap := make(map[string]interface{})
	cur, err := c.Mgo.Collection(m.BlockTxs).Find(context.Background(), filter, opts)
	if err != nil {
		g.JSON(http.StatusOK, gin.H{"err": "Error"})
		return
	}
	for cur.Next(context.Background()) {
		elem := new(TxInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}

		data.TransactionsList = append(data.TransactionsList, *elem)
	}

	optCount := new(options.CountOptions)
	optCount.Hint=Index_height_chainid

	count, _ := c.Mgo.Collection(m.BlockTxs).CountDocuments(context.Background(), filter,optCount)
	data.Txcount = int(count)
	rsp.Data = data
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return

}

func (c *Chain) GetBlockHeader(g *gin.Context) {
	heightstr := g.PostForm("height")
	height, _ := strconv.Atoi(heightstr)
	chainidstr := g.PostForm("chainId")
	chainid, _ := strconv.Atoi(chainidstr)

	rsp := new(Rsp)
	data := new(BlockInfo)
	filter := bson.D{{"height", height}, {"chainid", chainid}}
	err := c.Mgo.Collection(m.BlockData).FindOne(context.Background(), filter).Decode(&data)
	if err != nil {
		g.JSON(http.StatusOK, gin.H{"a": "b"})
		return
	}
	rsp.Data = BlockOne{*data}
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

func (c *Chain) GetTransactionsByAddress(g *gin.Context) {
	address := g.PostForm("address")
	page, _ := strconv.Atoi(g.PostForm("page"))
	rsp := new(Rsp)
	data := new(TxList)
	filter := bson.D{{"from", address}}

	opts := new(options.FindOptions)
	limit := int64(6)
	skip := int64((page - 1) * 6)
	opts.Limit = &limit
	sortMap := make(map[string]interface{})
	sortMap["timestamp"] = -1
	opts.Sort = sortMap
	opts.Skip = &(skip)
	cur, err := c.Mgo.Collection(m.BlockTxs).Find(context.Background(), filter, opts)
	if err != nil {
		g.JSON(http.StatusOK, gin.H{"err": "Error"})
		return
	}
	for cur.Next(context.Background()) {
		elem := new(TxInfo)
		err := cur.Decode(elem)
		if err != nil {
			log.Error(err.Error())
		}

		data.TransactionsList.DataList = append(data.TransactionsList.DataList, *elem)
	}


	optCount := new(options.CountOptions)
	optCount.Hint="_id"

	count, _ := c.Mgo.Collection(m.BlockTxs).CountDocuments(context.Background(), filter,optCount)
	data.TransactionsList.PageNum = page
	data.TransactionsList.PageSize = 6
	data.TransactionsList.Total = int(count)
	data.TransactionsList.Pages = data.TransactionsList.Total / 6
	rsp.Data = data
	rsp.Code = CODE_Success
	rsp.Msg = "success"
	g.JSON(http.StatusOK, rsp)
	return
}

*/

/* 查页面钱包用户(用户A,用户B) */
func (c *Chain) GetUserPayInfo(g *gin.Context) {
	rsp := new(Rsp)
	cur, err := c.Mgo.Collection(m.UserPayInfo).Find(context.Background(), bson.D{})
	if err != nil {
		rsp.Msg = err.Error()
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	var userpayinfos []UserPayInfo
	for cur.Next(context.Background()) {
		elme := new(UserPayInfo)
		err := cur.Decode(elme)
		if err == nil {
			userpayinfos = append(userpayinfos, *elme)
		}
	}
	rsp.Msg = "success"
	rsp.Data = userpayinfos
	rsp.Code = CODE_Success
	g.JSON(http.StatusOK, rsp)
	return
}

/* paylog  */
func (c *Chain) GetUserPayLog(g *gin.Context) {

	rsp := new(Rsp)
	opts := new(options.FindOptions)
	sortMap := make(map[string]interface{})
	sortMap["paytime"] = -1
	opts.Sort = sortMap

	cur, err := c.Mgo.Collection(m.UserPayInfoLog).Find(context.Background(), bson.D{}, opts.SetLimit(int64(6)))
	if err != nil {
		rsp.Msg = err.Error()
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	var paylogs []UserPayInfoLog
	for cur.Next(context.Background()) {
		log := new(UserPayInfoLog)
		err := cur.Decode(log)
		if err == nil {
			paylogs = append(paylogs, *log)
		}
	}
	rsp.Msg = "success"
	rsp.Data = paylogs
	rsp.Code = CODE_Success
	g.JSON(http.StatusOK, rsp)
	return
}

var Payqueue = list.New()

/* 页面转账 */
func (c *Chain) UserPay(g *gin.Context) {
	rsp := new(Rsp)
	fromname := g.PostForm("fromname")
	fromchainid := g.PostForm("fromchainid")
	fromaddress := g.PostForm("fromaddress")
	toaddress := g.PostForm("toaddress")
	paylines := g.PostForm("paylines")
	tochainid := g.PostForm("tochainid")
	//password:=g.PostForm("password")

	if fromchainid == "" || fromaddress == "" || paylines == "" || toaddress == "" || tochainid == "" {
		rsp.Msg = "参数有空值！"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	//1.pay
	res, err := rpc.HTTPSendTX(fromchainid, fromaddress, toaddress, paylines, "")
	fmt.Println("转账交易hash为:", res)
	fmt.Println("err:", err)
	if res == nil {
		rsp.Msg = "faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	paylog := new(UserPayInfoLog)
	for k, v := range res {
		if k == "errMsg" {
			rsp.Msg = "faild"
			rsp.Code = CODE_Error
			g.JSON(http.StatusOK, rsp)
			return
		}
		if k == "TXhash" {
			paylog.PayHash = v.(string)
		}
	}
	//2.save log
	paylog.ChainId = fromchainid
	paylog.PayLines = paylines
	paylog.PayTime = int32(time.Now().Unix())
	paylog.UserName = fromname
	paylog.ToAddress = toaddress
	paylog.FromChainId = fromchainid
	Payqueue.PushBack(*paylog) // 入对
	rsp.Msg = "success"
	rsp.Code = CODE_Success
	g.JSON(http.StatusOK, rsp)
	return
}

/* 获取金额 */
func (c *Chain) GetBalance(g *gin.Context) {
	rsp := new(Rsp)
	address := g.PostForm("address")
	chainid := g.PostForm("chainid")

	if address == "" || chainid == "" {
		rsp.Msg = "address or  chainid is empty"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	res, err := rpc.GetAcc(chainid, strings.TrimSpace(address))
	if err != nil {
		rsp.Msg = err.Error()
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	rsp.Msg = "success"
	rsp.Code = CODE_Success
	int64, err := strconv.ParseInt(res, 10, 64)
	rsp.Data = int64
	g.JSON(http.StatusOK, rsp)
	return
}

//mask 根据分片数运算 fromIndex 分片的起始chainId
func Shard2Chain(address string, mask uint, fromIndex int) string {
	bytes, err := hexutil.Decode(address)
	if err != nil {
		return ""
	}
	num := bytes[0] >> mask
	return strconv.Itoa(int(num) + fromIndex)
}

/* 获取金额   默认查询所有链的   */
func (c *Chain) GetAccountByAddress(g *gin.Context) {
	rsp := new(Rsp)
	var info PostParamTx
	err := g.BindJSON(&info)
	if err != nil {
		rsp.Msg = "json faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	address := info.Address

	chainInfo := getChainInfoStruct(c)
	fmt.Println(chainInfo)
	var queryIds []int
	for i := 0; i < len(chainInfo.Childrens); i++ {
		if len(chainInfo.Childrens[i].Childrens) > 0 { //已经分片的子链
			childrenCount := len(chainInfo.Childrens[i].Childrens)
			fromIndex := chainInfo.Childrens[i].Childrens[0].ChainId
			chainId := Shard2Chain(address, uint(childrenCount), fromIndex)
			fmt.Println(chainId)
			cid, err := strconv.Atoi(chainId)
			if err == nil {
				queryIds = append(queryIds, cid)
			}
		} else {
			queryIds = append(queryIds, chainInfo.Childrens[i].ChainId)
		}
	}

	//filter := bson.D{{}}
	filter := bson.M{"chainid": bson.M{"$in": queryIds}}
	chains := make([]ChainStat, 0)
	cur, err := c.Mgo.Collection(m.ChainStats).Find(context.Background(), filter)
	if err != nil {
		rsp.Msg = "faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	for cur.Next(context.Background()) {
		elme := new(ChainStat)
		err := cur.Decode(elme)
		if err == nil {
			chains = append(chains, *elme)
		}
	}

	var total decimal.Decimal
	infos := new(ResultBalance)
	for i := 0; i < len(chains); i++ {
		var info = new(ChainBalance)
		info.ChainId = int(chains[i].ChainId)
		chainid := strconv.Itoa(info.ChainId)
		res, err := rpc.GetAcc(chainid, strings.TrimSpace(address))
		if err == nil {
			info.Balance = res

			de2, err := decimal.NewFromString(res)

			if err == nil {
				total = total.Add(de2)
			}

		}

		infos.Details = append(infos.Details, *info)
	}
	infos.Balances = total
	rsp.Msg = "success"
	rsp.Code = CODE_Success
	rsp.Data = infos
	g.JSON(http.StatusOK, rsp)
	return
}

/* 获取金额 2 resp返回 */
func (c *Chain) GetAccount(g *gin.Context) {
	rsp := new(Rsp)
	address := g.PostForm("address")
	chainid := g.PostForm("chainid")

	if address == "" || chainid == "" {
		rsp.Msg = "address or chainid is empty "
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	res, err := rpc.GetAccRsp(chainid, strings.TrimSpace(address))
	if err != nil {
		rsp.Msg = err.Error()
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	rsp.Msg = "success"
	rsp.Code = CODE_Success
	//int64, err := strconv.ParseInt(res, 10, 64)
	rsp.Data = res
	g.JSON(http.StatusOK, rsp)
	return
}

/* 主链信息 获取*/
func (c *Chain) getMainInfo(g *gin.Context) {
	rsp := new(Rsp)
	id := g.Query("chainId")
	chainId, err := strconv.Atoi(id)
	if err != nil {
		rsp.Msg = "chainId is error!"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}
	fmt.Println(chainId)
	result := new(ChainInfo)
	filter := bson.D{{"chainid", chainId}}
	err = c.Mgo.Collection(m.ChainStats).FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		rsp.Msg = "faild"
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	} else {
		rsp.Msg = "success"
		rsp.Code = CODE_Success
		rsp.Data = result
		g.JSON(http.StatusOK, rsp)
		return
	}
}

func (c *Chain) GetFile(g *gin.Context) {
	//file:="http://192.168.1.108:8500/static/pdf.js/web/viewer.html?file=http://192.168.1.108:8500/static/file_1.pdf"
	//g.JSON(http.StatusOK, file)
	//return
	locale := g.Query("locale")
	if locale == US {
		fmt.Println(locale)
	} else {
		fmt.Println(locale)
	}

	rsp := new(Rsp)

	cur, err := c.Mgo.Collection(m.Files).Find(context.Background(), bson.D{})
	if err != nil {
		rsp.Msg = err.Error()
		rsp.Code = CODE_Error
		g.JSON(http.StatusOK, rsp)
		return
	}

	data := new(DownLoadFile)
	for cur.Next(context.Background()) {
		elme := new(Files)
		err := cur.Decode(elme)
		if err == nil {
			if strings.Contains(elme.FileName, "SDK") {
				data.SDK = append(data.SDK, *elme)
			} else {
				data.API = append(data.API, *elme)
			}
			//infos = append(infos, *elme)
		}
	}
	rsp.Msg = "success"
	rsp.Data = data
	rsp.Code = CODE_Success
	g.JSON(http.StatusOK, rsp)
	return
}

/*测试 */
func (c *Chain) GetTest(g *gin.Context) {
	input := "0x000000022c7536e3605d9c16a7a3d7b1898e529396a65c23000000000000000b000000039dbcadf3a1027042c52e0ea1cdf2d67a419eab4d00000000000079d90374cbb1"
	byte, err := hexutil.Decode(input)
	if err == nil {
		info := new(CashCheck)
		err := encoding.Unmarshal(byte, info)
		if err == nil {
			fmt.Println(info.Amount)
			c.Mgo.Collection("test").InsertOne(context.Background(), info)
		}
	}

	info := new(CashCheck)
	err1 := c.Mgo.Collection("test").FindOne(context.Background(), bson.D{}).Decode(info)
	if err1 != nil {
		log.Error(err1.Error())
	}
}

func (c *Chain) Test(g *gin.Context) {
	log.Error("Test log ....")

	g.JSON(http.StatusOK, "test")
}

func getAmount(input string) int64 {
	info := new(CashCheck)
	byte, err := hexutil.Decode(input)
	if err == nil {
		err := encoding.Unmarshal(byte, info)
		if err == nil {
			return info.Amount.Int64()
		}
	}
	return info.Amount.Int64()
}

const (
	AddressLength = 20
)

type (
	ChainID uint32
	Height  uint64

	Address [AddressLength]byte

	Addresser interface {
		Address() Address
	}
)

type CashCheck struct {
	FromChain    ChainID  `json:"FromChain"`    // 转出链
	FromAddress  Address  `json:"FromAddr"`     // 转出账户
	Nonce        uint64   `json:"Nonce"`        // 转出账户提交请求时的nonce
	ToChain      ChainID  `json:"ToChain"`      // 目标链
	ToAddress    Address  `json:"ToAddr"`       // 目标账户
	ExpireHeight Height   `json:"ExpireHeight"` // 过期高度，指的是当目标链高度超过（不含）这个值时，这张支票不能被支取，只能退回
	Amount       *big.Int `json:"Amount"`       // 金额
}

// 4字节FromChain + 20字节FromAddress + 8字节Nonce + 4字节ToChain + 20字节ToAddress +
// 8字节ExpireHeight + 1字节len(Amount.Bytes()) + Amount.Bytes()
// 均为BigEndian
func (c *CashCheck) Serialization(w io.Writer) error {
	buf4 := make([]byte, 4)
	buf8 := make([]byte, 8)

	binary.BigEndian.PutUint32(buf4, uint32(c.FromChain))
	_, err := w.Write(buf4)
	if err != nil {
		return err
	}

	_, err = w.Write(c.FromAddress[:])
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint64(buf8, uint64(c.Nonce))
	_, err = w.Write(buf8)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(buf4, uint32(c.ToChain))
	_, err = w.Write(buf4)
	if err != nil {
		return err
	}

	_, err = w.Write(c.ToAddress[:])
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint64(buf8, uint64(c.ExpireHeight))
	_, err = w.Write(buf8)
	if err != nil {
		return err
	}

	buf4 = buf4[:1]
	var mbytes []byte
	if c.Amount != nil {
		mbytes = c.Amount.Bytes()
	}
	buf4[0] = byte(len(mbytes))
	_, err = w.Write(buf4)
	if err != nil {
		return err
	}
	if buf4[0] > 0 {
		_, err = w.Write(mbytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CashCheck) Deserialization(r io.Reader) error {
	buf4 := make([]byte, 4)
	buf8 := make([]byte, 8)

	_, err := r.Read(buf4)
	if err != nil {
		return err
	}
	c.FromChain = ChainID(binary.BigEndian.Uint32(buf4))

	_, err = r.Read(c.FromAddress[:])
	if err != nil {
		return err
	}

	_, err = r.Read(buf8)
	if err != nil {
		return err
	}
	c.Nonce = binary.BigEndian.Uint64(buf8)

	_, err = r.Read(buf4)
	if err != nil {
		return err
	}
	c.ToChain = ChainID(binary.BigEndian.Uint32(buf4))

	_, err = r.Read(c.ToAddress[:])
	if err != nil {
		return err
	}

	_, err = r.Read(buf8)
	if err != nil {
		return err
	}
	c.ExpireHeight = Height(binary.BigEndian.Uint64(buf8))

	buf4 = buf4[:1]
	_, err = r.Read(buf4)
	if err != nil {
		return err
	}
	length := int(buf4[0])

	if length > 0 {
		mbytes := make([]byte, length)
		_, err = r.Read(mbytes)
		if err != nil {
			return err
		}
		c.Amount = new(big.Int)
		c.Amount.SetBytes(mbytes)
	} else {
		c.Amount = big.NewInt(0)
	}

	return nil
}
