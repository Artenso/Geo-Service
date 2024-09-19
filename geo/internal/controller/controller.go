package controller

import (
	"context"

	"github.com/Artenso/geo/internal/converter"
	"github.com/Artenso/geo/internal/service"
	desc "github.com/Artenso/geo/pkg/grpc_geo_provider"
)

type Controller struct {
	desc.UnimplementedGeoProviderServer

	service service.GeoProvider
}

func NewController(service service.GeoProvider) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) AddressSearch(ctx context.Context, req *desc.AddressSearchRequest) (*desc.AddressSearchResponse, error) {
	addresses, err := c.service.AddressSearch(req.Input)
	if err != nil {
		return nil, err
	}

	return converter.ToAddressSearchResponse(addresses), nil
}

func (c *Controller) GeoCode(ctx context.Context, req *desc.GeoCodeRequest) (*desc.GeoCodeResponse, error) {
	addresses, err := c.service.GeoCode(req.Lat, req.Lng)
	if err != nil {
		return nil, err
	}

	return converter.ToGeoCodeResponse(addresses), nil
}
