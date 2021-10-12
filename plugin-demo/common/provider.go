package provider

import (
	"net/rpc"

	plugin "github.com/hashicorp/go-plugin"
)

type Provider interface {
	Create(scheme *Scheme) (string, error)
	Update(scheme *Scheme) (string, error)
	Delete(scheme *Scheme) (string, error)
	Get(id string, resp *Scheme) error
}

type Scheme struct {
	Default interface{}
}

type ProviderRPC struct {
	client *rpc.Client
}

func (r *ProviderRPC) Create(scheme *Scheme) (string, error) {
	var resp string
	err := r.client.Call("Plugin.Create", scheme, &resp)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (r *ProviderRPC) Update(scheme *Scheme) (string, error) {
	var resp string
	err := r.client.Call("Plugin.Update", scheme, &resp)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (r *ProviderRPC) Delete(scheme *Scheme) (string, error) {
	var resp string
	err := r.client.Call("Plugin.Delete", scheme, &resp)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (r *ProviderRPC) Get(id string, resp *Scheme) error {
	err := r.client.Call("Plugin.Get", id, &resp)
	if err != nil {
		return err
	}

	return nil
}

type ProviderRPCServer struct {
	Impl Provider
}

func (s *ProviderRPCServer) Create(args *Scheme, resp *string) (err error) {
	*resp, err = s.Impl.Create(args)
	return
}

func (s *ProviderRPCServer) Update(args *Scheme, resp *string) (err error) {
	*resp, err = s.Impl.Update(args)
	return
}

func (s *ProviderRPCServer) Delete(args *Scheme, resp *string) (err error) {
	*resp, err = s.Impl.Delete(args)
	return
}

func (s *ProviderRPCServer) Get(id string, resp *Scheme) (err error) {
	err = s.Impl.Get(id, resp)
	return
}

type ProviderPlugin struct {
	Impl Provider
}

func (p *ProviderPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ProviderRPCServer{Impl: p.Impl}, nil
}

func (*ProviderPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ProviderRPC{client: c}, nil
}
