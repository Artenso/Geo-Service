package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	gRPCctrl "github.com/Artenso/user-service/internal/controller"
	"github.com/Artenso/user-service/internal/service"
	"github.com/Artenso/user-service/internal/storage"
	desc "github.com/Artenso/user-service/pkg/user_service"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	storage        storage.IStorage
	service        service.IService
	gRPCcontroller *gRPCctrl.Controller
	gRPCServer     *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {

	list, err := net.Listen("tcp", os.Getenv("GRPC_PORT"))
	if err != nil {
		return fmt.Errorf("failed to mapping port: %s", err.Error())
	}

	if err := a.gRPCServer.Serve(list); err != nil {
		return fmt.Errorf("failed to run server: %s", err.Error())
	}

	return nil

}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initStorage,
		a.initService,
		a.initGrpcController,
		a.initGRPCServer,
	}

	for _, init := range inits {
		err := init(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
		return err
	}

	return nil
}

func (a *App) initStorage(ctx context.Context) error {
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
		log.Fatalf("failed to init db connection: %s", err.Error())
	}

	a.storage = storage.NewStorage(conn)
	return nil
}

func (a *App) initService(_ context.Context) error {
	a.service = service.NewService(a.storage)
	return nil
}

func (a *App) initGrpcController(_ context.Context) error {
	a.gRPCcontroller = gRPCctrl.NewController(a.service)
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	s := grpc.NewServer()

	desc.RegisterUserServiceServer(s, a.gRPCcontroller)

	reflection.Register(s)

	a.gRPCServer = s

	return nil
}
