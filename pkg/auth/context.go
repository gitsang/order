package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
	RoleKey     contextKey = "role"
)

func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user_id not found in context")
	}
	return userID, nil
}

func SetUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, UsernameKey, username)
}

func GetUsername(ctx context.Context) (string, error) {
	username, ok := ctx.Value(UsernameKey).(string)
	if !ok {
		return "", errors.New("username not found in context")
	}
	return username, nil
}

func SetRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, RoleKey, role)
}

func GetRole(ctx context.Context) (string, error) {
	role, ok := ctx.Value(RoleKey).(string)
	if !ok {
		return "", errors.New("role not found in context")
	}
	return role, nil
}
