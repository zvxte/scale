package node

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
)

const (
	DefaultCACertFile = "/etc/scale/ca.crt"
	DefaultKeyFile    = "/etc/scale/node.key"
	DefaultCertFile   = "/etc/scale/node.crt"
)

// LoadMTLSConfig loads the TLS configuration for mutual authentication.
// Default file paths will be used for empty string arguments.
func LoadMTLSConfig(caCertFile, keyFile, certFile string) (*tls.Config, error) {
	if keyFile == "" {
		keyFile = DefaultKeyFile
	}
	if certFile == "" {
		certFile = DefaultCertFile
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load X509 key pair: %w", err)
	}

	if caCertFile == "" {
		caCertFile = DefaultCACertFile
	}
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("failed to append CA certificate to certificate pool")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}, nil
}
