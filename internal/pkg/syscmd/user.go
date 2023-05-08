package syscmd

import (
	"os/user"
)

func FindUser(username string) (string, error) {
	user, err := user.Lookup(username)
	if err != nil {
		return "", err
	}

	return user.Username, nil
}
