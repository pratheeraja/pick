package main

import (
	"github.com/docopt/docopt-go"
)

const (
	Version string = "pick version 0.1.0"
)

func main() {
	usage := `pick - minimal password manager
	
Usage:
    pick add 		[<alias>] [<username>] [<password>]
    pick cat 		<alias>
    pick cp 		<alias>
    pick rm 		<alias>
    pick ls
    pick -h | --help
    pick -v | --version

Options:
    -h, --help
    -v, --version

The most commonly used pick commands are:
    add 				Save a new credential to the safe
    cat					Print a credential to STDOUT
    cp					Copy a credential password to the clipboard
    ls					List credentials
    rm					Remove a credential
`

	args, _ := docopt.Parse(usage, nil, true, Version, true)

	dispatch(args)
}

func dispatch(args map[string]interface{}) {
	if args["cp"].(bool) {
		alias := args["<alias>"].(string)

		CopyCredential(alias)
	}

	if args["rm"].(bool) {
		alias := args["<alias>"].(string)

		DeleteCredential(alias)
	}

	if args["ls"].(bool) {

		ListCredentials()
	}

	if args["cat"].(bool) {
		alias := args["<alias>"].(string)

		ReadCredential(alias)
	}

	if args["add"].(bool) {
		alias, ok := args["<alias>"].(string)
		if !ok {
			alias = ""
		}

		username, ok := args["<username>"].(string)
		if !ok {
			username = ""
		}

		password, ok := args["<password>"].(string)
		if !ok {
			password = ""
		}

		WriteCredential(alias, username, password)
	}
}
