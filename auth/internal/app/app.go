package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Artenso/auth-service/internal/clients/user/client" // replase with import from user-service in normal case
	"github.com/Artenso/auth-service/internal/controller"
	"github.com/Artenso/auth-service/internal/service"
	desc "github.com/Artenso/auth-service/pkg/auth_service"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	service    service.Service
	userClient client.Client
	controller *controller.Controller
	gRPCServer *grpc.Server
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
		a.initUserClient,
		a.initController,
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
	tokenAuth := jwtauth.New(
		"HS256",
		[]byte(os.Getenv("JWTSECRET")),
		nil,
	)

	a.service = service.NewService(tokenAuth)
	return nil
}

func (a *App) initUserClient(_ context.Context) error {
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s%s",
			os.Getenv("USER_SERVICE_HOST"),
			os.Getenv("USER_SERVICE_PORT"),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	a.userClient = client.NewGRPCclient(conn)

	return nil
}

func (a *App) initController(_ context.Context) error {
	a.controller = controller.NewController(a.service, a.userClient)
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	s := grpc.NewServer()

	desc.RegisterAuthServiceServer(s, a.controller)

	reflection.Register(s)

	a.gRPCServer = s

	return nil
}
