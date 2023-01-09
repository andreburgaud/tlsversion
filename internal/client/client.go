package client

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

type Result int

const (
	ConnectionFailed Result = iota
	NotSupported
	Supported
)

// SupportedTls dials a remote server and retrieve the certificate to examine the expiration date
// the domain is either domain or domain:port. If the port is not provided, it defaults to 443.
func SupportedTls(domain string, tlsVersion uint16, timeout int) (Result, error) {

	var host, port string

	var err error
	if host, port, err = net.SplitHostPort(domain); err != nil {
		host = domain
		port = "443"
	}

	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", timeout))
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: duration},
		"tcp", net.JoinHostPort(host, port),
		&tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tlsVersion,
			MaxVersion:         tlsVersion,
		})

	if err != nil {
		if strings.Contains(err.Error(), "tls") {
			return NotSupported, err
		}
		return ConnectionFailed, err
	}

	defer func(conn *tls.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	return Supported, nil
}
