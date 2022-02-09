package parser

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// VaultAccount struct which contains
// an array of accounts
type VaultAccount struct {
	Accounts []Account `json:"accounts"`
}

// Account struct which contains private key and address
type Account struct {
	PrivateKey string `json:"privateKey"`
	Address    string `json:"address"`
}

// Contract struct which contains chain ID and SC address
type Contract struct {
	ChainIds []struct {
		Name    string `json:"name"`
		ChainId string `json:"chainId"`
		SC      struct {
			NitroAdj struct {
				Address string `json:"address"`
			} `json:"NitroAdjudicator"`
		} `json:"contracts"`
	} `json:"1337"`
}

func FromFileToAccounts(file string) (VaultAccount, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return VaultAccount{}, nil
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var accounts VaultAccount
	json.Unmarshal(byteValue, &accounts)

	return accounts, nil
}

func FromFileToContract(file string) (Contract, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return Contract{}, nil
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result Contract
	json.Unmarshal([]byte(byteValue), &result)

	return result, err
}
