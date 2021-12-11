package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	app "github.com/koen-vr/algo-prng-roller/shared/contract"
	exc "github.com/koen-vr/algo-prng-roller/shared/execute"
	net "github.com/koen-vr/algo-prng-roller/shared/network"
)

func init() {
	appStatusCmd.AddCommand(appGlobalStatusCmd)

	Status.AddCommand(appStatusCmd)
}

var Status = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the game",
	Long:  `Options to check the state and value of elements in the game.`,
	Run: func(cmd *cobra.Command, args []string) {
		// No args passed, fallback to help
		cmd.HelpFunc()(cmd, args)
	},
}

var appStatusCmd = &cobra.Command{
	Use:   "app",
	Short: "Create and fund accounts",
	Long:  `Creates and fund the required accounts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// No args passed, fallback to help
		cmd.HelpFunc()(cmd, args)
	},
}

var appGlobalStatusCmd = &cobra.Command{
	Use:   "global",
	Short: "Create and fund accounts",
	Long:  `Creates and fund the required accounts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := app.GetAppId("roller")
		if nil != err {
			return fmt.Errorf("status: get id: %s", err)
		}
		out, err := exc.List([]string{"-c", fmt.Sprintf(
			"goal app read -d %s --app-id %d --guess-format --global",
			net.NodePath(), id,
		)})
		if len(out) > 0 {
			fmt.Println(out)
		}
		if nil != err {
			return fmt.Errorf("status: app read: %s", err)
		}
		return nil
	},
}
