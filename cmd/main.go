package main

import (
	"flag"
	"os"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "bibirt-sock"
	// Version is the version of the compiled software.
	Version = "0.0.1-alpha"
	// flagconf is the config flag.
	flagconf = "../../configs"

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
}
