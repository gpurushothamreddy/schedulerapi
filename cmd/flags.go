package main

import (
	"fmt"
	"os"

	goflags "github.com/jessevdk/go-flags"
)

func parse() CommandlineArgs {
	args := CommandlineArgs{}

	parser := goflags.NewParser(&args, goflags.Default)
	parser.NamespaceDelimiter = "-"
	_, err := parser.Parse()
	if err != nil {
		fmt.Println("Could not parse args", err.Error())
		os.Exit(1)
	}

	return args
}

type CommandlineArgs struct {
	Server ServerArgs `group:"Server" namespace:"server"`
}

type ServerArgs struct {
	HostName string `long:"hostname" required:"true" description:"Name of the server/instance hosting this specific execution of the service binary"`
	Port     int    `long:"port" required:"true" description:"Port for the server to listen on for secure requests"`
}
