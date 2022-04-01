package cmd

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

type Wallet struct {
	account types.Account
	c       *client.Client
}

func CreateNewWallet(RPCEndpoint string) Wallet {
	// create a new wallet using types.NewAccount()
	newAccount := types.NewAccount()
	data := []byte(newAccount.PrivateKey)

	err := ioutil.WriteFile("data", data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return Wallet{
		newAccount,
		client.NewClient(RPCEndpoint),
	}
}

func ImportOldWallet(privateKey []byte, RPCEndpoint string) (Wallet, error) {
	wallet, err := types.AccountFromBytes(privateKey)
	if err != nil {
		log.Fatal("Failed to import:", err)
		return Wallet{}, err
	}

	return Wallet{
		wallet,
		client.NewClient(RPCEndpoint),
	}, nil
}

func GetBalance() (uint64, error) {
	privateKey, rErr := ioutil.ReadFile("data")
	wallet, _ := ImportOldWallet(privateKey, rpc.DevnetRPCEndpoint)
	balance, bErr := wallet.c.GetBalance(
		context.TODO(),                      // request context
		wallet.account.PublicKey.ToBase58(), // wallet to fetch balance for
	)

	// Check for read error
	if rErr != nil {
		log.Fatal("Failed to read wallet:", rErr)
		return 0, rErr
	}

	// Check for balance error
	if bErr != nil {
		log.Fatal("Failed to find balance:", bErr)
		return 0, bErr
	}

	return balance, nil
}

func RequestAirdrop(amount uint64) (string, error) {
	// request for SOL using RequestAirdrop()
	privateKey, kErr := ioutil.ReadFile("data")
	wallet, _ := ImportOldWallet(privateKey, rpc.DevnetRPCEndpoint)
	amount = amount * 1e9 // turning SOL into lamports
	txhash, hErr := wallet.c.RequestAirdrop(
		context.TODO(), // request context wallet.account.PublicKey.ToBase58(), // wallet address requesting airdrop
		wallet.account.PublicKey.ToBase58(),
		amount, // amount of SOL in lamport
	)

	// Check for privateKey error
	if kErr != nil {
		log.Fatal("Failed to request drop, key error:", kErr)
		return "", kErr
	}

	// Check for txhash error
	if hErr != nil {
		log.Fatal("Failed to request drop, hash error:", hErr)
		return "", hErr
	}
	return txhash, nil
}

func Transfer(receiver string, amount uint64) (string, error) {
	// 1. fetch the most recent blockhash
	privateKey, kErr := ioutil.ReadFile("data")
	wallet, _ := ImportOldWallet(privateKey, rpc.DevnetRPCEndpoint)
	response, resErr := wallet.c.GetRecentBlockhash(context.TODO())

	// 1a. Check for private Key error
	if kErr != nil {
		return "", kErr
	}

	// 1b. Check for response error
	if resErr != nil {
		return "", resErr
	}

	// 2. make a transfer message with the latest block hash
	sysProgTransParams := sysprog.TransferParam{wallet.account.PublicKey, // public key of the transaction sender
		common.PublicKeyFromString(receiver), // wallet address of the transaction receiver
		amount,                               // transaction amount in lamport
	}

	nMParams := types.NewMessageParam{
		wallet.account.PublicKey, // public key of the transaction signer
		[]types.Instruction{
			sysprog.Transfer(sysProgTransParams), // transaction amount in lamport
		},
		response.Blockhash, // recent block hash
	}

	message := types.NewMessage(nMParams)

	// 3. create a transaction with the message and TX signer
	newTransactionParams := types.NewTransactionParam{
		message,
		[]types.Account{wallet.account, wallet.account},
	}

	tx, err := types.NewTransaction(newTransactionParams)
	if err != nil {
		return "", err
	}

	// 4. send the transaction to the blockchain
	txhash, err := wallet.c.SendTransaction(context.TODO(), tx)
	if err != nil {
		return "", err
	}
	return txhash, nil
}
