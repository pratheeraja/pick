package safe

import (
	"cmd/crypto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Safe struct {
	CreatedOn int64                 `json:"createdOn"`
	CreatedBy string                `json:"createdBy"`
	Data      map[string]Credential `json:"data"`
}

type Credential struct {
	Alias     string `json:"alias"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedOn int64  `json:"createdOn"`
}

// Save encrypts the safe with the provided password and writes
// it to disk.
func (safe *Safe) Save(safePath string, password string) (err error) {
	encryptedSafe, err := crypto.EncryptText(safe.toJson(), password)

	err = ioutil.WriteFile(safePath, []byte(encryptedSafe), 0600)
	if err != nil {
		return
	}

	return
}

// Load loads the encrypted Safe file at safePath, decrypts the file, and
// returns the Safe.
func Load(safePath string, password string) (safe *Safe, err error) {
	encryptedSafe, err := ioutil.ReadFile(safePath)
	if err != nil {
		return
	}

	decryptedSafe, err := crypto.DecryptText(string(encryptedSafe), password)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal([]byte(decryptedSafe), &safe)
	if err != nil {
		log.Fatal(err)
	}

	if safe.Data == nil {
		safe.Data = make(map[string]Credential)
	}

	return
}

// AddCredential adds the provided credential to the safe.
func (safe *Safe) AddCredential(cred Credential) (err error) {
	if _, exists := safe.Data[cred.Alias]; exists {
		return fmt.Errorf("Credential with alias '%s' already exists!", cred.Alias)
	}

	safe.Data[cred.Alias] = cred

	return
}

// GetCredential returns the credential with the provided alias.
func (safe *Safe) GetCredential(alias string) (cred Credential, err error) {
	if _, ok := safe.Data[alias]; !ok {
		err = fmt.Errorf("Credential with alias '%s' does not exist!", alias)
		return Credential{}, err
	}

	return safe.Data[alias], nil
}

// RemoveCredential deletes the credential with the provided alias.
func (safe *Safe) RemoveCredential(alias string) (err error) {
	if _, ok := safe.Data[alias]; !ok {
		err = fmt.Errorf("Credential with alias '%s' does not exist!", alias)
		return err
	}

	delete(safe.Data, alias)
	return nil
}

// toJson marshals the safe to JSON.
func (safe *Safe) toJson() string {
	j, _ := json.Marshal(safe)

	return string(j)
}
