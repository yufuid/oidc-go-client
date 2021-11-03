package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
	"yufuid.com/oidc-go-client/config"
	"yufuid.com/oidc-go-client/pkg"
)

var (
	cfgFile string
	cfg     config.ClientConfig
	rootCmd = &cobra.Command{
		Use:   "oidc-client",
		Short: "oidc-client client",
		Long:  "yufu auth common client",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("rootCmd")
		},
	}
)

var serverCommand = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.InitClient(cfg);
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig();
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	serverCommand.Flags().StringVarP(&cfgFile, "config", "c", "", "config file");
	_ = serverCommand.MarkFlagRequired("config");
	rootCmd.AddCommand(serverCommand);
}

func initConfig() {
	if cfgFile == "" {
		fmt.Println("Error: missing config file")
		os.Exit(1)
	}

	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("failed to load config file %s: %v", cfgFile, err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		_ = fmt.Errorf("failed to marshal config file %s: %v", cfgFile, err)
	}
	fmt.Println()
}
