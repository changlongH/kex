/*
Copyright © 2023 hcl <hcl2685@gmail.com>
*/
package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kex",
	Short: "An interactive control tool for kube pods",
	Long: `An interactive control tool for kube pods.
You don't need to remember complex kubectl command.
Usage:
	$kex pods
show cluster pods list. You can select pod to describe or enter container.
Before that you need to initialize the configuration first。
	$touch ~/.kube/config
	$cp kex.yaml.template ~/.kex.yaml
	$vi ~/.kex.yaml`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

var cfgFile string
var Viper *viper.Viper

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	log.SetFlags(0)
	Viper = viper.New()
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kex.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		Viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		Viper.AddConfigPath(home)
		Viper.SetConfigType("yaml")
		Viper.SetConfigName(".kex")
	}

	Viper.AutomaticEnv()

	err := Viper.ReadInConfig()
	cobra.CheckErr(err)

	fillConfigAbsPath := func(key string, newkey string, required bool) {
		val := Viper.GetString(key)
		ret := filepath.IsAbs(val)
		if !ret {
			val = strings.Replace(val, "~/", "", 1)
			if len(val) > 0 {
				home, err := os.UserHomeDir()
				cobra.CheckErr(err)
				absPath := filepath.Join(home, val)
				Viper.Set(newkey, absPath)
			} else if required {
				err := errors.New("not set kubeConfig value in " + Viper.ConfigFileUsed())
				cobra.CheckErr(err)
			}
		}
	}
	fillConfigAbsPath("kubeConfig", "absKubeConfig", true)
	fillConfigAbsPath("kubeConfigsPath", "absKubeConfigsPath", false)
}

func GetKubeConfigFile() string {
	return Viper.GetString("absKubeConfig")
}

func GetKubeConfigsPath() string {
	return Viper.GetString("absKubeConfigsPath")
}
