package gitlib

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type User struct {
	username string
}

func GetGitUsers() ([]User, error) {
	//runCommand("cat /etc/group | grep gittk-users")

	groups, err := os.ReadFile("./group")
	if err != nil {
		return nil, err
	}

	return parseGroupFile(groups, "gittk-users"), nil
}

func runCommand(c string) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("\"%s\"", c))

	if err := cmd.Run(); err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			fmt.Println(err.Error())

		default:
			log.Panic(err)
		}
	}

	fmt.Println(cmd.Stdout)
}

func parseGroupFile(content []byte, groupname string) []User {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Split(bufio.ScanLines)

	var users []User

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ":")
		if line[0] == groupname {
			for _, u := range strings.Split(line[3], ",") {
				users = append(users, User{
					username: u,
				})
			}
			break
		}
	}

	return users
}

func getUserDataForPubkey(username string) {

}
