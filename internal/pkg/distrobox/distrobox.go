package distrobox

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func GetContainerId(containerName string) string {
	rawData := []byte{}
	for _, char := range fmt.Sprintf(`{"containerName":"%s"}`, containerName) {
		rawData = append(rawData, byte(rune(char)))
	}
	return hex.EncodeToString(rawData)
}

func GetContainers() ([]string, error) {
	out, err := exec.Command("distrobox", "ls", "--no-color").Output()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`\S{12} \| (\S+) .* \| .* \|`)
	groups := re.FindAllStringSubmatch(string(out), -1)

	containers := make([]string, len(groups))
	for i, container := range groups {
		containers[i] = container[1]
	}

	return containers, nil
}

func StartContainer(containerName string) error {
	out, err := exec.Command("distrobox", "enter", containerName, "-e", `"echo STARTED"`).Output()
	if err != nil {
		return err
	}

	if !strings.Contains(string(out), "STARTED") {
		return errors.New("container not started")
	}

	return nil
}
