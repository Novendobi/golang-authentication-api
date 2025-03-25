package models

import (
    "database/sql"
    "log"
    "os"
	"errors"

    "github.com/lib/pq"
    "github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)


type User struct {
	Id int
	Username string
	Email string
	Password string
}

var db *sql.DB

func connectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ERROR: Could not load enviroment, create a '.env' file", err)
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
        log.Fatal("DATABASE_URL not set in .env")
    }
	db, err = sql.Open("postgres", databaseUrl)
    if err != nil {
        log.Fatal("Error connecting to database:", err)
    }
	err = db.Ping()
    if err != nil {
        log.Fatal("Error pinging database:", err)
    }
    log.Println("Connected to PostgreSQL from models!")

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(50) NOT NULL UNIQUE,
            email VARCHAR(100) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL
        )
    `)
    if err != nil {
        log.Fatal("Error creating table:", err)
    }
    log.Println("Users table ready!")
}


func addUser(username string, email string, plainPassword string) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
    if err != nil {
        return User{}, err
    }

	var user User
    err = db.QueryRow(
        "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, username, email, password",
        username, email, string(hashedPassword),
    ).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	
    if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			log.Println("Duplicate username or email:", err)
			return User{}, errors.New("username or email already exists")
		}

		log.Println("Error inserting user:", err)
		return User{}, err
	}
    return user, nil
}


func updateUser(id int, username string, email string, plainPassword string) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
    if err != nil {
        return User{}, err
    }

	var user User
    err = db.QueryRow(
		`UPDATE users 
         SET username = $1, email = $2, password = $3 
         WHERE id = $4 
         RETURNING id, username, email, password`,
		username, email, string(hashedPassword), id,
	).Scan(&user.Id, &user.Username, &user.Email, &user.Password)

	
    if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			log.Println("Duplicate username or email:", err)
			return User{}, errors.New("username or email already exists")
		}

		log.Println("Error inserting user:", err)
		return User{}, err
	}
    return user, nil
}


func getAllUsers() ([]User, error) {
	rows, err := db.Query(
		"SELECT id, username, email, password FROM users",
	)
	if err != nil {
		log.Fatal("Couldn't Fetch data from Table users", err)
	}
	defer rows.Close()
	var users []User
	for rows.Next(){
		var u User
		if err := rows.Scan(&u.Id, &u.Username, &u.Email, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

type AuthType string
const (
	ByEmail AuthType = "email"
	ByUsername AuthType = "username"
)

func verifyUser(authType AuthType, emailOrusername string, plainPassword string) (bool, error) {
	var user User
	var err error
	switch authType {
	case ByEmail:
		err = db.QueryRow(
			"SELECT id, username, email, password FROM users WHERE email = $1", emailOrusername,
		).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	case ByUsername:
		err = db.QueryRow(
			"SELECT id, username, email, password FROM users WHERE username = $1", emailOrusername,
		).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	default:
		log.Fatal("Unkown Authentication Type", authType)
	}
	if err != nil {
		return false, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword)); err != nil {
		return false, err
	}

	return true, nil
}