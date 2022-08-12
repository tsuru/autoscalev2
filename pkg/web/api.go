// Copyright 2022 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	once sync.Once
	e    *echo.Echo
}

func (s *Server) Start(address string) error {
	s.once.Do(s.initWebServer)
	return s.e.Start(address)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.once.Do(s.initWebServer)
	return s.e.Shutdown(ctx)
}

func (s *Server) initWebServer() { s.e = newEcho() }

func newEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover(), middleware.Logger())

	e.GET("/resources/plans", listPlans)

	e.POST("/resources", createInstance)
	e.GET("/resources/{name}", getInstance)
	e.PUT("/resources/{name}", updateInstance)
	e.DELETE("/resources/{name}", removeInstance)
	e.GET("/resources/{name}/status", getInstanceStatus)

	e.POST("/resources/{name}/bind", bindUnit)
	e.DELETE("/resources/{name}/bind", unbindUnit)
	e.POST("/resources/{name}/app-bind", bindApp)
	e.DELETE("/resources/{name}/app-bind", unbindApp)

	return e
}

func listPlans(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func createInstance(c echo.Context) error {
	return errors.New("not implemented yet")
}

func getInstance(c echo.Context) error {
	return errors.New("not implemented yet")
}

func updateInstance(c echo.Context) error {
	return errors.New("not implemented yet")
}

func removeInstance(c echo.Context) error {
	return errors.New("not implemented yet")
}

func getInstanceStatus(c echo.Context) error {
	return errors.New("not implemented yet")
}

func bindUnit(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func unbindUnit(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func bindApp(c echo.Context) error {
	return errors.New("not implemented yet")
}

func unbindApp(c echo.Context) error {
	return errors.New("not implemented yet")
}
