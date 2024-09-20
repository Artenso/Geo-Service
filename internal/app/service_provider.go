package app

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/Artenso/Geo-Provider/client"
	gRPCclient "github.com/Artenso/Geo-Provider/client/grpc_geo_provider"
	jsonRPCclient "github.com/Artenso/Geo-Provider/client/json_rpc_geo_provider"
	"github.com/Artenso/Geo-Service/internal/controller"
	"github.com/Artenso/Geo-Service/internal/logger"
	"github.com/Artenso/Geo-Service/internal/responder"
	"github.com/Artenso/Geo-Service/internal/service"
	storage "github.com/Artenso/Geo-Service/internal/storage/pg"
	"github.com/jackc/pgx/v5"
	"github.com/ptflp/godecoder"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serviceProvider struct {
	dbConn     *pgx.Conn
	storage    storage.IStorage
	servise    service.IService
	responder  responder.Responder
	rpcClient  client.Client
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

func (s *serviceProvider) Service(ctx context.Context) service.IService {
	if s.servise == nil {
		s.servise = service.NewService(s.Storage(ctx))
	}

	return s.servise
}

func (s *serviceProvider) Responder(ctx context.Context) responder.Responder {
	if s.responder == nil {
		s.responder = responder.NewResponder(godecoder.NewDecoder())
	}

	return s.responder
}

func (s *serviceProvider) RPCclient(ctx context.Context) client.Client {
	if s.rpcClient == nil {
		switch os.Getenv("RPC_PROTOCOL") {
		case "grpc":
			conn, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				logger.Fatalf("Ошибка при подключении к серверу: %s", err)
			}

			s.rpcClient = gRPCclient.NewGRPCclient(conn)

		case "json-rpc":
			conn, err := net.Dial("tcp", "localhost:1234")
			if err != nil {
				logger.Fatalf("Ошибка при подключении к серверу: %s", err)
			}

			s.rpcClient = jsonRPCclient.NewJSONrpcClient(conn)
		}
	}

	return s.rpcClient
}

func (s *serviceProvider) Controller(ctx context.Context) *controller.Controller {
	if s.controller == nil {
		s.controller = controller.NewController(s.Responder(ctx), s.Service(ctx), s.RPCclient(ctx))
	}

	return s.controller
}
