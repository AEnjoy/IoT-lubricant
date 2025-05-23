package crypto

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql/driver"
	"errors"
	"os"

	"github.com/bytedance/sonic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	ErrInvalidRootCrt = errors.New("failed to parse root certificate")
)

type Tls struct {
	Enable     bool   `json:"enable" yaml:"enable" env:"TLS_ENABLE" envDefault:"false"`
	SkipVerify bool   `json:"skip_verify" yaml:"skip_verify" env:"TLS_SKIP_VERIFY"`
	FromFile   bool   `json:"from_file" yaml:"from_file" env:"TLS_FROM_FILE"`
	Key        string `json:"key" yaml:"key" env:"TLS_KEY" `   // server.key
	Cert       string `json:"cert" yaml:"cert" env:"TLS_CERT"` // client.crt
	CA         string `json:"ca" yaml:"ca" env:"TLS_CA"`       // server.crt
}

func (t *Tls) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return sonic.Unmarshal(bytes, t)
}

func (t Tls) Value() (driver.Value, error) {
	return sonic.MarshalString(t)
}

// GetTLSLinkConfig for client
func (t Tls) GetTLSLinkConfig() (credentials.TransportCredentials, error) {
	if !t.Enable {
		return nil, nil
	}
	var cert tls.Certificate
	var caCert []byte
	var err error
	if t.FromFile {
		cert, err = tls.LoadX509KeyPair(t.Cert, t.Key)
		if err != nil {
			return nil, err
		}
		caCert, err = os.ReadFile(t.CA)
		if err != nil {
			return nil, err
		}
	} else if len(t.Cert) != 0 && len(t.Key) != 0 {
		cert, err = tls.X509KeyPair([]byte(t.Cert), []byte(t.Key))
		if err != nil {
			return nil, err
		}
		caCert = []byte(t.CA)
	}

	caCertPool, _ := x509.SystemCertPool()
	if len(caCert) != 0 && !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, ErrInvalidRootCrt
	}
	cred := credentials.NewTLS(&tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	})
	return cred, nil
}

// GetServerTlsConfig for server
func (t Tls) GetServerTlsConfig() (grpc.ServerOption, error) {
	if !t.Enable {
		return nil, nil
	}
	var creds credentials.TransportCredentials
	var err error
	if t.FromFile {
		creds, err = credentials.NewServerTLSFromFile(t.CA, t.Key)
		if err != nil {
			return nil, err
		}
	} else {
		cert, err := tls.X509KeyPair([]byte(t.Cert), []byte(t.Key))
		if err != nil {
			return nil, err
		}
		creds = credentials.NewServerTLSFromCert(&cert)
	}

	return grpc.Creds(creds), nil
}
