# JTW credentials
```shell
openssl ecparam -name prime256v1 -genkey -noout -out privkey.pem
openssl ec -in privkey.pem -pubout -out pubkey.pem
```