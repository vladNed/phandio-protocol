package cmd

import (
	"encoding/hex"
	"fmt"
	"os"

	"bbogdan95/moneroutils/pkg/common"
	mcrypto "bbogdan95/moneroutils/pkg/crypto/monero"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "A tool to generate and sum monero keys",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var generateWalletCmd = &cobra.Command{
	Use:   "new",
	Short: "Generates a new wallet",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		privateKeys, err := mcrypto.GenerateKeys()
		if err != nil {
			panic(err)
		}

		publicKeys := privateKeys.PublicKeyPair()

		fmt.Printf("PrivateSendKey: %s\n", privateKeys.SpendKey().Hex())
		fmt.Printf("PublicSpendKey: %s\n", publicKeys.SpendKey().Hex())
		fmt.Printf("PrivateViewKey: %s\n", privateKeys.ViewKey().Hex())
		fmt.Printf("PublicViewKey:  %s\n", publicKeys.ViewKey().Hex())
		fmt.Println()
		fmt.Printf("Mainnet:     %s\n", publicKeys.Address(common.Mainnet))
		fmt.Printf("Stagenet:    %s\n", publicKeys.Address(common.Stagenet))
		fmt.Printf("Development: %s\n", publicKeys.Address(common.Development))
	},
}

var generatePubFromPriv = &cobra.Command{
	Use:   "pub",
	Short: "Pass a private key to derive the associated public key.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		pkBytes, err := hex.DecodeString(args[0])
		if err != nil {
			panic(err)
		}

		pvk, err := mcrypto.NewPrivateSpendKey(pkBytes)
		if err != nil {
			panic(err)
		}

		fmt.Printf("PrivateKey: %s\n", pvk.Hex())
		fmt.Printf("PublicKey : %s\n", pvk.Public().Hex())
	},
}

var generateViewFromSpend = &cobra.Command{
	Use:   "view",
	Short: "Pass a private spend key to generate a private view key.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		pkBytes, err := hex.DecodeString(args[0])
		if err != nil {
			panic(err)
		}

		pvk, err := mcrypto.NewPrivateSpendKey(pkBytes)
		if err != nil {
			panic(err)
		}
		keyPair, err := pvk.AsPrivateKeyPair()
		if err != nil {
			panic(err)
		}

		fmt.Printf("PrivateSpendKey: %s\n", keyPair.SpendKey().Hex())
		fmt.Printf("PrivateViewKey : %s\n", keyPair.ViewKey().Hex())
	},
}

var sumPrivateSpendKeys = &cobra.Command{
	Use:   "sumsk",
	Short: "Pass two private spend keys to generate their sum.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		pkABytes, err := hex.DecodeString(args[0])
		if err != nil {
			panic(err)
		}
		pvkA, err := mcrypto.NewPrivateSpendKey(pkABytes)
		if err != nil {
			panic(err)
		}

		pkBBytes, err := hex.DecodeString(args[1])
		if err != nil {
			panic(err)
		}
		pvkB, err := mcrypto.NewPrivateSpendKey(pkBBytes)
		if err != nil {
			panic(err)
		}

		sum := mcrypto.SumPrivateSpendKeys(pvkA, pvkB)

		isVerbose, _ := cmd.Flags().GetBool("verbose")

		if !isVerbose {
			fmt.Printf("SumPrivateSpendKey: %s\n", sum.Hex())
		} else {
			sumkp, err := sum.AsPrivateKeyPair()
			if err != nil {
				panic(err)
			}

			fmt.Printf("PrivateSendKey: %s\n", sumkp.SpendKey().Hex())
			fmt.Printf("PublicSpendKey: %s\n", sumkp.PublicKeyPair().SpendKey().Hex())
			fmt.Printf("PrivateViewKey: %s\n", sumkp.ViewKey().Hex())
			fmt.Printf("PublicViewKey:  %s\n", sumkp.PublicKeyPair().SpendKey().Hex())
			fmt.Println()
			fmt.Printf("Mainnet:     %s\n", sumkp.PublicKeyPair().Address(common.Mainnet))
			fmt.Printf("Stagenet:    %s\n", sumkp.PublicKeyPair().Address(common.Stagenet))
			fmt.Printf("Development: %s\n", sumkp.PublicKeyPair().Address(common.Development))
		}
	},
}

func init() {
	sumPrivateSpendKeys.PersistentFlags().BoolP("verbose", "v", false, "provides all keys when possible")

	rootCmd.AddCommand(generateWalletCmd)
	rootCmd.AddCommand(generatePubFromPriv)
	rootCmd.AddCommand(generateViewFromSpend)
	rootCmd.AddCommand(sumPrivateSpendKeys)
}
