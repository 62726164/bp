#!/bin/bash

# Generate a self-signed wildcard TLS cert. cd here and run script.
# use *.iso.local for Docker in ECS
# use *.bloomingpassword.fun for Linode

openssl ecparam -out ec-secp384r1.pem -name secp384r1

openssl req -nodes -newkey ec:ec-secp384r1.pem -keyout privkey.pem -new -out csr.pem

openssl x509 -sha256 -req -days 3650 -in csr.pem -signkey privkey.pem -out cert.pem 

cp cert.pem privkey.pem ../
rm *.pem
