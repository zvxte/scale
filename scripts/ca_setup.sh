#!/bin/bash
set -e

HOSTNAME="localhost"
KEY_SIZE="4096"
KEY_FILE="ca.key"
CERT_FILE="ca.crt"
OPENSSL_CONFIG_FILE="ca_openssl.conf"
DAYS_VALID="365"

cat > $OPENSSL_CONFIG_FILE <<EOL
[ req ]
default_bits = $KEY_SIZE
default_keyfile = $KEY_FILE
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[ req_distinguished_name ]
commonName = $HOSTNAME

[ v3_req ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = $HOSTNAME
EOL

openssl genpkey -algorithm RSA -out $KEY_FILE \
    -pkeyopt rsa_keygen_bits:$KEY_SIZE -quiet

openssl req -x509 -new -key $KEY_FILE -out $CERT_FILE -config \
    $OPENSSL_CONFIG_FILE -extensions v3_req -sha256 -days $DAYS_VALID

rm $OPENSSL_CONFIG_FILE