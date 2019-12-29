package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os/exec"
	"time"

	"github.com/sarpt/discordauth/pkg/discordauth"
)

var clientID *string
var clientSecret *string
var timeout *int
var address *string
var route *string

const defaultAddress string = "localhost:8080"
const defaultRoute string = "/oauth2/callback"

var discordScopes []string = []string{"identify", "guilds"}

func init() {
	clientID = flag.String("id", "", "specifies id of client application created in discord")
	clientSecret = flag.String("secret", "", "specifies secret of client application created in discord")
	timeout = flag.Int("timeout", 0, "specifies timeout in seconds of obtaining the authorization token - when no value is provided there is no timeout")
	address = flag.String("address", defaultAddress, "specifies address to which OAuth2 redirect should occur")
	route = flag.String("route", defaultRoute, "specifies route to which OAuth2 redirect should occur")
	flag.Parse()
}

func main() {
	var redirect = url.URL{
		Host:   *address,
		Path:   *route,
		Scheme: "http",
	}

	if *clientID == "" {
		// fetch client id from ./conf.json or ~/.discordauth/conf.json
	}

	if *clientSecret == "" {
		// fetch secret from ./conf.json or ~/.discordauth/conf.json
	}

	config, state := discordauth.GetAuthConfig(*clientID, *clientSecret, discordScopes, redirect)

	xdgCommand := exec.Command("xdg-open", config.AuthCodeURL(state))
	xdgCommand.Run()

	var ctx context.Context
	var cancel context.CancelFunc

	if *timeout == 0 {
		ctx = context.Background()
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
		defer cancel()
	}

	auth, err := discordauth.GetAuthorizedClient(ctx, config, state, redirect)

	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("The token expires: %s", auth.Token.Expiry))
}
