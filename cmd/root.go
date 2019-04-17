// Copyright © 2019 Mohammed Al-Ameen <mohammed.alameen@protonmail.ch>
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
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     configuration
)

var ver = "0.1"

type dbConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	SslModeOn bool
	DbName    string
}

type configuration struct {
	DB dbConfig
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "serverlessdb",
	Short: "Event-driven databases | ver: " + ver,
	Long:  `Event-driven databases | ver: ` + ver,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cfg.DB.DbName)
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.serverlessdb.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".serverlessdb" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".serverlessdb")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		cfg = parseConfig()
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		viper.SetEnvPrefix("NETLIFY")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
	}
}

func parseConfig() (c configuration) {
	c = configuration{
		DB: dbConfig{
			Host:      viper.GetString("database.hostname"),
			Port:      viper.GetInt("database.port"),
			Username:  viper.GetString("database.username"),
			Password:  viper.GetString("database.password"),
			SslModeOn: viper.GetBool("database.sslmodeon"),
			DbName:    viper.GetString("database.dbname"),
		},
	}
	return
}
