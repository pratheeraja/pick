package main

import (
	"bufio"
	"cmd/crypto"
	"cmd/safe"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"os/user"
	"strings"
	"time"
)

var masterPassword string

// CopyCredential copies a credential's password to the clipboard.
func CopyCredential(alias string) {
	cred, err := loadCredential(alias)
	if err != nil {
		log.Fatal(err)
	}

	err = clipboard.WriteAll(cred.Password)
	if err != nil {
		log.Fatal(err)
	}
}

// DeleteCredential deletes a credential from the safe.
func DeleteCredential(alias string) {
	safePath := getSafePath()

	// 1. Load the safe
	_safe, err := loadSafe(safePath)
	if err != nil {
		log.Fatal(err)
	}

	err = _safe.RemoveCredential(alias)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Save the safe
	safePassword := getMasterPassword(
		"Enter a master password to lock your safe")

	err = _safe.Save(safePath, safePassword)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("done")
}

// ListCredentials displays all of the credentials in the safe.
func ListCredentials() {
	_safe, err := loadSafe(getSafePath())
	if err != nil {
		log.Fatal(err)
	}

	for k := range _safe.Data {
		fmt.Println(k)
	}

	//prettyPrint(_safe.Data)
}

// ReadCredential displays a single credential from the safe.
func ReadCredential(alias string) {
	cred, err := loadCredential(alias)
	if err != nil {
		log.Fatal(err)
	}

	prettyPrint(cred)
}

// WriteCredential writes a new credential to the safe.
func WriteCredential(alias string, username string, password string) {

	// 1. Get a safe
	_safe, err := loadOrCreateSafe()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Collect any info that was not provided
	if alias == "" {
		alias = getInput("Enter an alias")
	}

	_, ok := _safe.Data[alias]
	if ok {
		log.Fatal(errors.New("Credential for " + alias + " already exists"))
	}

	if username == "" {
		username = getInput("Enter a username for " + alias)
	}

	if password == "" {
		if getAnswer("Generate password", "y") {
			password, err = crypto.GeneratePassword(50)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			password = getPassword("Enter your password for " + alias)
		}
	}

	// 3. Add the new credential to the safe
	cred := safe.Credential{alias, username, password, time.Now().Unix()}
	err = _safe.AddCredential(cred)
	if err != nil {
		log.Fatal(err)
	}

	// 4. Save the safe
	safePassword := getMasterPassword(
		"Enter a master password to lock your safe")

	err = _safe.Save(getSafePath(), safePassword)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("done")
}

func getPassword(prompt string) string {
	fmt.Printf("%s\n> ", prompt)
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n")

	return string(password)
}

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s\n> ", prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return string(text[:len(text)-1])
}

func getAnswer(question string, defaultChoice string) bool {
	prompt := question + "? (y/n)"

	yes := getInput(prompt)
	if yes == "" {
		yes = defaultChoice
	}

	return strings.Contains(yes, "y") || strings.Contains(yes, "yes")
}

func getMasterPassword(prompt string) string {
	if masterPassword == "" {
		masterPassword = getPassword(prompt)
	}

	return masterPassword
}

func getSafePath() string {
	// Use default safe location
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// ~/.pick.safe
	return usr.HomeDir + "/.pick.safe"
}

func prettyPrint(v interface{}) {
	j, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf(string(j))
}

func loadOrCreateSafe() (s *safe.Safe, err error) {
	safePath := getSafePath()
	if safeExists(safePath) {
		s, err = loadSafe(safePath)

	} else {
		if getAnswer("Unable to find safe, create new", "y") {
			s, err = createNewSafe(safePath)
		} else {
			// They chose not to create a new safe
			err = errors.New("You must create or provide a safe")
			s = &safe.Safe{}
		}
	}

	return
}

func safeExists(safePath string) bool {
	if _, err := os.Stat(safePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func loadSafe(safePath string) (*safe.Safe, error) {
	if !safeExists(safePath) {
		return nil, errors.New("Safe does not exist")
	}

	safePassword := getMasterPassword(
		"Enter a master password to unlock your safe")

	return safe.Load(safePath, safePassword)
}

func createNewSafe(safePath string) (s *safe.Safe, err error) {
	if safeExists(safePath) {
		err = errors.New("Safe already exists")
		s = &safe.Safe{}
		return
	}

	usr, _ := user.Current()
	s = &safe.Safe{
		time.Now().Unix(),
		usr.Name,
		make(map[string]safe.Credential),
	}

	return
}

func loadCredential(alias string) (safe.Credential, error) {
	_safe, err := loadOrCreateSafe()
	if err != nil {
		log.Fatal(err)
	}

	return _safe.GetCredential(alias)
}
