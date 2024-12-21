package node

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
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
	if _, err := os.Stat(caCertFile); err != nil {
		return nil, err
	}
	if _, err := os.Stat(keyFile); err != nil {
		return nil, err
	}
	if _, err := os.Stat(certFile); err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("failed to append CA certificate to pool")
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
