package eth

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
)

// Inscribe
//
//	@Description: 铭刻
//	@param msc20Json 铭文
func Inscribe(fromPrivateKey string, to string, msc20Json string) (string, error) {
	// 1. 得到私钥
	privateKey, err := crypto.HexToECDSA(fromPrivateKey[2:])
	if err != nil {
		return "", err
	}
	// 2. 得到公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// 3. 交易数据
	chainID, err := Client.ChainID(context.Background())
	if err != nil {
		return "", err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}
	gasFeeCap, err := Client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	gasTipCap, err := Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return "", err
	}
	var data = []byte("data:text/plain;rule=esip6," + msc20Json)

	toAddress := common.HexToAddress(to[2:])

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       23000,
		To:        &toAddress,
		Value:     big.NewInt(0),
		Data:      data,
	})

	// 4. 交易签名
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		return "", err
	}

	// 5. 发送交易
	err = Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}
