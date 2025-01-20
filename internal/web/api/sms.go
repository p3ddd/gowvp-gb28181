// Code generated by gowebx, DO AVOID EDIT.
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gowvp/gb28181/internal/conf"
	"github.com/gowvp/gb28181/internal/core/sms"
	"github.com/gowvp/gb28181/internal/core/sms/store/smsdb"
	"github.com/ixugo/goweb/pkg/web"
	"gorm.io/gorm"
)

type SmsAPI struct {
	smsCore sms.Core
}

func NewSMSCore(db *gorm.DB, cfg *conf.Bootstrap) sms.Core {
	core := sms.NewCore(smsdb.NewDB(db).AutoMigrate(true))
	if err := core.Run(
		cfg.Media.IP,
		cfg.Media.Secret,
		cfg.Media.RTPPortRange,
		cfg.Media.WebHookIP,
		cfg.Media.HTTPPort,
		cfg.Server.HTTP.Port,
	); err != nil {
		panic(err)
	}
	return core
}

func NewSmsAPI(core sms.Core) SmsAPI {
	return SmsAPI{smsCore: core}
}

func registerSms(g gin.IRouter, api SmsAPI, handler ...gin.HandlerFunc) {
	{
		group := g.Group("/media_servers", handler...)
		group.GET("", web.WarpH(api.findMediaServer))
		group.GET("/:id", web.WarpH(api.getMediaServer))
		group.PUT("/:id", web.WarpH(api.editMediaServer))
		group.POST("", web.WarpH(api.addMediaServer))
		group.DELETE("/:id", web.WarpH(api.delMediaServer))
	}
}

// >>> mediaServer >>>>>>>>>>>>>>>>>>>>

func (a SmsAPI) findMediaServer(c *gin.Context, in *sms.FindMediaServerInput) (any, error) {
	items, total, err := a.smsCore.FindMediaServer(c.Request.Context(), in)
	return gin.H{"items": items, "total": total}, err
}

func (a SmsAPI) getMediaServer(c *gin.Context, _ *struct{}) (any, error) {
	mediaServerID := c.Param("id")
	return a.smsCore.GetMediaServer(c.Request.Context(), mediaServerID)
}

func (a SmsAPI) editMediaServer(c *gin.Context, in *sms.EditMediaServerInput) (any, error) {
	mediaServerID := c.Param("id")
	return a.smsCore.EditMediaServer(c.Request.Context(), in, mediaServerID)
}

func (a SmsAPI) addMediaServer(c *gin.Context, in *sms.AddMediaServerInput) (any, error) {
	return a.smsCore.AddMediaServer(c.Request.Context(), in)
}

func (a SmsAPI) delMediaServer(c *gin.Context, _ *struct{}) (any, error) {
	mediaServerID := c.Param("id")
	return a.smsCore.DelMediaServer(c.Request.Context(), mediaServerID)
}
