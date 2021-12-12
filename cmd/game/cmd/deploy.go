package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/algorand/go-algorand-sdk/types"

	acc "github.com/koen-vr/algo-dice-roll/shared/account"
	app "github.com/koen-vr/algo-dice-roll/shared/contract"
	net "github.com/koen-vr/algo-dice-roll/shared/network"
)

var Deploy = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys the game to the network",
	Long:  `Deploys the game's contracts, assets and configuration on to the active network.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := buildContracts()
		cobra.CheckErr(err)
		err = deployContracts()
		cobra.CheckErr(err)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		err := verifyDeploy()
		cobra.CheckErr(err)
	},
}

func verifyDeploy() error {
	if !net.IsActive() {
		return fmt.Errorf("deploy: requires an active node")
	}
	if !acc.HasAccount("manager") {
		return fmt.Errorf("deploy: requires a funded accounts")
	}
	return nil
}

func buildContracts() error {
	if err := app.Build([]string{
		"clear", "roller",
	}); nil != err {
		return err
	}
	return nil
}

func deployContracts() error {
	// TODO Load main account
	account, err := acc.Load("manager")
	if nil != err {
		return err
	}

	if err := app.Deploy(app.Config{
		Manager:      account,
		ClearProg:    "clear",
		ApprovalProg: "roller",
		LocalSchema: types.StateSchema{
			NumUint:      0,
			NumByteSlice: 0,
		},
		GlobalSchema: types.StateSchema{
			NumUint:      1,
			NumByteSlice: 1,
		},
	}); nil != err {
		return err
	}
	return nil
}
