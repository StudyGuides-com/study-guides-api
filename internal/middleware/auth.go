package middleware

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const userIDKey contextKey = "userID"
const userRoleKey contextKey = "userRole"

// AuthUnaryInterceptor extracts user ID from JWT and stores it in context if present.
func AuthUnaryInterceptor(secret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		authHeader := md["authorization"]

		if len(authHeader) > 0 {
			tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")

			token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if sub, ok := claims["sub"].(string); ok {
					ctx = context.WithValue(ctx, userIDKey, sub)
				}
				if role, ok := claims["roles"].([]interface{}); ok {
					roleStrings := make([]string, len(role))
					for i, r := range role {
						if str, ok := r.(string); ok {
							roleStrings[i] = str
						}
					}
					roleStr := strings.Join(roleStrings, ",")
					ctx = context.WithValue(ctx, userRoleKey, roleStr)
				}
			}
		}

		return handler(ctx, req)
	}
}
