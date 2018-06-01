package main

import (
	"github.com/kabukky/httpscerts"
	"log"
)

const (
	sslCert = "etc/cert.pem"
	sslKey  = "etc/key.pem"
)

func setupSsl() {
	if settings.UseSsl == false {
		return
	}

	// need to edit isCA in httpscerts package to = false
	// not very clean

	err := httpscerts.Check(sslCert, sslKey)
	if err != nil {
		err = httpscerts.Generate(sslCert, sslKey, settings.SslHostname)
		if err != nil {
			log.Fatalf("Unable to create ssl certs: %v", err)
		}
	}
}
