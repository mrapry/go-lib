package graphqlserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"

	"github.com/mrapry/go-lib/codebase/factory"
	"github.com/mrapry/go-lib/config"
	"github.com/mrapry/go-lib/golibhelper"
	"github.com/mrapry/go-lib/golibshared"
	"github.com/mrapry/go-lib/logger"
	"github.com/mrapry/go-lib/tracer"

	"github.com/graph-gophers/graphql-go"
)

type graphqlServer struct {
	httpEngine  *http.Server
	httpHandler Handler
	service     factory.ServiceFactory
}

// NewServer create new GraphQL server
func NewServer(service factory.ServiceFactory) factory.AppServerFactory {
	return &graphqlServer{
		httpHandler: NewHandler(service),
	}
}

func (s *graphqlServer) Serve() {
	s.httpEngine = new(http.Server)

	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", s.httpHandler.ServeGraphQL)
	mux.HandleFunc("/graphql/playground", s.httpHandler.ServePlayground)

	s.httpEngine.Addr = fmt.Sprintf(":%d", config.BaseEnv().GraphQLPort)
	s.httpEngine.Handler = mux

	fmt.Println(golibhelper.StringYellow("[GraphQL] endpoint: /graphql"))
	fmt.Println(golibhelper.StringYellow("[GraphQL] playground: /graphql/playground"))
	fmt.Printf("\x1b[34;1m⇨ GraphQL server run at port [::]%s\x1b[0m\n\n", s.httpEngine.Addr)
	if err := s.httpEngine.ListenAndServe(); err != nil {
		switch e := err.(type) {
		case *net.OpError:
			panic(e)
		}
	}
}

func (s *graphqlServer) Shutdown(ctx context.Context) {
	deferFunc := logger.LogWithDefer("Stopping GraphQL HTTP server...")
	defer deferFunc()

	s.httpEngine.Shutdown(ctx)
}

// Handler interface
type Handler interface {
	ServeGraphQL(resp http.ResponseWriter, req *http.Request)
	ServePlayground(resp http.ResponseWriter, req *http.Request)
}

// NewHandler for create public graphql handler (maybe inject to rest handler)
func NewHandler(service factory.ServiceFactory) Handler {
	resolverModules := make(map[string]interface{})
	var resolverFields []reflect.StructField // for creating dynamic struct
	for _, m := range service.GetModules() {
		if name, handler := m.GraphQLHandler(); handler != nil {
			resolverModules[name] = handler
			resolverFields = append(resolverFields, reflect.StructField{
				Name: name,
				Type: reflect.TypeOf(handler),
			})
		}
	}

	resolverVal := reflect.New(reflect.StructOf(resolverFields)).Elem()
	for k, v := range resolverModules {
		val := resolverVal.FieldByName(k)
		val.Set(reflect.ValueOf(v))
	}

	gqlSchema := golibhelper.LoadAllFile(config.BaseEnv().GraphQLSchemaDir, ".graphql")
	resolver := resolverVal.Addr().Interface()

	schema := graphql.MustParseSchema(string(gqlSchema), resolver,
		graphql.UseStringDescriptions(),
		graphql.UseFieldResolvers(),
		graphql.Tracer(&tracer.GraphQLTracer{}))
	return &handlerImpl{
		schema: schema,
	}
}

type handlerImpl struct {
	schema *graphql.Schema
}

func (s *handlerImpl) ServeGraphQL(resp http.ResponseWriter, req *http.Request) {

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &params); err != nil {
		params.Query = string(body)
	}

	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = req.Header.Get("X-Real-IP")
		if ip == "" {
			ip, _, _ = net.SplitHostPort(req.RemoteAddr)
		}
	}
	req.Header.Set("X-Real-IP", ip)

	ctx := context.WithValue(req.Context(), golibshared.ContextKey("headers"), req.Header)
	response := s.schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(responseJSON)
}

func (s *handlerImpl) ServePlayground(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte(`<!DOCTYPE html>
	<html>
	<head>
		<meta charset=utf-8/>
		<meta name="viewport" content="user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, minimal-ui">
		<link rel="shortcut icon" href="https://graphcool-playground.netlify.com/favicon.png">
		<link rel="stylesheet" href="//cdn.jsdelivr.net/npm/graphql-playground-react@1.7.8/build/static/css/index.css"/>
		<link rel="shortcut icon" href="//cdn.jsdelivr.net/npm/graphql-playground-react@1.7.8/build/favicon.png"/>
		<script src="//cdn.jsdelivr.net/npm/graphql-playground-react@1.7.8/build/static/js/middleware.js"></script>
		<title>Playground</title>
	</head>
	<body>
	<style type="text/css">
		html { font-family: "Open Sans", sans-serif; overflow: hidden; }
		body { margin: 0; background: #172a3a; }
	</style>
	<div id="root"/>
	<script type="text/javascript">
		window.addEventListener('load', function (event) {
			const root = document.getElementById('root');
			root.classList.add('playgroundIn');
			const wsProto = location.protocol == 'https:' ? 'wss:' : 'ws:'
			GraphQLPlayground.init(root, {
				endpoint: location.protocol + '//' + location.host + '/graphql',
				subscriptionsEndpoint: wsProto + '//' + location.host + '/graphql',
				settings: {
					'request.credentials': 'same-origin'
				}
			})
		})
	</script>
	</body>
	</html>`))
}
