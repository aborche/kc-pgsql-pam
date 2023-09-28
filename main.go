package main

import (
	"log"
	"os"
	"strings"

	"github.com/aborche/kc-pgsql-pam/internal/auth"
	"github.com/aborche/kc-pgsql-pam/internal/conf"
	"github.com/aborche/kc-pgsql-pam/internal/flags"
	"github.com/aborche/kc-pgsql-pam/internal/utils"
)

var (
	version   string
	buildDate string
	commitSha string
)

func main() {
	// displayVersion()
	flags.DisplayHelp(version, buildDate, commitSha)
	c, err := conf.LoadConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	providerEndpoint := c.Endpoint + "/realms/" + c.Realm
	username := os.Getenv("PAM_USER")

	// Check user domain
	if len(c.AllowedDomains) > 0 {
		index := strings.Index(username, "@")
		if index < 1 {
			log.Fatalf("OIDC Auth: entered username '%s' not contains domain part", username)
			os.Exit(4)
		}
		domain := username[strings.LastIndex(username, "@")+1:]
		if !utils.CheckStringInArray(c.AllowedDomains, domain) {
			log.Fatalf("OIDC Auth: domain '%s' in username '%s' is not allowed here", domain, username)
			os.Exit(4)
		}
	}

	// Analyze the input from stdIn and split the password if it containcts "/"  return otp and pass
	password, otp, err := auth.ReadPasswordWithOTP()
	if err != nil {
		log.Fatal(err)
	}
	// Get provider configuration
	provider, err := auth.GetOIDCProvider(providerEndpoint)
	if err != nil {
		log.Fatalf("OIDC Auth: Failed to retrieve provider configuration: %v\n", err)
	}

	// get token endpint from the provider
	tokenUrl := provider.Endpoint().TokenURL

	// Retrieve an OIDC token using the password grant type
	accessToken, err := auth.RequestJWT(username, password, otp, tokenUrl, c.ClientID, c.ClientSecret, c.ClientScope)
	if err != nil {
		log.Fatalf("OIDC Auth: '%s' Failed to retrieve token: %v\n", username, err)
		os.Exit(2)
	}

	// Verify the token and retrieve the ID token
	if err := auth.VerifyToken(accessToken, c.ClientID, c.ClientSecret, c.Realm, c.Endpoint, c.GroupsClaim, c.AllowedGroups); err != nil {
		// handle the error
		log.Fatal(err)
		os.Exit(3)
	}
	log.Printf("OIDC Auth: '%s' Token acquired and verified Successfully.\n", username)
}
