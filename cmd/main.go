package main

import (
	"github.com/ell1jah/bmstu_web/cmd/server"
	jwtManager "github.com/ell1jah/bmstu_web/internal/pkg/jwt"
	postDelivery "github.com/ell1jah/bmstu_web/internal/post/delivery"
	postLogic "github.com/ell1jah/bmstu_web/internal/post/logic"
	postRepository "github.com/ell1jah/bmstu_web/internal/post/repository"
	rateRepository "github.com/ell1jah/bmstu_web/internal/rate/repository"
	userDelivery "github.com/ell1jah/bmstu_web/internal/user/delivery"
	userLogic "github.com/ell1jah/bmstu_web/internal/user/logic"
	userRepository "github.com/ell1jah/bmstu_web/internal/user/repository"

	_ "github.com/GoAdminGroup/go-admin/adapter/echo"
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/examples/datamodel"
	"github.com/GoAdminGroup/go-admin/modules/config"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/themes/adminlte"
	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-contrib/prometheus"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var testCfgPg = postgres.Config{DSN: "host=localhost user=postgres password=postgres port=13080"}

var prodCfgPg = postgres.Config{DSN: "host=cloth_pg user=postgres password=postgres port=5432"}
var jwtKey = []byte("sdoBsm#vpw,vsdS3902F,dvd]s")

func initAdmin(e *echo.Echo) {
	eng := engine.Default()

	cfg := config.Config{
		Databases: config.DatabaseList{},
		UrlPrefix: "admin", // The url prefix of the website.
		// Store must be set and guaranteed to have write access, otherwise new administrator users cannot be added.
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language: language.EN,
		Theme:    adminlte.Adminlte.ThemeName,
		IndexUrl: "/",
		Debug:    true,
	}

	cfg.Databases.Add("default", config.Database{
		Host:   "cloth_pg",
		Port:   "5432",
		User:   "postgres",
		Pwd:    "postgres",
		Name:   "postgres",
		Driver: "postgresql",
	})

	// Add configuration and plugins, use the Use method to mount to the web framework.
	_ = eng.AddConfig(&cfg).Use(e)

	eng.HTML("GET", "/admin", datamodel.GetContent)
}

func main() {
	govalidator.SetFieldsRequiredByDefault(true)

	db, err := gorm.Open(postgres.New(prodCfgPg),
		&gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("postgres connect success")

	userRepo := userRepository.NewPgRepo(db)
	postRepo := postRepository.NewPgRepo(db)
	rateRepo := rateRepository.NewPgRepo(db)

	userLogic := userLogic.NewLogic(userRepo)
	postLogic := postLogic.NewLogic(postRepo, userRepo, rateRepo)

	e := echo.New()
	initAdmin(e)

	e.Logger.SetHeader(`time=${time_rfc3339} level=${level} prefix=${prefix} ` +
		`file=${short_file} line=${line} message:`)
	e.Logger.SetLevel(log.INFO)

	p := prometheus.NewPrometheus("echo", nil)
	p.MetricsPath = "/prometheus"
	p.SetMetricsPath(e)
	p.Use(e)

	e.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format: `time=${time_custom} remote_ip=${remote_ip} ` +
			`host=${host} method=${method} uri=${uri} user_agent=${user_agent} ` +
			`status=${status} error="${error}" ` +
			`bytes_in=${bytes_in} bytes_out=${bytes_out}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	e.Use(echoMiddleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	sessionManager := jwtManager.NewJWTSessionsManager(jwtKey, jwt.SigningMethodHS256)

	authMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey:    jwtKey,
		SigningMethod: jwt.SigningMethodHS256.Alg(),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtManager.Claims)
		},
	})

	userDelivery.NewHandler(userLogic, sessionManager).SetRoutes(e, authMiddleware)
	postDelivery.NewHandler(postLogic).SetRoutes(e, authMiddleware)

	s := server.NewServer(e)
	if err := s.Start(); err != nil {
		e.Logger.Fatal(err)
	}
}
