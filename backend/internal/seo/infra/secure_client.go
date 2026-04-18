package infra

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

var ipValidator = func(ip net.IP) bool {
	return !ip.IsPrivate() && !ip.IsLoopback() && !ip.IsLinkLocalUnicast()
}

func CreateSecureClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 100,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, fmt.Errorf("failed to split host and port: %w", err)
				}

				ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
				if err != nil {
					return nil, fmt.Errorf("failed to lookup IP: %w", err)
				}

				if len(ips) == 0 {
					return nil, fmt.Errorf("no IP addresses found for host: %s", host)
				}

				resIP := ips[0]
				if !ipValidator(resIP) {
					return nil, fmt.Errorf("access to restricted IP denied: %s", resIP)
				}

				return dialer.DialContext(ctx, network, net.JoinHostPort(resIP.String(), port))
			},
		},
	}
}
