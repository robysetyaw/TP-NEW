package utils

import (
	"errors"
	"strings"
	model "trackprosto/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

func ExtractTokenFromAuthHeader(authHeader string) (string, error) {
	// Mengecek format header Authorization
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid authorization header format")

	}

	// Mengambil token dari header Authorization
	token := strings.TrimPrefix(authHeader, "Bearer ")

	return token, nil
}

// VerifyJWTToken memverifikasi token JWT dan mengembalikan klaim JWT jika token valid
func VerifyJWTToken(tokenString string) (jwt.MapClaims, error) {
	// Menentukan fungsi validasi token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifikasi menggunakan HMAC, RSA, atau algoritma validasi lainnya
		// Pastikan kunci rahasia (secret key) sesuai dengan yang digunakan saat pembuatan token
		// Misalnya, untuk validasi menggunakan HMAC dengan algoritma HS256:
		secretKey := []byte("secret-key")
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Mengecek apakah token valid dan memiliki klaim
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ParseJWTToken(tokenString string, secretKey []byte, claims jwt.Claims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return token, nil
	}

	return nil, errors.New("invalid token")
}

func SendResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, model.Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func GetUsernameFromContext(c *gin.Context) (string, error) {
	token, err := ExtractTokenFromAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	claims, err := VerifyJWTToken(token)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	username, ok := claims["username"].(string)
	if !ok {
		logrus.Error(err)
		return "", errors.New("username not found in claims")
	}

	return username, nil
}

func NonEmpty(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

func NonZero(value, defaultValue float64) float64 {
	if value != 0 {
		return value
	}
	return defaultValue
}
