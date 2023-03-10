package config

import (
	"fmt"

	"github.com/spf13/cobra"

	appCfg "github.com/cgxarrie-go/prq/config"
)

// azurePATCmd represents the azurePAT command
var azurePATCmd = &cobra.Command{
	Use:   "pat",
	Short: "set azure PAT",
	Long:  `Set the Azure PAT in the configuration file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := runAzurePATCmd(args)
		return err
	},
}

func init() {
	ConfigCmd.AddCommand(azurePATCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// azurePATCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// azurePATCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runAzurePATCmd(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}
	cfg := appCfg.GetInstance()
	cfg.Load()

	cfg.Azure.PAT = args[0]
	err := cfg.Save()
	if err != nil {
		return err
	}
	return nil
}
