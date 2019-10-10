package bifrost

import (
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
	LogFilter        regexp.Regexp
}

func Execute(config ManagerConfig) error {
	return execute(config)
}

func execute(config ManagerConfig) error {
	for i := 0; i < config.Repeat; i++ {
		command := exec.Command(config.AbsolutePath, config.ProgramArguments...)

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
	}

	return nil
}
