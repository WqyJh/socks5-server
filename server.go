package main

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/things-go/go-socks5"
)

type params struct {
	User            string   `env:"PROXY_USER" envDefault:""`
	Password        string   `env:"PROXY_PASSWORD" envDefault:""`
	Port            string   `env:"PROXY_PORT" envDefault:"1080"`
	AllowedDestFqdn string   `env:"ALLOWED_DEST_FQDN" envDefault:""`
	AllowedIPs      []string `env:"ALLOWED_IPS" envSeparator:"," envDefault:""`
}

func main() {
	// Working with app params
	cfg := params{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	var opts = []socks5.Option{
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "", log.LstdFlags))),
	}

	if cfg.User+cfg.Password != "" {
		creds := socks5.StaticCredentials{
			os.Getenv("PROXY_USER"): os.Getenv("PROXY_PASSWORD"),
		}
		cator := socks5.UserPassAuthenticator{Credentials: creds}
		opts = append(opts, socks5.WithAuthMethods([]socks5.Authenticator{cator}))
	}

	if cfg.AllowedDestFqdn != "" {
		opts = append(opts, socks5.WithRule(PermitDestAddrPattern(cfg.AllowedDestFqdn)))
	}

	server := socks5.NewServer(opts...)

	// // Set IP whitelist
	// if len(cfg.AllowedIPs) > 0 {
	// 	whitelist := make([]net.IP, len(cfg.AllowedIPs))
	// 	for i, ip := range cfg.AllowedIPs {
	// 		whitelist[i] = net.ParseIP(ip)
	// 	}
	// 	server.SetIPWhitelist(whitelist)
	// }

	log.Printf("Start listening proxy service on port %s\n", cfg.Port)
	if err := server.ListenAndServe("tcp", ":"+cfg.Port); err != nil {
		log.Fatal(err)
	}
}
