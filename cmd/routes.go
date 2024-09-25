package main

import (
	"net/http"
)

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	app.logger.Info().Ctx(ctx).Msg("/ping")
	w.Write([]byte("OK"))
}

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", app.ping)
	app.protoServer.Handle(mux, app.tracerProvider)

	return app.recoverPanic(app.enabledCORS(app.contextBuilder(mux)))
}
