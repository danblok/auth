#!/bin/bash

# Generates a self-signed CA certificate and a private key.
openssl req -new \
	-newkey rsa:4096 \
	-keyout ca.key \
	-subj /C=US/ST=NY/L=NY/O=auth/CN=localhost \
	-x509 \
	-sha256 \
	-days 365 \
	-out ca.crt

# Generates a private key for server.
openssl genrsa -out server.key 4096

# Generates a server CSR.
openssl req -new \
	-key server.key \
	-subj /C=US/ST=NY/L=NY/O=auth/CN=localhost \
	-config openssl.cnf \
	-out server.csr

# Generates a server certificate signed by the CA and a private key.
openssl x509 \
	-req \
	-in server.csr \
	-extfile openssl.cnf \
	-extensions server_ext \
	-CA ca.crt \
	-CAkey ca.key \
	-CAcreateserial \
	-days 365 \
	-sha256 \
	-out server.crt

# Verifies validity of certs.
openssl verify -verbose -CAfile ca.crt server.crt
