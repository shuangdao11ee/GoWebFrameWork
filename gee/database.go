package gee

import (
	"database/sql"
	"math/rand"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//this struct is used to contribute a database pointer to engine
//it will be used everytime a information from wechat server come in
type DB struct {
	DbPointer *sql.DB
}

//connect to target database and create a database pointer
func (db *DB) Init(DBName string) {
	db.DbPointer, _ = sql.Open("sqlite3", DBName)
}

//create the corresponding id for each openid in database
func (db *DB) IDCreated(openid string) {
	stmt, _ := db.DbPointer.Prepare("INSERT INTO user(openid, id) values(?,?)")
	ID := GetRandomString(3)
	for db.GetOPENID(ID) != "" {
		ID = GetRandomString(3)
	}
	_, _ = stmt.Exec(openid, ID)
}

//select the openid by id from database, if openid doesn't exist, it return ""
func (db *DB) GetOPENID(id string) string {
	rows, _ := db.DbPointer.Query("SELECT * FROM user")
	for rows.Next() {
		var openid string
		var cId string
		_ = rows.Scan(&openid, &cId)
		if cId == id {
			return openid
		}
	}
	return ""
}

//select the id by openid from database, if id doesn't exist, it return ""
func (db *DB) GetID(openid string) string {
	rows, _ := db.DbPointer.Query("SELECT * FROM user")
	for rows.Next() {
		var id string
		var cOpenid string
		_ = rows.Scan(&cOpenid, &id)
		if cOpenid == openid {
			return id
		}
	}
	return ""
}

//create a random string which length is confirmed
func GetRandomString(l int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
