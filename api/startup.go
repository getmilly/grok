package api

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/myheartz/grok/logging"
	"github.com/sarulabs/di"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/swaggo/swag"
)

//Server wraps API configurations.
type Server struct {
	Engine   *gin.Engine
	Settings *Settings

	DIBuilder *di.Builder
	Container di.Container

	Healthz HealthzHandler

	router      *gin.RouterGroup
	controllers []string

	mu sync.Mutex
}

var (
	containerKey       = "di-container"
	errContainerNotSet = errors.New("container not set in request scope")
)

//Configure creates a new API server
func Configure(generator SettingGenerator, healthz HealthzHandler) *Server {
	server := &Server{}
	server.Settings = generator()
	server.Healthz = healthz

	logging.LogWithApplication(server.Settings.ApplicationName)

	builder, err := di.NewBuilder()

	if err != nil {
		panic(err)
	}

	server.DIBuilder = builder
	server.Engine = gin.New()
	server.Engine.Use(LogMiddleware())
	server.Engine.Use(gin.Recovery())

	server.Engine.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	server.router = server.Engine.Group(server.Settings.BasePath)

	server.router.Use(server.containerHandler())
	server.router.GET("/healthz", server.healtz())

	doc := NewSwaggerDoc(server.Settings.SwaggerPath)
	swag.Register(swag.Name, doc)

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
	server.mu.Lock()
	defer server.mu.Unlock()
	return server.DIBuilder.Add(def)
}

//AddController register a new controller to be added to routes.
func (server *Server) AddController(def di.Def) error {
	server.mu.Lock()
	defer server.mu.Unlock()
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
