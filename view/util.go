package view

import (
	"context"

	"com.github.denisbytes.goimageai/models"
)

func AuthenticatedUser(ctx context.Context) models.AuthenticatedUser {
	user, ok := ctx.Value(models.UserContextKey).(models.AuthenticatedUser)
	if !ok {
		return models.AuthenticatedUser{}
	}
	return user
}
