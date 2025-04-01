package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/mcorrigan89/identity/internal/infrastructure/middleware"
	"github.com/mcorrigan89/identity/internal/interfaces/http/handlers"
)

func NewRouter(mux *http.ServeMux, middleware middleware.Middleware, userHandler *handlers.UserHandler) http.Handler {

	api := humago.New(mux, huma.DefaultConfig("OpenMic API", "1.0.0"))

	// User routes
	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/user/{id}",
		Summary:     "Get a user by ID",
		Tags:        []string{"User"},
	}, userHandler.GetUserByID)

	huma.Register(api, huma.Operation{
		OperationID: "create-password-user",
		Method:      http.MethodPost,
		Path:        "/user/password",
		Summary:     "Create password user",
		Tags:        []string{"User"},
	}, userHandler.CreatePasswordUser)

	huma.Register(api, huma.Operation{
		OperationID: "update-user",
		Method:      http.MethodPut,
		Path:        "/user/{id}",
		Summary:     "Update user",
		Tags:        []string{"User"},
	}, userHandler.UpdateUser)

	huma.Register(api, huma.Operation{
		OperationID: "oauth-link",
		Method:      http.MethodGet,
		Path:        "/user/oauth/{provider}",
		Summary:     "Get OAuth link",
		Tags:        []string{"User"},
	}, userHandler.OAuthLinkURL)

	huma.Register(api, huma.Operation{
		OperationID: "oauth-login",
		Method:      http.MethodPost,
		Path:        "/user/oauth/{provider}",
		Summary:     "OAuth login",
		Tags:        []string{"User"},
	}, userHandler.OAuthLogin)

	return middleware.RecoverPanic(middleware.EnabledCORS(middleware.ContextBuilder(mux)))
}
