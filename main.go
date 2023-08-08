package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")           //Comes from the joho/godotenv, used to load in the .env file that holds sensetive information
	portString := os.Getenv("PORT") // Comes from Go standard library, os is the package name, Getenv("Key") looks at the loaded .env file for a variable that matches the name passed in

	if portString == "" { //Error check, if port is blank log it and kill code
		log.Fatal("PORT is not found in the enviroment")
	}
	router := chi.NewRouter() // {Note Stamp 1} Comes from the go-chi/chi package, used to define a router

	router.Use(cors.Handler(cors.Options{ // {Note Stamp 2} Comes from the go-chi/cors package, these are configurations to the router that allow in this example allow a lot. They can be tightend up to security
		AllowedOrigins:   []string{"https://*", "http//*"},                    // Allow sending through these
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allow these methods
		AllowedHeaders:   []string{"*"},                                       // Allow any headers
		ExposedHeaders:   []string{"Link"},                                    //Hover over for more details and google search for better explanations {Note Stamp 2 End}
		AllowCredentials: false,
		MaxAge:           300,
	}))

	/*{Note Stamp 3}
	To hook up a handler, we need setup a router to a specific "path" and method
	 - below we are initializing a new chi router called v1Router - version 1 router
	 - we are then connecting our desired handler (handlerReadiness) to a specific "path" ("/healthz") in our router (v1Router) using the HandleFunc() method
	 	-- update changed from HandleFunc() to Get so that only Get request can hit this handler
	 - We are then Mounting - aka attaching - our v1Router to our base router (router) by specifing a "path" ("/v1") and the desired router (v1Router)
	 	- This is good practice because it lets us modularize routers and paths, if we later want to create a version 2 of the path, we can just change the mounted router and void major breaking
		- The full path will be "/v1/healthz"
			- "healthz" - a standard path naming convention for a path that just checks a server is alive an running
	*/
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr) //part of json.go {Note Stamp 2}

	router.Mount("/v1", v1Router)
	//{Note Stamp 3 End}

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server starting on port %v", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
