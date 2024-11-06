package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Genera un token JWT con un nombre de usuario y una clave secreta
func GenerateJWT(username string, jwtKey []byte) (string, error) {
	// Crea un nuevo token con el nombre de usuario y la fecha de expiración (72 horas)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Expira en 72 horas
	})

	// Firma el token con la clave secreta y retorna el token firmado
	return token.SignedString(jwtKey)
}

// Valida un token JWT usando la clave secreta
func ValidateJWT(tokenString string, jwtKey []byte) (bool, error) {
	// Parsea el token usando la clave secreta para verificar su validez
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Asegúrate de que el token use el algoritmo correcto
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo de firma no es válido")
		}
		return jwtKey, nil
	})

	// Verifica si hubo un error al parsear el token
	if err != nil {
		return false, fmt.Errorf("error al parsear el token: %v", err)
	}

	// Verifica si el token es válido
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Obtén la fecha de expiración del token (en formato float64)
		exp, ok := claims["exp"].(float64)
		if !ok {
			return false, fmt.Errorf("fecha de expiración no válida")
		}

		// Verifica si el token ha expirado
		if exp > float64(time.Now().Unix()) {
			return true, nil
		} else {
			return false, fmt.Errorf("el token ha expirado")
		}
	}

	// Si el token no es válido por alguna razón
	return false, fmt.Errorf("token no válido")
}
