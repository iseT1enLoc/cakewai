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
	exp := time.Now().Add(time.Minute * 5)
	claims := &domain.JwtCustomClaims{
		ID: user_id.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},

	}

	// Create the JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", claims, err
	}

	// Return the access token and claims
	return t, claims, nil
}

func CreateRefreshToken(user_id primitive.ObjectID, is_admin bool, secret string, expiry int) (refreshToken string, refresh_token_claim domain.JwtCustomRefreshClaims, err error) {

	exp := time.Now().UTC().Add(time.Minute * 2)
	fmt.Printf("\nCreate refreshtoken exp: %v\n", exp.Local().Format("Mon Jan 2 15:04:05 2006"))

	refresh_token_claim = domain.JwtCustomRefreshClaims{
		ID:      user_id.String(),
		IsAdmin: is_admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	fmt.Printf("\nCreate refreshtoken exp: %v\n", exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh_token_claim)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", refresh_token_claim, err
	}
	return rt, refresh_token_claim, nil
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

func ExtractIDAndRole(requestToken string, secretKey string) (string, *bool, error) {
	// Parse the token with the specified secret key
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return "", nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims as MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Print(claims)
	if !ok {
		return "", nil, errors.New("failed to parse claims")
	}

	// Retrieve the ID from the "name" claim
	idStr, ok := claims["_id"].(string)
	// Check the "exp" claim

	if !ok {
		return "", nil, errors.New("expiration time (exp) claim not found")
	}

	if strings.HasPrefix(idStr, "ObjectID(\"") && strings.HasSuffix(idStr, "\")") {
		idStr = idStr[10 : len(idStr)-2] // Remove `ObjectID("` and `")` around the hex
	}
	// Convert the string to an ObjectID
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse ObjectID: %w", err)
	}
	is_admin, ok := claims["is_admin"].(bool)

	fmt.Printf("ID is now %v\n", objID.Hex())
	return objID.Hex(), &is_admin, nil
}
