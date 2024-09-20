package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Artenso/Geo-Provider/client"
	"github.com/Artenso/Geo-Service/internal/converter"
	"github.com/Artenso/Geo-Service/internal/model"
	"github.com/Artenso/Geo-Service/internal/responder"
	"github.com/Artenso/Geo-Service/internal/service"
)

type Controller struct {
	responder responder.Responder
	service   service.IService
	rpcClient client.Client
}

func NewController(responder responder.Responder, service service.IService, rpcClient client.Client) *Controller {
	return &Controller{
		responder: responder,
		service:   service,
		rpcClient: rpcClient,
	}
}

// @Summary      Registration
// @Description  Saves your username and password in db
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input   body      model.RequestAuth  true  "registration data"
// @Success      200  {object}  responder.Response
// @Failure      400  {object}  responder.Response
// @Router       /api/register [post]
func (c *Controller) Registration(w http.ResponseWriter, r *http.Request) {
	requestBody, err := c.decodeRequestBody(r)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	input := &model.RequestAuth{}

	json.Unmarshal(requestBody, input)

	if input.Name == "" || input.Pass == "" {
		c.responder.ErrorBadRequest(w, fmt.Errorf("missing username or password"))
		return
	}

	user := converter.RequestAuthToUser(input)

	if err := c.service.RegistrateUser(r.Context(), user); err != nil {
		c.responder.ErrorInternal(w, err)
		return
	}

	c.responder.OutputJSON(w,
		responder.Response{
			Success: true,
			Message: fmt.Sprintf("Congratulations, %s, you are successfully registered", user.Name),
		},
	)
}

// @Summary      Log in
// @Description  returns JWT if you are registered user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input   body      model.RequestAuth  true  "registration data"
// @Success      200  {object}  responder.Response
// @Failure      400  {object}  responder.Response
// @Router       /api/login [post]
func (c *Controller) Authentication(w http.ResponseWriter, r *http.Request) {
	requestBody, err := c.decodeRequestBody(r)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	input := &model.RequestAuth{}

	json.Unmarshal(requestBody, input)

	if input.Name == "" || input.Pass == "" {
		c.responder.ErrorBadRequest(w, fmt.Errorf("missing username or password"))
		return
	}

	user := converter.RequestAuthToUser(input)

	token, err := c.service.AuthenticateUser(r.Context(), user)
	if err != nil {
		if errors.Is(err, model.ErrorUserNotFound) {
			c.responder.ErrorForbidden(w, err)
		} else {
			c.responder.ErrorInternal(w, err)
		}
		return
	}

	c.responder.OutputJSON(w,
		responder.Response{
			Success: true,
			Data: model.ResponseLogin{
				Token: token,
			},
		},
	)
}

// @Summary      Adress Search
// @Security     ApiKeyAuth
// @Description  Get full address info by its part
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        query   query      string  true  "part of address"
// @Success      200  {object}  responder.Response
// @Failure      400  {object}  responder.Response
// @Failure      403  {object}  responder.Response
// @Failure      401  {object}  responder.Response
// @Failure      500  {object}  responder.Response
// @Router       /api/address/search [post]
func (c *Controller) GetAddrByPart(w http.ResponseWriter, r *http.Request) {
	requestAddressSearch := model.RequestAddressSearch{
		Query: r.FormValue("query"),
	}

	if len(requestAddressSearch.Query) == 0 {
		c.responder.ErrorBadRequest(w, fmt.Errorf("empty query"))
		return
	}

	res, err := c.rpcClient.AddressSearch(r.Context(), requestAddressSearch.Query)
	if err != nil {
		c.responder.ErrorInternal(w, err)
		return
	}

	c.responder.OutputJSON(w,
		responder.Response{
			Success: true,
			Data:    res,
		},
	)
}

// @Summary      GeoCode
// @Security     ApiKeyAuth
// @Description  Get full address info by coordinates
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        lat   query      string  true  "latitude"
// @Param        lng   query      string  true  "longitude"
// @Success      200  {object}  responder.Response
// @Failure      400  {object}  responder.Response
// @Failure      403  {object}  responder.Response
// @Failure      401  {object}  responder.Response
// @Failure      500  {object}  responder.Response
// @Router       /api/address/geocode [post]
func (c *Controller) GetAddrByCoord(w http.ResponseWriter, r *http.Request) {
	requestAddressGeocode := model.RequestAddressGeocode{
		Lat: r.FormValue("lat"),
		Lng: r.FormValue("lng"),
	}

	if requestAddressGeocode.Lat == "" || requestAddressGeocode.Lng == "" {
		c.responder.ErrorBadRequest(w, fmt.Errorf("empty lat or lng parametr"))
		return
	}

	res, err := c.rpcClient.GeoCode(r.Context(), requestAddressGeocode.Lat, requestAddressGeocode.Lng)
	if err != nil {
		c.responder.ErrorInternal(w, err)
		return
	}

	c.responder.OutputJSON(w,
		responder.Response{
			Success: true,
			Data:    res,
		},
	)
}

func (c *Controller) decodeRequestBody(r *http.Request) ([]byte, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
