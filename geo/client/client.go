package client

import (
	"context"

	"github.com/Artenso/geo/internal/model"
	desc "github.com/Artenso/geo/pkg/grpc_geo_provider"
	"google.golang.org/grpc"
)

type Client interface {
	AddressSearch(ctx context.Context, input string) ([]*model.Address, error)
	GeoCode(ctx context.Context, lat, lng string) ([]*model.Address, error)
}

type client struct {
	client desc.GeoProviderClient
}

func NewGRPCclient(conn *grpc.ClientConn) Client {
	return &client{
		client: desc.NewGeoProviderClient(conn),
	}
}

func (c *client) AddressSearch(ctx context.Context, input string) ([]*model.Address, error) {
	req := &desc.AddressSearchRequest{
		Input: input,
	}

	resp, err := c.client.AddressSearch(ctx, req)
	if err != nil {
		return nil, err
	}

	addresses := make([]*model.Address, 0, len(resp.Addresses))

	for _, descAddr := range resp.Addresses {
		addr := &model.Address{
			City:   descAddr.City,
			Street: descAddr.Street,
			House:  descAddr.House,
			Lat:    descAddr.Lat,
			Lon:    descAddr.Lon,
		}

		addresses = append(addresses, addr)
	}

	return addresses, nil
}

func (c *client) GeoCode(ctx context.Context, lat, lng string) ([]*model.Address, error) {
	req := &desc.GeoCodeRequest{
		Lat: lat,
		Lng: lng,
	}

	resp, err := c.client.GeoCode(ctx, req)
	if err != nil {
		return nil, err
	}

	addresses := make([]*model.Address, 0, len(resp.Addresses))

	for _, descAddr := range resp.Addresses {
		addr := &model.Address{
			City:   descAddr.City,
			Street: descAddr.Street,
			House:  descAddr.House,
			Lat:    descAddr.Lat,
			Lon:    descAddr.Lon,
		}

		addresses = append(addresses, addr)
	}

	return addresses, nil
}
