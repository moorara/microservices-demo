package main

import (
	"github.com/caddyserver/caddy/caddy/caddymain"

	// plug in plugins here, for example:
	_ "github.com/btburke/caddy-jwt"
	_ "github.com/casbin/caddy-authz"
	_ "github.com/miekg/caddy-prometheus"
	_ "github.com/xuqingfeng/caddy-rate-limit"
)

func main() {
	// optional: disable telemetry
	// caddymain.EnableTelemetry = false
	caddymain.Run()
}
