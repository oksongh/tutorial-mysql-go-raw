package main

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func Connect() *sql.DB {
	// サーバー公開しないのでパスワード埋め込みで
	conf := mysql.Config{
		DBName:               "test",
		User:                 "eg_user",
		Passwd:               "p@ssW0rd",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func Ping(db *sql.DB) {

	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected")
}

type User struct {
	id      int
	name    string
	stmt    string
	otakuId int
}

func SqlInsert(db *sql.DB, user User) (int64, int64) {

	res, err := db.Exec("INSERT INTO user(name,stmt,otaku_id) VALUES(?,?,?)", user.name, user.stmt, user.otakuId)
	if err != nil {
		log.Fatal(err)
	}
	liid, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("inserted:%d, %d rows", liid, ra)
	return liid, ra
}

func SqlUpdate(db *sql.DB, user User) {

	res, err := db.Exec("UPDATE user SET stmt = ? WHERE id = ?", user.stmt, user.id)

	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.LastInsertId())
	log.Println(res.RowsAffected())

}

func SqlSelect(db *sql.DB) {

	rows, err := db.Query("SELECT * from user")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.id, &user.name, &user.stmt, &user.otakuId)
		if err != nil {
			log.Println(err)
		}
		log.Println(user)
	}
}

func SqlDelete(db *sql.DB, id string) {

	res, err := db.Exec("DELETE FROM user WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}

	res.LastInsertId()
	affected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("deleted: %d row\n", affected)

}

func updateDemo() {
	db := Connect()
	defer db.Close()

	SqlSelect(db)

	SqlUpdate(db, User{
		stmt:    "updated record",
		otakuId: 2,
	})

	SqlSelect(db)

}

func deleteDemo() {

	db := Connect()
	defer db.Close()

	Ping(db)

	// insert
	id, _ := SqlInsert(db, User{
		name:    "dummyman",
		stmt:    "I'll be deleted.",
		otakuId: 80,
	})

	// check insert
	SqlSelect(db)
	fmt.Println()

	// delete
	SqlDelete(db, fmt.Sprint(id))

	// check delete
	SqlSelect(db)

}

func main() {
	log.Println("delete demo")
	deleteDemo()

	fmt.Println()
	log.Println("update demo")
	updateDemo()

}
