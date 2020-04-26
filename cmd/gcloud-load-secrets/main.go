package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	gcloudSecrets "github.com/in4it/gcloud-load-secrets/pkg/gcloud/secrets"
)

func main() {

	var (
		secretsPrefix string
		secretsLabel  string
		cmd           string
		debug         bool
	)
	flag.StringVar(&secretsPrefix, "prefix", "", "prefix to filter on when retrieving secrets")
	flag.StringVar(&secretsLabel, "label", "", "label to filter on when retrieving secrets")
	flag.StringVar(&cmd, "cmd", "", "execute command")
	flag.BoolVar(&debug, "debug", false, "enable debug output")

	flag.Parse()

	if cmd == "" {
		fmt.Printf("missing required -cmd argument/flag\n")
		os.Exit(2)
	}

	readSecrets, err := gcloudSecrets.NewReadSecrets()

	if err != nil {
		panic(err)
	}

	secrets, err := readSecrets.ListSecrets(secretsPrefix, secretsLabel)
	if err != nil {
		panic(err)
	}
	secretsWithVersion, err := readSecrets.GetSecretsValue(secrets)
	if err != nil {
		panic(err)
	}
	if debug {
		for _, v := range secretsWithVersion {
			fmt.Printf("[debug] secrets: %s %s\n", v.Name, v.Payload)
		}
	}
	execCommand(cmd, readSecrets.GetKV(secrets))
}

func execCommand(input string, secrets []string) {
	var (
		args    []string
		env     []string
		command string
	)
	args = strings.Split(input, " ")
	command = args[0]
	env = append(os.Environ(), secrets...)

	err := syscall.Exec(command, args, env)
	if err != nil {
		log.Println(err)
	}
}
