package main

import (
	"fmt"
	"testing"

	"qiyetalk-server-go/db"
	"qiyetalk-server-go/models"
)

// TestUser ...
func TestUserSelection(m *testing.T) {
	db := db.GetDB()
	var users []models.User
	err := db.Model(&users).Select()
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user, user.Data["privilege"])
	}
}
