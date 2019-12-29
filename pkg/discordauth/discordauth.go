package discordauth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

const authURL string = "https://discordapp.com/api/oauth2/authorize"
const tokenURL string = "https://discordapp.com/api/oauth2/token"

// Notifier provides channels that inform about either code or errors
type Notifier struct {
	code chan string
	err  chan error
}

// Authorization groups authorized Discord client with a Token
type Authorization struct {
	Client *http.Client
	Token  *oauth2.Token
}

// GetAuthConfig returns the url for which to redirect user to insert credientials
func GetAuthConfig(clientID string, clientSecret string, scopes []string, redirect url.URL) (oauth2.Config, string) {
	state := strconv.FormatInt(time.Now().Unix(), 10)

	endpoint := oauth2.Endpoint{
		AuthURL:  authURL,
		TokenURL: tokenURL,
	}

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     endpoint,
		RedirectURL:  redirect.String(),
		Scopes:       scopes,
	}

	return config, state
}

// GetAuthorizedClient will listen for redirect requests, returning the client that is ready to use for doing further request to discord API
// make sure to call this after handling the url (by opening the browser or otherwise)
func GetAuthorizedClient(ctx context.Context, config oauth2.Config, state string, redirect url.URL) (Authorization, error) {
	authorization := Authorization{}
	var err error

	notifier := Notifier{
		code: make(chan string),
		err:  make(chan error),
	}

	go listenForRedirect(ctx, config.RedirectURL, notifier, state, redirect)

	select {
	case code := <-notifier.code:
		if code == "" {
			return authorization, err
		}

		token, err := config.Exchange(ctx, code)
		if err != nil {
			return authorization, err
		}

		authorization.Token = token
		authorization.Client = config.Client(ctx, token)
	case err = <-notifier.err:
	}

	return authorization, err
}

func listenForRedirect(ctx context.Context, callbackURL string, notifier Notifier, state string, redirect url.URL) {
	code := make(chan string)

	handler := http.NewServeMux()
	handler.HandleFunc(redirect.EscapedPath(), getCallbackHandler(code, state))

	server := http.Server{
		Addr:    redirect.Host,
		Handler: handler,
	}

	go func(ctx context.Context, code <-chan string) {
		var receivedCode string

		select {
		case receivedCode = <-code:
		case <-ctx.Done():
			fmt.Println("Server shutdown due to timeout")
		}

		if err := server.Shutdown(ctx); err != nil {
			notifier.err <- err
			return
		}

		notifier.code <- receivedCode
	}(ctx, code)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		notifier.err <- err
		return
	}
}

func getCallbackHandler(code chan<- string, state string) func(http.ResponseWriter, *http.Request) {
	return func(reswr http.ResponseWriter, req *http.Request) {
		queryVals := req.URL.Query()
		codeVal := queryVals.Get("code")
		stateVal := queryVals.Get("state")

		if codeVal == "" {
			fmt.Println("Received empty code - listening further...")
			return
		}

		if stateVal != state {
			fmt.Println("Received state is incorrect with expected - listening further...")
			return
		}

		fmt.Println(fmt.Sprintf("Received code: %s", codeVal))
		code <- codeVal
	}
}
