package main

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/vail130/devereux/devereux"
)

var (
	app = kingpin.New("devereux", "A command-line password manager.")

	new          = app.Command("new", "Create a new password repository.")
	repoName     = new.Arg("name", "Name of password repository.").Required().String()
	repoKey      = new.Flag("key", "Repository password.").Short('k').String()
	setAsDefault = new.Flag("default", "Make this the default password repository.").Short('d').Bool()

	set             = app.Command("set", "Add a password to a repository.")
	setPasswordName = set.Arg("name", "Name of password.").Required().String()
	setPassword     = set.Arg("password", "Actual password.").Required().String()
	setPasswordKey  = set.Flag("key", "Repository password.").Short('k').String()
	setPasswordRepo = set.Flag("repo", "Password repository to use.").Short('r').String()

	get             = app.Command("get", "Get a password from a repository.")
	getPasswordName = get.Arg("name", "Name of password.").Required().String()
	getPasswordKey  = get.Flag("key", "Repository password.").Short('k').String()
	getPasswordRepo = get.Flag("repo", "Password repository to use.").Short('r').String()
)

func main() {
	kingpin.Version("0.0.1")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	// Create new password repo
	case new.FullCommand():
		name, err := devereux.CreateRepository(*repoName, *repoKey, *setAsDefault)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("Created password repository \"%s\".\n", name)
		}

	// Add password to repo
	case set.FullCommand():
		err := devereux.SetPassword(*setPasswordName, *setPassword, *setPasswordKey, *setPasswordRepo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("Saved \"%s\".\n", *setPasswordName)
		}

	// Get password from repo
	case get.FullCommand():
		password, err := devereux.GetPassword(*getPasswordName, *getPasswordKey, *getPasswordRepo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("Copied \"%s\" to your clipboard.\n", *getPasswordName)
			clipboard.WriteAll(password)
		}
	}
}
