package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/budhip/common/postgre"
	"github.com/gofiber/fiber/v2"

	"google.golang.org/grpc"

	_ "github.com/lib/pq"

	commonGrpc "github.com/budhip/common/grpc"
	coreConfig "github.com/budhip/env-config"
	viperConfig "github.com/budhip/env-config/viper"
	deliveryGrpc "github.com/budhip/example-user/delivery/grpc"
	handler "github.com/budhip/example-user/delivery/http"

	conf "github.com/budhip/example-user/config"


	postgreUserRepo "github.com/budhip/example-user/repository"
	userSrv "github.com/budhip/example-user/service"
)

func getConfig() (coreConfig.Config, error) {
	return viperConfig.NewConfig("amarthacore", "config.json")
}

func main() {
	sugarLogger := conf.InitLogger()

	defer sugarLogger.Sync()

	config, err := getConfig()
	if err != nil {
		sugarLogger.Errorf("can not read config.json: %s", err)
		return
	}


	dbHost := config.GetString(`database.host`)
	dbPort := config.GetString(`database.port`)
	dbUser := config.GetString(`database.user`)
	dbPass := config.GetString(`database.pass`)
	dbName := config.GetString(`database.name`)

	dbConfig := postgre.Config{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPass,
		Name:     dbName,
	}

	db, err := postgre.DB(dbConfig)
	if err != nil {
		sugarLogger.Errorf("failed connect to database: %s", err)
		return
	}
	defer db.Close()

	sugarLogger.Infof("Successfully connected to database")

	grpcAddr := config.GetString("server.user_address")

	snow := conf.NewSnowflake()

	// service
	upr := postgreUserRepo.NewPostgreRepository(db)
	userService := userSrv.NewUserService(upr, snow)

	// grpc
	pbServer := grpc.NewServer(commonGrpc.WithDefault()...)
	deliveryGrpc.NewUserServerGRPC(pbServer, userService)

	// HTTP
	// Creates a new Fiber instance.
	app := fiber.New(fiber.Config{
		AppName:      "Fiber Example User Clean Architecture",
		ServerHeader: "Example User",
	})

	httpAddr := config.GetString("server.user_address_http")

	api := app.Group("/api")
	// Prepare our endpoints for the API.
	handler.NewUserHandler(api.Group("/v1/users"), userService)

	// Prepare an endpoint for 'Not Found'.
	app.All("*", func(c *fiber.Ctx) error {
		errorMessage := fmt.Sprintf("Route '%s' does not exist in this API!", c.OriginalURL())

		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"status":  "fail",
			"message": errorMessage,
		})
	})

	go func() {
		commonGrpc.Serve(grpcAddr, pbServer)
	}()

	go func() {
		sugarLogger.Fatal(app.Listen(httpAddr))
	}()

	sugarLogger.Infof("gRPC server started. Listening on port: %s", grpcAddr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	sugarLogger.Infof("All server stopped!")
}
