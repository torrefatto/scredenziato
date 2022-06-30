package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const progName = "scredenziato"

var Version = "dev"
var subCmds = map[string]func(args ...string) error{
	"list": listCmd,
	//"get":     getCmd,
	//"help":    helpCmd,
	"version": versionCmd,
}

func usage() {
	var verbs []string
	for verb := range subCmds {
		verbs = append(verbs, verb)
	}
	fmt.Printf("%s: <%s> [opts]\n", progName, strings.Join(verbs, "|"))
}

func main() {
	flag.Usage = usage
	flag.Parse()

	subcmd, ok := subCmds[flag.Arg(0)]
	if !ok {
		usage()
		os.Exit(-1)
	}

	if err := subcmd(flag.Args()[1:]...); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func listCmd(args ...string) error {
	var onlyServerURL bool
	flagset := flag.NewFlagSet("list", flag.ExitOnError)
	flagset.BoolVar(&onlyServerURL, "server-url", false, "Only show server urls")
	flagset.Usage = func() {
		fmt.Println("list: lists all the credentials stored if a credential store is found")
		flagset.PrintDefaults()
	}

	if err := flagset.Parse(args); err != nil {
		return err
	}

	helper, err := getHelper()
	if err != nil {
		return err
	}

	list, err := helper.List()
	if err != nil {
		return err
	}

	formatList(list, onlyServerURL)

	return nil
}

func formatList(list map[string]string, onlyServerURL bool) {
	for server, user := range list {
		if onlyServerURL {
			fmt.Println(server)
		} else {
			fmt.Printf("%s\t%s\n", server, user)
		}
	}
}

func versionCmd(...string) error {
	fmt.Printf("%s version: %s\n", progName, Version)
	return nil
}
