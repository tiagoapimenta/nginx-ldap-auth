package data

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

const (
	username1 = "Alice"
	username2 = "James"
	password1 = "master"
	password2 = "shadow"
	password3 = "qwerty"
	success   = time.Second / 2
	wrong     = time.Second / 5
)

func printPassMap(t *testing.T, storage *Storage, prefix string) {
	buffer := bytes.Buffer{}
	first := true
	for k, v := range storage.passwords {
		if first {
			first = false
		} else {
			buffer.WriteByte(',')
		}
		correct := "<nil>"
		if v.correct != nil {
			correct = v.correct.password
		}
		fmt.Fprintf(&buffer, "%s:{correct:%s,wrong:%+v}", k, correct, v.wrong)
	}
	t.Logf("%s passwords: %s\n", prefix, buffer.String())
}

func testCache(t *testing.T, storage *Storage, id int, username, password string, eok, efound bool) {
	ok, found := storage.Get(username, password)
	if ok != eok || found != efound {
		t.Errorf("Test %d expected (%v %v) given (%v %v)\n", id, eok, efound, ok, found)
	}
}

func TestPasswordTimeout(t *testing.T) {
	storage := NewStorage(success, wrong)

	testCache(t, storage, 0, username1, password1, false, false)
	testCache(t, storage, 1, username1, password2, false, false)
	testCache(t, storage, 2, username1, password3, false, false)
	testCache(t, storage, 3, username2, password1, false, false)
	testCache(t, storage, 4, username2, password2, false, false)
	testCache(t, storage, 5, username2, password3, false, false)

	storage.Put(username1, password1, true)
	storage.Put(username1, password3, false)
	printPassMap(t, storage, "add")

	testCache(t, storage, 6, username1, password1, true, true)
	testCache(t, storage, 7, username1, password2, false, false)
	testCache(t, storage, 8, username1, password3, false, true)
	testCache(t, storage, 9, username2, password1, false, false)
	testCache(t, storage, 10, username2, password2, false, false)
	testCache(t, storage, 11, username2, password3, false, false)

	time.Sleep(wrong + wrong/2)
	printPassMap(t, storage, "timed")

	testCache(t, storage, 12, username1, password1, true, true)
	testCache(t, storage, 13, username1, password2, false, false)
	testCache(t, storage, 14, username1, password3, false, false)
	testCache(t, storage, 15, username2, password1, false, false)
	testCache(t, storage, 16, username2, password2, false, false)
	testCache(t, storage, 17, username2, password3, false, false)

	time.Sleep(success - wrong)
	printPassMap(t, storage, "expired")

	testCache(t, storage, 18, username1, password1, false, false)
	testCache(t, storage, 19, username1, password2, false, false)
	testCache(t, storage, 20, username1, password3, false, false)
	testCache(t, storage, 21, username2, password1, false, false)
	testCache(t, storage, 22, username2, password2, false, false)
	testCache(t, storage, 23, username2, password3, false, false)
}
