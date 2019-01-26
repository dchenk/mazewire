package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const googleServiceAccountFile = "google-sa.json"

const dbType = "cockroach"

// Start starts up CockroachDB and builds and runs the app locally.
func Start(_ []string) error {

	cockroachCmd, err := startCockroach()
	if err != nil {
		return err
	}
	defer func() {
		if cockroachCmd.Process == nil {
			fmt.Println("cockroach process is nil upon exiting")
			return
		}
		err := cockroachCmd.Process.Signal(syscall.SIGHUP)
		if err != nil {
			fmt.Println("cockroach process terminate error:", err)
		}
		cockroachCmd.Wait() // no need to get any error
	}()

	// Let cockroach fully start up first (to not clog up the terminal output).
	time.Sleep(time.Millisecond * 300)
	dbPing := func(i int) bool {
		db, err := sql.Open("postgres", envVars["DB_CONNECTION"]+envVars["DB_PARAMS"])
		if err != nil {
			if i == 2 {
				fmt.Println(colorErr("Opening DB connection failed on third try: %v", err))
			}
			return false
		}
		if err := db.Ping(); err != nil {
			fmt.Println(colorErr("Could not ping the database: %v", err))
			return false
		}
		return true
	}
	for i := 0; i < 3; i++ {
		if dbPing(i) {
			break
		}
	}

	compileDaemonPath, err := exec.LookPath("CompileDaemon")
	if err != nil {
		return fmt.Errorf("could not find CompileDaemon executable; %v", err)
	}

	buildCmd := exec.Command(compileDaemonPath, fmt.Sprintf("-build=\"go build -tags=%v\"", dbType), "-command=./main/main", "-build-dir=main", "-recursive=true",
		"-exclude-dir=.git", "-exclude-dir=cockroach-data", "-exclude-dir=sql", "-exclude-dir=front", "-exclude-dir=run", "-exclude-dir=dns_certs", "-exclude-dir=try",
		"-include=main", "-include=data", "-include=users")
	buildCmd.Stderr = os.Stderr
	buildCmd.Stdout = os.Stdout
	buildCmd.Env = append(os.Environ(), "GOOGLE_APPLICATION_CREDENTIALS="+googleServiceAccountFile)

	if err := buildCmd.Start(); err != nil {
		return fmt.Errorf("could not start CompileDaemon; %v", err)
	}
	defer func() {
		if buildCmd.Process == nil {
			fmt.Println("CompileDaemon process is nil upon exiting")
			return
		}
		if err := buildCmd.Process.Kill(); err != nil {
			fmt.Println("CompileDaemon process kill error:", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until an interrupt.
	<-c

	return nil

}
