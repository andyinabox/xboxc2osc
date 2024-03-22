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

func (p *Publisher) Float(route []string, floats ...float32) error {
	r := p.path(route)
	msg := osc.NewMessage(r)
	for _, f := range floats {
		msg.Append(f)
	}
	log.Printf("%s: %v", r, floats)
	return p.client.Send(msg)
}

func (p *Publisher) Bool(route []string, bools ...bool) error {
	r := p.path(route)
	msg := osc.NewMessage(r)
	for _, b := range bools {
		msg.Append(b)
	}
	log.Printf("%s: %v", r, bools)
	return p.client.Send(msg)
}

func (p *Publisher) path(s []string) string {
	args := []string{"/", p.conf.RoutePrefix}
	args = append(args, s...)
	return path.Join(args...)
}
