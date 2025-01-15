package main

import (
	"slices"

	_ "github.com/comrade-coop/apocryph/backend/aapp-dns"
	_ "github.com/coredns/coredns/plugin/cache"
	_ "github.com/coredns/coredns/plugin/errors"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/health"
	_ "github.com/coredns/coredns/plugin/log"
	_ "github.com/coredns/coredns/plugin/ready"

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
	dnsserver.Directives = slices.Concat([]string{"aapp_dns"}, dnsserver.Directives)
}

func main() {
	coremain.Run()
}
