package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

func authRequired(c *fiber.Ctx) error {

	cookie := c.Cookies("jwt")

	jwtSecretKey := "TestSecret"

	token, err := jwt.ParseWithClaims(
		cookie,
		// &jwt.RegisteredClaims{},
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecretKey), nil
		})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claim := token.Claims.(jwt.MapClaims)

	fmt.Println(claim)
	return c.Next()
}

func main() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&Book{}, &User{})
	// InitData(db)

	app := fiber.New()
	app.Use("/books", authRequired)

	app.Get("/books", func(c *fiber.Ctx) error {
		return c.JSON(GetBooks(db))
	})

	app.Get("/books/:id", func(c *fiber.Ctx) error {
		bookId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.JSON(GetBook(db, bookId))
	})

	app.Post("/books", func(c *fiber.Ctx) error {
		book := new(Book)
		if err := c.BodyParser(book); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		if err := CreateBook(db, book); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusCreated)
	})

	app.Put("/books/:id", func(c *fiber.Ctx) error {
		bookId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		book := new(Book)
		if err := c.BodyParser(book); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		book.ID = uint(bookId)

		if err := UpdateBook(db, book); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := CreateUser(db, user); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendStatus(fiber.StatusCreated)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		token, err := LoginUser(db, user)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 72),
			HTTPOnly: true,
		})

		return c.JSON(fiber.Map{
			"message": "Login successful.",
		})
	})

	app.Listen(":8080")

}

func InitData(db *gorm.DB) {
	books := []Book{
		{Name: "The Go Programming Language", Author: "Alan A. A. Donovan", Description: "A comprehensive guide to Go programming.", Price: 4500},
		{Name: "Clean Code", Author: "Robert C. Martin", Description: "A handbook of agile software craftsmanship.", Price: 3500},
		{Name: "Introduction to Algorithms", Author: "Thomas H. Cormen", Description: "A detailed book on algorithms.", Price: 5000},
		{Name: "Design Patterns", Author: "Erich Gamma", Description: "Elements of reusable object-oriented software.", Price: 4000},
		{Name: "You Don't Know JS", Author: "Kyle Simpson", Description: "A deep dive into JavaScript.", Price: 3000},
	}

	// Insert books into the database
	for _, book := range books {
		CreateBook(db, &book)
	}
}
