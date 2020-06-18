package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	dsn := "server=%s;port=%d;database=%s;user id=%s;password=%s;log=127"
	cs := fmt.Sprintf(dsn, "localhost", 1433, "master", "sa", "Password1")

	conn, err := sql.Open("sqlserver", cs)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	err = loop(conn)
	if err != nil {
		panic(err)
	}
}

func loop(conn *sql.DB) error {
	ctx := context.Background()
	for true {
		tx, err := conn.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("Failed to Begintx: %w", err)
		}

		res, err := query(ctx, tx)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to query: %w", err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("Failed to Commit: %w", err)
		}

		fmt.Println("res: ", strings.Join(res, ", "))

		time.Sleep(1 * time.Second)
	}

	return nil
}

func query(ctx context.Context, tx *sql.Tx) ([]string, error) {
	rows, err := tx.QueryContext(ctx, "select name from sys.tables")
	if err != nil {
		return nil, fmt.Errorf("Failed to QueryContext: %w", err)
	}
	defer rows.Close()

	var names []string
	var name string
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("Failed to Scan: %w", err)
		}
		names = append(names, name)
	}
	return names, nil
}
