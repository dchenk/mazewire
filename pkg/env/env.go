// Package env helps determine the environment variables with which the application is running.
//
// There are two ways to configure the environment:
//
//   1. Define variables in a file with key=value pairs separated by a newline character.
//
//      Each variable's name ends when the first '=' character is reached, and everything following
//      that on the line is the variable's value. Empty lines and empty variable values are permitted.
//      A line is ignored if it begins with the '#' character.
//
//      To use this method, simply start the program with no arguments and have a "env.txt" file in
//      the working directory where the program is running. If you would like to specify the variables
//      in another file, provide a "-env-file=/path/to/vars-file" argument.
//
//   2. Define variables as normal environment variables as used in the operating system's shell.
//      To use this method, start the program with the argument "-env=vars".
//
// The following environment variables must be defined for server instance when it is started:
//
//   DB_CONNECTION
//     A connection string that the database system in use will use to connect to the database, without any
//     database name or other connection parameters.
//
//   DB_NAME
//     The name of the database to be used for storing all application data other than static media files.
//
//
// The following are standard recognized environment variables that may or may not be used depending on
// the way you're running the system:
//
//   DB_PARAMS
//     A set of key=value pairs used to configure how the connection to the database is made.
//
//   PLUGINS_DIR
//     The directory in which plugin binaries are stored. The default value is a directory named "plugins"
//     in the current working directory of the main application.
//
//   CHANGE_TOKEN_KEY
//     This variable, if set to the value "yes", will indicate to the server that the TOKEN_KEY variable that
//     is also set now must be used for user log-ins and that all server instances should be updated to use
//     this new value. Setting this variable to "yes" without setting the TOKEN_KEY variable is an error.
//
//   GCP_PROJECT
//     The ID of the Google Cloud Platform project within which the application is running and using the
//     various services used within the app.
//
//
// The following are variables that you need to set when initializing your cluster for the very first
// time. After that, new instances should not have these variables set at startup.
// Note that a root user is created upon initialization. This user has their first name and last name set
// to the default value of "root", which may be changed using the Admin UI.
//
//   INIT
//     A value of "yes" for this variable indicates that the cluster should be initialized.
//
//   TOKEN_KEY
//     A random string that is used to authenticate user log-ins. The same value should be used on all server
//     instances of the application. The value must be at least 20 characters long.
//     You may force all logged in users to be logged out by re-deploying your system with a new value for
//     this variable along with CHANGE_TOKEN_KEY set to "yes".
//
//   ADMIN_UI_ROOT
//     The full URL of the root where the Admin UI is accessed.
//     After the system is booted up the first time, new server instances don't need to set this variable
//     if it's been saved in the database, but new instances without this variable defined may be started
//     only after the value is already saved in the database. The value can be edited using the Admin UI.
//
//   ADMIN_SRC_ROOT
//     The full root URL from which static resources for the admin UI are downloaded.
//     After the system is booted up the first time, new server instances don't need to set this variable
//     if it's been saved in the database, but new instances without this variable defined may be started
//     only after the value is already saved in the database. The value can be edited using the Admin UI.
//
//   ROOT_USER
//     The username to use for the root user account. The values "admin" and "root" are not permitted.
//
//   ROOT_EMAIL
//     The email address by which the root user may be reached.
//
//   ROOT_PASS
//     The password to use for the root user account. The value must be at least 8 characters long.
//     For just this variable, any spaces immediately to the right of the '=' and trailing the value are
//     trimmed out (so the password will not have any leading or trailing spaces).
//
//
// Plugins may require environment variables to be set any way they like, though it's recommended to use
// one of the standard methods described above. Plugins at startup are given all environment variables
// that start with the string "WW_".
//
package env

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	envConfigFile   = "file"
	envConfigVars   = "vars"
	fileDefaultPath = "env.txt"

	// Required variables.
	VarDbConnection = "DB_CONNECTION"
	VarDbName       = "DB_NAME"

	// Other standard variables.
	VarDbParams       = "DB_PARAMS"
	VarPluginsDir     = "PLUGINS_DIR"
	VarChangeTokenKey = "CHANGE_TOKEN_KEY"
	VarGcpProject     = "GCP_PROJECT"

	// Variables used for initializing the cluster.
	VarInit         = "INIT"
	VarTokenKey     = "TOKEN_KEY"
	VarAdminUIRoot  = "ADMIN_UI_ROOT"
	VarAdminSrcRoot = "ADMIN_SRC_ROOT"
	VarRootUser     = "ROOT_USER"
	VarRootEmail    = "ROOT_EMAIL"
	VarRootPass     = "ROOT_PASS"

	// varDevMode is used for developing this system.
	// If set to the value "y", it indicates that the application should run in development mode.
	varDevMode = "DEV_MODE"
)

var requiredVars = []string{VarDbConnection, VarDbName}

var standardVars = append(requiredVars, VarDbParams, VarPluginsDir, VarChangeTokenKey, VarGcpProject)

var initVars = []string{VarInit, VarTokenKey, VarAdminUIRoot, VarAdminSrcRoot, VarRootUser, VarRootEmail, VarRootPass}

var (
	flagEnv = flag.String("env", envConfigFile,
		fmt.Sprintf("How the server environment is configured; either %q or %q", envConfigFile, envConfigVars))
	flagPath = flag.String("env-file", fileDefaultPath,
		"Where the environment variables are defined if using the file method")
)

// vars contains the environment variables that are set once upon program init.
var vars map[string]string

// Vars returns the environment variables with which the application is running.
func Vars() map[string]string {
	m := make(map[string]string, len(vars))
	for k, v := range vars {
		m[k] = v
	}
	return m
}

// checkForRequiredVars returns a list of all of the standard environment variables that have not been set.
func checkForRequiredVars(vs map[string]string) (missing []string) {
	for _, v := range requiredVars {
		if val, ok := vs[v]; !ok || strings.TrimSpace(val) == "" {
			missing = append(missing, v)
		}
	}
	return
}

// parseFile returns a map of all of the variables defined in the variables file located at filePath.
func parseFile(filePath string) (map[string]string, error) {
	envBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read env file: %v", err)
	}

	lines := bytes.Split(envBytes, []byte{'\n'})
	m := make(map[string]string, len(lines))

	key := make([]byte, 0, 24)
	for _, v := range lines {
		v = bytes.TrimSpace(v)
		ln := len(v)
		if ln == 0 || v[0] == '#' {
			continue
		}

		// Extract the name of the variable.
		key = key[:0]
		i := 0
		for i < ln {
			char := v[i]
			i++
			if char == '=' {
				break
			}
			key = append(key, char)
		}

		k := string(key)
		val := string(v[i:])
		if k == VarRootPass {
			// The password may not begin or end with spaces.
			val = strings.TrimSpace(val)
		}
		m[k] = val
	}

	// Set default values.
	for k, v := range defaults {
		if m[k] == "" {
			m[k] = v
		}
	}

	return m, nil
}

var defaults = map[string]string{
	VarPluginsDir: "plugins",
}

// parseVars returns a map of the standard and the initialization environment variables used by the
// system along with all of the variables that begin with the string "WW_".
func parseVars() (map[string]string, error) {
	toCheck := append(standardVars, initVars...)
	m := make(map[string]string, len(toCheck))
	for _, v := range toCheck {
		m[v] = os.Getenv(v)
	}
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "WW_") {
			split := strings.Split(v, "=")
			m[split[0]] = split[1]
		}
	}
	return m, nil
}

// Init says whether the INIT variable is set so that the cluster may be initialized.
func Init() bool {
	return vars[VarInit] == "yes"
}

// Prod says whether the server is running in production mode.
func Prod() bool {
	return vars[varDevMode] != "y"
}
