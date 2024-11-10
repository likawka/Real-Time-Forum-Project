package handlers

import (
	"net/http"
	
	httpSwagger "github.com/swaggo/http-swagger"
    _ "project-root/docs" // Імпортує згенеровані Swagger документи
)

// SwaggerHandler повертає http.HandlerFunc для обробки Swagger UI
func SwaggerHandler() http.HandlerFunc {
    return httpSwagger.WrapHandler
}