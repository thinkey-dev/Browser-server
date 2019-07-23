package rpc

import (
	"PublicChainBrowser-Server/utils/sha3"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"log/log"

	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-ini/ini"
)

const (
	Addr    string = "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	StepStr string = "0000000000000000000000000000000000000000000000000000000000000080"
)

var (
	addr         = "http://192.168.1.13:8090"
	ContractAddr = "0x44b9402f12402352409c05fb31a750e28e1b6d07"
	SaleAddr     = "0x1940170c07d69ee15eee9c9cf6780917ae873122"
	UserAddr     = "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
)

var (
	url string = ""
)

type SendTx struct {
	ChainId     string `json:"chainId"`
	FromChainId string `json:"fromChainId"`
	From        string `json:"from"`
	ToChainId   string `json:"toChainId"`
	To          string `json:"to"`
	Sig         string `json:"sig"`
	Pub         string `json:"pub"`
	Nonce       string `json:"nonce"`
	Value       string `json:"value"`
	Input       string `json:"input"`
}

type GetAccount struct {
	ChainId string `json:"chainId"`
	Address string `json:"address"`
}
type GetBlock struct {
	ChainId string `json:"chainId"`
	Hash    string `json:"hash"`
}
type GetChainNodeInfo struct {
	ChainId string `json:"chainIds"`
}
type GetTransactionByHash struct {
	Chainid string `json:"chainId"`
	Hash    string `json:"hash"`
}

type PostParams struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type GetAccountResult struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

var (
	key = "4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"
	pk  = "0x044e3b81af9c2234cad09d679ce6035ed1392347ce64ce405f5dcd36228a25de6e47fd35c4215d1edf53e6f83de344615ce719bdb0fd878f6ed76f06dd277956de"
)

func init() {
	cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		panic(err)
	}
	url = cfg.Section("rpc").Key("rpcaddr").String()
}

func rpcPost(params []byte) (map[string]interface{}, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, err
	}

	return result, nil
}
func rpcPostByte(params []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//var result map[string]interface{}
	//if err := json.Unmarshal([]byte(body), &result); err != nil {
	//	return nil, err
	//}

	return body, nil
}

func HTTPSendTX(chainId string, from string, to string, value string, input string) (map[string]interface{}, error) {
	nonce, err := GetNonce(chainId, from)
	if err != nil {
		return nil, nil
	}

	sendtx := SendTx{Pub: pk, Sig: "", ChainId: chainId, FromChainId: chainId, ToChainId: chainId, From: from, To: to, Nonce: nonce, Value: value, Input: input}
	HashSerialize(&sendtx)
	params := PostParams{Method: "SendTx", Params: sendtx}
	jsons, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	resp, err := rpcPost(jsons)
	log.Info(resp)
	return resp, err
}

func GetNonce(chainId string, addr string) (string, error) {
	getaccount := GetAccount{ChainId: chainId, Address: addr}
	params := PostParams{Method: "GetAccount", Params: getaccount}
	jsons, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	resp, err := rpcPost(jsons)
	if err != nil {
		return "", err
	}
	if _, ok := resp["nonce"]; ok {
		nonce := strconv.Itoa(int(resp["nonce"].(float64)))
		return nonce, nil
	} else {
		return "", errors.New("get nonce err")
	}
}

func GetAcc(chainId string, addr string) (string, error) {
	if chainId == "2" {
		fmt.Println(chainId)
	}
	getaccount := GetAccount{ChainId: chainId, Address: addr}
	params := PostParams{Method: "GetAccount", Params: getaccount}
	jsons, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	resp, err := rpcPost(jsons)
	if err != nil {
		return "", err
	}
	var numStr string
	for key, val := range resp {
		if key == "balance" {
			numStr = fmt.Sprint(val)
		}
	}
	decimalNum, err := decimal.NewFromString(numStr)
	if err != nil {
		//log.Errorf("decimal.NewFromString error, numStr:%s, err:%v", numStr, err)
		return "", err
	}
	return decimalNum.String(), nil
}

func GetAcc1(chainId string, addr string) (GetAccountResult, error) {
	var result GetAccountResult
	getaccount := GetAccount{ChainId: chainId, Address: addr}
	params := PostParams{Method: "GetAccount", Params: getaccount}
	jsons, err := json.Marshal(params)
	if err != nil {
		return result, err
	}
	resp, err := rpcPostByte(jsons)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetAccRsp(chainId string, addr string) (interface{}, error) {
	getaccount := GetAccount{ChainId: chainId, Address: addr}
	params := PostParams{Method: "GetAccount", Params: getaccount}
	jsons, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	resp, err := rpcPost(jsons)
	if err != nil {
		return nil, err
	}
	return resp, nil
	//if _, ok := resp["balance"]; ok {
	//	//nonce := strconv.Itoa(int(resp["balance"].(float64)))
	//	return resp, nil
	//} else {
	//	return "", errors.New("get balance err")
	//}
}

func HashSerialize(tx *SendTx) (int, error) {
	var err error
	var input string

	// val, _ := math.ParseBig256(tx.Value)
	if len(tx.Input) > 2 && tx.Input[0:2] == "0x" {
		input = tx.Input[2:]
	}
	str := []string{tx.ChainId, tx.From[2:], tx.To[2:], tx.Nonce, tx.Value, input}
	p := strings.Join(str, "")
	tmp := sha3.NewKeccak256()
	tmp.Write([]byte(p))
	hash := tmp.Sum(nil)

	println("hash=", hexutil.Encode(hash))

	pKey, err := crypto.HexToECDSA(key)

	sig, err := crypto.Sign(hash, pKey)
	if err != nil {
		return 0, err
	}
	tx.Sig = hexutil.Encode(sig)
	log.Info(tx.Sig)
	return 1, nil
}

func HashSerialize_Cat(tx *SendTx) (int, error) {
	var err error
	val, _ := math.ParseBig256(tx.Value)

	str := []string{tx.ChainId, tx.From[2:], tx.To[2:], tx.Nonce, val.String(), tx.Input[2:]}
	p := strings.Join(str, "")
	tmp := sha3.NewKeccak256()
	tmp.Write([]byte(p))
	hash := tmp.Sum(nil)

	pKey, err := crypto.HexToECDSA(key)
	sig, err := crypto.Sign(hash, pKey)
	if err != nil {
		return 0, err
	}
	tx.Sig = hexutil.Encode(sig)
	tx.Pub = pk
	log.Info(tx.Sig)
	return 1, nil
}

func GetBlockInfo(chainid string, hash string) (map[string]interface{}, error) {
	getblock := GetBlock{Hash: hash, ChainId: chainid}
	params := PostParams{Method: "GetTransactionByHash", Params: getblock}
	jsons, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	resp, err := rpcPost(jsons)
	if err == nil {
		log.Info(resp)
		return resp, err
	} else {
		fmt.Println(err.Error())
		return nil, err
	}
}

func HttpCallTransaction(sendtx SendTx) string {
	//sendtx := SendTx{ChainId: ChainId, From: from, FromChainId: FromChainId, ToChainId: ToChainId, To: to, Nonce: nonce, Value: val, Input: input}
	params := PostParams{Method: "CallTransaction", Params: sendtx}
	jsons, errs := json.Marshal(params)
	fmt.Printf("%v\n", string(jsons))
	if errs != nil {
		return ""
	}
	req, err := http.NewRequest("POST", addr, bytes.NewBuffer(jsons))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
func GenInputValue(value string, vposition string) string {
	tmplen := 64 - len(value)
	tmpstr := make([]rune, tmplen)
	for i := range tmpstr {
		tmpstr[i] = '0'
	}
	if vposition == "front" {
		return value + string(tmpstr)
	} else if vposition == "back" {
		return string(tmpstr) + value
	}
	return ""
}
func RpcPost(params []byte) string {
	req, err := http.NewRequest("POST", addr, bytes.NewBuffer(params))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func GetChainInfo(chainId string) ([]byte, error) {
	getchainnode := GetChainNodeInfo{ChainId: chainId}
	params := PostParams{Method: "GetChainInfo", Params: getchainnode}
	jsons, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url+"/chaininfo", bytes.NewBuffer(jsons))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	} else {
		return body, nil
	}

	//var result [] controllers.ChainNodeInfo
	//if err := json.Unmarshal([]byte(body), &result); err != nil {
	//	return nil, err
	//}

	//return result, nil
}
