package main

import (
	"context"
	"strconv"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

type DNS struct {
	swarm *Swarm
}

// ServeDNS implements the plugin.Handler interface. This method gets called when apocryphS3DNS is used
// in a Server.
func (e DNS) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	answer := &dns.Msg{}
	answer.SetReply(r)
	answer.Authoritative = true

	resultIps, err := e.swarm.FindBucket(state.QName())
	if err != nil {
		return 0, err
	}

	answer.Extra = []dns.RR{}

	for _, ip := range resultIps {
		switch state.Family() {
		case 1:
			ip4 := ip.To4()
			if ip4 != nil {
				answer.Extra = append(answer.Extra, &dns.A{
					Hdr: dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()},
					A:   ip4,
				})
			}
		case 2:
			answer.Extra = append(answer.Extra, &dns.AAAA{
				Hdr:  dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Class: state.QClass()},
				AAAA: ip,
			})
		}
	}

	port, _ := strconv.ParseUint(state.Port(), 10, 16)
	var name string
	if state.QName() == "." {
		name = "_" + state.Proto() + "." + state.QName()
	} else {
		name = "_" + state.Proto() + state.QName()
	}
	answer.Extra = append(answer.Extra, &dns.SRV{
		Hdr:    dns.RR_Header{Name: name, Rrtype: dns.TypeSRV, Class: state.QClass()},
		Port:   uint16(port),
		Target: ".",
	})

	w.WriteMsg(answer)

	return 0, nil
}
func (e DNS) Name() string { return "apocryphS3DNS" }

func init() { plugin.Register("apocryphS3DNS", setup) }

func setup(c *caddy.Controller) error {
	c.Next()
	if !c.NextArg() {
		return plugin.Error("apocryphS3DNS", c.ArgErr())
	}
	serfAddress := c.Val()
	if c.NextArg() {
		return plugin.Error("apocryphS3DNS", c.ArgErr())
	}

	swarm, err := NewSwarm(serfAddress)
	if err != nil {
		return err
	}

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return DNS{swarm}
	})

	return nil
}

func RunDNS(ctx context.Context, dnsAddress string, swarm *Swarm) error {
	server, err := dnsserver.NewServer(dnsAddress, []*dnsserver.Config{
		{
			Zone: ".",
			Plugin: []plugin.Plugin{
				func(next plugin.Handler) plugin.Handler {
					return DNS{swarm}
				},
			},
		},
	})
	if err != nil {
		return err
	}

	l, err := server.Listen()
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		l.Close()
	}()

	return nil
}
