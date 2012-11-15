package up

import (
	"fmt"
	"github.com/bmizerany/pq" // Postgres Database Driver
	"os"
)

// Get the Postgres Database Url.
func getDatabaseUrl() string {
	var database_url = os.Getenv("DATABASE_URL")
	// Set a default database url if there is nothing in the environemnt
	if database_url == "" {
		// Postgres.app uses this variable to set-up postgres.
		user := os.Getenv("USER")
		// Inject the username into the connection string
		database_url = "postgres://" + user + "@localhost/up?sslmode=disable"
		// Let the user know they are using a default
		fmt.Println("--- INFO: No DATABASE_URL env var detected, defaulting to " + database_url)
	}
	conn_str, err := pq.ParseURL(database_url)
	if err != nil {
		panic("Unable to Parse DB URL connection string: " + database_url)
	}
	return conn_str
}

// Get the Port from the environment so we can run on Heroku
func getPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("--- INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

// Get the secret key for interacting with the application
func getSecret() string {
	var secret = os.Getenv("SECRET")
	if secret == "" {
		panic("You must set a SECRET in your environment.\n\n $ export SECRET=foo")
	}
	return secret
}
