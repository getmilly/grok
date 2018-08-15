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
	Engine   *gin.Engine
	Settings *Settings

	DIBuilder *di.Builder
	Container di.Container

	router      *gin.RouterGroup
	controllers []string
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

	server.router = server.Engine.Group(server.Settings.BasePath)

	if server.Settings.SwaggerPath == "" {
		panic("Swagger path is needed.")
	}

	server.router.Use(server.containerHandler())
	server.router.Use(server.healtz())
	server.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if server.Settings.Authorize {
		server.router.Use(AuthMiddleware(
			NewAuthService(
				server.Settings.Authorization.JwksURI,
				server.Settings.Authorization.Issuer,
				server.Settings.Authorization.Audience,
			),
		))
	}

	return server
}

//AddDependency register a new dependency in DI container.
func (server *Server) AddDependency(def di.Def) error {
	return server.DIBuilder.Add(def)
}

//AddController register a new controller to be added to routes.
func (server *Server) AddController(def di.Def) error {
	server.controllers = append(server.controllers, def.Name)
	return server.DIBuilder.Add(def)
}

//Run starts the server.
func (server *Server) Run() {
	server.Container = server.DIBuilder.Build()

	for _, ctrl := range server.extractControllers() {
		ctrl.RegisterRoutes(server.router)
	}

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

func (server *Server) extractControllers() []Controller {
	var ctrls []Controller

	for _, name := range server.controllers {
		def := server.Container.Get(name)
		ctrl, ok := def.(Controller)

		if !ok {
			panic("Defs added in AddController must implements Controller")
		}

		ctrls = append(ctrls, ctrl)
	}

	return ctrls
}
