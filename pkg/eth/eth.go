package eth

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/storyicon/sigverify"
	"gohub/pkg/config"
	"gohub/pkg/logger"
	"math/big"
)

var (
	ErrSignatureMismatch = errors.New("签名不匹配")
)
var Client *ethclient.Client

func SetupEth() {
	// 客户端
	var err error
	Client, err = ethclient.Dial(config.GetString("eth.rpc_url"))
	if err != nil {
		panic(err)
	}
}

func getPublicAddress(privateKey *ecdsa.PrivateKey) common.Address {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logger.Error("eth", "getPublicAddress", "cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

// VerifySignature 验证以太坊签名
func VerifySignature(address, message, signature string) error {
	valid, err := sigverify.VerifyEllipticCurveHexSignatureEx(
		common.HexToAddress(address),
		[]byte(message),
		signature,
	)
	if err != nil {
		return errors.WithStack(err)
	}
	if !valid {
		return errors.New("签名验证失败")
	}
	return nil
}

// getTransactor
//
//	@Description: 获取事务发送者
//	@param value 单位为wei
//	@param gasLimit  gas限制，0为自动评估
//	@return *bind.TransactOpts
//	@return error
func getTransactor(value int64, gasLimit uint64) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(config.GetString("eth.private_Key"))
	if err != nil {
		return nil, err
	}
	chainID, err := Client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)

	var nonce uint64
	nonce, err = Client.PendingNonceAt(context.Background(), getPublicAddress(privateKey))
	transactor.Nonce = big.NewInt(int64(nonce))
	if err != nil {
		return nil, err
	}
	transactor.GasPrice, err = Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	transactor.GasTipCap, err = Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, err
	}
	transactor.Value = big.NewInt(value)
	transactor.GasLimit = gasLimit
	return transactor, nil
}
