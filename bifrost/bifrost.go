package bifrost

import (
	"bufio"
	"os"
	"os/exec"
	"regexp"
	"time"
)

type ManagerConfig struct {
	AbsolutePath     string
	ProgramArguments []string
	Timeout          time.Duration
	Repeat           int
	Verbose          bool
	Log              bool
	LogFilter        *regexp.Regexp
}

func Execute(config ManagerConfig) error {
	return execute(config)
}

func execute(config ManagerConfig) error {
	config.Repeat++

	for i := 0; i < config.Repeat; i++ {
		command := exec.Command(config.AbsolutePath, config.ProgramArguments...)
		log, err := os.Create("heimdall.log")
		if err != nil {
			return err
		}

		stdoutDone := make(chan interface{})
		stderrDone := make(chan interface{})

		stdout, _ := command.StdoutPipe()
		stderr, _ := command.StderrPipe()

		go func() {
			rd := bufio.NewReader(stdout)

			for {
				str, err := rd.ReadString('\n')
				if err != nil {
					break
				}

				if config.Verbose {
					os.Stdout.Write([]byte(str))
				}

				if config.Log {
					if config.LogFilter != nil {
						if config.LogFilter.MatchString(str) {
							log.Write([]byte(str))
						}
					} else {
						log.Write([]byte(str))
					}
				}

			}
			// do something with stdout
			close(stdoutDone)
		}()

		go func() {
			rd := bufio.NewReader(stderr)

			for {
				str, err := rd.ReadString('\n')
				if err != nil {
					break
				}

				if config.Verbose {
					os.Stdout.Write([]byte(str))
				}

				if config.Log {
					if config.LogFilter != nil {
						if config.LogFilter.MatchString(str) {
							log.Write([]byte(str))
						}
					} else {
						log.Write([]byte(str))
					}
				}
			}

			// do something with stderr
			close(stderrDone)
		}()

		// everything else

		if err := command.Start(); err != nil {
			return err
		}

		if config.Timeout > 0 {
			time.AfterFunc(config.Timeout, func() {
				command.Process.Kill()
			})
		}

		if err := command.Wait(); err != nil {
			return err
		}
		<-stdoutDone
		<-stderrDone
	}

	return nil
}
