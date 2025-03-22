package certificates

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
)

// CertificateProvider defines the interface for retrieving TLS certificates
type CertificateProvider interface {
	GetCertificate() (*tls.Certificate, error)
	GetTLSConfig() (*tls.Config, error)
}

// FileCertificateProvider implements CertificateProvider by loading certificates from files
type FileCertificateProvider struct {
	CertFile string
	KeyFile  string
}

// NewFileCertificateProvider creates a new provider that loads certificates from the specified files
func NewFileCertificateProvider(certFile, keyFile string) *FileCertificateProvider {
	return &FileCertificateProvider{
		CertFile: certFile,
		KeyFile:  keyFile,
	}
}

// GetCertificate loads and returns the TLS certificate from the files
func (p *FileCertificateProvider) GetCertificate() (*tls.Certificate, error) {
	// Load the certificate and key
	cert, err := tls.LoadX509KeyPair(p.CertFile, p.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	return &cert, nil
}

// GetTLSConfig returns a TLS configuration using the certificate
func (p *FileCertificateProvider) GetTLSConfig() (*tls.Config, error) {
	cert, err := p.GetCertificate()
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{*cert},
		MinVersion:   tls.VersionTLS12, // Require TLS 1.2 or higher for security
	}, nil
}

// DefaultCertificatePath returns the default path for certificates
func DefaultCertificatePath() string {
	// Return a path relative to the project root for the certificate files
	return "certs"
}

// GetDefaultCertificateProvider returns a provider configured with the default certificate paths
func GetDefaultCertificateProvider() (*FileCertificateProvider, error) {
	certDir := DefaultCertificatePath()

	// Check if the directory exists
	if _, err := os.Stat(certDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("[❌ERR] -> Certificate directory not found: %s", certDir)
	}

	certFile := filepath.Join(certDir, "server.crt")
	keyFile := filepath.Join(certDir, "server.key")

	// Check if the files exist
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("[❌ERR] -> Certificate file not found: %s", certFile)
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("[❌ERR] -> Key file not found: %s", keyFile)
	}

	fmt.Println("[🔐TLS] -> TLS Key and Cert loaded successfully.")
	return NewFileCertificateProvider(certFile, keyFile), nil
}
