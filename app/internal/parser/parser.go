package parser

import (
	"encoding/json"
	"io"
	"os"
)

// VaultAccount struct which contains an array of accounts with private keys and addresses
type VaultAccount struct {
	Accounts []struct {
		PrivateKey string `json:"privateKey"`
		Address    string `json:"address"`
	} `json:"accounts"`
}

// Contract struct which contains chain ID and SC address
type Contract struct {
	Name    string `json:"name"`
	ChainId string `json:"chainId"`
	SC      struct {
		NitroAdj struct {
			Address string `json:"address"`
		} `json:"NitroAdjudicator"`
	} `json:"contracts"`
}

func ToVaultAccount(file string) (VaultAccount, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return VaultAccount{}, nil
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var accounts VaultAccount
	json.Unmarshal(byteValue, &accounts)

	return accounts, nil
}

func ToContract(file string) (map[string][]Contract, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return map[string][]Contract{}, err
	}

	defer jsonFile.Close()

	contract := make(map[string][]Contract)
	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &contract)
	if err != nil {
		return map[string][]Contract{}, err
	}

	return contract, nil
}
