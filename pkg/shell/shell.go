package shell

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

const sh = "sh"

// Run a shell command.
func Run(command []byte, stdOut io.Writer) {
	cmd := exec.Command(sh, "-c", string(command))
	out, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Could not get pipe for stdout - %v", err)
		return
	}
	errOut, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Could not get pipe for stderr - %v", err)
		return
	}
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Could not start command - %v", err)
		return
	}
	go func() {
		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			stdOut.Write(scanner.Bytes())
		}
	}()
	go func() {
		scanner := bufio.NewScanner(errOut)
		for scanner.Scan() {
			stdOut.Write(scanner.Bytes())
		}
	}()
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Could not wait for command completion - %v", err)
	}
}
