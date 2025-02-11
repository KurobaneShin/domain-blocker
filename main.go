package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// DomainTime struct holds the domain name and its specific block time range
type DomainTime struct {
	Domain    string
	BlockFrom int // Block start hour (0-23)
	BlockTo   int // Block end hour (0-23)
}

var domainTimes = []DomainTime{
	{"x.com", 10, 18},
	{"asuracomic.net", 10, 18},
	{"reaperscans.com", 10, 18},
	{"flamecomics.xyz", 10, 18},
	{"mangakakalot.com", 10, 18},
	{"mangadex.org", 10, 18},
}

const (
	hostsFile = "/etc/hosts"
	localhost = "127.0.0.1"
)

func main() {
	err := createSystemdService()

	if err != nil {
		fmt.Printf("Error creating systemd service: %v\n", err)
		return
	}

	for {
		currentTime := time.Now()
		hour := currentTime.Hour()

		for _, domainTime := range domainTimes {
			if shouldBlock(domainTime, hour) {
				blockDomain(domainTime.Domain)
			} else {
				fmt.Printf("Unblocking domain: %s\n", domainTime.Domain)
				unblockDomain(domainTime.Domain)
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

// createSystemdService generates a systemd service file and enables the service
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

	// Replace placeholder with the correct path to your compiled Go binary
	compiledPath := myPath + "/block_domains" // Modify this with the actual path
	serviceContent = strings.Replace(serviceContent, "/path/to/your/block_domains", compiledPath, -1)

	// Write the service content to the file
	err = os.WriteFile(serviceFilePath, []byte(serviceContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write systemd service file: %v", err)
	}

	// Reload systemd daemon to recognize the new service
	cmd := exec.Command("sudo", "systemctl", "daemon-reload")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to reload systemd daemon: %v", err)
	}

	// Enable the service to start on boot
	cmd = exec.Command("sudo", "systemctl", "enable", "block_domains.service")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to enable systemd service: %v", err)
	}

	// Start the service
	cmd = exec.Command("sudo", "systemctl", "start", "block_domains.service")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start systemd service: %v", err)
	}

	fmt.Println("Systemd service created and enabled successfully.")
	return nil
}

// shouldBlock checks if the current hour falls within the block time for the domain
func shouldBlock(domainTime DomainTime, currentHour int) bool {
	shouldBlock := currentHour >= domainTime.BlockFrom && currentHour < domainTime.BlockTo
	return shouldBlock
}

// blockDomain adds the domain to /etc/hosts to block it
func blockDomain(domain string) {
	hostsContent, err := os.ReadFile(hostsFile)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", hostsFile, err)
		return
	}

	if strings.Contains(string(hostsContent), domain) {
		fmt.Printf("Domain %s is already blocked\n", domain)
		return
	}

	if !strings.EqualFold(string(hostsContent), domain) {
		hostsContent = append(hostsContent, []byte("\n"+localhost+"\t"+domain)...)
		fmt.Printf("Blocking domain: %s\n %s\n", domain, string(hostsContent))
		err = os.WriteFile(hostsFile, hostsContent, 0644)
		if err != nil {
			fmt.Printf("Error writing to %s: %v\n", hostsFile, err)
		}
	}
}

func unblockDomain(domain string) {
	hostsContent, err := os.ReadFile(hostsFile)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", hostsFile, err)
		return
	}

	hostsContent = removeDomain(hostsContent, domain)

	err = os.WriteFile(hostsFile, hostsContent, 0644)
	if err != nil {
		fmt.Printf("Error writing to %s: %v\n", hostsFile, err)
	}
}

func removeDomain(hostsContent []byte, domain string) []byte {
	contentStr := string(hostsContent)

	lines := strings.Split(contentStr, "\n")
	var newLines []string
	for _, line := range lines {
		if !strings.EqualFold(line, localhost+"\t"+domain) {
			newLines = append(newLines, line)
		}
	}

	return []byte(strings.Join(newLines, "\n"))
}
