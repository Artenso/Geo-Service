package model

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type RequestAuth struct {
	Name string `json:"username"`
	Pass string `json:"password"`
}

type ResponseLogin struct {
	Token string `json:"token"`
}

type ResponseAddress struct {
	Addresses []*Address `json:"addresses"`
}
