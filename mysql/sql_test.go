package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	// 注册mysql驱动
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserId int64
}

func Test_sql(t *testing.T) {
	db, err := sql.Open("mysql", "username:password@(localhost:3306)/database")
	if err != nil {
		t.Error(err)
		return
	}

	db.SetConnMaxLifetime(time.Second)
	// 执行sql
	ctx := context.Background()
	row := db.QueryRowContext(ctx, "select user_id from user where order by created_at desc limit 1")
	if row.Err() != nil {
		t.Error(row.Err())
		return
	}

	// 解析结果
	var u User
	if err = row.Scan(&u.UserId); err != nil {
		t.Error(err)
		return
	}
}

var ErrNotFound = errors.New("not found")

func Test_error(t *testing.T) {

	// Wrapping ErrNotFound
	err := fmt.Errorf("something went wrong: %w", ErrNotFound)

	// Checking if err contains ErrNotFound
	if errors.Is(err, ErrNotFound) {
		fmt.Println("Error is ErrNotFound")
	} else {
		fmt.Println("Error is not ErrNotFound")
	}
}
