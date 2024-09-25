package api

import (
	"net/http"
	"sync"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	identityv1connect "github.com/mcorrigan89/identity/internal/api/serviceapis/identity/v1/identityv1connect"
	"github.com/mcorrigan89/identity/internal/config"
	"github.com/mcorrigan89/identity/internal/services"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/sdk/trace"
)

type ProtoServer struct {
	config   *config.Config
	wg       *sync.WaitGroup
	logger   *zerolog.Logger
	services *services.Services

	identityV1Server *IdentityServerV1
}

func NewProtoServer(cfg *config.Config, logger *zerolog.Logger, wg *sync.WaitGroup, services *services.Services) *ProtoServer {

	identityV1Server := newIdentityProtoUrlServer(cfg, logger, wg, services)

	return &ProtoServer{
		config:           cfg,
		wg:               wg,
		logger:           logger,
		services:         services,
		identityV1Server: identityV1Server,
	}
}

func (s *ProtoServer) Handle(r *http.ServeMux, tracerProvider *trace.TracerProvider) {

	reflector := grpcreflect.NewStaticReflector(
		"serviceapis.identity.v1.IdentityService",
	)

	reflectPath, reflectHandler := grpcreflect.NewHandlerV1(reflector)
	r.Handle(reflectPath, reflectHandler)
	reflectPathAlpha, reflectHandlerAlpha := grpcreflect.NewHandlerV1Alpha(reflector)
	r.Handle(reflectPathAlpha, reflectHandlerAlpha)

	otelInterceptor, err := otelconnect.NewInterceptor(
		otelconnect.WithTracerProvider(tracerProvider),
	)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Failed to create OpenTelemetry interceptor")
	}

	identityV1Path, identityV1Handle := identityv1connect.NewIdentityServiceHandler(s.identityV1Server, connect.WithInterceptors(otelInterceptor))
	r.Handle(identityV1Path, identityV1Handle)
}
