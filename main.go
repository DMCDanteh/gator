package main

import (
	"log"
	"os"

	"database/sql"

	"github.com/DMCDanteh/gator/internal/config"
	"github.com/DMCDanteh/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	dbURL := cfg.DbURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	appState := state{
		cfg: &cfg,
		db:  database.New(db),
	}
	appCommands := commands{
		commands: make(map[string]func(*state, command) error),
	}
	appCommands.register("login", handlerLogin)
	appCommands.register("register", handlerRegister)
	appCommands.register("reset", handlerReset)
	appCommands.register("users", handlerUsers)
	appCommands.register("agg", handlerAgg)
	appCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	appCommands.register("feeds", handlerFeeds)
	appCommands.register("follow", middlewareLoggedIn(handlerFollow))
	appCommands.register("following", middlewareLoggedIn(handlerFollowing))
	appCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	appCommands.register("browse", middlewareLoggedIn(handlerBrowse))
	if len(os.Args) < 2 {
		log.Fatal("Error: no command provided. Usage: gator <command> [arguments]")
	}
	appCommand := command{
		name:      os.Args[1],
		arguments: os.Args[2:],
	}
	if err := appCommands.run(&appState, appCommand); err != nil {
		log.Fatal(err)
	}
}
