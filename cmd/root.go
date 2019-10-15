/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/dnoberon/heimdall/bifrost"

	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "heimdall",
	Short: "Heimdall is a quickly configured monitor for CLI applications.",
	Long: `Heimdall gives you a quick way to monitor, repeat, and selectively 
log a CLI application. Quick configuration options allow 
you to effectively test and monitor a CLI application in
development`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		absolutePath, err := filepath.Abs(args[0])
		if err != nil {
			log.Fatal("Unable to locate the executable provided " + err.Error())
		}

		repeat, _ := cmd.Flags().GetInt("repeat")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		parallelCount, _ := cmd.Flags().GetInt("parallelCount")

		toLog, _ := cmd.Flags().GetBool("log")
		logName, _ := cmd.Flags().GetString("logName")
		logOverwrite, _ := cmd.Flags().GetBool("logOverwrite")
		rawReg, _ := cmd.Flags().GetString("logFilter")
		verbose, _ := cmd.Flags().GetBool("verbose")

		config := bifrost.ManagerConfig{
			AbsolutePath:     absolutePath,
			Verbose:          verbose,
			Repeat:           repeat,
			ProgramArguments: args[1:],
			InParallelCount:  parallelCount,
			Log:              toLog,
			LogName:          logName,
			LogOverwrite:     logOverwrite,
			Timeout:          timeout}

		if rawReg != "" {
			regex, err := regexp.Compile(rawReg)
			if err != nil {
				log.Fatal(err)
			}

			config.LogFilter = regex
		}

		err = bifrost.Execute(config)
		if err != nil {
			log.Fatal(err)
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

	rootCmd.Flags().IntP("repeat", "r", 1, "Designate how many times to repeat your program with supplied arguments")
	rootCmd.Flags().DurationP("timeout", "t", 0, "Designate when to kill your provided program")

	rootCmd.Flags().IntP("parallelCount", "p", 1, "Designate how many instances of your should run in parallel at one time")

	rootCmd.Flags().BoolP("log", "l", false, "Toggle logging of provided program's stdout and stderr output to file, appends if file exists")
	rootCmd.Flags().String("logName", "heimdall.log", "Specify the log file name, defaults to heimdall.log")
	rootCmd.Flags().Bool("logOverwrite", false, "Toggle logging of provided program's stdout and stderr output to file")
	rootCmd.Flags().String("logFilter", "", "Allows for log filtering via regex string. Use only valid with log flag")
	rootCmd.Flags().BoolP("verbose", "v", false, "Toggle display of provided program's stdout and stderr output while heimdall runs")
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

		// Search config in home directory with name ".bifrost" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bifrost")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
