package svc

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewDataStoreServiceClientWithHost(host string, withTls bool) (serv DataStoreServiceClient, closeCall func() error, err error) {
	var (
		conn  *grpc.ClientConn
		creds credentials.TransportCredentials
	)

	if withTls {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, nil, err
		}

		tlsConfig := &tls.Config{
			RootCAs: certPool,
		}
		creds = credentials.NewTLS(tlsConfig)
	} else {
		creds = insecure.NewCredentials()
	}

	conn, err = grpc.NewClient(host,
		grpc.WithTransportCredentials(creds),
	)

	if err != nil {
		return nil, nil, err
	}
	serv = NewDataStoreServiceClient(conn)
	_, err = serv.Ping(context.Background(), &metapb.Ping{})
	if err != nil {
		return nil, nil, err
	}
	return serv, conn.Close, nil
}
