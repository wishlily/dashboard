package main

import (
	"net/http"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wishlily/dashboard/server/config"
	"github.com/wishlily/dashboard/server/controllers"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{Asset, AssetDir, AssetInfo, root}
	return &binaryFileSystem{
		fs,
	}
}

func main() {
	config.Parse()

	// set log level
	level := strings.ToLower(viper.GetString("log"))
	switch level {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default: // info
		log.SetLevel(log.InfoLevel)
	}
	if log.GetLevel() < log.InfoLevel {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// We can't use router.Static method to use '/' for static files.
	// see https://github.com/gin-gonic/gin/issues/75
	// r.StaticFS("/", assetFS())
	r.Use(static.Serve("/", BinaryFileSystem("assets")))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// add routes
	if err := controllers.Finance.Init(); err != nil {
		panic(err)
	}
	r.GET("/api/finance/record", controllers.Finance.Records)
	r.POST("/api/finance/record", controllers.Finance.Record)
	r.GET("/api/finance/account", controllers.Finance.Accounts)
	r.POST("/api/finance/account", controllers.Finance.Account)

	ip := viper.GetString("ip")
	port := viper.GetString("port")
	r.Run(ip + ":" + port)
}
