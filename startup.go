package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sarulabs/di"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

//Server wraps API configurations.
type Server struct {
	DIBuilder *di.Builder
	Container di.Container
	Engine    *gin.Engine
	Settings  *Settings
	Router    *gin.RouterGroup
}

//Settings stores some configs about how the API will woks.
type Settings struct {
	Host          string
	EnvsPath      string
	Authorize     bool
	Authorization struct {
		JwksURI  string
		Issuer   string
		Audience []string
	}
	SwaggerPath string
	BasePath    string
	Healthz     func() interface{}
}

//SettingGenerator creates a instance of Settings.
type SettingGenerator func() *Settings

var (
	containerKey       = "di-container"
	errContainerNotSet = errors.New("container not set in request scope")
)

//Configure creates a new API server
func Configure(generator SettingGenerator) *Server {
	server := &Server{}
	server.Settings = generator()

	builder, err := di.NewBuilder()

	if err != nil {
		panic(err)
	}

	server.DIBuilder = builder
	server.Engine = gin.Default()

	server.Engine.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	server.Router = server.Engine.Group(server.Settings.BasePath)

	if server.Settings.SwaggerPath == "" {
		panic("Swagger path is needed.")
	}

	server.Router.Use(server.containerHandler())
	server.Router.Use(server.healtz())
	server.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if server.Settings.Authorize {
		server.Router.Use(AuthMiddleware(
			NewAuthService(
				server.Settings.Authorization.JwksURI,
				server.Settings.Authorization.Issuer,
				server.Settings.Authorization.Audience,
			),
		))
	}

	return server
}

//Run starts the server.
func (server *Server) Run() {
	server.Container = server.DIBuilder.Build()
	err := server.Engine.Run(server.Settings.Host)

	if err != nil {
		panic(err)
	}
}

//Container return DI Container defined in request scope.
func Container(c *gin.Context) (di.Container, error) {
	container, ok := c.Get(containerKey)

	if !ok {
		return nil, errContainerNotSet
	}

	di, ok := container.(di.Container)

	if !ok {
		return nil, errContainerNotSet
	}

	return di, nil
}

//DotEnv generates settings using environment variables.
func DotEnv(files ...string) SettingGenerator {
	err := godotenv.Load(files...)

	if err != nil {
		panic(err)
	}

	return func() *Settings {
		//TODO: define pattern and build envs here.
		return &Settings{}
	}
}

func (server *Server) containerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if server.Container != nil {
			container, err := server.Container.SubContainer()

			if err != nil {
				panic(err)
			}

			defer container.Delete()
			c.Set(containerKey, container)
		}
		c.Next()
	}
}

func (server *Server) healtz() gin.HandlerFunc {
	return func(c *gin.Context) {
		var healthz interface{}
		if server.Settings.Healthz != nil {
			healthz = server.Settings.Healthz()
		}
		c.JSON(http.StatusOK, healthz)
	}
}
