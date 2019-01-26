package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// startCockroach starts cockroach in a goroutine and returns the exec.Cmd to control it.
func startCockroach() (*exec.Cmd, error) {
	cockroachPath, err := exec.LookPath("cockroach")
	if err != nil {
		return nil, fmt.Errorf("could not find cockroach executable; %v", err)
	}
	cockroachStart := exec.Command(cockroachPath, "start", "--host=localhost", "--port="+CockroachPort(), "--insecure")
	cockroachStart.Stdout = os.Stdout
	cockroachStart.Stderr = os.Stderr
	return cockroachStart, cockroachStart.Start()
}

func StartWaitCRDB(_ []string) error {
	cmd, err := startCockroach()
	if err != nil {
		return err
	}
	defer func() {
		if cmd.Process == nil {
			fmt.Println("cockroach process is nil upon exiting")
			return
		}
		err := cmd.Process.Signal(syscall.SIGHUP)
		if err != nil {
			fmt.Println("cockroach process terminate error:", err)
		}
		cmd.Wait()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until an interrupt.
	<-c
	return nil
}
