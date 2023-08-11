package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"
)

const (
	defaultServerIdleTimeout  = 12 * time.Second
	defaultServerReadTimeout  = 12 * time.Second
	defaultServerWriteTimeout = 12 * time.Second
	defaultServerAddr         = "0.0.0.0:3000"
)

type Server struct {
	*http.Server
}

type Option func(*Server)

func New(opts ...Option) *Server {
	defaultServer := &http.Server{
		Handler:      http.DefaultServeMux,
		Addr:         defaultServerAddr,
		TLSConfig:    getDefaultTLSConfig(),
		IdleTimeout:  defaultServerIdleTimeout,
		ReadTimeout:  defaultServerReadTimeout,
		WriteTimeout: defaultServerWriteTimeout,
	}

	server := &Server{
		defaultServer,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (s *Server) Start() error {
	return s.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.Shutdown(ctxTimeout)
}

func WithHandler(handler http.Handler) Option {
	return func(s *Server) {
		s.Handler = handler
	}
}

func WithAddress(addr string) Option {
	return func(s *Server) {
		s.Addr = addr
	}
}

func WithTLSConfig(config *tls.Config) Option {
	return func(s *Server) {
		s.TLSConfig = config
	}
}

func getDefaultTLSConfig() *tls.Config {
	tlsConfig := &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	return tlsConfig
}
