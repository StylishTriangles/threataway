package env

import (
	"threataway/source/shared/database"
	"time"
)

// Env is environment variable
type Env struct {
	ID    uint32
	name  string
	value string
}

// GetEnv returns value of environment variable
func GetEnv(name string) string {
	stmt, err := database.DB.Prepare("SELECT value FROM env WHERE name = ?")
	if err != nil {
		return ""
	}
	defer stmt.Close()
	ret := ""
	err = stmt.QueryRow(name).Scan(&ret)
	if err != nil {
		return ""
	}
	return ret
}

// SetEnv sets environment variable
func SetEnv(name, value string) {
	DeleteEnv(name)
	stmt, _ := database.DB.Prepare("INSERT INTO env(name, value) VALUES(?, ?)")
	defer stmt.Close()
	stmt.Exec(name, value)
}

// DeleteEnv deletes environment variable
func DeleteEnv(name string) {
	stmt, _ := database.DB.Prepare("DELETE FROM env WHERE name = ?")
	defer stmt.Close()
	stmt.Exec(name)
}

var timedFunc func()
var index = 0

// ChangeTimeout changes timeout of
func ChangeTimeout(minutes string, f func()) {
	index++
	timedFunc = func() {
		ix := index
		for ix == index {
			f()
			t, _ := time.ParseDuration(minutes + "m")
			time.Sleep(t)
		}
	}
	go timedFunc()
}
