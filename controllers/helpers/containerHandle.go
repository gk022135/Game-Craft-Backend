package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const defaultDB = "gamecraft"

func CreateMySQLContainer(userId string) (string, error) {
	containerName := "mysql_user_" + userId
	password := os.Getenv("MYSQL_CONTAINER_PASSWORD")
	if strings.TrimSpace(password) == "" {
		return "", fmt.Errorf("MYSQL_CONTAINER_PASSWORD is not set")
	}

	// docker run command with default DB to attempt initial creation
	cmd := exec.Command("docker", "run", "--name", containerName,
		"-e", "MYSQL_ROOT_PASSWORD="+password,
		"-e", "MYSQL_DATABASE="+defaultDB,
		"-d", "mysql:8")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error creating container: %v, output: %s", err, strings.TrimSpace(string(out)))
	}

	containerId := strings.TrimSpace(string(out))
	fmt.Println("Container created:", containerId)

	// Wait for mysqld to become ready (use password injection via docker exec -e)
	if err := waitForMySQL(containerId, password, 45*time.Second); err != nil {
		// include docker logs to help debugging
		logsCmd := exec.Command("docker", "logs", containerId)
		logsOut, _ := logsCmd.CombinedOutput()
		return containerId, fmt.Errorf("mysql not ready: %w; docker logs:\n%s", err, string(logsOut))
	}

	// Ensure DB exists (in case this container reused existing volume)
	if _, err := ensureDatabaseExists(containerId, password, defaultDB); err != nil {
		return containerId, fmt.Errorf("failed to ensure database exists: %w", err)
	}

	return containerId, nil
}

func DeleteContainer(containerId string) error {
	cmd := exec.Command("docker", "rm", "-f", containerId)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error deleting container: %v, output: %s", err, strings.TrimSpace(string(out)))
	}
	fmt.Println("Container deleted:", strings.TrimSpace(string(out)))
	return nil
}

// waitForMySQL pings the server inside the container until it responds or timeout elapses.
// NOTE: we pass MYSQL_PWD into the container exec using docker exec -e so mysqladmin can authenticate.
func waitForMySQL(containerId, password string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		// pass -e to docker exec so mysqladmin sees MYSQL_PWD inside container
		cmd := exec.Command("docker", "exec", "-e", "MYSQL_PWD="+password, containerId, "mysqladmin", "-uroot", "ping", "--silent")
		out, err := cmd.CombinedOutput()
		if err == nil {
			_ = out
			return nil
		}
		// optional: print short debug (comment out in prod)
		// fmt.Printf("mysqladmin ping failed: %s\n", strings.TrimSpace(string(out)))
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("mysql did not become ready in %s", timeout)
}

// ensureDatabaseExists runs CREATE DATABASE IF NOT EXISTS <db>;
// it runs without -D to avoid Unknown database when the DB doesn't yet exist.
func ensureDatabaseExists(containerId, password, dbName string) (string, error) {
	create := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", dbName)
	cmd := exec.Command("docker", "exec", "-i", "-e", "MYSQL_PWD="+password, containerId,
		"mysql", "-uroot", "-e", create)
	out, err := cmd.CombinedOutput()
	outStr := string(out)
	if err != nil {
		return outStr, fmt.Errorf("create database failed: %v, output: %s", err, strings.TrimSpace(outStr))
	}
	return outStr, nil
}

// RunQuery executes query in the given container against defaultDB (gamecraft).
// If the DB doesn't exist, ensureDatabaseExists will create it first.
func RunQuery(containerId, query string) (string, error) {
	password := os.Getenv("MYSQL_CONTAINER_PASSWORD")
	if strings.TrimSpace(password) == "" {
		return "", fmt.Errorf("MYSQL_CONTAINER_PASSWORD is not set")
	}

	// Wait for MySQL to be ready.
	if err := waitForMySQL(containerId, password, 30*time.Second); err != nil {
		return "", fmt.Errorf("mysql not ready: %w", err)
	}

	// Make sure defaultDB exists (safe no-op if already present).
	if _, err := ensureDatabaseExists(containerId, password, defaultDB); err != nil {
		return "", err
	}

	// Now run the provided query against defaultDB using -D
	cmd := exec.Command("docker", "exec", "-i", "-e", "MYSQL_PWD="+password, containerId,
		"mysql", "-uroot", "-D", defaultDB, "-e", query)

	out, err := cmd.CombinedOutput()
	outStr := string(out)
	fmt.Println(outStr)
	if err != nil {
		return outStr, fmt.Errorf("query failed: %v, output: %s", err, strings.TrimSpace(outStr))
	}
	return outStr, nil
}
