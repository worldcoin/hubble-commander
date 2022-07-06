package main

import (
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/ybbus/jsonrpc/v2"
)

func getNetworkInfo(client jsonrpc.RPCClient) *dto.NetworkInfo {
	resp, err := client.Call("hubble_getNetworkInfo")
	if err != nil {
		log.Fatal(err)
	}

	var networkInfo dto.NetworkInfo
	err = resp.GetObject(&networkInfo)
	if err != nil {
		log.Fatal(err)
	}

	return &networkInfo
}

func getUserState(client jsonrpc.RPCClient, treeIndex uint32) *dto.UserStateWithID {
	resp, err := client.Call("hubble_getUserState", treeIndex)
	if err != nil {
		log.Fatal(err)
	}

	var userState dto.UserStateWithID
	err = resp.GetObject(&userState)
	if err != nil {
		log.Fatal(err)
	}

	return &userState
}

func getPublicKey(client jsonrpc.RPCClient, treeIndex uint32) *models.PublicKey {
	resp, err := client.Call("hubble_getPublicKeyByStateID", treeIndex)
	if err != nil {
		log.Fatal(err)
	}

	var publicKey models.PublicKey
	err = resp.GetObject(&publicKey)
	if err != nil {
		log.Fatal(err)
	}

	return &publicKey
}

func getToStateID(ctx *cli.Context) uint32 {
	toStr := ctx.String("to")
	toInt, err := strconv.Atoi(toStr)
	if err != nil {
		log.Fatal(err)
	}

	return uint32(toInt)
}

func decodeHexString(asString string) []byte {
	asString = strings.TrimPrefix(asString, "0x")

	privateKey, err := hex.DecodeString(asString)
	if err != nil {
		log.Fatal(err)
	}

	return privateKey
}

func decodePublicKey(asString string) *models.PublicKey {
	asBytes := decodeHexString(asString)
	if len(asBytes) != models.PublicKeyLength {
		log.Fatal("public keys must be exactly 128 bytes long")
	}

	var key models.PublicKey
	copy(key[:], asBytes)
	return &key
}

func getTxType(ctx *cli.Context) txtype.TransactionType {
	if ctx.IsSet("type") {
		switch strings.ToLower(ctx.String("type")) {
		case "transfer":
			return txtype.Transfer
		case "c2t", "create2transfer":
			return txtype.Create2Transfer
		default:
			log.Fatal("unrecognized transaction type: ", ctx.String("type"))
		}
	}

	toStr := ctx.String("to")

	if strings.HasPrefix(toStr, "0x") {
		// a hexstring is probably a pubkey
		return txtype.Create2Transfer
	}

	_, err := strconv.Atoi(toStr)
	if err == nil {
		// it is technically possible for a pubkey to have no a-f chars but it is
		// exceedingly unlikely: (10/16)^64
		return txtype.Transfer
	}

	_, err = hex.DecodeString(toStr)
	if err == nil {
		return txtype.Create2Transfer
	}

	log.Fatal("failed to infer intended transaction type, --to must be either a stateid or a pubkey")
	return 0
}

//nolint:funlen
func sendTransaction(ctx *cli.Context) error {
	txType := getTxType(ctx)

	var toStateID uint32
	var toPublicKey *models.PublicKey

	switch txType {
	case txtype.Transfer:
		toStateID = getToStateID(ctx)
		log.Info("Sending transfer to stateID: ", toStateID)
	case txtype.Create2Transfer:
		toPublicKey = decodePublicKey(ctx.String("to"))
		log.Info("Sending create2transfer to pubkey: ", toPublicKey)
	default:
		panic("unreachable")
	}

	client := jsonrpc.NewClient(ctx.String("rpcurl"))

	networkInfo := getNetworkInfo(client)
	chainIDAsString := networkInfo.ChainID.String()

	log.Infof(
		"Connected to remote hubble. ChainID=%s BlockNumber=%d",
		chainIDAsString,
		networkInfo.BlockNumber,
	)

	var fromStateID uint32 = uint32(ctx.Int("from"))

	userState := getUserState(client, fromStateID)
	log.Infof(
		"Found sender account. ID=%d Token=%s Balance=%s Nonce=%s",
		fromStateID,
		userState.TokenID.String(),
		userState.Balance.String(),
		userState.Nonce.String(),
	)

	var amount models.Uint256 = models.MakeUint256(ctx.Uint64("amount"))
	var fee models.Uint256 = models.MakeUint256(ctx.Uint64("fee"))

	privateKeyAsString := ctx.String("privateKey")
	privateKey := decodeHexString(privateKeyAsString)

	wallet, err := bls.NewWallet(privateKey, networkInfo.SignatureDomain)
	if err != nil {
		log.Fatal(err)
	}

	serverPublicKey := getPublicKey(client, fromStateID)
	walletPublicKey := wallet.PublicKey()
	if serverPublicKey.String() != walletPublicKey.String() {
		log.Error("derived public key: ", walletPublicKey)
		log.Error("expected public key: ", serverPublicKey)
		log.Fatal("cannot sign transaction, provided private key is incorrect")
	}

	var resp *jsonrpc.RPCResponse

	switch txType {
	case txtype.Transfer:
		transfer := &dto.Transfer{
			FromStateID: &fromStateID,
			ToStateID:   &toStateID,
			Amount:      &amount,
			Nonce:       &userState.Nonce,
			Fee:         &fee,
			Signature:   nil,
		}

		transfer, err = api.SignTransfer(wallet, *transfer)
		if err != nil {
			log.Fatal(err)
		}

		resp, err = client.Call(
			"hubble_sendTransaction",
			[]dto.Transfer{*transfer},
		)
		if err != nil {
			log.Fatal(err)
		}
	case txtype.Create2Transfer:
		create2Transfer := &dto.Create2Transfer{
			FromStateID: &fromStateID,
			ToPublicKey: toPublicKey,
			Amount:      &amount,
			Nonce:       &userState.Nonce,
			Fee:         &fee,
			Signature:   nil,
		}

		create2Transfer, err = api.SignCreate2Transfer(wallet, *create2Transfer)
		if err != nil {
			log.Fatal(err)
		}

		resp, err = client.Call(
			"hubble_sendTransaction",
			[]dto.Create2Transfer{*create2Transfer},
		)
		if err != nil {
			log.Fatal(err)
		}
	default:
		panic("unreachable")
	}

	if resp.Error != nil {
		log.Error("Error: ", resp.Error)
	} else {
		txHash, err := resp.GetString()
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("Submitted transaction. TxHash=%s", txHash)
	}

	return nil
}
