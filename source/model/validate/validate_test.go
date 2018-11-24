package validate

import "testing"

type multiTest struct {
	pass string
	ok   bool
	err  error
}

var passArray = []multiTest{
	// common passwords
	{"$andmann", false, ErrCommonPassword},
	{"000000000000", false, ErrCommonPassword},
	{"zxcvbnm123456789", false, ErrCommonPassword},
	{"%E2%82%AC", false, ErrCommonPassword},
	{"%%passwo", false, ErrCommonPassword},
	{"helloworld", false, ErrCommonPassword},
	// similar to common password but using unicode
	{"pąssword", true, nil},
	{"ilovemöm", true, nil},
	// short passwords
	{"", false, ErrEmpty},
	{"a", false, ErrTooShort},
	{"asdfg", false, ErrTooShort},
	{"ääää", true, nil}, // 8 bytes
	// long passwords
	{"qwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnm0123", true, nil},              // 56 characters
	{"qwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnm01234", false, ErrTooLong},     // 57 characters
	{"АаБбВвГгДдЕеЁёЖжЗзИиЙйКкЛлКк", true, nil},                                          // 28 characters 56 bytes
	{"АаБбВв\000ГгДдЕеЁёЖжЗзИиЙйКкЛлМмНнОоПпРрСсТтУуФфХхЦцЧчШшЩщЪъЫ", false, ErrTooLong}, // 68 characters (with null)
	// random passwords
	{"correct horse battery staple", true, nil},
	{"Battle_Ounce_Holy_Cross_5", true, nil},
	{"h3LL0world!", true, nil}, // if you use passwords like this you deserve to be hacked
	{"aAaAaAaAAa", true, nil},
}

var usernameArray = []multiTest{
	// unallowed characters
	{"$andmann", false, ErrInvalidFormat},
	{"%E2%82%AC", false, ErrInvalidFormat},
	{"%%passwo", false, ErrInvalidFormat},
	{"hello'world", false, ErrInvalidFormat},
	{"EtTuBrute?", false, ErrInvalidFormat},
	{"ääää", false, ErrInvalidFormat},
	// short
	{"", false, ErrEmpty},
	{"a", true, nil},
	{"asdfg", true, nil},
	{".", true, nil},
	// long
	{"qwertyuiopasdfghjklzxcvbnm", true, nil},                                   // 26 chars
	{"qwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnm", false, ErrTooLong}, // 52 chars
	// random
	{"000000000000", true, nil},
	{"zxcvbnm123456789", true, nil},
}

var emailArray = []multiTest{
	// incorrect
	{"test", false, ErrInvalidFormat},
	{"", false, ErrEmpty},
	{"aaa@", false, ErrInvalidFormat},
	{"aaa@a.a", false, ErrInvalidFormat},
	// random
	{"aaa@a.aa", true, nil},
	{"test@test.com", true, nil},
}

var nameArray = []multiTest{
	// incorrect names
	{"", false, ErrEmpty},
	{"abcdefghijklmnopqrstuvwqyzabcdefghijklmnopqrstuvwqyz", false, ErrTooLong},
	// correct names
	{"Steve", true, nil},
}

func multiTestWrapper(t *testing.T, m []multiTest, f func(string) (bool, error)) {
	for _, pair := range m {
		v, e := f(pair.pass)
		if pair.ok != v {
			t.Error(
				"For", pair.pass,
				"expected", pair.ok,
				"got", v)
		}
		if pair.err != e {
			t.Error(
				"For", pair.pass,
				"expected", pair.err,
				"got", e)
		}
	}
}

func TestPassword(t *testing.T) {
	// change password path to be able to load passwords file if test is run from
	// directory of this file
	commonPasswordsPath = "../../../" + commonPasswordsPath
	multiTestWrapper(t, passArray, Password)
}

func TestUsername(t *testing.T) {
	multiTestWrapper(t, usernameArray, Username)
}

func TestEmail(t *testing.T) {
	// create long string
	str := ""
	for i := 0; i < 256; i++ {
		str += "a"
	}
	emailArray = append(emailArray, multiTest{str, false, ErrTooLong})
	multiTestWrapper(t, emailArray, Email)
}

func TestName(t *testing.T) {
	multiTestWrapper(t, nameArray, Name)
}
