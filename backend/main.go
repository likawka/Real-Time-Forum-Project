package main

import (
	"log"
	"net/http"
	"project-root/config"
	_ "project-root/docs" // This imports the generated Swagger docs
	"project-root/pkg/db"
	"project-root/pkg/handlers"
	"project-root/pkg/repositories"
	"project-root/pkg/websockets"

	"github.com/gorilla/mux"
)

// @swagger 2.0
// @title           kood-rt-forum API
// @version         1.0
// @description     This is a sample server for kood-rt-forum.
// @host            localhost:8080
// @BasePath        /api
func main() {
	// Load configuration
	config.LoadConfig()

	// Create a new DBHandler instance
	handler := db.NewDBHandler()

	// Initialize main database
	err := handler.InitMainDB(config.AppConfig.Database.Path, config.AppConfig.Database.InitScript)
	if err != nil {
		log.Println("Error initializing main database:", err)
		return
	}

	// Initialize messages database
	err = handler.InitMsgDB(config.AppConfig.MessagesDatabase.Path)
	if err != nil {
		log.Println("Error initializing messages database:", err)
		return
	}

	// Set the global DBHandler for repositories
	repositories.SetDBHandler(handler)

	// Create a new Gorilla Mux router instance
	r := mux.NewRouter()

	// Serve Swagger UI
	r.PathPrefix("/swagger/").Handler(handlers.SwaggerHandler())

	// Serve static files
	r.PathPrefix("/frontend/").Handler(http.StripPrefix("/frontend/", http.FileServer(http.Dir("../frontend"))))

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()

	// Posts
	api.HandleFunc("/posts", handlers.HandleGetPosts).Methods("GET")
	api.HandleFunc("/posts", handlers.HandleCreatePost).Methods("POST")
	api.HandleFunc("/posts/{postId:[0-9]+}", handlers.HandleGetPostAndComments).Methods("GET")
	api.HandleFunc("/posts/{postId:[0-9]+}/comments", handlers.HandleCreateComment).Methods("POST")

	// Ratings
	api.HandleFunc("/rate", handlers.HandleRate).Methods("PUT")

	// Users
	api.HandleFunc("/users", handlers.HandleGetUsers).Methods("GET")
	api.HandleFunc("/users/{nickname}/{type:posts|comments}", handlers.HandleGetUser).Methods("GET")

	// Authentication
	api.HandleFunc("/auth/register", handlers.HandleRegister).Methods("POST")
	api.HandleFunc("/auth/login", handlers.HandleLogin).Methods("POST")
	api.HandleFunc("/auth/logout", handlers.HandleLogout).Methods("DELETE")

	// Chats
	api.HandleFunc("/chats", handlers.HandleGetChats).Methods("GET")
	api.HandleFunc("/chats", handlers.HandleCreateChat).Methods("POST")
	api.HandleFunc("/chats/{chatHash}", handlers.HandleGetChat).Methods("GET")

	// Websockets
	manager := websockets.NewWebSocketManager()
	api.HandleFunc("/ws", manager.WebSocketHandler).Methods("GET")

	// Handle index path
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/public/index.html")
	})

	// Start the server
	log.Printf("Server running on http://localhost%s\n", config.AppConfig.PortNumber)
	log.Println("To stop the server press `Ctrl + C`")
	log.Fatal(http.ListenAndServe(config.AppConfig.PortNumber, r))
}
