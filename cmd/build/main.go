package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func main() {
	// Parse command line arguments for protocol
	protocolFlag := flag.String("protocol", "h1c", "Protocol to build for (h1c, h1tls, h2c, h2tls, h3)")
	flag.Parse()

	// Validate the protocol
	protocol := *protocolFlag
	validProtocols := map[string]bool{
		"h1c":   true,
		"h1tls": true,
		"h2c":   true,
		"h2tls": true,
		"h3":    true,
	}

	if !validProtocols[protocol] {
		fmt.Printf("Error: Invalid protocol '%s'\n", protocol)
		fmt.Println("Valid protocols: h1c, h1tls, h2c, h2tls, h3")
		os.Exit(1)
	}

	// Generate a unique ID for this build
	agentUUID := uuid.New().String()
	fmt.Printf("Building agent with UUID: %s\n", agentUUID)

	// Get current time for build timestamp
	buildTime := time.Now().UTC().Format(time.RFC3339)

	// Ensure bin directory exists
	binDir := "bin"
	if err := os.MkdirAll(binDir, 0755); err != nil {
		fmt.Printf("Error creating bin directory: %v\n", err)
		os.Exit(1)
	}

	// Target binary name based on selected protocol
	binaryName := filepath.Join(binDir, fmt.Sprintf("agent_%s", protocol))

	fmt.Printf("Building %s for protocol: %s\n", binaryName, protocol)

	// Construct the build command with the UUID and build time injected
	cmd := exec.Command("go", "build",
		"-o", binaryName,
		"-ldflags", fmt.Sprintf("-X main.embeddedUUID=%s -X main.buildTime=%s -X main.buildProtocol=%s",
			agentUUID, buildTime, protocol),
		"cmd/agent/main.go")

	// Connect command's stdout and stderr to our process
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the build command
	if err := cmd.Run(); err != nil {
		fmt.Printf("Build failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Build successful! Executable: %s\n", binaryName)
	fmt.Printf("Agent UUID: %s\n", agentUUID)
}
