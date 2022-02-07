package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-coord/controller"
	_ "github.com/shark/minigame-coord/docs"
	"github.com/shark/minigame-coord/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(addr string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.Cors())
	accountController := controller.NewAccountController()
	accountGroup := r.Group("/account")
	{
		accountGroup.POST("/login", accountController.Login)
		accountGroup.GET("/test", accountController.Test)
	}
	r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "SWAGGER"))
	if conf.Ini.Coord.UseTLS {
		if err := r.RunTLS(addr, conf.Ini.Coord.TLSPem, conf.Ini.Coord.TLSKey); err != nil {
			log.Fatalln(err)
		}
	} else {
		if err := r.Run(addr); err != nil {
			log.Fatalln(err)
		}
	}

}
