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
	"list":    listCmd,
	"get":     getCmd,
	"version": versionCmd,
}

func init() {
	subCmds["help"] = helpCmd
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

func getCmd(args ...string) error {
	var noUsername, noSecret bool
	flagset := flag.NewFlagSet("get", flag.ExitOnError)
	flagset.BoolVar(&noUsername, "no-username", false, "Omit the username from the return value")
	flagset.BoolVar(&noSecret, "no-secret", false, "Omit the secret from the return value")
	flagset.Usage = func() {
		fmt.Println("get <server-url>")
		fmt.Println("Given a server url, returns the associated username and/or password")
		flagset.PrintDefaults()
	}

	if err := flagset.Parse(args); err != nil {
		return err
	}

	if flagset.NArg() != 1 {
		flagset.Usage()
		return fmt.Errorf("get: wrong syntax")
	}

	helper, err := getHelper()
	if err != nil {
		return err
	}

	user, secret, err := helper.Get(flagset.Arg(0))
	if err != nil {
		return err
	}

	if !noUsername {
		fmt.Println(user)
	}

	if !noSecret {
		fmt.Println(secret)
	}

	return nil
}

func versionCmd(...string) error {
	fmt.Printf("%s version: %s\n", progName, Version)
	return nil
}

func helpCmd(arg ...string) error {
	if len(arg) > 0 {
		switch arg[0] {
		case "list":
			listCmd("-h")
		case "get":
			getCmd("-h")
		}
	}
	usage()
	return nil
}
