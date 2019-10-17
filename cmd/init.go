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
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/objconv/json"

	"github.com/dnoberon/heimdall/bifrost"

	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Args:  cobra.NoArgs,
	Short: "Create a configuration for heimdall to replace command flag arguments",
	Long: `Init creates a "heimdall_config.json" in the current directory. This configuration
will be read in when using the "heimdall run" command and
should take the place of any command flags you would use
normally`,
	Run: func(cmd *cobra.Command, args []string) {
		config := bifrost.ManagerConfig{}

		promptProgramPath(&config)
		promptProgramArguments(&config)
		promptRepeat(&config)
		promptTimeout(&config)
		promptParallel(&config)
		promptVerbose(&config)
		promptLog(&config)

		configFile, err := json.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile("heimdall_config.json", configFile, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func promptProgramPath(config *bifrost.ManagerConfig) {
	prompt := promptui.Prompt{
		Label:    "Executable path",
		Validate: emptyValidate,
	}

	path, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	config.AbsolutePath = absolutePath
}

func promptProgramArguments(config *bifrost.ManagerConfig) {
	prompt := promptui.Prompt{
		Label: "Program arguments separated by comma",
	}

	arguments, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	if arguments == "" {
		return
	}

	config.ProgramArguments = strings.Split(arguments, ",")
}

func promptTimeout(config *bifrost.ManagerConfig) {
	prompt := promptui.Prompt{
		Label:    "How long should we wait before killing your program? - e.g 10s, 1m, 1h ",
		Default:  "5m",
		Validate: timeValidate,
	}

	timeout, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	_, err = time.ParseDuration(timeout)
	if err != nil {
		log.Fatal(err)
	}

	config.TimeoutString = timeout
}

func promptRepeat(config *bifrost.ManagerConfig) {
	prompt := promptui.Prompt{
		Label:    "How many times should we repeat execution?",
		Default:  "1",
		Validate: intValidate,
	}

	repeat, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	repeatAmount, err := strconv.Atoi(repeat)
	if err != nil {
		log.Fatal(err)
	}

	config.Repeat = repeatAmount
}

func promptParallel(config *bifrost.ManagerConfig) {
	prompt := promptui.Prompt{
		Label:    "How many instances of your program should we run at the same time?",
		Default:  "1",
		Validate: intValidate,
	}

	instances, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	parallelCount, err := strconv.Atoi(instances)
	if err != nil {
		log.Fatal(err)
	}

	config.InParallelCount = parallelCount
}

func promptVerbose(config *bifrost.ManagerConfig) {
	prompt := promptui.Prompt{
		Label:    "Print all output to console? [Y/n] ",
		Validate: confirmValidate,
		Default:  "y",
	}

	confirm, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	config.Verbose = isYes(confirm)
}

func promptLog(config *bifrost.ManagerConfig) {
	prompt := promptui.Prompt{
		Label:    "Log the output generated by your program? [Y/n] ",
		Default:  "y",
		Validate: confirmValidate,
	}

	confirm, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	if !isYes(confirm) {
		return
	}

	prompt = promptui.Prompt{
		Label:   "What should we call the log file?",
		Default: "heimdall.log",
	}

	logName, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	config.LogName = logName

	prompt = promptui.Prompt{
		Label:    "Overwrite existing logs? [y/N] ",
		Validate: confirmValidate,
		Default:  "n",
	}

	confirm, err = prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	config.LogOverwrite = isYes(confirm)

	prompt = promptui.Prompt{
		Label:    "Filter incoming logs? [y/N] ",
		Validate: confirmValidate,
		Default:  "n",
	}

	confirm, err = prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	if isYes(confirm) {
		prompt = promptui.Prompt{
			Label:    "Regex expression for filtering logs",
			Default:  "",
			Validate: regexValidate,
		}

		expression, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		compiled, err := regexp.Compile(expression)
		if err != nil {
			log.Fatal(err)
		}

		config.LogFilter = compiled
	}

}

func emptyValidate(input string) error {
	if input == "" {
		return errors.New("Invalid text")
	}

	return nil
}

func confirmValidate(input string) error {
	if input == "y" || input == "Y" || input == "N" || input == "n" || input == "yes" || input == "no" {
		return nil
	}

	return errors.New("Bad confirmation (yes, y, no, n)")
}

func isYes(input string) bool {
	return input == "y" || input == "yes" || input == "Y"
}

func intValidate(input string) error {
	_, err := strconv.Atoi(input)

	if err != nil {
		return errors.New("Invalid number")
	}

	return nil
}

func regexValidate(input string) error {
	_, err := regexp.Compile(input)
	return err
}

func timeValidate(input string) error {
	_, err := time.ParseDuration(input)

	return err
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
