package flags

import (
	"flag"
	"fmt"
	"os"
)

func printHelpMessage() {
	fmt.Printf(`Usage: %s USERNAME PASSWORD/[OTP]

Generates a password grant token from Keycloak for the given user.

Options:
  -h, --help              Show this help message and exit
  -v, --version           Show version information

Notes:
  For the program to function properly, it needs to locate a configuration file called 'config.toml'.
  The program will search for this file in the current directory, default install '/opt/kc-pgsql-pam', '/etc/config.toml',
  and '$HOME/.config/config.toml', in that specific order.

  In addition to defaults, all configuration parameters can also be provided through environment variables.

  $KC_PGSQL_REALM $KC_PGSQL_ENDPOINT $KC_PGSQL_CLIENTID $KC_PGSQL_CLIENTSECRET $KC_PGSQL_CLIENTSCOPE
  
  To use the program, you must create a client in Keycloak and provide the following 
  information in the configuration file: realm, endpoint, client ID, client secret, and 
  client scope is optional.

Arguments:
  USERNAME                The username of the user is taken from $PAM_USER environment variable
  PASSWORD                The password of the user is taken from stdIn
  OTP                     (Optional) The OTP code if two-factor authentication is enabled i.e (password/otp)

  EXAMPLE                 (With otp): echo testpass/717912 | kc-pgsql-pam (Only Password): echo testpass | kc-pgsql-pam
`, os.Args[0])
}

// displayVersion displays build version information
func DisplayHelp(version, buildDate, commitSha string) {
	helpFlag := flag.Bool("help", false, "Show this help message and exit")
	hFlag := flag.Bool("h", false, "Show this help message and exit")
	versionFlag := flag.Bool("version", false, "Display version information")
	vFlag := flag.Bool("v", false, "Display version number (shorthand)")

	flag.Parse()

	if *helpFlag || *hFlag {
		printHelpMessage()
		os.Exit(0)

	}

	if *versionFlag || *vFlag {
		printVersionInfo(version, buildDate, commitSha)
		os.Exit(0)

	}
}

func printVersionInfo(version, buildDate, commitSha string) {
	fmt.Println("Version:", version)
	fmt.Println("Build Date:", buildDate)
	fmt.Println("Commit SHA:", commitSha)
}
