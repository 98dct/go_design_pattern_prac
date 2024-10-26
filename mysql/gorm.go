package main

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var (
	db   *gorm.DB
	once sync.Once
	dsn  = "root:root@(192.168.8.100:3306)/test101?charset=utf8&parseTime=true&loc=Asia%2fShanghai"
)

func getDB() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		//cfg := &gorm.Config{}
		//cfg.Logger = logger.Default.LogMode(logger.Info)

		cfg := &gorm.Config{}
		cfg.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Warn, // Log level
				IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Disable color
			},
		)

		db, err = gorm.Open(mysql.Open(dsn), cfg)
	})

	return db, err
}

type Reward struct {
	gorm.Model               // gorm默认开启软删除机制
	Amount     sql.NullInt64 `gorm:"column:amount"`
	Tp         string        `gorm:"not null"`
	UserId     int64         `gorm:"not null"`
}

func (r *Reward) TableName() string {
	return "reward"
}

// 查询
func Query() {
	db, _ := getDB()

	var r Reward
	if err := db.First(&r).Error; err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r)
}

// 创建
func create() {
	db, _ := getDB()

	r := Reward{

		Amount: sql.NullInt64{
			Int64: 0,
			Valid: true,
		},
		Tp:     "jurge",
		UserId: 456,
	}
	if err := db.Create(&r).Error; err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r)
}

// 删除
func delete() {
	db, _ := getDB()

	var r Reward
	if err := db.Delete(&r, 8).Error; err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r)
}

// 更新
func update() {
	db, _ := getDB()

	r := Reward{
		Amount: sql.NullInt64{
			Int64: 1000,
			Valid: true,
		},
	}

	if err := db.Where("id = ?", 2).Updates(&r).Error; err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r)
}

// 事务
func tx() {
	db, _ := getDB()

	f := func(tx *gorm.DB) error {
		return nil
	}
	if err := db.Transaction(f).Error; err != nil {
		fmt.Println(err)
		return
	}

}

var sli []int64

func main() {
	//Query()
	//create()
	//delete()
	//test()
	sli = make([]int64, 1024)
	//fmt.Println(s, s1)
}

func test() {
	s := make([]int64, 1024)
	s1 := make([]int64, 1023)
	fmt.Println(s, s1)
}
