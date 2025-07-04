package main

import (
	"log"
	"os"
	"pojok-baca-api/config"
	_ "pojok-baca-api/docs"
	"pojok-baca-api/handler"
	"pojok-baca-api/middleware"
	"pojok-baca-api/repository"
	"pojok-baca-api/service"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title           Pojok Baca API
// @version         1.0
// @description     This is an Pojok Baca API
// @termsOfService  http://swagger.io/terms/

// @contact.name   Tatakarai
// @contact.email  tatakarai@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	config.LoadEnv()
	db := config.DBInit()

	//if err := db.AutoMigrate(&model.User{}, &model.Book{}, &model.Rental{}, &model.DepositTransaction{}); err != nil {
	//	panic("Auto migrate fail : " + err.Error())
	//}

	log.Println("Auto migrate success")

	//USER
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	//BOOK
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewProductHandler(bookService)

	rentalRepo := repository.NewRentalRepository(db)
	rentalService := service.NewRentalService(rentalRepo)
	rentalHandler := handler.NewRentalHandler(rentalService, bookService, userService)

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	//group api
	api := e.Group("/api")

	//group users
	user := api.Group("/user")
	productGroup := api.Group("/products")
	rentalGroup := api.Group("/rentals")

	//User register & login
	user.POST("/register", userHandler.CreateUser)
	user.POST("/login", userHandler.Login)

	productGroup.GET("", bookHandler.GetBooks)
	productGroup.GET("/:id", bookHandler.GetBookByID)

	jwtSecret := os.Getenv("JWT_SECRET")

	user.Use(middleware.JWTMiddleware(jwtSecret))
	user.GET("/me", userHandler.GetDataByID)

	productGroup.Use(middleware.JWTMiddleware(jwtSecret))
	productGroup.POST("", bookHandler.CreateBook)
	productGroup.PUT("/:id", bookHandler.UpdateBookByID)
	productGroup.DELETE("/:id", bookHandler.DeleteBookByID)

	rentalGroup.Use(middleware.JWTMiddleware(jwtSecret))
	rentalGroup.POST("", rentalHandler.CreateRental)
	rentalGroup.GET("/report", rentalHandler.GetRentalByUserID)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
