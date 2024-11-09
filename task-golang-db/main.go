package main

import (
	"log"
	"os"
	"task-golang-db/handlers"
	"task-golang-db/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Initialize the database
	db := NewDatabase()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get DB from GORM:", err)
	}
	defer sqlDB.Close()

	// Get signing key from environment
	signingKey := os.Getenv("SIGNING_KEY")

	// Initialize Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allowed origin for your Vue app
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // Allow "Authorization" header
		AllowCredentials: true,
	}))

	// Homepage route
	router.GET("/homepage", handlers.Homepage(db))

	// Grouping routes under /auth
	authhandlers := handlers.NewAuth(db, []byte(signingKey))
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authhandlers.Login)
		authRoutes.POST("/upsert", authhandlers.Upsert)
	}

	// Grouping routes under /account
	accounthandlers := handlers.NewAccount(db)
	accountRoutes := router.Group("/account")
	{
		accountRoutes.POST("/create", accounthandlers.Create)
		accountRoutes.GET("/read/:id", accounthandlers.Read)
		accountRoutes.PATCH("/update/:id", accounthandlers.Update)
		accountRoutes.DELETE("/delete/:id", accounthandlers.Delete)
		accountRoutes.GET("/list", accounthandlers.List)
		accountRoutes.POST("/topup", accounthandlers.TopUp)
		accountRoutes.GET("/my", middleware.AuthMiddleware(signingKey), accounthandlers.My)
		accountRoutes.GET("/balance", middleware.AuthMiddleware(signingKey), accounthandlers.Balance)
		accountRoutes.POST("/transfer", middleware.AuthMiddleware(signingKey), accounthandlers.Transfer)
		accountRoutes.GET("/mutation", middleware.AuthMiddleware(signingKey), accounthandlers.Mutation)
	}

	// Grouping routes under /transaction-category
	transCatHandlers := handlers.NewTransCat(db)
	transCatRoutes := router.Group("/transaction-category")
	{
		transCatRoutes.POST("/create", transCatHandlers.Create)
		transCatRoutes.GET("/read/:id", transCatHandlers.Read)
		transCatRoutes.PATCH("/update/:id", transCatHandlers.Update)
		transCatRoutes.DELETE("/delete/:id", transCatHandlers.Delete)
		transCatRoutes.GET("/list", transCatHandlers.List)
		transCatRoutes.GET("/my", middleware.AuthMiddleware(signingKey), transCatHandlers.My)
	}

	// Grouping routes under /transaction
	transactionHandlers := handlers.NewTrans(db)
	transactionRoutes := router.Group("/transaction")
	{
		transactionRoutes.POST("/new", transactionHandlers.NewTransaction)
		transactionRoutes.GET("/list", transactionHandlers.TransactionList)
	}

	// Start the server
	router.Run(":8080")
}

// NewDatabase initializes the database connection
func NewDatabase() *gorm.DB {
	dsn := os.Getenv("DATABASE")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get DB object:", err)
	}

	var currentDB string
	err = sqlDB.QueryRow("SELECT current_database()").Scan(&currentDB)
	if err != nil {
		log.Fatal("Failed to query current database:", err)
	}

	log.Printf("Connected to Database: %s\n", currentDB)

	return db
}
