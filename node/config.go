package node

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
)

const (
	caCertFile = "/etc/scale/ca.crt"
	keyFile    = "/etc/scale/node.key"
	certFile   = "/etc/scale/node.crt"
)

type Config struct {
	TLSConfig *tls.Config
}

func LoadConfig() (*Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load X509 key pair: %w", err)
	}

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New(
			"failed to append CA certificate to the certificate pool",
		)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return &Config{
		TLSConfig: tlsConfig,
	}, nil
}
