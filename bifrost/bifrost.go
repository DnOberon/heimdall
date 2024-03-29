// Copyright 2019 John Darrington johnw.darrington@gmail.com

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

// Package bifrost encompasses all functions related to the short-lived process manager by the same name. Heimdall
// was the ever-vigilant guardian of the gods' stronghold, Asgard - now he will be the guardian of whichever program you choose.
// Heimdall is designed as both launcher and monitor of short-lived CLI tools and programs. Heimdall provides the ability
// to automatically repeat a process, kill a hung process started with the tool, and log the programs output (filtering logs
// is also possible). It is hoped that heimdall and bifrost will be a tool you reach for again and again when developing your CLI tool.
package bifrost

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"time"
)

// ManagerConfig manages configuration for the bifrost execution function
type ManagerConfig struct {
	AbsolutePath     string
	ProgramArguments []string

	Timeout       time.Duration `json:"-"`
	TimeoutString string
	Repeat        int

	InParallelCount int

	Log          bool
	LogName      string
	LogOverwrite bool
	LogFilter    *regexp.Regexp
	Verbose      bool

	logFile *os.File
	lock    *sync.Mutex
	wg      sync.WaitGroup
}

// Execute accepts a configuration and attempts to run the provided program and its arguments.
func Execute(config ManagerConfig) error {
	return execute(config)
}

func execute(config ManagerConfig) error {
	config.lock = &sync.Mutex{}
	config.wg = sync.WaitGroup{}

	if config.Log && config.LogOverwrite {
		f, err := os.Create(config.LogName)
		if err != nil {
			return err
		}

		defer f.Close()

		config.logFile = f
	} else if config.Log {
		f, err := os.OpenFile(config.LogName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		defer f.Close()

		config.logFile = f
	}

	for p := 0; p < config.InParallelCount; p++ {
		config.wg.Add(1)

		go func() {
			for i := 0; i < config.Repeat; i++ {
				command := exec.Command(config.AbsolutePath, config.ProgramArguments...)

				// attachment of reader/writers to command execution
				stdoutDone, stderrDone := attachLogger(command, config)

				if err := command.Start(); err != nil {
					//return err
				}

				if config.Timeout > 0 {
					time.AfterFunc(config.Timeout, func() {
						command.Process.Kill()
					})
				}

				if err := command.Wait(); err != nil {
					//	return err
				}

				// we're not going to wait on these to finish if we don't need them. Helps us avoid infinite loops if we somehow
				// screwed up the reader/write initiation
				if config.LogFilter != nil || config.Log || config.Verbose {
					<-stdoutDone
					<-stderrDone
				}

			}

			config.wg.Done()
		}()

	}

	config.wg.Wait()

	return nil
}

func attachLogger(cmd *exec.Cmd, config ManagerConfig) (stdoutDone chan interface{}, stderrDone chan interface{}) {
	stdoutDone = make(chan interface{})
	stderrDone = make(chan interface{})

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// while we could make a single function and simply assign the output writer, I choose to keep stdout and stderr
	// separate for ease of reading and understanding by those new to Go.
	go func() {
		rd := bufio.NewReader(stdout)

		for {
			str, err := rd.ReadString('\n')
			if err != nil {
				break
			}

			out := fmt.Sprintf("[H-PID:%d %s]  %s", cmd.Process.Pid, time.Now().UTC().Format("06-01-02 15:04:05"), str)

			if config.Verbose {
				os.Stdout.Write([]byte(out))
			}

			if config.LogFilter != nil {
				if config.LogFilter.MatchString(str) {
					config.lock.Lock()
					config.logFile.Write([]byte(out))
					config.lock.Unlock()
				}

			} else if config.Log {
				config.lock.Lock()
				config.logFile.Write([]byte(out))
				config.lock.Unlock()
			}

		}
		close(stdoutDone)
	}()

	go func() {
		rd := bufio.NewReader(stderr)

		for {
			str, err := rd.ReadString('\n')
			if err != nil {
				break
			}

			out := fmt.Sprintf("[H-PID:%d %s]  %s", cmd.Process.Pid, time.Now().UTC().Format("06-01-02 15:04:05"), str)

			if config.Verbose {
				os.Stdout.Write([]byte(out))
			}

			if config.LogFilter != nil {
				if config.LogFilter.MatchString(str) {
					config.lock.Lock()
					config.logFile.Write([]byte(out))
					config.lock.Unlock()
				}

			} else if config.Log {
				config.lock.Lock()
				config.logFile.Write([]byte(out))
				config.lock.Unlock()
			}
		}

		close(stderrDone)
	}()

	return stdoutDone, stderrDone
}
