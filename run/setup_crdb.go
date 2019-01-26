package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const databaseSchemaFile = "sql/db.sql" // CWD is where Makefile is.

func SetupCockroach(_ []string) error {

	schemaFile, err := os.Open(databaseSchemaFile)
	if err != nil {
		return err
	}
	defer schemaFile.Close()

	var schemas []dbSchema
	rdr := bufio.NewReader(schemaFile)

	inStatement := false // For the loop, indicate if we're inside a statement.
	for {
		line, err := rdr.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("reading the schema file: %v", err)
		}
		if strings.HasPrefix(line, "--") {
			if strings.HasPrefix(line, "-- END") {
				break
			}
			inStatement = false
			continue
		}
		if strings.TrimSpace(line) == "" {
			inStatement = false
			continue
		}
		line = strings.TrimSuffix(line, ";\n") // Remove trailing semicolons
		if inStatement {
			schemas[len(schemas)-1].statement += line
		} else {
			schemas = append(schemas, dbSchema{
				lineStart: line,
				statement: line,
			})
		}
		inStatement = true
	}

	cockroachCmd, err := startCockroach()
	if err != nil {
		return fmt.Errorf("could not start cockroach; %v", err)
	}
	defer func() {
		if cockroachCmd.Process == nil {
			fmt.Println("cockroach process is nil upon exiting")
			return
		}
		err := cockroachCmd.Process.Kill()
		if err != nil {
			fmt.Println("cockroach kill error:", err)
		}
	}()

	// Make sure the "cockroach start" command has started cockroach.
	time.Sleep(time.Second * 2)

	db, err := sql.Open("postgres", envVars["DB_CONNECTION"]+envVars["DB_PARAMS"])
	if err != nil {
		return fmt.Errorf("opening first DB connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping the database: %v", err)
	}

	// The first statement must be made manually because it depends on the database name.
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + envVars["DB_NAME"])
	if err != nil {
		return fmt.Errorf("statement to create database failed with error: %v", err)
	}

	if err = db.Close(); err != nil {
		// Simply output the error, do not return.
		fmt.Println("error closing first DB: ", err)
	}

	// Now we can use the database created, and we need a new sql.DB with the database name specified.
	db, err = sql.Open("postgres", envVars["DB_CONNECTION"]+"/"+envVars["DB_NAME"]+envVars["DB_PARAMS"])
	if err != nil {
		return fmt.Errorf("opening second DB connection: %v", err)
	}

	tx, txErr := db.Begin()
	if txErr != nil {
		return txErr
	}
	defer func() {
		if txErr == nil {
			commitErr := tx.Commit()
			if commitErr != nil {
				fmt.Printf("Could not commit the transaction: %v\n", commitErr)
			}
		} else {
			tx.Rollback()
		}
	}()

	for i := range schemas {
		_, txErr = tx.Exec(schemas[i].statement)
		if txErr != nil {
			return fmt.Errorf("the statement beginning with %q failed: %v", schemas[i].statement[:15], txErr)
		}
	}

	return nil
}

type dbSchema struct {
	lineStart string
	statement string
}
