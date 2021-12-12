package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/koen-vr/algo-dice-roll/cmd/game/cmd"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "Run development and test utilities",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("./contracts"); err != nil {
			fmt.Fprintln(os.Stderr, "contracts folder not available")
			os.Exit(1)
		}
		if _, err := os.Stat("./network.json"); err != nil {
			fmt.Fprintln(os.Stderr, "network setup file not available")
			os.Exit(1)
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(cmd.Create)
	rootCmd.AddCommand(cmd.Deploy)
	rootCmd.AddCommand(cmd.Network)
	rootCmd.AddCommand(cmd.Status)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.algo-cfg)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".algo-cfg"
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".algo-cfg")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
