package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Configuración del "upgrader" para WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Permitir conexiones desde cualquier origen
		return true
	},
}

// Manejador para WebSocket
func WsHandler(c *gin.Context) {
	// Asegurarse de que el paquete 'websocket' se está utilizando
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// En caso de error, devolvemos un mensaje HTTP 500
		http.Error(c.Writer, "No se pudo establecer la conexión WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Ciclo de comunicación con el cliente WebSocket
	for {
		// Leer mensaje del cliente
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// Si hubo un error al leer el mensaje, cerramos la conexión
			break
		}

		// Enviar el mismo mensaje de vuelta al cliente (eco)
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			// Si hubo un error al escribir el mensaje, cerramos la conexión
			break
		}
	}
}
