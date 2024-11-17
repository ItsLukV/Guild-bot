package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ItsLukV/Guild-bot/src/guildData"
	"github.com/ItsLukV/Guild-bot/src/guildEvent"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// -----------------------------------------------
// -------------- Struct and const --------------
// -----------------------------------------------

var (
	host            = os.Getenv("DB_HOST")
	port     uint16 = 5432
	user            = os.Getenv("DB_USER")
	password        = os.Getenv("DB_PASSWORD")
	dbname          = os.Getenv("DB")
)

type Database struct {
	pool *pgxpool.Pool
}

// -----------------------------------------------
// ------------------ Functions ------------------
// -----------------------------------------------

var lock = &sync.Mutex{}

var singleInstance *Database

func GetInstance() *Database {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")
			singleInstance = &Database{}
			singleInstance.init()
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singleInstance
}

func (d *Database) init() {
	if d == nil {
		return
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	dbpool, err := pgxpool.New(context.Background(), psqlInfo)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	d.pool = dbpool

	fmt.Println("Successfully connected!")
}

func (d *Database) Close() {
	d.pool.Close()
}

func (d *Database) Save(bot guildData.GuildBot) {
	err := d.saveUsers(bot.Users)
	if err != nil {
		return
	}
}

func (d *Database) AddUser(user guildData.GuildUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a timeout
	defer cancel()                                                          // Ensure the context is cancelled

	query := `
        INSERT INTO users (discord_snowflake, discord_username, minecraft_username, minecraft_uuid)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (discord_snowflake) DO NOTHING;
    `

	_, err := d.pool.Exec(ctx, query,
		user.Snowflake,       // $1: discord_snowflake
		user.DiscordUsername, // $2: discord_username
		user.McUsername,      // $3: minecraft_username
		user.McUUID,          // $4: minecraft_uuid
	)
	if err != nil {
		log.Printf("Failed to insert user %s: %v", user.Snowflake, err)
		return fmt.Errorf("failed to insert user: %w", err)
	}

	log.Printf("User %s inserted successfully.", user.Snowflake)
	return nil
}

func (d *Database) saveUsers(users map[string]guildData.GuildUser) error {
	if len(users) == 0 {
		log.Println("No users to insert")
		return nil
	}

	ctx := context.Background()

	// Begin a transaction
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start database transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `
        CREATE TEMP TABLE temp_users (
            discord_snowflake TEXT,
            discord_username TEXT,
            minecraft_username TEXT,
            minecraft_uuid TEXT
        ) ON COMMIT DROP;
    `)
	if err != nil {
		return fmt.Errorf("failed to create temporary table: %w", err)
	}

	rows := make([][]interface{}, 0, len(users))

	for _, user := range users {
		rows = append(rows, []interface{}{
			user.Snowflake,
			user.DiscordUsername,
			user.McUsername,
			user.McUUID,
		})
	}

	// Perform bulk insert into the temporary table
	copyCount, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"temp_users"},
		[]string{"discord_snowflake", "discord_username", "minecraft_username", "minecraft_uuid"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("failed to copy data into temporary table: %w", err)
	}

	// Insert into the main table with conflict handling
	_, err = tx.Exec(ctx, `
        INSERT INTO users (discord_snowflake, discord_username, minecraft_username, minecraft_uuid)
        SELECT discord_snowflake, discord_username, minecraft_username, minecraft_uuid FROM temp_users
        ON CONFLICT (discord_snowflake) DO NOTHING;
    `)
	if err != nil {
		return fmt.Errorf("failed to insert from temporary table: %w", err)
	}

	log.Printf("Batch insert completed successfully! Inserted %d users.", copyCount)
	return nil
}

func (d *Database) fetchUsers() (map[string]guildData.GuildUser, error) {
	// Create a context for the query
	ctx := context.Background()

	// Define the query to fetch all users
	query := `
        SELECT discord_snowflake, discord_username, minecraft_username, minecraft_uuid
        FROM users
    `

	// Execute the query
	rows, err := d.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Create a map to store the results
	users := make(map[string]guildData.GuildUser)

	// Iterate through the rows and populate the map
	for rows.Next() {
		var user guildData.GuildUser
		var snowflake string

		// Scan the row into the GuildUser struct and snowflake
		err := rows.Scan(&snowflake, &user.DiscordUsername, &user.McUsername, &user.McUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Store the user in the map with the snowflake as the key
		users[snowflake] = user
	}

	// Check for any errors that occurred during iteration
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}

	return users, nil
}

func LoadGuildBot() (guildData.GuildBot, error) {
	users, err := GetInstance().fetchUsers()
	if err != nil {
		return guildData.GuildBot{}, fmt.Errorf("failed to fetch users: %w", err)
	}
	return guildData.GuildBot{
		Users:  users,
		Events: make(map[int]guildEvent.Event),
	}, nil
}
