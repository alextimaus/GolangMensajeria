package main

import (
	auth "backend/Validacion/auto"
	"backend/Validacion/ws"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY")) // Obtén la clave secreta desde las variables de entorno

func main() {
	// Crear una instancia del router Gin
	r := gin.Default()

	// Endpoint de autenticación
	r.POST("/login", func(c *gin.Context) {
		var login struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		// Validar los datos de entrada
		if err := c.BindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de login inválidos"})
			return
		}

		// Generar un token JWT
		tokenString, err := auth.GenerateJWT(login.Username, jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	})

	// Endpoint para WebSocket
	r.GET("/ws", func(c *gin.Context) {
		token := c.Query("token") // Token JWT de la URL

		// Validar el token JWT
		isValid, err := auth.ValidateJWT(token, jwtKey)
		if err != nil || !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		ws.WsHandler(c)
	})

	// Ejecutar el servidor en todas las interfaces de red en el puerto 8080
	r.Run("0.0.0.0:8080")
}
