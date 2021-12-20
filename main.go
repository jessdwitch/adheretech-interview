// Script to retrieve tokens from a token source and store them in a DB
package main

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

const INSERT_TOKEN_SQL = `INSERT INTO secret_tokens (data) VALUES ($1)`
const DEFAULT_MAX_TOKENS = 455902

func main() {
	// Get the number of tokens to generate from the command line
	qty := os.Args[1]
	if val, err := strconv.Atoi(qty); err != nil {
		log.Fatal("Invalid quantity: ", qty)
	} else {
		if val < 1 {
			log.Fatal("Quantity must be greater than 0")
		}
		max, err := strconv.Atoi(os.Getenv("MAX_TOKENS"))
		if err != nil || max < 1 {
			max = DEFAULT_MAX_TOKENS
		}
		if val > max {
			log.Fatal("Quantity must be less than or equal to ", max)
		}
	}
	db := getDB()
	tokens, err := getAndStoreTokens(db, qty)
	if err != nil {
		log.Fatal(err)
	}
	for token, success := range tokens {
		if success {
			log.Printf("OK : %s", token)
		} else {
			log.Printf("ERR: %s", token)
		}
	}
}

// getAndInsertTokens: Get qty tokens from the token source and insert them into the database.
func getAndStoreTokens(db *sql.DB, qty string) (map[string]bool, error) {
	tokens, err := getTokens(qty)
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool)
	for _, token := range tokens {
		result[token] = insertToken(db, token)
	}
	return result, nil
}

// getTokens: Request qty tokens from the token source
func getTokens(qty string) ([]string, error) {
	token_host := os.Getenv("TOKEN_HOST")
	if token_host == "" {
		return nil, errors.New("TOKEN_HOST not set")
	}
	url := token_host + "?size=" + qty
	resp, err := http.Post(url, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("Token source returned " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result, last := []string{}, 0
	for i, r := range body {
		if r == '\n' {
			target := body[last:i]
			if len(target) > 0 {
				result = append(result, string(target))
			}
			last = i + 1
		}
	}
	return result, nil
}

// insertToken: Insert a token into the database. Return true if successful.
func insertToken(db *sql.DB, token string) bool {
	_, err := db.Exec(INSERT_TOKEN_SQL, token)
	return err == nil
}

// getDB: Get a database connection
func getDB() *sql.DB {
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "require"
	}
	connstr := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASS") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=" + sslmode
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
