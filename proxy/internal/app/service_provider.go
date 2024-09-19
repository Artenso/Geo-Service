package app

import (
	"context"
	"fmt"
	"os"

	authServiceClient "github.com/Artenso/proxy-service/internal/clients/auth/client" // replase with import from user-service in normal case
	geoServiceClient "github.com/Artenso/proxy-service/internal/clients/geo/client"   // replase with import from user-service in normal case
	userServiceClient "github.com/Artenso/proxy-service/internal/clients/user/client" // replase with import from user-service in normal case

	"github.com/Artenso/proxy-service/internal/controller"
	"github.com/Artenso/proxy-service/internal/logger"

	"github.com/jackc/pgx/v5"
	"github.com/ptflp/godecoder"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serviceProvider struct {
	dbConn     *pgx.Conn
	responder  controller.Responder
	userClient userServiceClient.Client
	geoClient  geoServiceClient.Client
	authClient authServiceClient.Client
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

func (s *serviceProvider) Responder(ctx context.Context) controller.Responder {
	if s.responder == nil {
		s.responder = controller.NewResponder(godecoder.NewDecoder())
	}

	return s.responder
}

func (s *serviceProvider) UserClient(ctx context.Context) userServiceClient.Client {
	if s.userClient == nil {

		conn, err := grpc.NewClient(
			fmt.Sprintf(
				"%s%s",
				os.Getenv("USER_SERVICE_HOST"),
				os.Getenv("USER_SERVICE_PORT"),
			),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatalf("Ошибка при подключении к серверу: %s", err)
		}

		s.userClient = userServiceClient.NewGRPCclient(conn)
	}

	return s.userClient
}

func (s *serviceProvider) GeoClient(ctx context.Context) geoServiceClient.Client {
	if s.geoClient == nil {

		conn, err := grpc.NewClient(
			fmt.Sprintf(
				"%s%s",
				os.Getenv("GEO_SERVICE_HOST"),
				os.Getenv("GEO_SERVICE_PORT"),
			),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatalf("Ошибка при подключении к серверу: %s", err)
		}

		s.geoClient = geoServiceClient.NewGRPCclient(conn)
	}

	return s.geoClient
}

func (s *serviceProvider) AuthClient(ctx context.Context) authServiceClient.Client {
	if s.authClient == nil {

		conn, err := grpc.NewClient(
			fmt.Sprintf(
				"%s%s",
				os.Getenv("AUTH_SERVICE_HOST"),
				os.Getenv("AUTH_SERVICE_PORT"),
			),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatalf("Ошибка при подключении к серверу: %s", err)
		}

		s.authClient = authServiceClient.NewGRPCclient(conn)
	}

	return s.authClient
}

func (s *serviceProvider) Controller(ctx context.Context) *controller.Controller {
	if s.controller == nil {
		s.controller = controller.NewController(
			s.Responder(ctx),
			s.UserClient(ctx),
			s.GeoClient(ctx),
			s.AuthClient(ctx),
		)
	}

	return s.controller
}
