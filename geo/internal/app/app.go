package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	gRPCctrl "github.com/Artenso/geo/internal/controller"
	"github.com/Artenso/geo/internal/service"
	desc "github.com/Artenso/geo/pkg/grpc_geo_provider"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	service        service.GeoProvider
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

func (a *App) initService(_ context.Context) error {
	a.service = service.NewGeoService(os.Getenv("DADATA_APIKEY"), os.Getenv("DADATA_SECRETKEY"))
	return nil
}

func (a *App) initGrpcController(_ context.Context) error {
	a.gRPCcontroller = gRPCctrl.NewController(a.service)
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	s := grpc.NewServer()

	desc.RegisterGeoProviderServer(s, a.gRPCcontroller)

	reflection.Register(s)

	a.gRPCServer = s

	return nil
}
