package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type CommandList []string

func (cl CommandList) contains(e string) bool {
	for _, a := range cl {
		if a == e {
			return true
		}
	}
	return false
}

var grantedCommands = CommandList{
	"git-receive-pack",
	"git-upload-pack",
}

func runCommand(cmd string) {
	home := os.Getenv("PWD")
	//cmd := os.Getenv("SSH_ORIGINAL_COMMAND")
	s := strings.Split(cmd, " ")
	switch {
	case grantedCommands.contains(s[0]):
		var ss []string
		for _, str := range s {
			ss = append(ss, strings.Trim(str, "'"))
		}

		subproc := exec.Command(s[0], ss[1:]...)
		subproc.Dir = home
		subproc.Stderr = os.Stderr
		subproc.Stdout = os.Stdout
		subproc.Stdin = os.Stdin

		if err := subproc.Run(); err != nil {
			e := err.(*exec.ExitError)
			os.Exit(e.ExitCode())
		} else {
			os.Exit(0)
		}

	case cmd == "":
		os.Exit(1)

	default:
		log.Println("Access denied")
		os.Exit(1)
	}
}

func main() {
	runCommand(os.Getenv("SSH_ORIGINAL_COMMAND"))
}
