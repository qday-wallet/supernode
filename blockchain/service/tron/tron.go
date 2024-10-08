package tron

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/sunjiangjun/supernode/blockchain"
	"github.com/sunjiangjun/supernode/blockchain/chain"
	"github.com/sunjiangjun/supernode/blockchain/config"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type Tron struct {
	log              *xlog.XLog
	nodeCluster      []*config.NodeCluster
	blockChainClient blockchain.ChainConn
}

func (t *Tron) Token(chainCode int64, contractAddr string, abi string, eip string) (string, error) {

	cluster := t.BalanceCluster(false)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}

	var resp map[string]any
	var err error
	if eip == "721" {
		url := fmt.Sprintf("%v/%v", cluster.NodeUrl, "wallet/triggerconstantcontract")
		resp, err = t.blockChainClient.GetToken721(url, cluster.NodeToken, contractAddr, contractAddr)
	} else if eip == "1155" {
		resp, err = t.blockChainClient.GetToken1155(cluster.NodeUrl, cluster.NodeToken, contractAddr, contractAddr)
	} else if eip == "20" {
		url := fmt.Sprintf("%v/%v", cluster.NodeUrl, "wallet/triggerconstantcontract")
		resp, err = t.blockChainClient.GetToken20ByHttp(url, cluster.NodeToken, contractAddr, contractAddr)
	} else {
		return "", fmt.Errorf("unknow the eip:%v", eip)
	}

	if err != nil {
		cluster.ErrorCount += 1
	}

	bs, _ := json.Marshal(resp)
	return string(bs), err
}

func (t *Tron) StartWDT() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for {
			<-ticker.C
			for _, v := range t.nodeCluster {
				v.ErrorCount = 0
			}
		}
	}()
}

func (t *Tron) MonitorCluster() any {
	return t.nodeCluster
}

func (t *Tron) GetCode(chainCode int64, address string) (string, error) {
	req := `{ "value": "%v", "visible": true}`
	req = fmt.Sprintf(req, address)
	return t.SendReq(chainCode, req, "wallet/getcontract")
}

func (t *Tron) GetAddressType(chainCode int64, address string) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("GetAddressType,Duration=%v", time.Since(start))
	}()
	req := `{ "value": "%v", "visible": true}`
	req = fmt.Sprintf(req, address)
	resp, err := t.SendReq(chainCode, req, "wallet/getcontract")
	if err != nil {
		return "", err
	}

	if gjson.Parse(resp).Get("code_hash").Exists() {
		//合约地址
		return "0x12", nil
	} else {
		//外部地址
		return "0x11", nil
	}
}

func (t *Tron) SubscribePendingTx(chainCode int64, receiverCh chan string, sendCh chan string) (string, error) {
	return "", fmt.Errorf("blockchain:%v,the method has not been implemented", chainCode)
}

func (t *Tron) SubscribeLogs(chainCode int64, address string, topics []string, receiverCh chan string, sendCh chan string) (string, error) {
	return "", fmt.Errorf("blockchain:%v,the method has not been implemented", chainCode)
}

func (t *Tron) UnSubscribe(chainCode int64, subId string) (string, error) {
	return "", fmt.Errorf("blockchain:%v,the method has not been implemented", chainCode)
}

func (t *Tron) GetBlockReceiptByBlockNumber(chainCode int64, number string) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("GetBlockReceiptByBlockNumber,Duration=%v", time.Since(start))
	}()
	req := `{"num": %v}`

	n, err := strconv.ParseInt(number, 0, 64)
	if err != nil {
		return "", err
	}
	req = fmt.Sprintf(req, n)
	return t.SendReq(chainCode, req, "wallet/gettransactioninfobyblocknum")
}

func (t *Tron) GetBlockReceiptByBlockHash(chainCode int64, hash string) (string, error) {
	return "", nil
}

func (t *Tron) GetTransactionReceiptByHash(chainCode int64, hash string) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("GetTransactionReceiptByHash,Duration=%v", time.Since(start))
	}()
	req := `{ "value": "%v"}`
	req = fmt.Sprintf(req, hash)
	return t.SendReq(chainCode, req, "wallet/gettransactioninfobyid")
}

func NewTron(cluster []*config.NodeCluster, blockchain int64, xlog *xlog.XLog) blockchain.API {
	blockChainClient := chain.NewChain(blockchain, xlog)
	if blockChainClient == nil {
		return nil
	}
	t := &Tron{
		log:              xlog,
		nodeCluster:      cluster,
		blockChainClient: blockChainClient,
	}
	t.StartWDT()
	return t
}

func NewTron2(cluster []*config.NodeCluster, blockchain int64, xlog *xlog.XLog) blockchain.ExApi {
	blockChainClient := chain.NewChain(blockchain, xlog)
	if blockChainClient == nil {
		return nil
	}
	t := &Tron{
		log:              xlog,
		nodeCluster:      cluster,
		blockChainClient: blockChainClient,
	}
	return t
}

func (t *Tron) GetAccountResourceForTron(chainCode int64, address string) (string, error) {
	req := `{
			  "address": "%v",
			  "visible": true
			}`
	req = fmt.Sprintf(req, address)
	res, err := t.SendReq(chainCode, req, "wallet/getaccountresource")
	if err != nil {
		return "", err
	}

	return res, nil
}

func (t *Tron) EstimateGasForTron(chainCode int64, from, to, functionSelector, parameter string) (string, error) {
	req := `{
			"owner_address": "%v",
			"contract_address": "%v",
			"function_selector": "%v",
			"parameter": "%v",
			"visible": true
			}`
	req = fmt.Sprintf(req, from, to, functionSelector, parameter)
	res, err := t.SendReq(chainCode, req, "wallet/triggerconstantcontract")
	if err != nil {
		return "", err
	}

	return res, nil
}

func (t *Tron) EstimateGas(chainCode int64, from, to, data string) (string, error) {
	req := `
	{
		  "id": 1,
		  "jsonrpc": "2.0",
		  "method": "eth_estimateGas",
		  "params": [
			{
			  "from":"%v",
			  "to": "%v",
			  "data": "%v"
			}
		  ]
	}`
	req = fmt.Sprintf(req, from, to, data)
	res, err := t.SendJsonRpc(chainCode, req)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (t *Tron) GasPrice(chainCode int64) (string, error) {
	req := `
		{
		"id": 1,
		"jsonrpc": "2.0",
		"method": "eth_gasPrice"
		}
		`

	//req = fmt.Sprintf(req, from, to, functionSelector, parameter)
	res, err := t.SendJsonRpc(chainCode, req)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (t *Tron) GetBlockByHash(chainCode int64, hash string, flag bool) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("GetBlockByHash,Duration=%v", time.Since(start))
	}()

	req := `{"id_or_num": "%v","detail":%v}`
	req = fmt.Sprintf(req, hash, flag)
	res, err := t.SendReq(chainCode, req, "wallet/getblock")
	if err != nil {
		return "", err
	}

	//var delTx bool = true
	//if delTx {
	//mp := gjson.Parse(res).Map()
	//delete(mp, "transactions")
	//r, _ := json.Marshal(mp)
	//return string(r), nil
	//} else {
	return res, nil
	//}
}

func (t *Tron) GetBlockByNumber(chainCode int64, number string, flag bool) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("GetBlockByNumber,Duration=%v", time.Since(start))
	}()

	req := `{"id_or_num": "%v","detail":%v}`

	n, err := strconv.ParseInt(number, 0, 64)
	if err != nil {
		return "", err
	}
	req = fmt.Sprintf(req, n, flag)
	res, err := t.SendReq(chainCode, req, "wallet/getblock")
	if err != nil {
		return "", err
	}

	//var delTx bool = true
	//if delTx {
	//mp := gjson.Parse(res).Map()
	//delete(mp, "transactions")
	//r, _ := json.Marshal(mp)
	//return string(r), nil
	//} else {
	return res, nil
	//}
}

func (t *Tron) GetTxByHash(chainCode int64, hash string) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("GetTxByHash,Duration=%v", time.Since(start))
	}()
	req := `{ "value": "%v"}`
	req = fmt.Sprintf(req, hash)
	return t.SendReq(chainCode, req, "wallet/gettransactionbyid")
}

func (t *Tron) SendJsonRpc(chainCode int64, req string) (string, error) {
	cluster := t.BalanceCluster(false)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}
	url := fmt.Sprintf("%v/%v", cluster.NodeUrl, "jsonrpc")
	return t.blockChainClient.SendRequestToChain(url, cluster.NodeToken, req)
}

func (t *Tron) Balance(chainCode int64, address string, tag string) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("Balance,Duration=%v", time.Since(start))
	}()
	req := `{"address":"%v",  "visible": true}`
	req = fmt.Sprintf(req, address)
	res, err := t.SendReq(chainCode, req, "wallet/getaccount")
	if err != nil {
		return "", err
	}

	r := gjson.Parse(res)
	if r.Get("Error").Exists() {
		return "", errors.New(r.Get("Error").String())
	} else {
		returnStr := `{"balance":%v}`
		balance := r.Get("balance").Int()
		returnStr = fmt.Sprintf(returnStr, balance)
		return returnStr, nil
	}
}

func (t *Tron) TokenBalance(chainCode int64, address string, contractAddr string, abi string) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("TokenBalance,Duration=%v", time.Since(start))
	}()
	cluster := t.BalanceCluster(false)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}

	url := fmt.Sprintf("%v/%v", cluster.NodeUrl, "wallet/triggerconstantcontract")
	mp, err := t.blockChainClient.GetToken20ByHttp(url, cluster.NodeToken, contractAddr, address)
	if err != nil {
		return "", err
	}
	rs, _ := json.Marshal(mp)
	return string(rs), nil
}

func (t *Tron) Nonce(chainCode int64, address string, tag string) (string, error) {
	return "", fmt.Errorf("blockchain:%v,the method has not been implemented", chainCode)
}

func (t *Tron) LatestBlock(chainCode int64) (string, error) {
	res, err := t.SendReq(chainCode, "", "wallet/getnowblock")
	if err != nil {
		return "", err
	}

	r := gjson.Parse(res)

	blockId := r.Get("blockID").String()
	number := r.Get("block_header.raw_data.number").Int()

	returnStr := `{"blockId":"%v","number":%v}`
	returnStr = fmt.Sprintf(returnStr, blockId, number)
	return returnStr, nil
}

func (t *Tron) SendRawTransaction(chainCode int64, signedTx string) (string, error) {
	req := `{
			  "transaction": "%v"
			}`
	req = fmt.Sprintf(req, signedTx)
	return t.SendReq(chainCode, req, "wallet/broadcasthex")
}

func (t *Tron) TraceTransaction(chainCode int64, address string) (string, error) {
	return "", fmt.Errorf("blockchain:%v,the method has not been implemented", chainCode)
}

func (t *Tron) GetLogs(chainCode int64, contracts string, fromBlock, toBlock string, topics ...string) (string, error) {
	req := `
			{
		  "id": 1,
		  "jsonrpc": "2.0",
		  "method": "eth_getLogs",
		  "params": [
			{
			  "address": [
				"%v"
			  ],
			  "fromBlock": "%v",
			  "toBlock": "%v",
			  "topics": [
				"%v"
			  ]
			}
		  ]
		}`
	req = fmt.Sprintf(req, contracts, fromBlock, toBlock, topics[0])
	res, err := t.SendJsonRpc(chainCode, req)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (t *Tron) SendReq(blockChain int64, reqBody string, url string) (resp string, err error) {
	reqBody = strings.Replace(reqBody, "\t", "", -1)
	reqBody = strings.Replace(reqBody, "\n", "", -1)
	defer func() {
		if err != nil {
			t.log.Errorf("method:%v,blockChain:%v,req:%v,err:%v,uri:%v", "SendReq", blockChain, reqBody, err, url)
		} else {
			t.log.Printf("method:%v,blockChain:%v,req:%v,resp:%v", "SendReq", blockChain, reqBody, "ok")
		}
	}()
	cluster := t.BalanceCluster(false)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}

	url = fmt.Sprintf("%v/%v", cluster.NodeUrl, url)
	resp, err = t.blockChainClient.SendRequestToChainByHttp(url, cluster.NodeToken, reqBody)
	if err != nil {
		cluster.ErrorCount += 1
	}
	return resp, err
}

func (t *Tron) BalanceCluster(trace bool) *config.NodeCluster {

	var resultCluster *config.NodeCluster
	l := len(t.nodeCluster)

	if l > 1 {
		//如果有多个节点，则根据权重计算
		mp := make(map[string][]int64, 0)
		originCluster := make(map[string]*config.NodeCluster, 0)

		var sum int64
		for _, v := range t.nodeCluster {
			if v.Weight == 0 {
				//如果没有设置weight,则默认设定5
				v.Weight = 5
			}
			sum += v.Weight
			key := fmt.Sprintf("%v/%v", v.NodeUrl, v.NodeToken)
			mp[key] = []int64{v.Weight, sum}
			originCluster[key] = v
		}

		f := math.Mod(float64(time.Now().Unix()), float64(sum))
		var nodeId string

		for k, v := range mp {
			if len(v) == 2 && f <= float64(v[1]) && f >= float64(v[1]-v[0]) {
				nodeId = k
				break
			}
		}
		resultCluster = originCluster[nodeId]
	} else if l == 1 {
		//如果 仅有一个节点，则只能使用该节点
		resultCluster = t.nodeCluster[0]
	} else {
		return nil
	}
	return resultCluster
}
