package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	kitlog "github.com/go-kit/log"
	"github.com/ppeymann/top-app.git/auth"
	"github.com/ppeymann/top-app.git/config"
	"github.com/ppeymann/top-app.git/docs"
	"github.com/ppeymann/top-app.git/env"
	"github.com/redis/go-redis/v9"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	Router        *gin.Engine
	Config        *config.Configuration
	Logger        kitlog.Logger
	paseto        auth.TokenMaker
	instrumenting serviceInstrumenting
	redis         *redis.Client
}

// EnvMode specified the running env 'release' represents production mode and ” represents development.
// it depended on gin GIN_MODE env for unifying and simplicity of setting.
var EnvMode = ""

func NewServer(logger kitlog.Logger, conf *config.Configuration, redis *redis.Client) *Server {
	svr := &Server{
		Logger:        logger,
		Config:        conf,
		instrumenting: newServiceInstrumenting(),
		redis:         redis,
	}

	router := gin.New()
	router.Use(gin.Recovery())

	EnvMode = os.Getenv("GIN_MODE")

	// setting swagger info if not in production mode
	if env.GetEnv("SWAGGER_ENABLED", "false") == "true" {
		docs.SwaggerInfo.Title = fmt.Sprintf("OTP App Backend [ AuthMode: %s ]", "Paseto")
		docs.SwaggerInfo.Description = "The Swagger Documentation For OtpApps Backend API server."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = env.GetEnv("HOST_URL", "localhost:8080")
		docs.SwaggerInfo.BasePath = "/api/v1"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	}

	// binding global
	router.Use(svr.metrics())

	// enable api rate limit if expected and enabled in config file
	if conf.RateLimit.Enabled {
		router.Use(svr.redisRateLimit())
	}

	if env.GetEnv("CORS_ENABLE", "false") == "true" {
		router.Use(svr.cors())
	}

	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Fatalln(err)
	}

	svr.Router = router

	svr.Router.GET("/metrics", svr.prometheus())

	svr.paseto, err = auth.NewPasetoMaker(env.GetEnv("JWT", ""))
	if err != nil {
		log.Fatal(err)
	}

	return svr
}

func (s *Server) Listen() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	if env.GetEnv("SWAGGER_ENABLED", "false") == "true" {
		s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	srv := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		Addr:              fmt.Sprintf("%s:%d", s.Config.Listener.Host, s.Config.Listener.Port),
		Handler:           s.Router,
	}

	// start https server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http listener stopped : %s", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully otp app server, press Ctrl+C again to force")

	// The context is used to inform the server it has 30 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}

	log.Println("otp service exiting")
}
