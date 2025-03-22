package identity

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	// IdentityFileName is the name of the file that stores the agent UUID
	IdentityFileName = ".agent_identity"
)

// GetOrCreateAgentUUID retrieves the agent's UUID from storage or creates a new one
func GetOrCreateAgentUUID() (string, error) {
	// Check if identity file exists
	identityPath := getIdentityFilePath()

	// Try to read existing UUID
	existingUUID, err := readUUIDFromFile(identityPath)
	if err == nil && isValidUUID(existingUUID) {
		// Valid UUID found
		fmt.Printf("Loaded existing agent UUID: %s\n", existingUUID)
		return existingUUID, nil
	}

	// Need to generate a new UUID
	newUUID := uuid.New().String()

	// Save the new UUID
	if err := saveUUIDToFile(identityPath, newUUID); err != nil {
		return "", fmt.Errorf("failed to save new UUID: %w", err)
	}

	fmt.Printf("Generated new agent UUID: %s\n", newUUID)
	return newUUID, nil
}

// getIdentityFilePath returns the path to the identity file
func getIdentityFilePath() string {
	// In a real C2 agent, you might want to use a less obvious location
	// or store in the registry on Windows for better stealth
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fall back to current directory if home dir can't be determined
		return IdentityFileName
	}

	return filepath.Join(homeDir, IdentityFileName)
}

// readUUIDFromFile reads the UUID from a file
func readUUIDFromFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

// saveUUIDToFile saves the UUID to a file
func saveUUIDToFile(path string, uuid string) error {
	return ioutil.WriteFile(path, []byte(uuid), 0600) // Restrictive permissions
}

// isValidUUID checks if a string is a valid UUID
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
