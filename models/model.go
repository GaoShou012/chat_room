package models

import "time"

type Model struct {
	ID        *uint64    `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"default:NULL" sql:"index"`
}