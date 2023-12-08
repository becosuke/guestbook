package repository

import (
	domain "github.com/becosuke/guestbook/api/internal/domain/post"
)

type Generator interface {
	GenerateSerial() *domain.Serial
}
