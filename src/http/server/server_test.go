package server

import "testing"

func TestServerOptionsDefaultHTTPMode(t *testing.T) {
	httpServer := NewServer(Config{
		Host: "127.0.0.1",
		Port: "8080",
	})

	options := httpServer.GetOptions()
	if options["Mode"] != "http" {
		t.Fatalf("expected default mode http, got %s", options["Mode"])
	}
	if options["Port"] != "8080" {
		t.Fatalf("expected http port 8080, got %s", options["Port"])
	}
}

func TestServerOptionsHTTPAndTLSMode(t *testing.T) {
	httpServer := NewServer(Config{
		Host: "127.0.0.1",
		Port: "8000",
		TLS: TLSConfig{
			Enable:   true,
			Port:     "8443",
			CertFile: "server.crt",
			KeyFile:  "server.key",
		},
	})

	options := httpServer.GetOptions()
	if options["Mode"] != "http+https" {
		t.Fatalf("expected tls mode http+https, got %s", options["Mode"])
	}
	if options["TLSPort"] != "8443" {
		t.Fatalf("expected tls port 8443, got %s", options["TLSPort"])
	}
}

func TestServerOptionsTLSOnlyMode(t *testing.T) {
	httpServer := NewServer(Config{
		Host: "127.0.0.1",
		TLS: TLSConfig{
			Enable:   true,
			Port:     "8443",
			CertFile: "server.crt",
			KeyFile:  "server.key",
		},
	})

	options := httpServer.GetOptions()
	if options["Mode"] != "https" {
		t.Fatalf("expected tls only mode https, got %s", options["Mode"])
	}
	if _, exists := options["Port"]; exists {
		t.Fatalf("expected no http port, got %s", options["Port"])
	}
	if options["TLSPort"] != "8443" {
		t.Fatalf("expected tls port 8443, got %s", options["TLSPort"])
	}
}

func TestServerOptionsLegacyTLSOnlyMode(t *testing.T) {
	httpServer := NewServer(Config{
		Host: "127.0.0.1",
		Port: "8443",
		TLS: TLSConfig{
			Enable:   true,
			CertFile: "server.crt",
			KeyFile:  "server.key",
		},
	})

	options := httpServer.GetOptions()
	if options["Mode"] != "https" {
		t.Fatalf("expected legacy tls only mode https, got %s", options["Mode"])
	}
	if _, exists := options["Port"]; exists {
		t.Fatalf("expected no http port, got %s", options["Port"])
	}
	if options["TLSPort"] != "8443" {
		t.Fatalf("expected tls port 8443, got %s", options["TLSPort"])
	}
}

func TestValidateConfigRequiresHTTPOrTLS(t *testing.T) {
	httpServer := NewServer(Config{})

	err := httpServer.validateConfig()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestValidateHTTPSConfig(t *testing.T) {
	tests := []struct {
		port      string
		name      string
		tlsConfig TLSConfig
		wantErr   bool
	}{
		{
			port: "8000",
			name: "valid config",
			tlsConfig: TLSConfig{
				Port:     "8443",
				CertFile: "server.crt",
				KeyFile:  "server.key",
			},
		},
		{
			port: "",
			name: "missing tls port",
			tlsConfig: TLSConfig{
				CertFile: "server.crt",
				KeyFile:  "server.key",
			},
			wantErr: true,
		},
		{
			port: "8000",
			name: "same http and tls port",
			tlsConfig: TLSConfig{
				Port:     "8000",
				CertFile: "server.crt",
				KeyFile:  "server.key",
			},
			wantErr: true,
		},
		{
			port: "8000",
			name: "legacy tls port from http port",
			tlsConfig: TLSConfig{
				CertFile: "server.crt",
				KeyFile:  "server.key",
			},
		},
		{
			port: "8000",
			name: "missing cert file",
			tlsConfig: TLSConfig{
				Port:    "8443",
				KeyFile: "server.key",
			},
			wantErr: true,
		},
		{
			port: "8000",
			name: "missing key file",
			tlsConfig: TLSConfig{
				Port:     "8443",
				CertFile: "server.crt",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpServer := NewServer(Config{
				Port: tt.port,
				TLS:  tt.tlsConfig,
			})

			err := httpServer.validateHTTPSConfig()
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}
