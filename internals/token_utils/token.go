package tokenutil

import (
	"errors"
	"fmt"
	"time"

	apperror "cakewai/cakewai.com/component/apperr"
	"cakewai/cakewai.com/domain"

	"github.com/golang-jwt/jwt/v4"
)

func CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	//exp := time.Now().Add(time.Duration(expiry))
	exp := time.Now().Add(time.Second * 3600)
	claims := &domain.JwtCustomClaims{
		Name:     user.Name,
		GoogleId: user.GoogleId,
		Email:    user.Email,
		ID:       user.Id.Hex(),
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

func CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	claimsRefresh := &domain.JwtCustomRefreshClaims{
		ID:       user.Id.Hex(),
		Name:     user.Name,
		GoogleId: user.GoogleId,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expiry))),
		},
	}
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
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", apperror.ErrInvalidToken
	}

	sid := claims["id"].(string)

	return sid, nil
}

func ExtractID(requestToken string, secretKey string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])

		}

		return []byte(secretKey), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)

	fmt.Printf("\nExtractID in extractID function %s\n", claims["id"])
	if err != nil {
		fmt.Printf("\nerror extractID %s\n\n", err)
		return "", err
	}
	fmt.Printf("\nExtractID in extractID function %s\n", token)
	//claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", errors.ErrUnsupported
	}

	id := claims["_id"].(string)

	idString := string(id)
	fmt.Printf("id String: %d", len(idString))
	return idString, nil
	// token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	// 	}
	// 	return []byte(secretKey), nil
	// })

	// if err != nil {
	// 	return "", err
	// }

	// claims, ok := token.Claims.(jwt.MapClaims)

	// if !ok && !token.Valid {
	// 	fmt.Println("In invalid token extract id")
	// 	return "", fmt.Errorf("Invalid Token")
	// }

	// return claims["id"].(string), nil

}
