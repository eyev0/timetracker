package http

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	config "github.com/spf13/viper"

	"github.com/eyev0/timetracker/internal/http/controllers"
	"github.com/eyev0/timetracker/internal/http/middleware"
	"github.com/eyev0/timetracker/internal/log"
)

func InitHttpServer() {
	server := gin.Default()

	addr := fmt.Sprintf("%s:%d", config.GetString("SERVER_IP"), config.GetInt("SERVER_PORT"))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	apiRouter := server.Group("/api")
	apiRouter.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Timetracker backend"})
	})

	authRouter := apiRouter.Group("/auth")
	authRouter.POST("/register", controllers.SignUpUser)
	authRouter.POST("/login", controllers.SignInUser)
	authRouter.GET("/logout", middleware.DeserializeUser, controllers.LogoutUser)

	googleTokenRouter := authRouter.Group("/google_token")
	googleTokenRouter.POST("/refresh", middleware.DeserializeUser, controllers.RefreshToken)
	googleTokenRouter.GET("/", middleware.DeserializeUser, controllers.GetToken)

	apiRouter.GET("/sessions/oauth/google", controllers.GoogleOAuth)

	usersRouter := apiRouter.Group("/users", middleware.DeserializeUser)
	usersRouter.GET("/me", controllers.GetMe)
	usersRouter.PATCH("/settings", controllers.PatchSettings)

	entriesRouter := apiRouter.Group("/entries", middleware.DeserializeUser)
	entriesRouter.POST("/create", controllers.CreateEntry)
	entriesRouter.POST("/update", controllers.UpdateEntry)
	entriesRouter.GET("/current", controllers.GetEntry)

	// router.StaticFS("/images", http.Dir("public"))
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Route Not Found"})
	})

	log.Logger.Infof("Starting http server on %s", addr)
	log.Logger.Fatal(server.Run(addr))
}
