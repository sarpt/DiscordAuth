package oauth

import (
	"encoding/json"
	"io"

	"golang.org/x/oauth2"
)

// GetTokenFromJSON returns Token from JSON content
func GetTokenFromJSON() (oauth2.Token, error) {
	var token oauth2.Token

	return token, nil
}

// WriteTokenToJSON encodes Token as JSON conent
func WriteTokenToJSON(writer io.Writer, token oauth2.Token) error {
	return json.NewEncoder(writer).Encode(token)
}
