// +build ignore

package main

import "os"
import "os/exec"

func main() {
	cmd := exec.Command("go", "test", "-c", "./test/integration")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
