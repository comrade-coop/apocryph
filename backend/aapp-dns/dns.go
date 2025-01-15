package dns

import (
	"context"
	"net"

	"github.com/comrade-coop/apocryph/backend/swarm"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("aapp_dns")

type AappDNS struct {
	Next       plugin.Handler
	swarm      *swarm.Swarm
	baseDomain string
}

func (e AappDNS) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	answer := &dns.Msg{}
	answer.SetReply(r)
	answer.Authoritative = true

	domainName := state.QName()
	if !dns.IsSubDomain(e.baseDomain, domainName) {
		return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
	}

	names := dns.SplitDomainName(domainName)

	subNames := names[:len(names)-dns.CountLabel(e.baseDomain)]
	log.Info("Subnames are:", len(subNames), subNames)

	if len(subNames) != 1 {
		return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
	}

	bucketId := subNames[len(subNames)-1]
	log.Info("bucketId is:", bucketId)

	resultHostnames, err := e.swarm.FindBucketBestNodes(bucketId)
	if err != nil {
		return 0, err
	}

	answer.Extra = []dns.RR{}

	for _, hostname := range resultHostnames {
		switch state.Family() {
		case 1:
			ip4, err := net.ResolveIPAddr("ip4", hostname)
			if err != nil {
				continue
			}
			answer.Extra = append(answer.Extra, &dns.A{
				Hdr: dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()},
				A:   ip4.IP,
			})
		case 2:
			ip6, err := net.ResolveIPAddr("ip6", hostname)
			if err != nil {
				continue
			}
			answer.Extra = append(answer.Extra, &dns.AAAA{
				Hdr:  dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Class: state.QClass()},
				AAAA: ip6.IP,
			})
		}
	}

	w.WriteMsg(answer)

	return 0, nil
}
func (e AappDNS) Name() string { return "aapp_dns" }

func init() { plugin.Register("aapp_dns", setup) }

func setup(c *caddy.Controller) error {
	c.Next()
	if !c.NextArg() {
		return plugin.Error("aapp_dns", c.ArgErr())
	}
	baseDomain := c.Val()
	if !c.NextArg() {
		return plugin.Error("aapp_dns", c.ArgErr())
	}
	serfAddress := c.Val()
	if c.NextArg() {
		return plugin.Error("aapp_dns", c.ArgErr())
	}

	swarm, err := swarm.NewSwarm(serfAddress, "")
	if err != nil {
		return err
	}

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return AappDNS{next, swarm, baseDomain}
	})

	return nil
}
