package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func installChrony() error {
	_, err := exec.LookPath("apt")
	if err == nil {
		cmd := exec.Command("sudo", "apt", "update")
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to update apt: %v", err)
		}
		cmd = exec.Command("sudo", "apt", "install", "-y", "chrony")
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to install chrony: %v", err)
		}
		return nil
	}

	_, err = exec.LookPath("yum")
	if err == nil {
		cmd := exec.Command("sudo", "yum", "install", "-y", "chrony")
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to install chrony: %v", err)
		}
		return nil
	}

	return fmt.Errorf("unsupported package manager")
}

func startChronyService() error {
	cmd := exec.Command("sudo", "systemctl", "enable", "--now", "chronyd")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start chrony service: %v", err)
	}
	return nil
}

func getChronySyncTime() (float64, error) {
	cmd := exec.Command("chronyc", "tracking")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("error executing chronyc command: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Reference ID") {
			tokens := strings.Fields(line)
			for i, token := range tokens {
				if token == "RefTime" {
					syncTime, err := time.Parse(time.RFC3339, tokens[i+1])
					if err != nil {
						return 0, err
					}
					elapsed := time.Since(syncTime).Minutes()
					return elapsed, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("could not find Chrony sync time")
}

func main() {
	err := installChrony()
	if err != nil {
		fmt.Println("Error installing Chrony:", err)
		return
	}

	err = startChronyService()
	if err != nil {
		fmt.Println("Error starting Chrony service:", err)
		return
	}

	syncTime, err := getChronySyncTime()
	if err != nil {
		fmt.Println("Error getting Chrony sync time:", err)
		return
	}
	fmt.Printf("Chrony sync time: %.2f minutes\n", syncTime)
}
