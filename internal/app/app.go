package app

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/Artenso/Geo-Service/internal/token"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	return a.httpServer.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initTokenAuth,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initTokenAuth(_ context.Context) error {
	token.Init()
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	r.Post("/api/login", a.serviceProvider.Controller(ctx).Authentication)
	r.Post("/api/register", a.serviceProvider.Controller(ctx).Registration)

	r.Group(func(r chi.Router) {
		r.Use(token.Verifier())
		r.Use(jwtauth.Authenticator)

		r.Post("/api/address/search", a.serviceProvider.Controller(ctx).GetAddrByPart)
		r.Post("/api/address/geocode", a.serviceProvider.Controller(ctx).GetAddrByCoord)

		r.Route("/metrics/pprof", func(r chi.Router) {
			r.HandleFunc("/*", pprof.Index)
			r.HandleFunc("/cmdline", pprof.Cmdline)
			r.HandleFunc("/profile", pprof.Profile)
			r.HandleFunc("/symbol", pprof.Symbol)
			r.HandleFunc("/trace", pprof.Trace)
		})
	})

	a.httpServer = &http.Server{
		Addr:         os.Getenv("PORT"),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return nil
}
