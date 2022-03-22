/*
@Time    : 3/22/22 21:38
@Author  : Neil
@File    : cert_config
*/

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

func GetClientConfig() (*tls.Config, error) {
	// set tls config
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	certBytes, err := ioutil.ReadFile("certs/client.pem")
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

func GetServerConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	certBytes, err := ioutil.ReadFile("certs/client.pem")
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
