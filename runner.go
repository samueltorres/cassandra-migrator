package main

import (
	"os"
	"os/exec"
	"strings"
)

type CqlExecutor struct {
	Hosts      string
	Keyspace   string
	Username   string
	Password   string
	CQLVersion string
	SSLEnabled bool
}

func (e *CqlExecutor) Execute(file string) error {
	cmd := exec.Command("cqlsh", e.buildArgs(file)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (e *CqlExecutor) buildArgs(file string) []string {
	var args []string

	args = append(args, "-f", file)

	if e.Username != "" {
		args = append(args, "-u", e.Username)
	}

	if e.Password != "" {
		args = append(args, "-p", e.Password)
	}

	if e.Keyspace != "" {
		args = append(args, "-k", e.Keyspace)
	}

	if e.CQLVersion != "" {
		args = append(args, "--cqlversion", e.CQLVersion)
	}

	if e.SSLEnabled {
		args = append(args, "--ssl")
	}

	if e.Hosts != "" {
		args = append(args, strings.Split(e.Hosts, ":")...)
	}

	return args
}
