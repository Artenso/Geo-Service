package converter

import (
	"github.com/Artenso/geo/internal/model"
	desc "github.com/Artenso/geo/pkg/grpc_geo_provider"
)

func ToAddressSearchResponse(addresses []*model.Address) *desc.AddressSearchResponse {
	resp := &desc.AddressSearchResponse{}

	for _, addr := range addresses {
		descAddr := &desc.Address{
			City:   addr.City,
			Street: addr.Street,
			House:  addr.House,
			Lat:    addr.Lat,
			Lon:    addr.Lon,
		}

		resp.Addresses = append(resp.Addresses, descAddr)
	}

	return resp
}

func ToGeoCodeResponse(addresses []*model.Address) *desc.GeoCodeResponse {
	resp := &desc.GeoCodeResponse{}

	for _, addr := range addresses {
		descAddr := &desc.Address{
			City:   addr.City,
			Street: addr.Street,
			House:  addr.House,
			Lat:    addr.Lat,
			Lon:    addr.Lon,
		}

		resp.Addresses = append(resp.Addresses, descAddr)
	}

	return resp
}
