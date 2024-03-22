package oscpub

import (
	"log"
	"path"

	"github.com/hypebeast/go-osc/osc"
)

type Config struct {
	HostDomain  string
	HostPort    int
	RoutePrefix string
}

type Publisher struct {
	conf   *Config
	client *osc.Client
}

func New(conf *Config) *Publisher {
	return &Publisher{conf, osc.NewClient(conf.HostDomain, conf.HostPort)}
}

func (p *Publisher) Float(route []string, f float32) error {
	r := p.path(route)
	msg := osc.NewMessage(r)
	msg.Append(f)
	log.Printf("%s: %f", r, f)
	return p.client.Send(msg)
}

func (p *Publisher) Bool(route []string, b bool) error {
	r := p.path(route)
	msg := osc.NewMessage(r)
	msg.Append(b)
	log.Printf("%s: %v", r, b)
	return p.client.Send(msg)
}

func (p *Publisher) path(s []string) string {
	args := []string{"/", p.conf.RoutePrefix}
	args = append(args, s...)
	return path.Join(args...)
}
