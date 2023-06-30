package taskcreate

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	kafkaClient "github.com/uduncloud/easynode/common/kafka"
	"github.com/uduncloud/easynode/common/util"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service"
	"github.com/uduncloud/easynode/task/service/taskcreate/db"
	"github.com/uduncloud/easynode/task/service/taskcreate/ether"
	"github.com/uduncloud/easynode/task/service/taskcreate/tron"
	"time"
)

type Service struct {
	config      *config.Config
	store       service.StoreTaskInterface
	log         *xlog.XLog
	nodeId      string
	kafkaClient *kafkaClient.EasyKafka
	sendCh      chan []*kafka.Message
	//receiverCh  chan []*kafka.Message
}

func (s *Service) Start(ctx context.Context) {

	go s.startKafka(ctx)

	go s.updateLatestBlock(ctx)

	go func(ctx context.Context) {
		interrupt := true
		for interrupt {
			//处理自产生区块任务，包括：区块
			s.startCreateBlockProc()
			//<-time.After(20 * time.Second)
			select {
			case <-ctx.Done():
				interrupt = false
				break
			default:
				time.Sleep(13 * time.Second)
				continue
			}
		}
	}(ctx)
}

func (s *Service) startKafka(ctx context.Context) {
	broker := fmt.Sprintf("%v:%v", s.config.TaskKafka.Host, s.config.TaskKafka.Port)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	s.kafkaClient.Write(kafkaClient.Config{Brokers: []string{broker}}, s.sendCh, nil, ctx)
}

func (s *Service) updateLatestBlock(ctx context.Context) {
	log := s.log.WithFields(logrus.Fields{
		"model": "updateLatestBlock",
		"id":    time.Now().UnixMilli(),
	})
	blockConfigs := s.config.BlockConfigs
	if len(blockConfigs) < 1 {
		log.Warnf("config.BlockConfigs|info=%v", "chain config is null")
		return
	}

	interrupt := true
	for interrupt {
		for _, v := range blockConfigs {
			lastNumber, err := GetLastBlockNumber(v.BlockChainCode, s.log, v)
			if err != nil {
				log.Errorf("GetLastBlockNumber|err=%v", err)
				continue
			}
			if lastNumber > 1 {
				_ = s.store.UpdateLastNumber(v.BlockChainCode, lastNumber)
			}
		}

		select {
		case <-ctx.Done():
			interrupt = false
			break
		default:
			time.Sleep(10 * time.Second)
			continue
		}
	}
}
func (s *Service) startCreateBlockProc() {
	log := s.log.WithFields(logrus.Fields{
		"model": "startCreateBlockProc",
		"id":    time.Now().UnixMilli(),
	})
	blockConfigs := s.config.BlockConfigs
	if len(blockConfigs) < 1 {
		log.Warnf("config.BlockConfigs|info=%v", "chain config is null")
		return
	}

	for _, v := range blockConfigs {
		//生成最新的区块任务
		err := s.NewBlockTask(*v, log)
		if err != nil {
			log.Errorf("NewBlockTask|err=%v", err.Error())
		}

	}
}

func (s *Service) GetLastBlockNumberForEther(v *config.BlockConfig) error {
	lastNumber, err := ether.NewEther(s.log).GetLatestBlockNumber(v)
	if err != nil {
		return err
	}
	if lastNumber > 1 {
		return s.store.UpdateLastNumber(v.BlockChainCode, lastNumber)
	}
	return nil
}

func (s *Service) GetLastBlockNumberForTron(v *config.BlockConfig) error {
	lastNumber, err := tron.NewTron(s.log).GetLatestBlockNumber(v)
	if err != nil {
		return err
	}
	if lastNumber > 1 {
		return s.store.UpdateLastNumber(v.BlockChainCode, lastNumber)
	}
	return nil
}

func (s *Service) NewBlockTask(v config.BlockConfig, log *logrus.Entry) error {
	if v.BlockMin < 1 {
		panic("Min blockNumber is not less 1")
	}

	//已经下发的最新区块高度
	UsedMaxNumber, lastBlockNumber, err := s.store.GetRecentNumber(v.BlockChainCode)
	if err != nil {
		log.Errorf("GetRecentNumber|err=%v", err)
		return err
	}

	//如果从未下发该链区块任务，则 使用配置的最小区块高度
	if UsedMaxNumber == 0 {
		UsedMaxNumber = v.BlockMin
	}

	//获取指定区块高度
	//UsedMaxNumber~BlockMax

	//如果没有配置最大高度，则最大高度 时时读取链上最新高度
	if v.BlockMax < 1 {
		v.BlockMax = lastBlockNumber
	}

	if UsedMaxNumber >= v.BlockMax {
		return errors.New("UsedMaxNumber > BlockMax")
	}
	list := make([]*service.NodeTask, 0)

	UsedMaxNumber++

	for UsedMaxNumber <= v.BlockMax {
		//ns := &service.NodeSource{BlockChain: v.BlockChainCode, BlockNumber: fmt.Sprintf("%v", UsedMaxNumber), SourceType: 2}

		task := &service.NodeTask{
			NodeId:      s.nodeId,
			BlockNumber: fmt.Sprintf("%v", UsedMaxNumber),
			BlockChain:  v.BlockChainCode,
			TaskType:    2,
			TaskStatus:  0,
			CreateTime:  time.Now(),
			LogTime:     time.Now(),
			Id:          time.Now().UnixNano(),
		}

		list = append(list, task)
		UsedMaxNumber++
	}

	recentNumber := UsedMaxNumber - 1

	//生成任务并保存
	msgList, err := s.store.AddNodeTask(list)
	if err != nil {
		return err
	}
	if len(msgList) > 0 {
		s.sendCh <- msgList
	}

	//更新最新下发的区块高度
	err = s.store.UpdateRecentNumber(v.BlockChainCode, recentNumber)
	if err != nil {
		return err
	}

	return nil
}

func NewService(config *config.Config) *Service {
	xg := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFile("./log/task/create_task", 24*time.Hour)
	kf := kafkaClient.NewEasyKafka(xg)
	sendCh := make(chan []*kafka.Message, 5)
	//receiverCh := make(chan []*kafka.Message, 5)
	store := db.NewFileTaskCreateService(config, xg)
	nodeId, err := util.GetLocalNodeId()
	if err != nil {
		panic(err)
	}
	return &Service{
		config:      config,
		log:         xg,
		store:       store,
		nodeId:      nodeId,
		kafkaClient: kf,
		sendCh:      sendCh,
		//receiverCh:  receiverCh,
	}
}
