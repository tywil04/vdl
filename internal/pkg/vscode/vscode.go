package vscode

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func findExecutable() ([]string, error) {
	out, err := exec.Command("which", "code").Output()
	if err == nil {
		return []string{string(out)}, nil
	}

	fOut, err := exec.Command("flatpak", "list").Output()
	if err != nil {
		return nil, err
	}

	if strings.Contains(string(fOut), "com.visualstudio.code") {
		return []string{"flatpak", "run", "com.visualstudio.code"}, nil
	}

	return nil, errors.New("vscode not found")
}

func OpenContainer(containerId, directory string) error {
	commands, err := findExecutable()
	if err != nil {
		return err
	}

	commands = append(
		commands,
		"--folder-uri",
		fmt.Sprintf(`vscode-remote://attached-container+%s/%s`, containerId, directory),
	)
	cmd := exec.Command(
		commands[0],
		commands[1:]...,
	)
	cmd.Start()

	return nil
}
