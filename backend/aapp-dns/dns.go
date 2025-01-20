package dns

import (
	"context"
	"net"
	"strconv"

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
	Next        plugin.Handler
	swarm       *swarm.Swarm
	baseDomains []string
}

func (e AappDNS) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	var subNames []string
	found := false

	for _, baseDomain := range e.baseDomains {
		domainName := state.QName()
		if dns.IsSubDomain(baseDomain, domainName) {
			found = true
			names := dns.SplitDomainName(domainName)
			subNames = names[:len(names)-dns.CountLabel(baseDomain)]
		}
	}
	if !found {
		return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
	}

	log.Debug("Subnames are:", subNames)

	// if len(subNames) != 1 {
	// 	return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
	// }

	bucketId := subNames[len(subNames)-1]
	log.Debug("bucketId is:", bucketId)

	answer := &dns.Msg{}
	answer.SetReply(r)
	answer.Authoritative = true
	answer.Compress = true

	resultHostnames, err := e.swarm.FindBucketBestNodes(bucketId)
	if err != nil {
		log.Warningf("Failed to find node: %v", err)
		return dns.RcodeServerFailure, err
	}

	for _, hostname := range resultHostnames {
		switch state.Family() {
		case 1:
			ip4, err := net.ResolveIPAddr("ip4", hostname) // TODO: Ugly recursion
			if err != nil {
				log.Warningf("Failed to resolve node %s: %v", hostname, err)
				continue
			}
			log.Infof("Success A! %s: %v", hostname, ip4.IP)
			answer.Answer = append(answer.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass(), Ttl: 100},
				A:   ip4.IP,
			})
		case 2:
			ip6, err := net.ResolveIPAddr("ip6", hostname) // TODO: Ugly recursion
			if err != nil {
				log.Warningf("Failed to resolve node %s: %v", hostname, err)
				continue
			}
			log.Infof("Success AAAA! %s: %v", hostname, ip6.IP)
			answer.Answer = append(answer.Answer, &dns.AAAA{
				Hdr:  dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Class: state.QClass(), Ttl: 100},
				AAAA: ip6.IP,
			})
		}
	}

	port, _ := strconv.ParseUint(state.Port(), 10, 16)
	answer.Extra = append(answer.Extra, &dns.SRV{
		Hdr:    dns.RR_Header{Name: "_" + state.Proto() + "." + state.QName(), Rrtype: dns.TypeSRV, Class: state.QClass()},
		Port:   uint16(port),
		Target: ".",
	})

	err = w.WriteMsg(answer)

	return dns.RcodeSuccess, err
}
func (e AappDNS) Name() string { return "aapp_dns" }

func init() { plugin.Register("aapp_dns", setup) }

func setup(c *caddy.Controller) error {
	c.Next()
	baseDomains := []string{}
	var serfAddress string
	for c.NextBlock() {
		switch c.Val() {
		case "base":
			if !c.NextArg() {
				print(1)
				return plugin.Error("aapp_dns", c.ArgErr())
			}
			baseDomains = append(baseDomains, c.Val())
			if c.NextArg() {
				return plugin.Error("aapp_dns", c.ArgErr())
			}
		case "serf":
			if !c.NextArg() {
				return plugin.Error("aapp_dns", c.ArgErr())
			}
			serfAddress = c.Val()
			if c.NextArg() {
				return plugin.Error("aapp_dns", c.ArgErr())
			}
		default:
			return plugin.Error("aapp_dns", c.ArgErr())
		}
	}

	swarm, err := swarm.NewSwarm(serfAddress, "")
	if err != nil {
		return err
	}

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return AappDNS{next, swarm, baseDomains}
	})

	return nil
}
