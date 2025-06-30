package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fallenkarma/wasatext/internal/handlers"
	"github.com/fallenkarma/wasatext/internal/repository/postgres"
	"github.com/fallenkarma/wasatext/internal/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	} else {
		log.Println("Successfully loaded .env file")
	}

    // Get environment variables
    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080" 
    }
    
    dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	log.Print("DB_CONNECTION_STRING: ", dbConnectionString)

	const UPLOADS_BASE_PATH = "/app/uploads"
	repo,err := postgres.NewPostgresRepository(dbConnectionString,UPLOADS_BASE_PATH)
	if err != nil {
		log.Fatalf("Connection to database failed: %v", err)
	}

	// Initialize service with repository
	svc := service.New(repo)

	// Initialize handlers with service
	handler := handlers.New(svc)

	// Initialize router
	r := mux.NewRouter()


    // This serves files from the UPLOADS_BASE_PATH on the /uploads/ endpoint.
    // So if a file is saved at /uploads/user_photos/user123_12345.jpg
    // it will be accessible at http://your-backend-ip:port/uploads/user_photos/user123_12345.jpg
    fileServer := http.FileServer(http.Dir(UPLOADS_BASE_PATH))
    r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fileServer))

	
	// Add API prefix
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Public routes (no auth required)
	apiRouter.HandleFunc("/session", handler.Login).Methods("POST")

	// Protected routes (auth required)
	protected := apiRouter.NewRoute().Subrouter()
	protected.Use(handler.AuthMiddleware)

	// User routes
	protected.HandleFunc("/users", handler.GetUsers).Methods("GET")
	protected.HandleFunc("/users/me", handler.GetMyUser).Methods("GET")
	protected.HandleFunc("/users/me/username", handler.SetMyUserName).Methods("PUT")
	protected.HandleFunc("/users/me/photo", handler.SetMyPhoto).Methods("PUT")

	// Conversation routes
	protected.HandleFunc("/conversations", handler.CreateConversation).Methods("POST")
	protected.HandleFunc("/conversations", handler.GetMyConversations).Methods("GET")
	protected.HandleFunc("/conversations/{id}", handler.GetConversation).Methods("GET")

	// Message routes
	protected.HandleFunc("/messages", handler.SendMessage).Methods("POST")
	protected.HandleFunc("/messages/forward", handler.ForwardMessage).Methods("POST")
	protected.HandleFunc("/messages/{id}/reaction", handler.CommentMessage).Methods("POST")
	protected.HandleFunc("/messages/{id}/reaction", handler.UncommentMessage).Methods("DELETE")
	protected.HandleFunc("/messages/{id}", handler.DeleteMessage).Methods("DELETE")
	protected.HandleFunc("/messages/{id}", handler.UpdateMessage).Methods("PUT")

	// Group routes
	protected.HandleFunc("/groups/{id}/members", handler.AddToGroup).Methods("POST")
	protected.HandleFunc("/groups/{id}/leave", handler.LeaveGroup).Methods("POST")
	protected.HandleFunc("/groups/{id}/name", handler.SetGroupName).Methods("PUT")
	protected.HandleFunc("/groups/{id}/photo", handler.SetGroupPhoto).Methods("PUT")

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Wrap the router with CORS handler
	handlerWithCORS := crs.Handler(r)	

	// Create server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handlerWithCORS,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Starting server on :" + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}