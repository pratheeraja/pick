package crypto

import (
	"bytes"
	"errors"
	"github.com/golang/crypto/openpgp"
	"github.com/golang/crypto/openpgp/armor"
	"io/ioutil"
)

// GeneratePassword generates a password of length.
func GeneratePassword(length int) (password string, err error) {
	return "", errors.New("crypto.GeneratePassword NOT IMPLEMENTED!")
}

// EncryptText uses PGP to symmetrically encrypt and armor text with the
// provided password.
func EncryptText(text string, password string) (encryptedText string, err error) {
	encbuf := bytes.NewBuffer(nil)

	w, err := armor.Encode(encbuf, "PGP SIGNATURE", nil)
	if err != nil {
		return
	}

	plaintext, err := openpgp.SymmetricallyEncrypt(w, []byte(password), nil, nil)
	if err != nil {
		return
	}

	_, err = plaintext.Write([]byte(text))

	plaintext.Close()
	w.Close()

	encryptedText = encbuf.String()

	return
}

// DecryptText uses PGP to decrypt symmetrically encrypted and armored text
// with the provided password.
func DecryptText(text string, password string) (decryptedText string, err error) {
	decbuf := bytes.NewBuffer([]byte(text))

	armorBlock, err := armor.Decode(decbuf)
	if err != nil {
		return
	}

	failed := false
	prompt := func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		// If the decrypted private key or given passphrase isn't correct,
		// the function will be called again, forever. This method will fail fast.
		// Ref: https://godoc.org/golang.org/x/crypto/openpgp#PromptFunction
		if failed {
			return nil, errors.New("Unable to unlock safe with provided password")
		}

		failed = true

		return []byte(password), nil
	}

	md, err := openpgp.ReadMessage(armorBlock.Body, nil, prompt, nil)

	if err != nil {
		return
	}

	decryptedBuf, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return
	}

	decryptedText = string(decryptedBuf)

	return
}
