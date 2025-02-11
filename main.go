package main

import (
	"fmt"
	"os"
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

	for _, line := range strings.Split(string(hostsContent), "\n") {
		if strings.EqualFold(line, localhost+"\t"+domain) {
			fmt.Printf("Domain %s is already blocked\n", domain)
			return
		}
	}

	hostsContent = append(hostsContent, []byte("\n"+localhost+"\t"+domain)...)
	fmt.Printf("Blocking domain: %s\n %s\n", domain, string(hostsContent))
	err = os.WriteFile(hostsFile, hostsContent, 0644)
	if err != nil {
		fmt.Printf("Error writing to %s: %v\n", hostsFile, err)
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
