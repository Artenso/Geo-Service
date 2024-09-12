package app

import (
	"context"
	"fmt"
	"os"

	"github.com/Artenso/Geo-Service/internal/controller"
	"github.com/Artenso/Geo-Service/internal/logger"
	"github.com/Artenso/Geo-Service/internal/responder"
	"github.com/Artenso/Geo-Service/internal/service"
	storage "github.com/Artenso/Geo-Service/internal/storage/pg"
	"github.com/jackc/pgx/v5"
	"github.com/ptflp/godecoder"
)

type serviceProvider struct {
	dbConn *pgx.Conn

	storage storage.IStorage

	geoServ service.GeoProvider

	servise service.IService

	responder responder.Responder

	controller *controller.Controller
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) DbConn(ctx context.Context) *pgx.Conn {
	if s.dbConn == nil {
		dbDSN := fmt.Sprintf(
			`postgres://%s:%s@%s:%s/%s`,
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
		conn, err := pgx.Connect(ctx, dbDSN)
		if err != nil {
			logger.Fatalf("failed to init db connection: %s", err.Error())
		}

		s.dbConn = conn
	}

	return s.dbConn
}

func (s *serviceProvider) Storage(ctx context.Context) storage.IStorage {
	if s.storage == nil {
		s.storage = storage.NewStorage(s.DbConn(ctx))
	}

	return s.storage
}

func (s *serviceProvider) GeoServ(ctx context.Context) service.GeoProvider {
	if s.geoServ == nil {
		s.geoServ = service.NewGeoService(os.Getenv("DADATA_APIKEY"), os.Getenv("DADATA_SECRETKEY"))
	}

	return s.geoServ
}

func (s *serviceProvider) Service(ctx context.Context) service.IService {
	if s.servise == nil {
		s.servise = service.NewService(s.Storage(ctx), s.GeoServ(ctx))
	}

	return s.servise
}

func (s *serviceProvider) Responder(ctx context.Context) responder.Responder {
	if s.responder == nil {
		s.responder = responder.NewResponder(godecoder.NewDecoder())
	}

	return s.responder
}

func (s *serviceProvider) Controller(ctx context.Context) *controller.Controller {
	if s.controller == nil {
		s.controller = controller.NewController(s.Responder(ctx), s.Service(ctx))
	}

	return s.controller
}
