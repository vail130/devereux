package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/atotto/clipboard"

	"github.com/vail130/devereux/devereux"
)

var (
	app = kingpin.New("devereux", "A command-line password manager.")

	new = app.Command("new", "Create a new password repository.")
	repoName = new.Arg("name", "Name of password repository.").Required().String()
	setAsDefault = new.Flag("default", "Make this the default password repository.").Short('d').Bool()

	set = app.Command("set", "Add a password to a repository.")
	setPasswordName = set.Arg("name", "Name of password.").Required().String()
	setPassword = set.Arg("password", "Actual password.").Required().String()
	setPasswordRepo = set.Flag("repo", "Password repository to use.").Short('r').String()

	get = app.Command("get", "Get a password from a repository.")
	getPasswordName = get.Arg("name", "Name of password.").Required().String()
	getPasswordRepo = get.Flag("repo", "Password repository to use.").Short('r').String()
)

func main() {
	kingpin.Version("0.0.1")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	// Create new password repo
	case new.FullCommand():
		_, err := devereux.CreateRepository(*repoName, *setAsDefault)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	// Add password to repo
	case set.FullCommand():
		err := devereux.SetPassword(*setPasswordName, *setPassword, *setPasswordRepo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("Saved %s.\n", *setPasswordName)
		}

	// Get password from repo
	case get.FullCommand():
		password, err := devereux.GetPassword(*getPasswordName, *getPasswordRepo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("Copied %s to your clipboard.\n", *getPasswordName)
			clipboard.WriteAll(password)
		}
	}
}
