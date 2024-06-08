package models

import "time"

type Device struct {
	Name     string
	LastUsed time.Time
	Type     string
	TypeId   int
}
