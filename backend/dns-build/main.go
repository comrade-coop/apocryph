package main

import (
	_ "github.com/comrade-coop/apocryph/backend/aapp-dns"
	_ "github.com/coredns/coredns/plugin/cache"
	_ "github.com/coredns/coredns/plugin/errors"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/health"
	_ "github.com/coredns/coredns/plugin/log"
	_ "github.com/coredns/coredns/plugin/ready"
	_ "github.com/coredns/coredns/plugin/reload"
	_ "github.com/coredns/coredns/plugin/rewrite"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
)

// var directives = []string{
// 	"example",
// 	"whoami",
// 	"aapp_dns",
// 	"startup",
// 	"shutdown",
// }

func init() {
	dnsserver.Directives = []string{
		"root",
		"metadata",
		"geoip",
		"cancel",
		"tls",
		"timeouts",
		"multisocket",
		"reload",
		"nsid",
		"bufsize",
		"bind",
		"debug",
		"trace",
		"ready",
		"health",
		"pprof",
		"prometheus",
		"errors",
		"log",
		"rewrite",
		"cache",
		"aapp_dns",
		"forward",
	}
}

func main() {
	coremain.Run()
}
