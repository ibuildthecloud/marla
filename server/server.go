package server

import (
	"fmt"
	"net"
	"net/http"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/server/router"
	"github.com/docker/docker/api/server/router/container"
	"github.com/docker/docker/api/server/router/image"
	"github.com/docker/docker/cmd/dockerd/hack"
	"github.com/docker/docker/pkg/listeners"
	"github.com/docker/docker/runconfig"
	"github.com/gorilla/mux"
	"github.com/ibuildthecloud/marla/daemon"
)

const (
	// TODO: Configurable
	socketGroup = "docker"
	socket      = "/var/run/marla.sock"
)

// versionMatcher defines a variable matcher to be parsed by the router
// when a request is about to be served.
const versionMatcher = "/v{version:[0-9.]+}"

type Server struct {
	daemon   *daemon.Daemon
	listener net.Listener
	setupErr error
	router   *mux.Router
}

func New(daemon *daemon.Daemon) (*Server, error) {
	s := &Server{
		daemon: daemon,
	}
	s.setupRoutes()
	return s, nil
}

func (s *Server) setupRoutes() {
	decoder := runconfig.ContainerDecoder{}

	routers := []router.Router{
		container.NewRouter(s.daemon, decoder),
		image.NewRouter(s.daemon, decoder),
	}

	m := mux.NewRouter()

	logrus.Debugf("Registering routers")
	for _, apiRouter := range routers {
		for _, r := range apiRouter.Routes() {
			f := s.makeHTTPHandler(r.Handler())

			logrus.Debugf("Registering %s, %s", r.Method(), r.Path())
			m.Path(versionMatcher + r.Path()).Methods(r.Method()).Handler(f)
			m.Path(r.Path()).Methods(r.Method()).Handler(f)
		}
	}

	//notFoundHandler := //httputils.MakeErrorHandler(errors.NewRequestNotFoundError(fmt.Errorf("!! page not found")))
	notFoundHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println("!!!!!", req.URL)
	})
	m.HandleFunc(versionMatcher+"/{path:.*}", notFoundHandler)
	m.NotFoundHandler = notFoundHandler

	s.router = m
}

func (s *Server) makeHTTPHandler(handler httputils.APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		vars := mux.Vars(r)
		if vars == nil {
			vars = make(map[string]string)
		}

		if err := handler(ctx, w, r, vars); err != nil {
			logrus.Errorf("Handler for %s %s returned error: %v", r.Method, r.URL.Path, err)
			httputils.MakeErrorHandler(err)(w, r)
		}
	}
}

func (s *Server) Listen() error {
	ls, err := listeners.Init("unix", socket, socketGroup, nil)
	if err != nil {
		return err
	}

	s.listener = &hack.MalformedHostHeaderOverride{ls[0]}
	return nil
}

func (s *Server) Serve() error {
	logrus.Infof("API listen on %s", s.listener.Addr())
	return http.Serve(s.listener, s.router)
}

func (s *Server) ListenAndServe() error {
	if err := s.Listen(); err != nil {
		return err
	}

	return s.Serve()
}
