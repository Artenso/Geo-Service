package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	authServiceClient "github.com/Artenso/proxy-service/internal/clients/auth/client" // replase with import from user-service in normal case
	geoServiceClient "github.com/Artenso/proxy-service/internal/clients/geo/client"   // replase with import from user-service in normal case
	userServiceClient "github.com/Artenso/proxy-service/internal/clients/user/client" // replase with import from user-service in normal case
	"github.com/Artenso/proxy-service/internal/model"

	"github.com/go-chi/jwtauth"
)

type Controller struct {
	responder  Responder
	userClient userServiceClient.Client
	geoClient  geoServiceClient.Client
	authClient authServiceClient.Client
}

func NewController(
	responder Responder,
	userClient userServiceClient.Client,
	geoClient geoServiceClient.Client,
	athClient authServiceClient.Client) *Controller {
	return &Controller{
		responder:  responder,
		userClient: userClient,
		geoClient:  geoClient,
		authClient: athClient,
	}
}

// @Summary      Registration
// @Description  Saves your username and password in db
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input   body      model.RequestAuth  true  "registration data"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
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

	id, err := c.authClient.Register(r.Context(), input.Name, input.Pass)
	if err != nil {
		c.responder.ErrorInternal(w, err)
		return
	}

	c.responder.OutputJSON(w,
		Response{
			Success: true,
			Message: fmt.Sprintf("Congratulations, %s, you are successfully registered!\nYour id: %d", input.Name, id),
		},
	)
}

// @Summary      Log in
// @Description  returns JWT if you are registered user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input   body      model.RequestAuth  true  "registration data"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
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

	token, err := c.authClient.LogIn(r.Context(), input.Name, input.Pass)
	if err != nil {
		if errors.Is(err, model.ErrorUserNotFound) {
			c.responder.ErrorForbidden(w, err)
		} else {
			c.responder.ErrorInternal(w, err)
		}
		return
	}

	c.responder.OutputJSON(w,
		Response{
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
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      403  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
// @Router       /api/address/search [post]
func (c *Controller) GetAddrByPart(w http.ResponseWriter, r *http.Request) {
	requestAddressSearch := model.RequestAddressSearch{
		Query: r.FormValue("query"),
	}

	if len(requestAddressSearch.Query) == 0 {
		c.responder.ErrorBadRequest(w, fmt.Errorf("empty query"))
		return
	}

	res, err := c.geoClient.AddressSearch(r.Context(), requestAddressSearch.Query)
	if err != nil {
		c.responder.ErrorInternal(w, err)
		return
	}

	c.responder.OutputJSON(w,
		Response{
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
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      403  {object}  Response
// @Failure      401  {object}  Response
// @Failure      500  {object}  Response
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

	res, err := c.geoClient.GeoCode(r.Context(), requestAddressGeocode.Lat, requestAddressGeocode.Lng)
	if err != nil {
		c.responder.ErrorInternal(w, err)
		return
	}

	c.responder.OutputJSON(w,
		Response{
			Success: true,
			Data:    res,
		},
	)
}

func (c *Controller) Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := jwtauth.TokenFromHeader(r)

		if err := c.authClient.Verify(r.Context(), token); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		// Token is verified, pass it through
		next.ServeHTTP(w, r)
	})
}

func (c *Controller) decodeRequestBody(r *http.Request) ([]byte, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
