package main

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var cmdKill = &Command{
	Exec:        runKill,
	UsageLine:   "kill [OPTIONS] SERVER",
	Description: "Kill a running server",
	Help:        "Kill a running server.",
}

func init() {
	cmdKill.Flag.BoolVar(&killHelp, []string{"h", "-help"}, false, "Print usage")
	// FIXME: add --signal option
}

// Flags
var killHelp bool // -h, --help flag

func runKill(cmd *Command, args []string) {
	if killHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	serverId := cmd.GetServer(args[0])
	command := "halt"
	server, err := cmd.API.GetServer(serverId)
	if err != nil {
		log.Fatalf("Failed to get server information for %s: %v", serverId, err)
	}

	execCmd := append(NewSshExecCmd(server.PublicAddress.IP, true, []string{command}))

	log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))

	spawn := exec.Command("ssh", execCmd...)
	spawn.Stdout = os.Stdout
	spawn.Stdin = os.Stdin
	spawn.Stderr = os.Stderr
	err = spawn.Run()
	if err != nil {
		log.Fatal(err)
	}
}