/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/portto/solana-go-sdk/rpc"
	"github.com/spf13/cobra"
)

// importWalletCmd represents the importWallet command
var importWalletCmd = &cobra.Command{
	Use:   "importWallet",
	Short: "Imports an existing wallet",
	Long:  "Imports an existing wallet from a given private key.",
	Run: func(cmd *cobra.Command, args []string) {
		privateKey, _ := ioutil.ReadFile("data")
		fmt.Println("Importing wallet from the 'key_data' file.")
		wallet, _ := ImportOldWallet(privateKey, rpc.DevnetRPCEndpoint)
		fmt.Println("Public Key: " + wallet.account.PublicKey.ToBase58())
		balance, _ := GetBalance()
		fmt.Println("Wallet balance: " + strconv.Itoa(int(balance/1e9)) + " SOL")
	},
}

func init() {
	rootCmd.AddCommand(importWalletCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importWalletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importWalletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
