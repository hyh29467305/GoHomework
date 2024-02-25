//go:build wireinject

package main

import (
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//第三方依赖
		ioc.InitDB, ioc.InitRedis,

		//dao 部分
		dao.NewUserDao,

		//cache 部分
		cache.NewCodeCache, cache.NewUserCache,

		//repository 部分
		repository.NewCachedUserRegistory, repository.NewCachedCodeRepository,

		// Service 部分
		ioc.InitSMSService, service.NewCodeService, service.NewUserService,

		// handler 部分
		web.NewUserHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
