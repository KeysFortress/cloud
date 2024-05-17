package models

import "github.com/pquerna/otp"

type Algorithm struct {
	Id      int
	Name    string
	Related otp.Algorithm
}
