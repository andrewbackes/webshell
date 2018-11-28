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
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Could not get pipe for stdout - %v", err)
		return
	}
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Could not start command - %v", err)
		return
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			stdOut.Write(scanner.Bytes())
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading stdout - %v", err)
		}
	}()
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Could not wait for command completion - %v", err)
	}
}
