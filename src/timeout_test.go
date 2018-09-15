package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"
)

const (
	username1 = "Alice"
	username2 = "James"
	password1 = "master"
	password2 = "shadow"
	password3 = "qwerty"
)

func TestMain(m *testing.M) {
	config.Timeout.Success = time.Second / 2
	config.Timeout.Wrong = time.Second / 5
	os.Exit(m.Run())
}

func printPassMap(t *testing.T, prefix string) {
	buffer := bytes.Buffer{}
	first := true
	for k, v := range passwords {
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

func testCache(t *testing.T, id int, username, password string, eok, efound bool) {
	ok, found := getCache(username, password)
	if ok != eok || found != efound {
		t.Errorf("Test %d expected (%v %v) given (%v %v)\n", id, eok, efound, ok, found)
	}
}

func TestPasswordTimeout(t *testing.T) {
	testCache(t, 0, username1, password1, false, false)
	testCache(t, 1, username1, password2, false, false)
	testCache(t, 2, username1, password3, false, false)
	testCache(t, 3, username2, password1, false, false)
	testCache(t, 4, username2, password2, false, false)
	testCache(t, 5, username2, password3, false, false)

	putCache(username1, password1, true)
	putCache(username1, password3, false)
	printPassMap(t, "add")

	testCache(t, 6, username1, password1, true, true)
	testCache(t, 7, username1, password2, false, false)
	testCache(t, 8, username1, password3, false, true)
	testCache(t, 9, username2, password1, false, false)
	testCache(t, 10, username2, password2, false, false)
	testCache(t, 11, username2, password3, false, false)

	time.Sleep(config.Timeout.Wrong + config.Timeout.Wrong/2)
	printPassMap(t, "timed")

	testCache(t, 12, username1, password1, true, true)
	testCache(t, 13, username1, password2, false, false)
	testCache(t, 14, username1, password3, false, false)
	testCache(t, 15, username2, password1, false, false)
	testCache(t, 16, username2, password2, false, false)
	testCache(t, 17, username2, password3, false, false)

	time.Sleep(config.Timeout.Success - config.Timeout.Wrong)
	printPassMap(t, "expired")

	testCache(t, 18, username1, password1, false, false)
	testCache(t, 19, username1, password2, false, false)
	testCache(t, 20, username1, password3, false, false)
	testCache(t, 21, username2, password1, false, false)
	testCache(t, 22, username2, password2, false, false)
	testCache(t, 23, username2, password3, false, false)
}
