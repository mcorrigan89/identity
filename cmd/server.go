package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func (server *appServer) serve(routes http.Handler) error {

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	routesHandler := h2c.NewHandler(routes, &http2.Server{})

	srv := &http.Server{
		Addr:        fmt.Sprintf(":%d", server.config.Port),
		Handler:     routesHandler,
		IdleTimeout: time.Minute,
		// ReadTimeout:  5 * time.Second,
		// WriteTimeout: 10 * time.Second,
		TLSConfig: tlsConfig,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		server.logger.Info().Str("signal", s.String()).Msg("caught signal")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			server.logger.Err(err).Ctx(ctx).Msg("failed to shutdown server")
			shutdownError <- err
		}

		server.logger.Info().Str("addr", srv.Addr).Msg("completing background tasks")

		server.wg.Wait()
		shutdownError <- nil
	}()

	server.logger.Info().Str("addr", srv.Addr).Str("env", server.config.Env).Msg("listening on")

	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	server.logger.Info().Str("addr", srv.Addr).Msg("stopped server")

	return nil
}
