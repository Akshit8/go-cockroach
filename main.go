package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/jackc/pgx/v4"
)

var (
	dbUser string
	dbHost string
	dbPort string
)

func init() {
	dbUser = os.Getenv("DB_USER")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
}

func main() {
	config, err := pgx.ParseConfig(fmt.Sprintf("postgresql://%s@%s:%s/bank?sslmode=disable", dbUser, dbHost, dbPort))
	if err != nil {
		log.Fatal("error parsing db source: ", err)
	}

	// connect to the "bank" database
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal("error connecting to db: ", err)
	}
	log.Print("successfully connected to cockroach db")
	defer conn.Close(context.Background())

	// create the "accounts" table
	if _, err = conn.Exec(context.Background(),
		"CREATE TABLE IF NOT EXISTS accounts (id INT PRIMARY KEY, balance INT)",
	); err != nil {
		log.Fatal("error creating table accounts: ", err)
	}

	// Insert two rows into the "accounts" table.
	// createAccounts(conn)

	fmt.Println("Initial balance")
	printBalances(conn)

	// run a transfer in transaction.
	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return transferFunds(context.Background(), tx, 1, 2, 250)
	})
	if err == nil {
		fmt.Println("Transfer Successful")
		fmt.Println("Final balance")
		printBalances(conn)
	} else {
		fmt.Println("Transfer Failed: ", err)
	}
}

func transferFunds(ctx context.Context, tx pgx.Tx, fromID int, toID int, amount int) error {
	// Read the balance
	var fromBalance int
	err := tx.QueryRow(ctx,
		"SELECT balance FROM accounts WHERE id = $1", fromID).Scan(&fromBalance)
	if err != nil {
		return err
	}

	if fromBalance < amount {
		return fmt.Errorf("insufficient balance in account: %d", fromID)
	}

	// Perform the transfer
	_, err = tx.Exec(ctx,
		"UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx,
		"UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
	if err != nil {
		return err
	}

	return nil
}

func printBalances(conn *pgx.Conn) {
	rows, err := conn.Query(context.Background(),
		"SELECT id, balance FROM accounts",
	)
	if err != nil {
		log.Println("error retrieving rows: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, balance int
		if err := rows.Scan(&id, &balance); err != nil {
			log.Println("error scanning rows: ", err)
		}
		fmt.Printf("ID: %d Balance: %d\n", id, balance)
	}
}

func createAccounts(conn *pgx.Conn) {
	if _, err := conn.Exec(context.Background(),
		"INSERT INTO accounts (id, balance) VALUES (1, 1000), (2, 250)",
	); err != nil {
		log.Println("error inserting accounts: ", err)
	}
}
