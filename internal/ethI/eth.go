package ethI

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"gohub/internal/dao"
	"gohub/internal/enum"
	"gohub/internal/model"
	"gohub/internal/service"
	"gohub/pkg/config"
	"gohub/pkg/eth"
	"gohub/pkg/logger"
	"math/big"
	"strings"
)

type Subscription struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

var confirmBlock int

// 订阅事件
var execTopic = crypto.Keccak256Hash([]byte("MSC20OrderExecuted(address,address,bytes32,string,uint256,uint256,uint16)"))
var cancelTopic = crypto.Keccak256Hash([]byte("MSC20OrderCanceled(address,bytes32)"))

func ListeningOrderLog() {
	// 起始区块
	var startNum = config.GetInt64("start_block")
	confirmBlock = config.GetInt("confirm_block")
	// 数据库中最新的区块号
	var orderLog *model.OrderLogDO
	if err := dao.OrderLog.Model().Order("id desc").Find(&orderLog).Error; err != nil {
		panic(err)
	}

	if orderLog != nil {
		startNum = orderLog.BlockNumber
	}

	address := strings.ToLower(config.Get("eth.contract_address"))

	// 同步到最新区块
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(address)},
		Topics:    [][]common.Hash{{execTopic, cancelTopic}},
		FromBlock: big.NewInt(startNum),
	}
	filterLogs, err := eth.WsClient.FilterLogs(context.Background(), query)
	if err != nil {
		panic(err)
	}
	filterLogs = filterLogs[1:]
	for _, vLog := range filterLogs {
		handleLog(vLog)
	}

	if len(filterLogs) > 0 {
		query.FromBlock = big.NewInt(int64(filterLogs[len(filterLogs)-1].BlockNumber))
	}

	// 订阅事件
	logs := make(chan types.Log)
	sub, err := eth.WsClient.SubscribeFilterLogs(context.Background(), query, logs)

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				logger.Errorv(errors.Wrap(err, "监听订单日志出错，退出"))
				return
			case vLog := <-logs:
				handleLog(vLog)
			}
		}
	}()
}

func handleLog(vLog types.Log) {
	logger.Debugf("Log:%+v", vLog)
	status := enum.OrderLogStatusSuccess
	switch vLog.Topics[0].Hex() {
	case execTopic.Hex(): // 订单执行
		if errCode := service.Order.Execute(vLog.Topics[1].String()); errCode != enum.OrderLogStatusSuccess {
			logger.Error("订单执行错误：" + errCode.Name)
			status = errCode
			break
		}
	case cancelTopic.Hex(): // 订单取消
		if errCode := service.Order.Cancel(vLog.Topics[1].String()); errCode != enum.OrderLogStatusSuccess {
			logger.Error("取消订单执行错误：" + errCode.Name)
			status = errCode
			break
		}
	}
	// 保存日志
	err := savaLog(vLog, status.Code)
	if err != nil {
		logger.Errorv(err)
	}
}

var queue LogQueue = make([]types.Log, confirmBlock)

func savaLog(vLog types.Log, status string) error {
	// 等待一段区块确认，避免区块重组
	vLog, ok := queue.Enqueue(vLog)
	if ok {
		return nil
	}
	topics := make([]string, len(vLog.Topics))
	for i, v := range vLog.Topics {
		topics[i] = v.Hex()
	}

	err := dao.OrderLog.Model().Create(&model.OrderLogDO{
		Address:     vLog.Address.Hex(),
		Topics:      topics,
		BlockNumber: int64(vLog.BlockNumber),
		TxHash:      vLog.TxHash.Hex(),
		TxIndex:     vLog.TxIndex,
		Index:       vLog.Index,
		Status:      status,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

type LogQueue []types.Log

// Enqueue
//
//	@Description: 确认队列
//	@receiver q
//	@param n
//	@return types.Log
//	@return bool 入队是否成功
func (q *LogQueue) Enqueue(n types.Log) (types.Log, bool) {
	if confirmBlock <= 0 {
		return n, false
	}
	if len(*q) < confirmBlock {
		*q = append(*q, n)
		return types.Log{}, true
	} else {
		element := (*q)[0]
		*q = (*q)[1:]
		return element, false
	}
}
