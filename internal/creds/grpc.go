package creds

import "google.golang.org/grpc/credentials"

func GetTransportCredentials() (credentials.TransportCredentials, error) {
	certFile := "/misc/cert/cert.crt"
	keyFile := "/misc/cert/cert.key"
	return credentials.NewServerTLSFromFile(certFile, keyFile)
}
