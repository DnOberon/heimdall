package bifrost

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBifrostExecuteAndRepeat(t *testing.T) {
	timer := time.Now()
	path, _ := filepath.Abs("../heimdall")

	hProcess := exec.Command(path, "--repeat=2", "tester")

	err := hProcess.Start()
	assert.Nil(t, err)

	err = hProcess.Wait()
	assert.Nil(t, err)

	assert.True(t, hProcess.ProcessState.Success())
	assert.True(t, time.Now().Sub(timer) > 10*time.Second)
}

func TestBifrostExecuteTimeout(t *testing.T) {
	timer := time.Now()
	path, _ := filepath.Abs("heimdall")

	hProcess := exec.Command(path, "--repeat=1", "--timeout=10s", "tester", "30")

	err := hProcess.Start()
	assert.Nil(t, err)

	hProcess.Wait()

	assert.True(t, hProcess.ProcessState.Exited())
	assert.True(t, time.Now().Sub(timer) > 10*time.Second)
	assert.True(t, time.Now().Sub(timer) < 30*time.Second)
}

func TestBifrostLogging(t *testing.T) {
	path, _ := filepath.Abs("../heimdall")

	hProcess := exec.Command(path, "--repeat=2", "tester")

	pipe, err := hProcess.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	err = hProcess.Start()
	assert.Nil(t, err)

	rd := bufio.NewReader(pipe)

	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			break
		}

		fmt.Println(str)
	}

	err = hProcess.Wait()
	assert.Nil(t, err)
}
