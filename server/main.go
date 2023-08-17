package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type User struct {
	ID       int
	Username string
	Role     string // "service_provider" or "service_consumer"
}

type Service struct {
	ID          int
	Name        string
	Description string
	ProviderID  int
	ImageURL    string // New field for the image URL
}

type Order struct {
	ID         int
	ServiceID  int
	ConsumerID int
	Status     string // "pending", "completed"
}

var (
	users        = make(map[int]User)
	services     = make(map[int]Service)
	orders       = make(map[int]Order)
	userCount    = 0
	serviceCount = 0
	orderCount   = 0

	// Store user credentials (for simplicity, use a map, in a real-world scenario, use a database)
	userCredentials = make(map[string]string)
	userSessions    = make(map[string]int) // Map user sessions to user IDs
	sessionCounter  = 0
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Endpoint to create a new user
	app.Post("/users", func(c *fiber.Ctx) error {
		userCount++
		user := User{
			ID:       userCount,
			Username: c.FormValue("username"),
			Role:     c.FormValue("role"),
		}
		users[user.ID] = user
		return c.JSON(user)
	})

	// Endpoint to list all services
	app.Get("/services", func(c *fiber.Ctx) error {
		var serviceList []Service
		for _, service := range services {
			serviceList = append(serviceList, service)
		}
		return c.JSON(serviceList)
	})

	// Endpoint to create a new service
	app.Post("/services", func(c *fiber.Ctx) error {
		serviceCount++
		providerID, _ := strconv.Atoi(c.FormValue("provider_id"))
		service := Service{
			ID:          serviceCount,
			Name:        c.FormValue("name"),
			Description: c.FormValue("description"),
			ProviderID:  providerID,
			ImageURL:    c.FormValue("image_url"), // Get image URL from form
		}
		services[service.ID] = service
		return c.JSON(service)
	})

	// Endpoint to place an order
	app.Post("/orders", func(c *fiber.Ctx) error {
		orderCount++
		serviceID, _ := strconv.Atoi(c.FormValue("service_id"))
		consumerID, _ := strconv.Atoi(c.FormValue("consumer_id"))
		order := Order{
			ID:         orderCount,
			ServiceID:  serviceID,
			ConsumerID: consumerID,
			Status:     "pending",
		}
		orders[order.ID] = order
		return c.JSON(order)
	})

	// Endpoint to manage orders
	app.Put("/orders/:id", func(c *fiber.Ctx) error {
		orderID, _ := strconv.Atoi(c.Params("id"))
		order, found := orders[orderID]
		if !found {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Order not found",
			})
		}
		order.Status = "completed"
		orders[orderID] = order
		return c.JSON(order)
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		if _, exists := userCredentials[username]; exists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Username already exists",
			})
		}

		userCount++
		userCredentials[username] = password
		users[userCount] = User{
			ID:       userCount,
			Username: username,
			Role:     c.FormValue("role"),
		}
		return c.JSON(fiber.Map{
			"message": "Registration successful",
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		storedPassword, exists := userCredentials[username]
		if !exists || storedPassword != password {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid credentials",
			})
		}

		sessionCounter++
		sessionID := strconv.Itoa(sessionCounter)
		userSessions[sessionID] = userCount // Map session to user ID

		return c.JSON(fiber.Map{
			"message":    "Login successful",
			"session_id": sessionID,
			"user_id":    userCount,
			"username":   username,
			"user_role":  users[userCount].Role,
		})
	})

	// Run the server
	log.Fatal(app.Listen("0.0.0.0:3000"))
}
