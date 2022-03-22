openssl pkcs12 -export -in client.pem -inkey client.key -out client.p12
base64 client.p12 > client.txt
