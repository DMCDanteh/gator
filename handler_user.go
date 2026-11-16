package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/DMCDanteh/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("The login handler expects a single argument, the username")
	}
	username := cmd.arguments[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}
	if err := s.cfg.SetUser(cmd.arguments[0]); err != nil {
		return err
	}
	fmt.Printf("The user has been set.\n")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("usage: register <username>")
	}
	username := cmd.arguments[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Printf("Username: %s.\nID: %v\n", user.Name, user.ID)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.ResetUsers(context.Background()); err != nil {
		return err
	}
	fmt.Printf("Successfully reset database.\n")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
