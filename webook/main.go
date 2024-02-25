package main

import (
	"time"
	"webook/config"
	"webook/internal/repository/dao"
	"webook/internal/service/sms"
	"webook/internal/service/sms/localsms"
	"webook/internal/web/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisSession "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	server := InitWebServer()
	// db := initDB()

	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr: config.Config.Redis.Addr,
	// })

	// server := initWebServer()
	// codeSvc := initCodeSvc(redisClient)

	// initUserHandler(db, redisClient, codeSvc, server)

	// server := gin.Default()
	// server.GET("/hello", func(ctx *gin.Context) {
	// 	ctx.String(http.StatusOK, "hello, 启动成功了！")
	// })

	server.Run(":8081")
}

// func initUserHandler(db *gorm.DB, redisClient redis.Cmdable, codeSvc *service.CodeService, server *gin.Engine) {
// 	ud := dao.NewUserDao(db)
// 	uc := cache.NewUserCache(redisClient)
// 	ur := repository.NewUserRegistory(ud, uc)
// 	us := service.NewUserService(ur)
// 	hdl := web.NewUserHandler(us, codeSvc)
// 	hdl.RegisterRoutes(server)
// }

// func initCodeSvc(redisClient redis.Cmdable) *service.CodeService {
// 	cc := cache.NewCodeCache(redisClient)
// 	crepo := repository.NewCodeRepository(cc)
// 	return service.NewCodeService(crepo, initMemorySms())
// }

func initMemorySms() sms.Service {
	return localsms.NewService()
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}

	db = db.Debug()

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		// AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
		// AllowOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"x-jwt-token"},
	}))

	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr: config.Config.Redis.Addr,
	// })
	// redisLimiter := limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 100)

	// server.Use(ratelimit.NewBuilder(redisLimiter).Build())

	useJWT(server)
	// userSession(server)

	return server
}

func useJWT(server *gin.Engine) {
	login := middleware.LoginMiddlewareJWTBuilder{}
	server.Use(login.CheckLoginJWT())
}

func userSession(server *gin.Engine) {
	login := &middleware.LoginMiddlewareBuilder{}
	// store := cookie.NewStore([]byte("secret"))
	store, err := redisSession.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("x88Uf1yeybFEQA1yedOtpYJ9TibkH2vY"),
		[]byte("x88Uf1yeybFEQA1yedOtpYJ9TibkH2vX"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
