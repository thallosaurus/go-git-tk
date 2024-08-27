package gitlib

import "os"

type sshkey_entry struct {
}

func OpenAuthorizedKeys(path string) {
	if path == "" {
		home := os.Getenv("HOME")
		path = home + "/.ssh/authorized_keys"
	}
}
