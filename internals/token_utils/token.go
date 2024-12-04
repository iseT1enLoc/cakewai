package tokenutil

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	apperror "cakewai/cakewai.com/component/apperr"
	"cakewai/cakewai.com/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang-jwt/jwt/v4"
)

func CreateAccessToken(user_id primitive.ObjectID, secret string, expiry int) (accessToken string, err error) {
	//exp := time.Now().Add(time.Duration(expiry))
	exp := time.Now().Add(time.Minute * 2)
	claims := &domain.JwtCustomClaims{
		ID: user_id.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func CreateRefreshToken(user_id primitive.ObjectID, secret string, expiry int) (refreshToken string, err error) {

	exp := time.Now().UTC().Add(time.Minute * 60)
	claimsRefresh := &domain.JwtCustomRefreshClaims{
		ID: user_id.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	fmt.Printf("\nCreate refreshtoken exp: %v\n", exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return rt, nil
}

func Is_authorized(requestToken string, secretkey string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			print(token)
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])

		}
		return []byte(secretkey), nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.ErrUnexpectedSigningMethod
		}
		return []byte(secret), nil
	})
	fmt.Println("Line 71 thot not cuoc tinh")
	if err != nil {
		log.Fatalf(err.Error())
		return "", err
	}
	claims, ok := token.Claims.(*domain.JwtCustomClaims)
	fmt.Println("after")
	if !ok && !token.Valid {
		return "", apperror.ErrInvalidToken
	}

	sid := claims.ID
	fmt.Printf("id string type %v", sid)
	return sid, nil
}

// func ExtractID(requestToken string, secretKey string) (string, error) {
// 	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

// 			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])

// 		}

// 		return []byte(secretKey), nil
// 	})
// 	fmt.Print(token.Claims)
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	idStr, ok := claims["_id"].(string)
// 	if !ok {
// 		return "", errors.New("ID not found or is not a string")
// 	}

// 	// Check if idStr includes "ObjectID(...)"; if so, strip it
// 	if strings.HasPrefix(idStr, "ObjectID(\"") && strings.HasSuffix(idStr, "\")") {
// 		idStr = idStr[10 : len(idStr)-2] // Remove `ObjectID("` and `")` around the hex
// 	}

// 	// Convert the string to ObjectID
// 	objID, err := primitive.ObjectIDFromHex(idStr)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to parse ObjectID: %w", err)
// 	}

// 	fmt.Printf("ID is now %v\n", objID.Hex())
// 	return objID.Hex(), nil
// }

func ExtractID(requestToken string, secretKey string) (string, error) {
	// Parse the token with the specified secret key
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims as MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Print(claims)
	if !ok {
		return "", errors.New("failed to parse claims")
	}

	// Retrieve the ID from the "name" claim
	idStr, ok := claims["_id"].(string)
	// Check the "exp" claim

	if !ok {
		return "", errors.New("expiration time (exp) claim not found")
	}

	if strings.HasPrefix(idStr, "ObjectID(\"") && strings.HasSuffix(idStr, "\")") {
		idStr = idStr[10 : len(idStr)-2] // Remove `ObjectID("` and `")` around the hex
	}
	// Convert the string to an ObjectID
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse ObjectID: %w", err)
	}

	fmt.Printf("ID is now %v\n", objID.Hex())
	return objID.Hex(), nil
}
