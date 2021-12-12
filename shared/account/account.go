package account

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/algorand/go-algorand-sdk/crypto"

	net "github.com/koen-vr/algo-dice-roll/shared/network"
)

// Note: Do not use in production as it is insecure! Why?
// Storing private keys in plain text to disk will cost you.
// This is just designed to ease development, consult experts.

func Load(name string) (crypto.Account, error) {
	account, err := loadFromFile(name)
	if nil != err {
		return crypto.Account{}, fmt.Errorf(
			"account: load from file: %s", err,
		)
	}
	return account, nil
}

func Create(name string) (crypto.Account, error) {
	account := crypto.GenerateAccount()
	if err := saveToFile(name, account); nil != err {
		return crypto.Account{}, fmt.Errorf(
			"account: save to file: %s", err,
		)
	}
	return account, nil
}

func HasAccount(name string) bool {
	if _, err := os.Stat(getPath(name)); err != nil {
		return false
	}
	return true
}

func getPath(name string) string {
	return fmt.Sprintf("%s/%s.key", net.Path(), name)
}

func saveToFile(name string, account crypto.Account) error {
	if err := ioutil.WriteFile(getPath(name),
		account.PrivateKey, os.ModePerm,
	); err != nil {
		return fmt.Errorf(
			"write file: %s", err,
		)
	}
	return nil
}

func loadFromFile(name string) (crypto.Account, error) {
	data, err := os.ReadFile(getPath(name))
	if err != nil {
		return crypto.Account{}, fmt.Errorf(
			"read file: %s", err,
		)
	}
	account, err := crypto.AccountFromPrivateKey(data)
	if err != nil {
		return crypto.Account{}, fmt.Errorf(
			"from private key: %s", err,
		)
	}
	return account, nil
}
