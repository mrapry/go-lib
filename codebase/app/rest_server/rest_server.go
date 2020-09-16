package restserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	echoMidd "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	graphqlserver "github.com/mrapry/go-lib/codebase/app/graphql_server"
	"github.com/mrapry/go-lib/codebase/factory"
	"github.com/mrapry/go-lib/config"
	"github.com/mrapry/go-lib/golibhelper"
	"github.com/mrapry/go-lib/logger"
	"github.com/mrapry/go-lib/tracer"
	"github.com/mrapry/go-lib/wrapper"
)

type restServer struct {
	serverEngine   *echo.Echo
	service        factory.ServiceFactory
	graphqlHandler graphqlserver.Handler
}

// NewServer create new REST server
func NewServer(service factory.ServiceFactory) factory.AppServerFactory {
	server := &restServer{
		serverEngine: echo.New(),
		service:      service,
	}

	// inject graphql handler, delete/comment this code if you want separate graphql server from echo rest server
	if config.BaseEnv().UseGraphQL {
		server.graphqlHandler = graphqlserver.NewHandler(service)
	}

	return server
}

func (h *restServer) Serve() {

	h.serverEngine.HTTPErrorHandler = wrapper.CustomHTTPErrorHandler
	h.serverEngine.Use(echoMidd.Recover(), echoMidd.CORS())

	h.serverEngine.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message":   fmt.Sprintf("Service %s up and running", h.service.Name()),
			"timestamp": time.Now().Format(time.RFC3339Nano),
		})
	})

	useSWAGGER, ok := os.LookupEnv("USE_SWAGGER")
	if !ok {
		panic("missing USE_SWAGGER environment")
	}
	statusSwagger, _ := strconv.ParseBool(useSWAGGER)

	if statusSwagger {
		fs := http.FileServer(http.Dir("/docs"))
		h.serverEngine.GET("/*", echo.WrapHandler(http.StripPrefix("/", fs)))
		h.serverEngine.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.URL(os.Getenv("SWAGGER_ADDRESS_DOC"))))
	}

	restRootPath := h.serverEngine.Group("",
		tracer.EchoRestTracerMiddleware, echoMidd.Logger(),
	)
	for _, m := range h.service.GetModules() {
		if h := m.RestHandler(); h != nil {
			h.Mount(restRootPath)
		}
	}

	if h.graphqlHandler != nil {
		h.serverEngine.POST("/graphql", echo.WrapHandler(h.graphqlHandler.ServeGraphQL()))
		h.serverEngine.GET("/graphql/playground", echo.WrapHandler(http.HandlerFunc(h.graphqlHandler.ServePlayground)))
		h.serverEngine.GET("/graphql/voyager", echo.WrapHandler(http.HandlerFunc(h.graphqlHandler.ServeVoyager)))
	}

	var routes strings.Builder
	httpRoutes := h.serverEngine.Routes()
	sort.Slice(httpRoutes, func(i, j int) bool {
		return httpRoutes[i].Path < httpRoutes[j].Path
	})
	for _, route := range httpRoutes {
		if !strings.Contains(route.Name, "(*Group)") {
			routes.WriteString(golibhelper.StringGreen(fmt.Sprintf("[REST-ROUTE] %-6s %-30s --> %s\n", route.Method, route.Path, route.Name)))
		}
	}
	fmt.Print(routes.String())

	h.serverEngine.HideBanner = true
	h.serverEngine.HidePort = true
	port := fmt.Sprintf(":%d", config.BaseEnv().RESTPort)
	fmt.Printf("\x1b[34;1m⇨ REST server run at port [::]%s\x1b[0m\n\n", port)
	if err := h.serverEngine.Start(port); err != nil {
		switch e := err.(type) {
		case *net.OpError:
			panic(e)
		}
	}
}

func (h *restServer) Shutdown(ctx context.Context) {
	deferFunc := logger.LogWithDefer("Stopping REST HTTP server...")
	defer deferFunc()

	if err := h.serverEngine.Shutdown(ctx); err != nil {
		panic(err)
	}
}
