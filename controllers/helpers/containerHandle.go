package helpers

import (
	"fmt"
	"os"
	"os/exec"
)

func CreateMySQLContainer(userId string) (string, error) {
	containerName := "mysql_user_" + userId
	password := os.Getenv("MYSQL_CONTAINER_PASSWORD")

	// docker run command
	cmd := exec.Command("docker", "run", "--name", containerName,
		"-e", "MYSQL_ROOT_PASSWORD="+password,
		"-d", "mysql:8")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error creating container: %v, output: %s", err, out)
	}

	containerId := string(out) // ye container ID return karega
	return containerId, nil
}

func DeleteContainer(containerId string) error {
	// stop + remove
	cmd := exec.Command("docker", "rm", "-f", containerId)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error deleting container: %v, output: %s", err, out)
	}
	fmt.Println("Container deleted:", string(out))
	return nil
}

func RunQuery(containerId, query string) (string, error) {
	cmd := exec.Command("docker", "exec", "-i", containerId,
		"mysql", "-uroot", "-psecret", "-e", query)

	out, err := cmd.CombinedOutput()
	return string(out), err
}