# JTW credentials
```shell
# EdDSA private key
openssl genpkey -algorithm Ed25519 -out privkey.pem
# ECDSA private key
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:prime256v1 -out privkey.pem
# public key from private key
openssl pkey -in privkey.pem -pubout -out pubkey.pem
```
