package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	acc "github.com/koen-vr/algo-prng-roller/shared/account"
	net "github.com/koen-vr/algo-prng-roller/shared/network"
)

func init() {
	Create.AddCommand(accountsCmd)
}

var Create = &cobra.Command{
	Use:   "create",
	Short: "Create required accounts and assets",
	Long:  `Creates the required accounts and game assets.`,
	Run: func(cmd *cobra.Command, args []string) {
		// No args passed, fallback to help
		cmd.HelpFunc()(cmd, args)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		err := verifyCreate()
		cobra.CheckErr(err)
	},
}

func verifyCreate() error {
	if !net.IsActive() {
		return fmt.Errorf("create: requires an active node")
	}
	return nil
}

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Create and fund accounts",
	Long:  `Creates and fund the required accounts.`,
	Run: func(cmd *cobra.Command, args []string) {
		account, err := acc.Create("manager")
		cobra.CheckErr(err)
		err = acc.FundFromSeed(10, account.Address.String())
		cobra.CheckErr(err)
	},
}
