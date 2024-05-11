package eth

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/storyicon/sigverify"
	"gohub/pkg/config"
)

var Client *ethclient.Client
var WsClient *ethclient.Client

func SetupEth() {
	// 客户端
	var err error
	Client, err = ethclient.Dial(config.GetString("eth.rpc_url"))
	if err != nil {
		panic(err)
	}
	WsClient, err = ethclient.Dial(config.Get("eth.ws_rpc_url"))
	if err != nil {
		panic(err)
	}

}

// VerifySignature 验证以太坊签名
func VerifySignature(address, message, signature string) error {
	valid, err := sigverify.VerifyEllipticCurveHexSignatureEx(
		common.HexToAddress(address),
		[]byte(message),
		signature,
	)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("签名验证失败")
	}
	return nil
}

func GetLastBlockNum() (uint64, error) {
	header, err := Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}

	return header.Number.Uint64(), nil
}
