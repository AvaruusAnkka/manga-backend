package postgresql

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func getDb() (db *sql.DB, err error) {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("user='%s' password=%s host=%s dbname='%s'", user, pass, host, name)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetUser(c *gin.Context) {
	user := User{}
	var users []User

	db, err := getDb()
	if err != nil {
		c.String(400, err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		c.String(400, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&user.ID, &user.Name)
		fmt.Println("ID:" + user.ID)
		fmt.Println("Name: " + user.Name)
		users = append(users, user)
	}

	c.JSON(200, users)
}

func AddUser(c *gin.Context) {
	name := c.PostForm("name")

	if name == "" {
		c.String(400, "Missing name")
		return
	}

	db, err := getDb()
	if err != nil {
		c.String(400, err.Error())
		return
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO users (name) VALUES ($1)", name)
	if err != nil {
		c.String(400, "Insert problem")
	}

	c.JSON(200, res)
}
