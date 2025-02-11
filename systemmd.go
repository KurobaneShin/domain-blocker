package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func createSystemdService() error {
	myPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	serviceFilePath := "/etc/systemd/system/block_domains.service"
	serviceContent := `[Unit]
Description=Block domains during specific time range

[Service]
ExecStart=/path/to/your/block_domains
Restart=always
User=root
Group=root

[Install]
WantedBy=multi-user.target
`

	compiledPath := myPath + "/domain-blocker"
	serviceContent = strings.Replace(serviceContent, "/path/to/your/block_domains", compiledPath, -1)

	err = os.WriteFile(serviceFilePath, []byte(serviceContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write systemd service file: %v", err)
	}

	cmd := exec.Command("sudo", "systemctl", "daemon-reload")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to reload systemd daemon: %v", err)
	}

	cmd = exec.Command("sudo", "systemctl", "enable", "block_domains.service")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to enable systemd service: %v", err)
	}

	cmd = exec.Command("sudo", "systemctl", "start", "block_domains.service")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start systemd service: %v", err)
	}

	fmt.Println("Systemd service created and enabled successfully.")
	return nil
}
