// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/SantoDE/varaco/configuration"
	"github.com/SantoDE/varaco/rancher"
	"github.com/SantoDE/varaco/render"
	"github.com/SantoDE/varaco/types"
	"os"
	"github.com/mitchellh/go-homedir"
	"github.com/SantoDE/varaco/websocket"
)

var cfgFile string

var config configuration.Rancher

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "varaco",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			rancher := new(rancher.Rancher)
			rancher.Config = &config

			path := viper.GetString("file.path")
			serviceToWatch := viper.GetString("SERVICE_TO_WATCH")
			additionalVCL := viper.GetString("ADDITIONAL_VCL")
			servicePort := viper.GetString("SERVICE_PORT")

			fmt.Printf("Path %s \n", path)
			fmt.Printf("Service to Watch %s \n", serviceToWatch)
			fmt.Printf("Additional VCL %s \n", additionalVCL)
			fmt.Printf("Service Port %s\n", servicePort)

			varnishConfig := make(chan types.VarnishConfiguration)

			rancher.Provide(serviceToWatch, varnishConfig)

			for {
				select {
				case cfg := <-varnishConfig:
					fmt.Printf("Render Config \n")
					render.RenderConfig(path, cfg.Backends, additionalVCL, servicePort)
					fmt.Printf("Execute Command\n")
					erc := rancher.ExecuteReloadCommand(cfg.Host)
					websocket.DoSocketCall(erc)
					fmt.Printf("done")
				}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() { 
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.varaco.yaml)")

	rootCmd.PersistentFlags().String("cattle.url", "", "config file (default is $HOME/.varaco.yaml)")
	rootCmd.PersistentFlags().String("cattle.access.key", "", "config file (default is $HOME/.varaco.yaml)")
	rootCmd.PersistentFlags().String("cattle.secret.key", "", "config file (default is $HOME/.varaco.yaml)")
	rootCmd.PersistentFlags().String("file.path", "", "config file (default is $HOME/.varaco.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.BindPFlag("cattle.url", rootCmd.PersistentFlags().Lookup("cattle.url"))
	viper.BindPFlag("cattle.access.key", rootCmd.PersistentFlags().Lookup("cattle.access.key"))
	viper.BindPFlag("cattle.secret.key", rootCmd.PersistentFlags().Lookup("cattle.secret.key"))
	viper.BindPFlag("file.path", rootCmd.PersistentFlags().Lookup("file.path"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".varaco" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".varaco")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initRancherConfig() {
	rancherConfig := new(configuration.Rancher)

	if viper.IsSet("cattle.url") && viper.GetString("cattle.url") != "" {
		rancherConfig.URL = viper.GetString("cattle.url")
	} else if viper.IsSet("CATTLE_URL"){
		rancherConfig.URL = viper.GetString("CATTLE_URL");
	}

	if viper.IsSet("cattle.access.key")  && viper.GetString("cattle.access.key") != "" {
		rancherConfig.AccessKey = viper.GetString("cattle.access.key")
	} else if viper.IsSet("CATTLE_ACCESS_KEY"){
		rancherConfig.AccessKey = viper.GetString("CATTLE_ACCESS_KEY");
	}

	if viper.IsSet("cattle.secret.key")  && viper.GetString("cattle.secret.key") != "" {
		rancherConfig.SecretKey = viper.GetString("cattle.secret.key")
	} else if viper.IsSet("CATTLE_SECRET_KEY"){
		rancherConfig.SecretKey = viper.GetString("CATTLE_SECRET_KEY");
	}

	config = *rancherConfig
}