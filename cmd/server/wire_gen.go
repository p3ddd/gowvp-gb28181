// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gowvp/gb28181/internal/conf"
	"github.com/gowvp/gb28181/internal/data"
	"github.com/gowvp/gb28181/internal/web/api"
	"github.com/gowvp/gb28181/pkg/gbs"
	"log/slog"
	"net/http"
)

// Injectors from wire.go:

func wireApp(bc *conf.Bootstrap, log *slog.Logger) (http.Handler, func(), error) {
	db, err := data.SetupDB(bc, log)
	if err != nil {
		return nil, nil, err
	}
	core := api.NewVersion(db)
	versionAPI := api.NewVersionAPI(core)
	smsCore := api.NewSMSCore(db, bc)
	smsAPI := api.NewSmsAPI(smsCore)
	uniqueidCore := api.NewUniqueID(db)
	mediaCore := api.NewMediaCore(db, uniqueidCore)
	webHookAPI := api.NewWebHookAPI(smsCore, mediaCore, bc)
	mediaAPI := api.NewMediaAPI(mediaCore, smsCore, bc)
	gb28181API := api.NewGb28181API(db, uniqueidCore)
	proxyAPI := api.NewProxyAPI(db, uniqueidCore)
	configAPI := api.NewConfigAPI(db, bc)
	gb28181 := api.NewGB28181(db, uniqueidCore)
	server, cleanup := gbs.NewServer(bc, gb28181, smsCore)
	usecase := &api.Usecase{
		Conf:       bc,
		DB:         db,
		Version:    versionAPI,
		SMSAPI:     smsAPI,
		WebHookAPI: webHookAPI,
		UniqueID:   uniqueidCore,
		MediaAPI:   mediaAPI,
		GB28181API: gb28181API,
		ProxyAPI:   proxyAPI,
		ConfigAPI:  configAPI,
		SipServer:  server,
	}
	handler := api.NewHTTPHandler(usecase)
	return handler, func() {
		cleanup()
	}, nil
}
