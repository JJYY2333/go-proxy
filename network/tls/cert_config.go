/*
@Time    : 3/22/22 21:38
@Author  : Neil
@File    : cert_config
*/

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
)

type KeyPair struct {
	pem string
	key string
}

func GetClientConfig(pair *KeyPair) (*tls.Config, error) {
	// set tls config
	fmt.Println(pair.pem, pair.key)
	cert, err := tls.LoadX509KeyPair(pair.pem, pair.key)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	certBytes, err := ioutil.ReadFile(pair.pem)
	if err != nil {
		log.Println("Unable to read cert.pem")
		return nil, err
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		log.Println("failed to parse root certificate")
		return nil, err
	}
	conf := &tls.Config{
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	return conf, nil
}

func GetServerConfig(clientKeyPair *KeyPair, serverKeyPair *KeyPair) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(serverKeyPair.pem, serverKeyPair.key)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	certBytes, err := ioutil.ReadFile(clientKeyPair.pem)
	if err != nil {
		log.Println("Unable to read cert.pem")
		return nil, err
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		log.Println("failed to parse root certificate")
		return nil, err
	}
	conf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}

	return conf, nil
}
