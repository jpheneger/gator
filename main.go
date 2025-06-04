package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jpheneger/gator/internal/config"
	"github.com/jpheneger/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("unable to get logged in user: %v", err)
		}
		err = handler(s, cmd, user)
		if err != nil {
			return fmt.Errorf("unable to run command '%s' for user '%s': %w", cmd.Name, user.Name, err)
		}
		return nil
	}
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error opening connection to database: %v", err)
	}
	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("help", handlerHelp)
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollowFeed))
	cmds.register("following", handlerFollowing)
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}

func handlerHelp(s *state, cmd command) error {
	fmt.Printf("login: Log in the currently active user.\nUsage: login <name>\n\n")
	fmt.Printf("register: Register a new user.\nUsage: register <name>\n\n")
	fmt.Printf("reset: Clear/empty the database.\nUsage: reset\n\n")
	fmt.Printf("users: List the users, highlight current user.\nUsage: users\n\n")
	fmt.Printf("agg: Aggregate the stored feeds at the provided interval/duration (ex: 1s, 1m, 1h15m, etc).\nUsage: agg <duration>\n\n")
	fmt.Printf("addfeed: Add a new RSS feed.\nUsage: addfeed <name> <url>\n\n")
	fmt.Printf("feeds: List feeds owned by logged in user\nUsage: feeds\n\n")
	fmt.Printf("follow: Follow a feed (logged in user)\nUsage: follow <url>\n\n")
	fmt.Printf("following: List feeds followed by logged in user \nUsage: following\n\n")
	fmt.Printf("unfollow: Unfollow a feed (logged in user)\nUsage: unfollow <url>\n\n")
	fmt.Printf("browse: List posts for logged in user. Default limit is 20.\nUsage: browse <limit>\n\n")
	return nil
}
