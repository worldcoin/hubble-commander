package main

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
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

func decodeHexString(asString string) []byte {
	asString = strings.TrimPrefix(asString, "0x")

	privateKey, err := hex.DecodeString(asString)
	if err != nil {
		log.Fatal(err)
	}

	return privateKey
}

func randomPublicKey() *models.PublicKey {
	privateKey := make([]byte, 32)
	_, err := rand.Read(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	wallet, err := bls.NewWallet(privateKey, bls.Domain{0x00, 0x00, 0x00, 0x00})
	if err != nil {
		log.Fatal(err)
	}

	return wallet.PublicKey()
}

func sendC2T(client jsonrpc.RPCClient, wallet *bls.Wallet, from uint32, nonce models.Uint256) string {
	toPublicKey := randomPublicKey()
	amount := models.MakeUint256(1)
	fee := models.MakeUint256(1)

	transfer := dto.Create2Transfer{
		FromStateID: &from,
		ToPublicKey: toPublicKey,
		Amount:      &amount,
		Nonce:       &nonce,
		Fee:         &fee,
		Signature:   nil,
	}

	create2Transfer, err := api.SignCreate2Transfer(wallet, transfer)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Call(
		"hubble_sendTransaction",
		[]dto.Create2Transfer{*create2Transfer},
	)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Error != nil {
		log.Fatal(resp.Error)
	}

	txHash, err := resp.GetString()
	if err != nil {
		log.Fatal(err)
	}

	return txHash
}

// runs commands against a running hubble instance
// e2e/bench/bench_transactions_test.go
func benchmarkHubble(ctx *cli.Context) error {
	// must be run against a state which has funds in this account
	//  the genesis I deployed earlier has already funded this account

	// privateKey
	privateKey := decodeHexString("4c7c9af5de8b5e5a5877d706d8a41faaf84aa3dd6e98b4a07d4eb7e44daf9c78")

	// stateID
	fromStateID := uint32(6)

	rpcURL := "http://localhost:8080"
	client := jsonrpc.NewClient(rpcURL)

	networkInfo := getNetworkInfo(client)
	chainIDAsString := networkInfo.ChainID.String()

	log.Infof(
		"Connected to remote hubble. ChainID=%s BlockNumber=%d",
		chainIDAsString,
		networkInfo.BlockNumber,
	)

	// check that we have the correct private key for this state

	userState := getUserState(client, fromStateID)
	log.Infof(
		"Found sender account. ID=%d Token=%s Balance=%s Nonce=%s",
		fromStateID,
		userState.TokenID.String(),
		userState.Balance.String(),
		userState.Nonce.String(),
	)

	// TODO: check that the balance is sufficient to create our transactions

	wallet, err := bls.NewWallet(privateKey, networkInfo.SignatureDomain)
	if err != nil {
		log.Fatal(err)
	}

	serverPublicKey := getPublicKey(client, fromStateID)
	walletPublicKey := wallet.PublicKey()
	if serverPublicKey.String() != walletPublicKey.String() {
		log.Error("derived public key: ", walletPublicKey)
		log.Error("expected public key: ", serverPublicKey)
		log.Fatal("cannot run benchmark, provided private key is incorrect")
	}

	nonce := &userState.Nonce

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		txHash := sendC2T(client, wallet, fromStateID, *nonce)
		log.Infof(
			"Sent C2T txHash=%s", txHash,
		)

		one := models.MakeUint256(1)
		nonce = nonce.Add(&one)
	}

	return nil
}
