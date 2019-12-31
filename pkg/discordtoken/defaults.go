package discordtoken

// DefaultAddress is address of a server listening to redirect from browser
const DefaultAddress string = "localhost:8080"

// DefaultRoute is a server route which should handle redirect from browser
const DefaultRoute string = "/oauth2/callback"

// DefaultScopes specify permissions that client using the token will have when accessing Discord API
var DefaultScopes []string = []string{"identify"}
