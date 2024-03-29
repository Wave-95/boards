// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Board struct {
	ID          pgtype.UUID
	Name        pgtype.Text
	Description pgtype.Text
	UserID      pgtype.UUID
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
}

type BoardInvite struct {
	ID         pgtype.UUID
	BoardID    pgtype.UUID
	SenderID   pgtype.UUID
	ReceiverID pgtype.UUID
	Status     pgtype.Text
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
}

type BoardMembership struct {
	ID        pgtype.UUID
	UserID    pgtype.UUID
	BoardID   pgtype.UUID
	Role      pgtype.Text
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type EmailVerification struct {
	ID         pgtype.UUID
	Code       string
	UserID     pgtype.UUID
	IsVerified pgtype.Bool
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
}

type Post struct {
	ID          pgtype.UUID
	UserID      pgtype.UUID
	Content     pgtype.Text
	Color       pgtype.Text
	Height      pgtype.Int4
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	PostOrder   pgtype.Float8
	PostGroupID pgtype.UUID
}

type PostGroup struct {
	ID        pgtype.UUID
	BoardID   pgtype.UUID
	Title     pgtype.Text
	PosX      pgtype.Int4
	PosY      pgtype.Int4
	ZIndex    pgtype.Int4
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type User struct {
	ID         pgtype.UUID
	Name       pgtype.Text
	Email      pgtype.Text
	Password   pgtype.Text
	IsGuest    pgtype.Bool
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
	IsVerified pgtype.Bool
}
