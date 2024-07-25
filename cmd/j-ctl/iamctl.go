package main

import (
	"j-iam/internal/iamctl/cmd"
	"os"
)

func main() {
	command := cmd.NewDefaultIAMCtlCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
