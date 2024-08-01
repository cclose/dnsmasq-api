package main

import (
	"github.com/cclose/dnsmasq-api/cmd/dnsMasqAPI/cmd"
)

var (
	BuildTimeStr string
	Commit       string
	Version      string
)

func main() {
	cmd.Execute(BuildTimeStr, Commit, Version)
}
