package user

import (
	"gowebapp/source/shared/database"
	"gowebapp/source/shared/password"
	"gowebapp/source/shared/testutil"
	"testing"
)

var randStr = testutil.RandAlphaStr

func init() {
	mi := &database.MySQLInfo{
		Username: "gomarket",
		Password: "password",
		Name:     "gomarket",
		Hostname: "localhost",
		Port:     3306,
	}
	database.Connect(mi)
}

func BenchmarkRegisterBenchmark32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Register(randStr(32), randStr(32), randStr(32))
	}
}

func TestRegister(t *testing.T) {
	l, e, p := randStr(8), randStr(16), randStr(16)
	err := Register(l, e, p)
	if err != nil {
		t.Error("Failed with error:", err)
	}
	u, err := GetByHandle(e)
	if err != nil {
		t.Error("Failed with error:", err)
	}
	if u.Login != l {
		t.Error(
			"For login",
			"expected", l,
			"got", u.Login)
	}
	if u.Email != e {
		t.Error(
			"For email",
			"expected", e,
			"got", u.Email)
	}
	if !password.Compare(p, u.PasswordHash) {
		t.Error("Password hash retrived from database is different from the one provided")
	}

	err = Register(l, e, p)
	if err == nil {
		t.Error(
			"For", l, e, p,
			"expected", ErrAlreadyRegistered,
			"got", err)
	}
}

func TestGetByHandle(t *testing.T) {
	l := randStr(16) // non existent login
	_, err := GetByHandle(l)
	if err == nil {
		t.Error(
			"For", l,
			"expected", ErrNotFound,
			"got", err)
	}
}
