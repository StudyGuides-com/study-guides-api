package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
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
		fmt.Printf("[DEBUG] AuthUnaryInterceptor called for method: %s\n", info.FullMethod)

		md, _ := metadata.FromIncomingContext(ctx)
		fmt.Printf("[DEBUG] Metadata from context: %+v\n", md)

		authHeader := md["authorization"]
		fmt.Printf("[DEBUG] Authorization header: %+v\n", authHeader)

		if len(authHeader) > 0 {
			tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")
			fmt.Printf("[DEBUG] Token string after Bearer removal: %s\n", tokenStr)
			fmt.Printf("[DEBUG] JWT Secret length: %d\n", len(secret))

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				fmt.Printf("[DEBUG] JWT Parse callback called, signing method: %v\n", token.Method)
				return []byte(secret), nil
			})
			if err != nil {
				fmt.Printf("[DEBUG] JWT Parse error: %v\n", err)
			}
			fmt.Printf("[DEBUG] Token valid: %v\n", token.Valid)

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// Log decoded JWT claims for debugging
				fmt.Printf("[DEBUG] JWT Claims: %+v\n", claims)

				if sub, ok := claims["sub"].(string); ok {
					fmt.Printf("[DEBUG] User ID from JWT: %s\n", sub)
					ctx = context.WithValue(ctx, userIDKey, sub)
				}

				// Check for roles in different possible formats
				fmt.Printf("[DEBUG] Raw roles claim: %+v (type: %T)\n", claims["roles"], claims["roles"])

				if role, ok := claims["roles"].([]interface{}); ok {
					fmt.Printf("[DEBUG] Processing roles as []interface{}: %+v\n", role)
					var roles []sharedpb.UserRole
					for _, r := range role {
						if str, ok := r.(string); ok {
							// Handle simple string roles (legacy format)
							fmt.Printf("[DEBUG] Processing role string: %s\n", str)
							switch strings.ToLower(str) {
							case "admin":
								fmt.Printf("[DEBUG] Mapped admin role to USER_ROLE_ADMIN\n")
								roles = append(roles, sharedpb.UserRole_USER_ROLE_ADMIN)
							case "user":
								roles = append(roles, sharedpb.UserRole_USER_ROLE_USER)
							case "freelancer":
								roles = append(roles, sharedpb.UserRole_USER_ROLE_FREELANCER)
							case "tester":
								roles = append(roles, sharedpb.UserRole_USER_ROLE_TESTER)
							default:
								fmt.Printf("[DEBUG] Unknown role: %s, mapping to UNSPECIFIED\n", str)
								roles = append(roles, sharedpb.UserRole_USER_ROLE_UNSPECIFIED)
							}
						} else if roleObj, ok := r.(map[string]interface{}); ok {
							// Handle complex role objects (current format)
							fmt.Printf("[DEBUG] Processing role object: %+v\n", roleObj)
							if roleInfo, ok := roleObj["role"].(map[string]interface{}); ok {
								if roleName, ok := roleInfo["name"].(string); ok {
									fmt.Printf("[DEBUG] Found role name in object: %s\n", roleName)
									switch strings.ToLower(roleName) {
									case "admin":
										fmt.Printf("[DEBUG] Mapped admin role object to USER_ROLE_ADMIN\n")
										roles = append(roles, sharedpb.UserRole_USER_ROLE_ADMIN)
									case "user":
										roles = append(roles, sharedpb.UserRole_USER_ROLE_USER)
									case "freelancer":
										roles = append(roles, sharedpb.UserRole_USER_ROLE_FREELANCER)
									case "tester":
										roles = append(roles, sharedpb.UserRole_USER_ROLE_TESTER)
									default:
										fmt.Printf("[DEBUG] Unknown role object name: %s, mapping to UNSPECIFIED\n", roleName)
										roles = append(roles, sharedpb.UserRole_USER_ROLE_UNSPECIFIED)
									}
								} else {
									fmt.Printf("[DEBUG] Role object missing 'name' field\n")
								}
							} else {
								fmt.Printf("[DEBUG] Role object missing 'role' field\n")
							}
						} else {
							fmt.Printf("[DEBUG] Role is neither string nor object: %+v (type: %T)\n", r, r)
						}
					}
					fmt.Printf("[DEBUG] Final mapped roles: %+v\n", roles)
					ctx = context.WithValue(ctx, userRoleKey, roles)
				} else if roleStr, ok := claims["roles"].(string); ok {
					// Handle single role as string
					fmt.Printf("[DEBUG] Processing single role as string: %s\n", roleStr)
					var role sharedpb.UserRole
					switch strings.ToLower(roleStr) {
					case "admin":
						fmt.Printf("[DEBUG] Mapped single admin role to USER_ROLE_ADMIN\n")
						role = sharedpb.UserRole_USER_ROLE_ADMIN
					case "user":
						role = sharedpb.UserRole_USER_ROLE_USER
					case "freelancer":
						role = sharedpb.UserRole_USER_ROLE_FREELANCER
					case "tester":
						role = sharedpb.UserRole_USER_ROLE_TESTER
					default:
						fmt.Printf("[DEBUG] Unknown single role: %s, mapping to UNSPECIFIED\n", roleStr)
						role = sharedpb.UserRole_USER_ROLE_UNSPECIFIED
					}
					roles := []sharedpb.UserRole{role}
					fmt.Printf("[DEBUG] Final mapped single role: %+v\n", roles)
					ctx = context.WithValue(ctx, userRoleKey, roles)
				} else {
					fmt.Printf("[DEBUG] No valid roles found in JWT claims\n")
				}
			}
		}

		return handler(ctx, req)
	}
}
