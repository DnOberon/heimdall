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
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/dnoberon/heimdall/bifrost"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Args:  cobra.NoArgs,
	Short: `Run heimdall using the "heimdall_config.json" file in the current directory `,
	Long: `You must first generate a "heimdall_config.json" file before running
this command. If a file does not already exist please run
"heimdall init"`,
	Run: func(cmd *cobra.Command, args []string) {
		config := bifrost.ManagerConfig{}

		file, err := ioutil.ReadFile("heimdall_config.json")
		if err != nil {
			log.Fatal(err)
			return
		}

		err = json.Unmarshal(file, &config)
		if err != nil {
			log.Fatal(err)
			return
		}

		timeout, err := time.ParseDuration(config.TimeoutString)
		if err != nil {
			log.Fatal(err)
			return
		}

		config.Timeout = timeout
		bifrost.Execute(config)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
