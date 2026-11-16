package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/DMCDanteh/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.arguments) > 1 {
		return fmt.Errorf("Usage: browse [limit]")
	}
	if len(cmd.arguments) == 1 {
		nmbr, err := strconv.Atoi(cmd.arguments[0])
		if err != nil {
			return err
		}
		limit = nmbr
	}
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("* %s\n* %s\n", post.Title, post.Url)
	}
	return nil
}
