package http

import (
	"github.com/go-chi/chi/v5"
)

// Every Handler layer should implement this interface, to make it easier for registering to the router.
type Handler interface {
	// This method is used to define the routes handled by the specific handler.
	// It takes a chi.Router as a parameter, allowing the handler to register its routes with the router.
	// Do not register prefix url pattern eg: "/users" on the handler func.
	Routes(r chi.Router)
	// This method should return the URL prefix for the handler's routes.
	// For example, if the handler handles routes related to users and all its routes should start with /users,
	// So this method would return string "/users".
	URLPattern() string
}

type RoutesConfiguration struct {
	router   *chi.Mux
	handlers []Handler
}

func SetupAPIRoutes(r *chi.Mux) *RoutesConfiguration {
	return &RoutesConfiguration{
		router: r,
	}
}

// Register All handler, should use this only once because it's directly assigning handlers field.
func (c *RoutesConfiguration) RegisterHandlers(handlers ...Handler) {
	c.handlers = handlers
}

func (c *RoutesConfiguration) Init() {
	c.router.Route("/v1", func(subRouter chi.Router) {
		// Iterate handlers and register to the router
		for _, eachHandler := range c.handlers {
			subRouter.Route(eachHandler.URLPattern(), func(r chi.Router) {
				eachHandler.Routes(r)
			})
		}
	})
}
