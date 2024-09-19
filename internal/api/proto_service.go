package api

import (
	"net/http"
	"sync"

	"connectrpc.com/grpcreflect"
	identityv1connect "github.com/mcorrigan89/identity/internal/api/serviceapis/identity/v1/identityv1connect"
	"github.com/mcorrigan89/identity/internal/config"
	"github.com/mcorrigan89/identity/internal/services"
	"github.com/rs/zerolog"
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

func (s *ProtoServer) Handle(r *http.ServeMux) {

	reflector := grpcreflect.NewStaticReflector(
		"serviceapis.identity.v1.IdentityService",
	)

	reflectPath, reflectHandler := grpcreflect.NewHandlerV1(reflector)
	r.Handle(reflectPath, reflectHandler)
	reflectPathAlpha, reflectHandlerAlpha := grpcreflect.NewHandlerV1Alpha(reflector)
	r.Handle(reflectPathAlpha, reflectHandlerAlpha)

	scheduleV1Path, scheduleV1Handle := identityv1connect.NewIdentityServiceHandler(s.identityV1Server)
	r.Handle(scheduleV1Path, scheduleV1Handle)
}
