#!/bin/bash
set -e

HOSTNAME="localhost"
KEY_SIZE="4096"
DAYS_VALID="365"

KEY_FILE="node.key"
CSR_FILE="node.csr"
CERT_FILE="node.crt"
OPENSSL_CONFIG_FILE="node_openssl.conf"

CA_KEY_FILE="ca.key"
CA_CERT_FILE="ca.crt"

cat > $OPENSSL_CONFIG_FILE <<EOL
[ req ]
default_bits = $KEY_SIZE
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

openssl req -new -key $KEY_FILE -out $CSR_FILE -config $OPENSSL_CONFIG_FILE \
     -noenc -quiet

openssl x509 -req -in $CSR_FILE -CA $CA_CERT_FILE -CAkey $CA_KEY_FILE \
    -CAcreateserial -out $CERT_FILE -extensions v3_req \
    -extfile $OPENSSL_CONFIG_FILE -days $DAYS_VALID -sha256

rm $CSR_FILE
rm $OPENSSL_CONFIG_FILE
