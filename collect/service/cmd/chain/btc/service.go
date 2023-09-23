package btc

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/0xcregis/easynode/blockchain"
	chainConfig "github.com/0xcregis/easynode/blockchain/config"
	chainService "github.com/0xcregis/easynode/blockchain/service"
	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type Service struct {
	log                *xlog.XLog
	chain              *config.Chain
	nodeId             string
	store              collect.StoreTaskInterface
	txChainClient      blockchain.API
	blockChainClient   blockchain.API
	receiptChainClient blockchain.API
}

func (s *Service) GetMultiBlockByNumber(blockNumber string, log *logrus.Entry, flag bool) ([]*collect.BlockInterface, []*collect.TxInterface) {
	return nil, nil
}

func (s *Service) Monitor() {
	go func() {
		for {
			<-time.After(30 * time.Second)
			if s.txChainClient != nil {
				txCluster := s.txChainClient.MonitorCluster()
				_ = s.store.StoreClusterNode(int64(s.chain.BlockChainCode), "tx", txCluster)
			}

			if s.blockChainClient != nil {
				blockCluster := s.blockChainClient.MonitorCluster()
				_ = s.store.StoreClusterNode(int64(s.chain.BlockChainCode), "block", blockCluster)
			}

			if s.receiptChainClient != nil {
				receiptCluster := s.receiptChainClient.MonitorCluster()
				_ = s.store.StoreClusterNode(int64(s.chain.BlockChainCode), "receipt", receiptCluster)
			}
		}
	}()
}

func (s *Service) GetBlockByHash(blockHash string, eLog *logrus.Entry, flag bool) (*collect.BlockInterface, []*collect.TxInterface) {
	start := time.Now()
	defer func() {
		eLog.Printf("GetBlockByHash.Duration =%v,blockHash:%v", time.Since(start), blockHash)
	}()
	//调用接口
	resp, err := s.blockChainClient.GetBlockByHash(int64(s.chain.BlockChainCode), blockHash, flag)
	if err != nil {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, err.Error(), blockHash)
		return nil, nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, "block is empty", blockHash)
		return nil, nil
	}

	resp = gjson.Parse(resp).Get("result").String()

	//解析数据
	block, txList := GetBlockFromJson(resp)

	if len(block.BlockHash) < 1 {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, resp, blockHash)
		return nil, nil
	}

	r := &collect.BlockInterface{BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, Block: block}
	if !flag { //仅区块，不涉及交易
		return r, nil
	}

	//fix prevout
	txList = s.getTxsOut(txList, eLog)

	addressList, err := s.store.GetMonitorAddress(int64(s.chain.BlockChainCode))
	if err != nil {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, err, blockHash)
		return nil, nil
	}
	addressMp := rebuildAddress(addressList)
	txs := make([]*collect.TxInterface, 0, len(txList))
	for _, tx := range txList {
		bs, _ := json.Marshal(tx)
		if s.CheckAddress(bs, addressMp) {
			t := &collect.TxInterface{TxHash: tx.TxHash, Tx: tx}
			txs = append(txs, t)
		}
	}
	//r := &service.BlockInterface{BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, Block: block}
	return r, txs
}

func (s *Service) GetBlockByNumber(blockNumber string, eLog *logrus.Entry, flag bool) (*collect.BlockInterface, []*collect.TxInterface) {

	start := time.Now()
	defer func() {
		eLog.Printf("GetBlockByNumber.Duration =%v,blockNumber:%v", time.Since(start), blockNumber)
	}()

	//调用接口
	resp, err := s.blockChainClient.GetBlockByNumber(int64(s.chain.BlockChainCode), blockNumber, flag)
	//resp, err := ether.Eth_GetBlockByNumber(cluster.Host, cluster.Key, blockNumber, s.log)
	if err != nil {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, err.Error(), blockNumber)
		return nil, nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, "block is empty", blockNumber)
		return nil, nil
	}

	resp = gjson.Parse(resp).Get("result").String()

	//解析数据
	block, txList := GetBlockFromJson(resp)

	if len(block.BlockHash) < 1 {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, resp, blockNumber)
		return nil, nil
	}

	r := &collect.BlockInterface{BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, Block: block}
	if !flag { //仅区块数据，不涉及交易
		return r, nil
	}

	//fix prevout
	txList = s.getTxsOut(txList, eLog)

	addressList, err := s.store.GetMonitorAddress(int64(s.chain.BlockChainCode))
	if err != nil {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, err, blockNumber)
		return nil, nil
	}

	addressMp := rebuildAddress(addressList)
	txs := make([]*collect.TxInterface, 0, len(txList))
	for _, tx := range txList {
		bs, _ := json.Marshal(tx)
		if s.CheckAddress(bs, addressMp) {
			t := &collect.TxInterface{TxHash: tx.TxHash, Tx: tx}
			txs = append(txs, t)
		} else {
			eLog.Warnf("the tx is ignored,hash=%v", tx.TxHash)
		}
	}
	//r := &service.BlockInterface{BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, Block: block}
	return r, txs
}

func (s *Service) GetTx(txHash string, eLog *logrus.Entry) *collect.TxInterface {

	//调用接口
	resp, err := s.txChainClient.GetTxByHash(int64(s.chain.BlockChainCode), txHash)
	//resp, err := ether.Eth_GetTransactionByHash(cluster.Host, cluster.Key, txHash, s.log)
	if err != nil {
		eLog.Errorf("GetTx|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetTx|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "tx is empty", txHash)
		return nil
	}
	resp = gjson.Parse(resp).Get("result").String()

	//解析数据
	tx := GetTxFromJson(resp)

	if len(tx.TxHash) < 1 {
		eLog.Errorf("GetTx|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, resp, txHash)
		return nil
	}

	//fix tx.blockNumber if blockNumber is null
	if len(tx.BlockNumber) < 1 {
		resp, err = s.txChainClient.GetBlockByHash(int64(s.chain.BlockChainCode), tx.BlockHash, false)
		if err != nil {
			eLog.Errorf("GetTx|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, resp, txHash)
			return nil
		}
		tx.BlockNumber = gjson.Parse(resp).Get("result.height").String()
	}

	addressList, err := s.store.GetMonitorAddress(int64(s.chain.BlockChainCode))
	if err != nil {
		eLog.Errorf("GetTx|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err, txHash)
		return nil
	}

	addressMp := rebuildAddress(addressList)
	bs, _ := json.Marshal(tx)
	if s.CheckAddress(bs, addressMp) {
		r := &collect.TxInterface{TxHash: tx.TxHash, Tx: tx}
		return r
	} else {
		eLog.Warnf("the tx is ignored,hash=%v", tx.TxHash)
	}

	//r := &collect.TxInterface{TxHash: tx.TxHash, Tx: tx}
	return nil
}

func (s *Service) GetReceiptByBlock(blockHash, number string, eLog *logrus.Entry) ([]*collect.ReceiptInterface, error) {
	return nil, nil
}

func (s *Service) GetReceipt(txHash string, eLog *logrus.Entry) (*collect.ReceiptInterface, error) {
	return nil, nil
}

func (s *Service) getTxsOut(list []*collect.Tx, eLog *logrus.Entry) []*collect.Tx {
	wg := sync.WaitGroup{}
	wg.Add(len(list))
	for _, v := range list {
		go func(v *collect.Tx) {
			defer wg.Done()

			resp, err := s.txChainClient.GetTxByHash(int64(s.chain.BlockChainCode), v.TxHash)
			//resp, err := ether.Eth_GetTransactionByHash(cluster.Host, cluster.Key, txHash, s.log)
			if err != nil {
				eLog.Errorf("getTxsOut|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), v.TxHash)
				return
			}

			//处理数据
			if resp == "" {
				eLog.Errorf("getTxsOut|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "tx is empty", v.TxHash)
				return
			}
			from := gjson.Parse(resp).Get("result.vin").String()
			v.FromAddr = from
		}(v)
	}

	wg.Wait()
	return list
}

func getCoreAddr(addr string) string {
	addr = strings.ToLower(addr)
	if strings.HasPrefix(addr, "0x") {
		return strings.Replace(addr, "0x", "", 1) //去丢0x
	}
	return addr
}

func rebuildAddress(addrList []string) map[string]int64 {
	mp := make(map[string]int64, len(addrList))
	for _, v := range addrList {
		addr := getCoreAddr(v)
		mp[addr] = 1
	}
	return mp
}

func (s *Service) CheckAddress(tx []byte, addrList map[string]int64) bool {

	//if len(addrList) < 1 || len(tx) < 1 {
	//	return false
	//}

	txAddressList := make(map[string]int64, 10)
	root := gjson.ParseBytes(tx)

	fromAddr := root.Get("from").String()

	fromList := gjson.Parse(fromAddr).Array()
	for _, v := range fromList {
		addr := v.Get("prevout.scriptPubKey.address").String()
		txAddressList[getCoreAddr(addr)] = 1
	}

	toAddr := root.Get("to").String()
	toList := gjson.Parse(toAddr).Array()
	for _, v := range toList {
		addr := v.Get("scriptPubKey.address").String()
		txAddressList[getCoreAddr(addr)] = 1
	}

	has := false
	for k := range txAddressList {
		if _, ok := addrList[k]; ok {
			has = true
			break
		}
	}
	return has
}

func NewService(c *config.Chain, x *xlog.XLog, store collect.StoreTaskInterface, nodeId string) collect.BlockChainInterface {

	var blockClient blockchain.API
	if c.BlockTask != nil {
		list := make([]*chainConfig.NodeCluster, 0, 4)
		for _, v := range c.BlockTask.FromCluster {
			temp := &chainConfig.NodeCluster{
				NodeUrl:   v.Host,
				NodeToken: v.Key,
				Weight:    v.Weight,
			}
			list = append(list, temp)
		}
		blockClient = chainService.NewApi(int64(c.BlockChainCode), list, x)
	}

	var txClient blockchain.API
	if c.TxTask != nil {
		list := make([]*chainConfig.NodeCluster, 0, 4)
		for _, v := range c.TxTask.FromCluster {
			temp := &chainConfig.NodeCluster{
				NodeUrl:   v.Host,
				NodeToken: v.Key,
				Weight:    v.Weight,
			}
			list = append(list, temp)
		}

		txClient = chainService.NewApi(int64(c.BlockChainCode), list, x)
	}

	var receiptClient blockchain.API
	if c.ReceiptTask != nil {
		list := make([]*chainConfig.NodeCluster, 0, 4)
		for _, v := range c.ReceiptTask.FromCluster {
			temp := &chainConfig.NodeCluster{
				NodeUrl:   v.Host,
				NodeToken: v.Key,
				Weight:    v.Weight,
			}
			list = append(list, temp)
		}
		receiptClient = chainService.NewApi(int64(c.BlockChainCode), list, x)
	}

	return &Service{
		log:                x,
		chain:              c,
		store:              store,
		txChainClient:      txClient,
		blockChainClient:   blockClient,
		receiptChainClient: receiptClient,
		nodeId:             nodeId,
	}
}
