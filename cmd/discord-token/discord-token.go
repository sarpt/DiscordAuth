package main

import (
	"flag"
	"fmt"

	"github.com/sarpt/discord-token/pkg/discordtoken"
)

var clientID *string
var clientSecret *string
var timeout *int
var address *string
var route *string
var forceRegenerate *bool
var tokenPath *string

func init() {
	clientID = flag.String("id", "", "Id of client application created in discord. When not provided, XDG_CONFIG_HOME with fallback to $HOME/.config/discordauth will be checked for client.json.")
	clientSecret = flag.String("secret", "", "Secret of client application created in discord. When not provided, XDG_CONFIG_HOME with fallback to $HOME/.config/discordauth will be checked for client.json.")
	timeout = flag.Int("timeout", 0, "Allowed time in seconds to obtain the authorization token. Prevents spawned server from waiting forever to OAuth2 callback redirection if user does not take action. When no value is provided server (and the command) waits indifinetely for redirection.")
	address = flag.String("address", discordtoken.DefaultAddress, "Address to which OAuth2 redirect should occur.")
	route = flag.String("route", discordtoken.DefaultRoute, "Route to which OAuth2 redirect should occur.")
	tokenPath = flag.String("token-path", "", "Path to a token file. When not provided, XDG_CONFIG_HOME with fallback to $HOME/.config/discordauth will be checked for token.json during reading, and XDG_CONFIG_HOME with fallback to $HOME/.config/discordauth will be used for writing.")

	flag.Parse()
}

func main() {
	client, err := discordtoken.NewClientInfo(*clientID, *clientSecret)
	if err != nil {
		panic(err)
	}

	config := discordtoken.Config{
		Client:   client,
		Redirect: discordtoken.GetRedirect(*address, *route),
		Scopes:   discordtoken.DefaultScopes,
	}

	ctx, cancel := discordtoken.GetContext(*timeout)
	if cancel != nil {
		defer cancel()
	}

	token, err := discordtoken.GenerateToken(ctx, config)
	if err != nil {
		panic(err)
	}

	err = discordtoken.WriteTokenFile(*tokenPath, *token)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("The token expires: %s", token.Expiry))
}
