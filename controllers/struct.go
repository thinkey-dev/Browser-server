package controllers

import (
	"github.com/shopspring/decimal"
	//"math/big"
)

/* indexname */
const (
	Index_id      = "_id_" //(这个自带的)
	Index_txtype  = "txtype_index"
	Index_hash    = "hash_index"
	Index_to      = "to_index"
	Index_chainid = "chainid_index"

	Index_height_chainid = "height_chainid" //唯一索引(块表是唯一 交易表不唯一)
	Index_timestamp      = "timestamp_index"
	Index_hash_chainid   = "hash_chainid" //唯一索引
	Index_from_to        = "from_to"
	Index_txtype_chainid = "txtype_chainid"

	//块高表 只有timestamp_index &  height_chainid 唯一索引

	Index_from = "from_index"

	Index_from_chainid = "from_chainid"
	Index_to_chainid   = "to_chainid"
)

type HeightTx struct {
	ChainId int `json:"chainId"`
	Height  int `json:"height"`
	TxCount int `json:"txCount"`
}

type Rsp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type StatData struct {
	Chainid int `json:"chainid"`
	Data    interface{}
	TmpData map[string]interface{}
}

type ChainList struct {
	Chainlist AllTxChainData `json:"chainList"`
}

type ChainInfoList struct {
	ChainInfolist AllTxChainData `json:"chainInfoList"`
}

type TxList struct {
	TransactionsList AllTxData `json:"transactionsList"`
}

type BlockTxs struct {
	TransactionsList []TxInfo `json:"transactionsList"`
	Txcount          int      `json:"txcount"`
}

type BlockList struct {
	Blocklist []BlockInfo `json:"blockList"`
}

type AllTxChainData struct {
	PageNum  int         `json:"pageNum"`
	PageSize int         `json:"pageSize"`
	Pages    int         `json:"pages"`
	Total    int         `json:"total"`
	DataList []ChainInfo `json:"dataList"`
}

type AllTxData struct {
	PageNum  int      `json:"pageNum"`
	PageSize int      `json:"pageSize"`
	Pages    int      `json:"pages"`
	Total    int      `json:"total"`
	DataList []TxInfo `json:"dataList"`
}

type AllTxChainCount struct {
	AccountCount  int32   `json:"accountCount"`
	CurrentHeight float64 `json:"currentHeight"`
	MemberCount   int32   `json:"memberCount"`
	Tps           int32   `json:"tps"`
	TxChainCount  int32   `json:"txChainCount"`
	TxCount       int32   `json:"txCount"`
}

type BlockOne struct {
	Block BlockInfo `json:"block"`
}

type ChainInfo struct {
	ChainId           int32    `json:"chainId"`
	Currentheight     float64  `json:"currentheight"`
	Txcount           int32    `json:"txcount"`
	Tps               int32    `json:"tps"`
	TpsLastEpoch      int32    `json:"tpsLastEpoch"`
	Lives             int32    `json:"lives"`
	Accountcount      int32    `json:"accountcount"`
	Epochlength       int32    `json:"epochlength"`
	Epochduration     int32    `json:"epochduration"`
	Lastepochduration int32    `json:"lastepochduration"`
	Currentcomm       []string `json:"currentcomm"`
}

type ChainStat struct {
	ChainId           int32    `json:"chainId"`
	Currentheight     float64  `json:"currentheight"`
	Txcount           int32    `json:"txcount"`
	ChildrenTxcount   int64    `json:"ChildrenTxcount"` //旗下交易量
	ChildrenCount     int      `json:"ChildrenCount"`   //分片数量
	Tps               int32    `json:"tps"`
	TpsLastEpoch      int32    `json:"tpsLastEpoch"`
	Lives             int32    `json:"lives"`
	Accountcount      int32    `json:"accountcount"`
	Epochlength       int32    `json:"epochlength"`
	Epochduration     int32    `json:"epochduration"`
	Lastepochduration int32    `json:"lastepochduration"`
	Currentcomm       []string `json:"currentcomm"`
	//CurrentcommCount      int `json:"currentcommCount"`
	Miner string `json:"miner"`
}
type ChainBalance struct {
	ChainId int    `json:"chainId"`
	Balance string `json:"balance"`
}

type ResultBalance struct {
	Balances decimal.Decimal `json:"balances"`
	Details  []ChainBalance  `json:"details"`
}

type AllBlockData struct {
	PageNum  int         `json:"pageNum"`
	PageSize int         `json:"pageSize"`
	Pages    int         `json:"pages"`
	Total    int         `json:"total"`
	DataList []BlockInfo `json:"dataList"`
}

type BlockInfo struct {
	Hash          string `json:"hash"`
	Previoushash  string `json:"previousHash"`
	Chainid       int    `json:"chainId"`
	Height        int    `json:"height"`
	Mergeroot     string `json:"mergeRoot"`
	Deltaroot     string `json:"deltaRoot"`
	Stateroot     string `json:"stateRoot"`
	Txcount       int    `json:"txcount"`
	Timestamp     int    `json:"timeTamp"`
	TheDateString string `json:"theDateString"`
	/* 添加的字段 */
	Miner       string `json:"miner"`
	MinerReward int    `json:"minerReward"`
	//NextHeight         int `json:"nextHeight"`
	ChainIdCommitteeCount int `json:"chainIdCommitteeCount"`
}

type TxInfo struct {
	ChainId int    `json:"chainId"`
	Height  int    `json:"height"`
	From    string `json:"from"`
	To      string `json:"to"`
	Nonce   int    `json:"nonce"`

	Timestamp     int    `json:"timestamp"`
	Thedatestring string `json:"theDateString"`
	Input         string `json:"input"`
	Hash          string `json:"hash"`
	Status        int    `json:"status"`
	TxType        string `json:"txType"`
	TxCost        int    `json:"txCost"`
	//ConnectedTx   interface{}
	ConnectedTx ConnectedTx `json:"connectedTx"`
	Value       float64     `json:"value"`
	//FromId        int    `json:"fromId"`
	//ToId          int    `json:"toId"`
}

type ConnectedTx struct {
	ChainId int    `json:"chainId"`
	Hash    string `json:"hash"`
	TxType  string `json:"txType"`
}

type Param struct {
	Chainid string `json:"chainid"`
	Address string `json:"address"`
}

type PostData struct {
	Method string `json:"method"`
	Params Param  `json:"params"`
}

type OnLine struct {
	ChainId   int32  `json:"chainId"`
	Address   string `json:"address"`
	MId       int32  `json:"mid"`
	UserId    int32  `json:"userId"`
	WorkTime  int32  `json:"workTime"`
	LeaveTime int32  `json:"leaveTime"`
}

type AllOnlineData struct {
	PageNum  int      `json:"pageNum"`
	PageSize int      `json:"pageSize"`
	Pages    int      `json:"pages"`
	Total    int      `json:"total"`
	DataList []OnLine `json:"dataList"`
}

type AllDownloadRecords struct {
	PageNum  int              `json:"pageNum"`
	PageSize int              `json:"pageSize"`
	Pages    int              `json:"pages"`
	Total    int              `json:"total"`
	DataList []DownloadRecord `json:"dataList"`
}
type AccountRes struct {
	Address     string `json:"address"`
	Balance     int64  `json:"balance"`
	CodeHash    int    `json:"codeHash"`
	Nonce       string `json:"nonce"`
	StorageRoot int    `json:"storageRoot"`
}
type DownloadRecord struct {
	Name         string `json:"name"`
	Uid          int    `json:"uid"`
	ChainId      int    `json:"chainid"`
	Hash         string `json:"hash"`
	Height       int    `json:"height"`
	DataChainId  string `json:"datachainid"`
	DownloadTime int32  `json:"downloadTime"`
}

type OffLine struct {
	Name          string `json:"name"`
	ChainId       int32  `json:"chainId"`
	Uid           int32  `json:"uid"`
	Objid         int32  `json:"objid"`
	Address       string `json:"address"`
	Hash          string `json:"hash"`
	CreateTime    int32  `json:"workTime"`
	License       int    `json:"license"`
	DownloadTimes int32  `json:"downloadtimes"`
	Datainfo      string `json:"datainfo"`
	Datatype      string `json:"datatype"`
	DataFrom      string `json:"datafrom"`
	Datachainid   string `json:"datachainid"`
	Height        int32  `json:"height"`
}

type AllOfflineData struct {
	PageNum  int       `json:"pageNum"`
	PageSize int       `json:"pageSize"`
	Pages    int       `json:"pages"`
	Total    int       `json:"total"`
	DataList []OffLine `json:"dataList"`
}

type RightInfo struct {
	Address string `json:"address"`
	Status  int    `json:"status"`
}

type UserChainInfo struct {
	ChainId     int32  `json:"chainid"`
	Address     string `json:"address"`
	ChainType   string `json: "chaintype"`
	DataNumbers int32  `json:"datanumbers"`
	Uid         int    `json:"uid"`
	Name        string `json:"name"`
	ChainName   string `json:"chainname"`
	WorkTime    int32  `json:"worktime"`
}

type UserPayInfo struct {
	UserId  int32  `json:"userid"`
	ChainId int32  `json:"chainid"`
	Address string `json:"address"`
	Name    string `json:"name"`
}
type TransactionCatUser struct {
	ChainId int32  `json:"chainid"`
	Address string `json:"address"`
	Name    string `json:"name"`
	Id      int    `json:"id"`
}

//转THK额度 日志
type UserPayInfoLog struct {
	UserName    string `json:"username"`
	FromChainId string `json:"fromchainid"`
	ToAddress   string `json:"toaddress"`
	PayTime     int32  `json:"paytime"`
	PayHash     string `json:"payhash"`
	ChainId     string `json:"chainid"`
	LogInfo     string `json:"loginfo"`
	PayLines    string `json:"paylines"`
}

//卖卖猫日志
type TransactionCatLog struct {
	UserName    string  `json:"username"`
	ToAddress   string  `json:"toaddress"`
	PayTime     int32   `json:"paytime"`
	PayHash     string  `json:"payhash"`
	InChainId   int32   `json:"chainid"`
	LogInfo     string  `json:"loginfo"`
	PayLines    float64 `json:"paylines"`
	BlockHeight string  `json:"blockheight"`
}

type AllUserChain struct {
	PageNum  int             `json:"pageNum"`
	PageSize int             `json:"pageSize"`
	Pages    int             `json:"pages"`
	Total    int             `json:"total"`
	DataList []UserChainInfo `json:"dataList"`
}

type AllUserRightInfo struct {
	PageNum  int         `json:"pageNum"`
	PageSize int         `json:"pageSize"`
	Pages    int         `json:"pages"`
	Total    int         `json:"total"`
	DataList []RightInfo `json:"dataList"`
}
type CatId struct {
	Id     int
	Saling bool
}
type PostParams struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

/**  节点信息*/
type ChainNodeInfo struct {
	ChainId      int    `json:"chainId"`
	DataNodeId   string `json:"dataNodeId"`
	DataNodeIp   string `json:"dataNodeIp"`
	DataNodePort int    `json:"dataNodePort"`
	Mode         int    `json:"mode"`
	Parent       int    `json:"parent"`
}

/* 父节点关联这种 */
type ChainInfoParents struct {
	ChainId      int                `json:"chainId"`
	DataNodeId   string             `json:"dataNodeId"`
	DataNodeIp   string             `json:"dataNodeIp"`
	DataNodePort int                `json:"dataNodePort"`
	Mode         int                `json:"mode"`
	Parent       int                `json:"parent"`
	Childrens    []ChainInfoParents `json:"childrens"`
}

/**  节点信息 层级 */
type ChainNodeInfoStruct struct {
	ChainId      int                   `json:"chainId"`
	DataNodeId   string                `json:"dataNodeId"`
	DataNodeIp   string                `json:"dataNodeIp"`
	DataNodePort int                   `json:"dataNodePort"`
	Mode         int                   `json:"mode"`
	Parent       int                   `json:"parent"`
	Childrens    []ChainNodeInfoStruct `json:"childrens"`
}

/**  节点信息*/
type ChainNodeInfos struct {
	ChinInfos []ChainNodeInfo
}
type Committee struct {
	ChainId       int      `json:"chainId"`
	Epoch         int      `json:"epoch"`
	MemberDetails []string `json:"memberdetails"`
}
type MainPageInfo struct {
	ThkPrice    float32 `json:"thkPrice"`
	MarketValue float32 `json:"marketValue"`
	TxCount     int64   `json:"txCount"`
	Tps         int64   `json:"tps"`
}

/*  每个交易类型对应的 交易数量 */
type TxTypeInfo struct {
	TxType  string `json:"txType"`
	TxCount int    `json:"txCount"`
}

/*   高度对应的交易信息  */
type HeightTxInfo struct {
	TxtSumCount  int `json:"txtSumCount"` //总数量
	TxByTypeInfo []TxTypeInfo
}

//txType1 合约发布 1
//txType2 合约交易  2
//txType3 链内交易 3
//txType4 跨链转账取款  4
//txType5 跨链转账存款  5
//txType6 跨链转账取消   6

const (
	TxType1 string = "1"
	TxType2 string = "2"
	TxType3 string = "3"
	TxType4 string = "4"
	TxType5 string = "5"
	TxType6 string = "6"
)

type PostParamTx struct {
	Page     int `json:"page"`
	Pagesize int `json:"pagesize"`

	ChainId  string `json:"chainId"`
	Hash     string `json:"hash"`
	Height   string `json:"height"`
	Epoch    string `json:"epoch"`
	Address  string `json:"address"`
	TxType   string `json:"txType"`
	To       string `json:"to"`
	Contract string `json:"contract"` //合约的hash

}

type Files struct {
	FileName    string `json:"fileName"`
	Version     string `json:"version"`
	FileSize    string `json:"fileSize"`
	CreateTime  int32  `json:"createTime"`
	PreviewUrl  string `json:"previewUrl"`
	DownloadUrl string `json:"downloadUrl"`
}

type CashCheckHash struct {
	Hashqk string `json:"hashqk"` //取款hash
	Hashcx string `json:"hashcx"` //存款hash 或者 撤销 Hash
}

type DownLoadFile struct {
	SDK []Files
	API []Files
}

type ResultData struct {
	Data interface{} `json:"data"`
}

const (
	QK string = "0x0000000000000000000000000000000000020000" //取款
	CK string = "0x0000000000000000000000000000000000030000" //存款
	QX string = "0x0000000000000000000000000000000000040000" //撤销
)

const (
	US string = "en-US"
	CN string = "zh-CN"
)

const (
	CODE_Error   int = 201
	CODE_Success int = 200
)
