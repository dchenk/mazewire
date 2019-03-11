// Package env helps determine the environment variables with which the application is running.
//
// See the docs on the env file under /docs.
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
	// Required variables.
	VarDbConnection = "DB_CONNECTION"
	VarDbName       = "DB_NAME"

	// Other standard variables.
	VarDbParams       = "DB_PARAMS"
	VarPluginsDir     = "PLUGINS_DIR"
	VarChangeTokenKey = "CHANGE_TOKEN_KEY"
	VarGcpProject     = "GCP_PROJECT"

	// Variables used to initialize the cluster.
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
	flagEnvFile = flag.String("env-file", "",
		"Where the environment variables are defined if using the file method")
)

// vars contains the environment variables that are set at program init.
var vars map[string]string

// Vars returns a copy of the environment variables with which the application is running.
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
// system along with all of the variables that begin with the string "MW_".
func parseVars() (map[string]string, error) {
	toCheck := append(standardVars, initVars...)
	m := make(map[string]string, len(toCheck))
	for _, v := range toCheck {
		m[v] = os.Getenv(v)
	}
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "MW_") {
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
