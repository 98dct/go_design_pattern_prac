package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"
)

func TestQuerySql(t *testing.T) {
	DB, err := sql.Open("mysql", "root:root@tcp(192.168.8.100:3306)/user?charset=utf8&parseTime=true&loc=Asia%2fShanghai")
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	ctx := context.Background()
	err = DB.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	rows, err := DB.QueryContext(ctx, "select id, name from users")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	var res []struct {
		id   sql.NullInt64
		name sql.NullString
	}

	for rows.Next() {
		var row struct {
			id   sql.NullInt64
			name sql.NullString
		}
		err := rows.Scan(&row.id, &row.name)
		if err != nil {
			log.Println(err)
			return
		}
		res = append(res, row)
	}
	for _, item := range res {
		fmt.Println(item.id.Int64, item.name.String)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return
	}
}

// prepare statement查询多行
// 相比于直接查询的好处：
// 效率：语句不需要被重新编译（编译意味着解析+优化+转译），使用binary proto 传输更加轻量，紧凑
// 安全性：降低了sql注入攻击的风险
func TestQueryMultiRowSql(t *testing.T) {
	DB, err := sql.Open("mysql", "root:root@tcp(192.168.8.100:3306)/user?charset=utf8&parseTime=true&loc=Asia%2fShanghai")
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	ctx := context.Background()
	stmt, err := DB.PrepareContext(ctx, "select id, name from users where id in (?,?,?)")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, 1, 2, 3)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
	var res []struct {
		id   sql.NullInt64
		name sql.NullString
	}

	for rows.Next() {
		var row struct {
			id   sql.NullInt64
			name sql.NullString
		}
		err := rows.Scan(&row.id, &row.name)
		if err != nil {
			log.Println(err)
			return
		}
		res = append(res, row)
	}
	for _, item := range res {
		fmt.Println(item.id.Int64, item.name.String)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return
	}

}

// prepare statement查询单行
func TestQuerySingleRowSql(t *testing.T) {
	DB, err := sql.Open("mysql", "root:root@tcp(192.168.8.100:3306)/user?charset=utf8&parseTime=true&loc=Asia%2fShanghai")
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	ctx := context.Background()
	stmt, err := DB.PrepareContext(ctx, `select id, name from users where id = ?`)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	var res struct {
		id   sql.NullInt64
		name sql.NullString
	}
	err = stmt.QueryRowContext(ctx, 2).Scan(&res.id, &res.name)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(res.id.Int64, res.name.String)

}

// 更新行
func TestUpdateRowSql(t *testing.T) {
	DB, err := sql.Open("mysql", "root:root@tcp(192.168.8.100:3306)/user?charset=utf8&parseTime=true&loc=Asia%2fShanghai")
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	ctx := context.Background()
	stmt, err := DB.PrepareContext(ctx, `update users set name = ?, updated_at = ? where id = ?`)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, "lisi", time.Now(), "222")
	if err != nil {
		log.Println(err)
		return
	}
	rowsAffected, _ := res.RowsAffected()
	lastInsertId, _ := res.LastInsertId()
	fmt.Println(rowsAffected, lastInsertId)
}

// 插入行
func TestInsertRowSql(t *testing.T) {
	DB, err := sql.Open("mysql", "root:root@tcp(192.168.8.100:3306)/user?charset=utf8&parseTime=true&loc=Asia%2fShanghai")
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	ctx := context.Background()
	stmt, err := DB.PrepareContext(ctx, `insert into users(name , phone, created_at, updated_at) values(?,?,?,?)`)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, "张三", "88999922311", time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return
	}
	rowsAffected, _ := res.RowsAffected()
	lastInsertId, _ := res.LastInsertId() // mysql服务器给本次插入的记录生成的主键id
	fmt.Println(rowsAffected, lastInsertId)
}

// 事务情况下执行sql
func TestTransSql(t *testing.T) {
	DB, err := sql.Open("mysql", "root:root@tcp(192.168.8.100:3306)/user?charset=utf8&parseTime=true&loc=Asia%2fShanghai")
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	ctx := context.Background()
	tx, _ := DB.Begin()
	stmt, err := tx.PrepareContext(ctx, `insert into users(name , phone, created_at, updated_at) values(?,?,?,?)`)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, "王麻子", "9445238387", time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	rowsAffected, _ := res.RowsAffected()
	lastInsertId, _ := res.LastInsertId() // mysql服务器给本次插入的记录生成的主键id
	fmt.Println(rowsAffected, lastInsertId)
}

func TestUnEscape(t *testing.T) {
	res, err := url.QueryUnescape("charset=utf8&parseTime=true&loc=Asia%2fShanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

func TestByte(t *testing.T) {
	fmt.Println(byte(1 >> 8))
}
