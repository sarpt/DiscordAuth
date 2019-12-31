package discordtoken

import (
	"encoding/json"
	"os"
)

const clientFileName string = "client.json"

// ClientInfo includes information about client credientials
type ClientInfo struct {
	ID     string
	Secret string
}

// NewClientInfo returns instance with filled id and secret. If id and secret not provided, the function tries read it from configuration files.
func NewClientInfo(id string, secret string) (ClientInfo, error) {
	var clientInfo ClientInfo

	if id != "" && secret != "" {
		clientInfo = ClientInfo{
			ID:     id,
			Secret: secret,
		}
		return clientInfo, nil
	}

	clientInfoFromFile, err := getClientInfoFromFile()
	if err != nil {
		return clientInfo, err
	}

	if id == "" {
		clientInfo.ID = clientInfoFromFile.ID
	}

	if secret == "" {
		clientInfo.Secret = clientInfoFromFile.Secret
	}

	return clientInfo, nil
}

func getClientInfoFromFile() (ClientInfo, error) {
	var clientInfo ClientInfo

	clientFilePath, err := getFilePath(clientFileName)
	if err != nil {
		return clientInfo, err
	}

	clientFile, err := os.Open(clientFilePath)
	if err != nil {
		return clientInfo, err
	}

	err = json.NewDecoder(clientFile).Decode(&clientInfo)
	return clientInfo, err
}
