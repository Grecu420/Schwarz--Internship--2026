package main

import (
	"Schwarz--Internship--2026/services/common/auth"
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var publicMethods = map[string]bool{
	"/gateway.GatewayService/Login":       true,
	"/gateway.GatewayService/CreateUser":  true,
	"/gateway.GatewayService/GetProperty": true,
	"/proto.UserService/GetUser":          true,
	"/proto.UserService/CreateUser":       true,
	"/userbase.UserService/Login":         true,
	"/userbase.UserService/CreateUser":    true,
	"/userbase.UserService/GetUser":       true,
}

type authInterceptor struct {
	secret []byte
}

func (i *authInterceptor) AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	// Skip verification for public methods
	if publicMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	// 1. Extract metadata from incoming context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization header missing")
	}

	// 2. Extract Bearer token
	tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")
	slog.Info("token", "token", tokenStr)
	if tokenStr == authHeader[0] {
		return nil, status.Error(codes.Unauthenticated, "invalid token format")
	}

	// 3. Parse & Validate signature and expiration
	loginClaims, err := auth.ParseToken(tokenStr, i.secret)

	if errors.Is(err, jwt.ErrTokenExpired) {
		slog.Error("expired token", "error", err)
		return nil, status.Error(codes.Unauthenticated, "expired token")
	} else if err != nil {
		slog.Error("invalid token", "error", err)

		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// 4. Store claims/user in context and pass to the endpoint handler
	newCtx := context.WithValue(ctx, "user_id", loginClaims.ID)
	return handler(newCtx, req)
}
