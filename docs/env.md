# Configuring the Environment

Package `env` helps determine the environment variables with which the application is running.

There are two ways to configure the environment:

2. Define variables as normal environment variables for a process. To use this method, start the
main program without an `env` argument.

1. Define variables in a file with key=value pairs, one per line. To use this method, start the
program with the argument `-env` to have the configuration read in from the file named "env.txt" in
the main program's working directory, or give the `-env` flag and the name of the file with
configuration with the `-env-file` argument (for example, `-env-file=~/mazewire-env.txt`). In such a
file, each variable's name ends when the first '=' character is reached, and everything between the
equals character and the newline character is taken as the variable's value. Empty lines are
skipped, and empty values are permitted. Lines that begin with the '#' character are skipped.


The following environment variables must be defined for each server instance to be started:

- DB_CONNECTION

  A connection string that the system will use to connect to the database, without the database name
  or other connection parameters.

- DB_NAME

  The name of the database to be used for storing application data other than static media files.


The following are standard recognized environment variables that may or may not be used depending on
the way you're running Mazewire:

- DB_PARAMS

  A set of key=value pairs used to configure how the connection to the database is made.

- PLUGINS_DIR

  The directory in which plugin binaries are stored. The default value is a directory named
  "plugins" in the current working directory of the main application. This directory will be
  created by Mazewire if it does not already exist.

- CHANGE_TOKEN_KEY

  This variable, if set to "yes", will indicate to the server that the TOKEN_KEY variable that is
  also set now must be used for user log-ins and that all server instances should be updated to use
  this new value. Setting this variable to "yes" without setting the TOKEN_KEY variable is an error.

- GCP_PROJECT

  The ID of the Google Cloud Platform project in which the application is running.


The following are variables that you need to set when initializing your cluster for the first time.
After the first initialization, new instances should not have these variables set at startup. The
root user's first name and last name set to the value "root" and can be changed using the Admin UI.

- INIT

  A value of "yes" for this variable indicates that the cluster should be initialized.

- TOKEN_KEY

  A random string that is used by the authentication system. The same value should be used on all
  server instances. The value must be at least 20 characters long. You cay force all logged in users
  to be logged out by re-deploying your system with a new value for this variable along with the
  value of CHANGE_TOKEN_KEY set to "yes".

- ADMIN_UI_ROOT

  The full URL of the root where the Admin UI is accessed. After the system is booted up the first
  time, new server instances don't need to set this variable if it's been saved in the database, but
  new instances without this variable defined may be started only after the value is already saved
  in the database. The value can be edited using the Admin UI.

- ADMIN_SRC_ROOT

  The full root URL from which static resources for the admin UI are downloaded. After the system is
  booted up the first time, new server instances don't need to set this variable if it's been saved
  in the database, but new instances without this variable defined may be started only after the
  value is already saved in the database. The value can be edited using the Admin UI.

- ROOT_USER

  The username to use for the root user account. The values "admin" and "root" are not permitted.

- ROOT_EMAIL

  The email address by which the root user may be reached.

- ROOT_PASS

  The password to use for the root user account. The value must be at least 8 characters long.
  For just this variable, any spaces immediately to the right of the '=' and trailing the value are
  trimmed out (so the password will not have any leading or trailing spaces).


Plugins may require additional environment variables to be set. Plugin processes are started with
all of the environment variables that start with the string "MW_".

