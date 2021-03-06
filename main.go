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

	list = app.Command("list", "List password repositories.")

	new          = app.Command("new", "Create a new password repository.")
	repoName     = new.Arg("name", "Name of password repository.").Required().String()
	setAsDefault = new.Flag("default", "Make this the default password repository.").Short('d').Bool()
	repoKey      = new.Flag("key", "Repository password.").Short('k').String()

	set             = app.Command("set", "Add a password to a repository.")
	setPasswordName = set.Arg("name", "Name of password.").Required().String()
	setPasswordRepo = set.Flag("repo", "Password repository to use.").Short('r').String()
	setPasswordKey  = set.Flag("key", "Repository password.").Short('k').String()
	setPassword     = set.Flag("password", "Actual password.").Short('p').String()

	get             = app.Command("get", "Get a password from a repository.")
	getPasswordName = get.Arg("name", "Name of password.").Required().String()
	getPasswordRepo = get.Flag("repo", "Password repository to use.").Short('r').String()
	getPasswordKey  = get.Flag("key", "Repository password.").Short('k').String()

	delete         = app.Command("delete", "Get a password from a repository.")
	deleteRepoName = delete.Arg("repo", "Password repository to use.").Required().String()
)

func main() {
	kingpin.Version("0.0.1")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	// List existing password repos
	case list.FullCommand():
		names, err := devereux.GetRepositories()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println("Password repositories:")
			for _, name := range names {
				fmt.Println(name)
			}
		}

	// Create new password repo
	case new.FullCommand():
		name, err := devereux.CreateRepository(*repoName, *setAsDefault, *repoKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("Created password repository \"%s\".\n", name)
		}

	// Add password to repo
	case set.FullCommand():
		err := devereux.SetPassword(*setPasswordName, *setPasswordRepo, *setPasswordKey, *setPassword)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("Saved \"%s\".\n", *setPasswordName)
		}

	// Get password from repo
	case get.FullCommand():
		password, err := devereux.GetPassword(*getPasswordName, *getPasswordRepo, *getPasswordKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("Copied \"%s\" to your clipboard.\n", *getPasswordName)
			clipboard.WriteAll(password)
		}

	// Delete a password repo
	case delete.FullCommand():
		err := devereux.DeleteRepository(*deleteRepoName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("Deleted password repository \"%s\".\n", *deleteRepoName)
		}
	}
}
