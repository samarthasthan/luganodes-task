package main

import (
	"github.com/samarthasthan/luganodes-task/api"
	"github.com/samarthasthan/luganodes-task/internal/store/controller"
	"github.com/samarthasthan/luganodes-task/internal/store/database"
	"github.com/samarthasthan/luganodes-task/pkg/env"
)

var (
	REST_API_PORT       string
	MYSQL_PORT          string
	MYSQL_ROOT_PASSWORD string
	MYSQL_HOST          string
	REDIS_PORT          string
	REDIS_HOST          string
)

func init() {
	REST_API_PORT = env.GetEnv("REST_API_PORT", "8000")
	MYSQL_PORT = env.GetEnv("MYSQL_PORT", "3306")
	MYSQL_ROOT_PASSWORD = env.GetEnv("MYSQL_ROOT_PASSWORD", "password")
	MYSQL_HOST = env.GetEnv("MYSQL_HOST", "localhost")
}

func main() {
	// Create mysql database
	sql := database.NewMySQL()
	err := sql.Connect("root:" + MYSQL_ROOT_PASSWORD + "@tcp(" + MYSQL_HOST + ":" + MYSQL_PORT + ")/luganodes")
	if err != nil {
		panic(err)
	}

	defer sql.Close()

	// Initialize the Controller from the store
	c := controller.NewController(sql, nil)
	// Start the REST API Server
	h := api.NewHandler(c)
	h.Handle()
	h.Logger.Fatal(h.Start(":" + REST_API_PORT))
}
