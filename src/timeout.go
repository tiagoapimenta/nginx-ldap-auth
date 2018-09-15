package main

import (
	"sort"
	"sync"
	"time"
)

type passtimer struct {
	password string
	timer    *time.Timer
}

type userpass struct {
	correct *passtimer
	wrong   []passtimer
}

var (
	passwords = map[string]*userpass{}
	mutex     = sync.RWMutex{}
)

func containsWrongPassword(data *userpass, password string) (int, bool) {
	size := len(data.wrong)
	if size == 0 {
		return 0, false
	}

	pos := sort.Search(size, func(i int) bool {
		return data.wrong[i].password >= password
	})

	return pos, pos < size &&
		data.wrong[pos].password == password
}

func getCache(username, password string) (bool, bool) {
	defer mutex.RUnlock()
	mutex.RLock()

	data, found := passwords[username]
	if !found {
		return false, false
	}

	if data.correct != nil && (*data.correct).password == password {
		return true, true
	}

	_, found = containsWrongPassword(data, password)

	return false, found
}

func putCache(username, password string, ok bool) {
	defer mutex.Unlock()
	mutex.Lock()

	data, found := passwords[username]
	if !found {
		data = &userpass{}
		passwords[username] = data
	}

	timeout := config.Timeout.Wrong
	if ok {
		timeout = config.Timeout.Success
	}

	pass := passtimer{
		password: password,
		timer: time.AfterFunc(timeout, func() {
			removeCache(username, password, ok)
		}),
	}

	if ok {
		if data.correct != nil {
			data.correct.timer.Stop()
		}
		data.correct = &pass
	} else {
		pos, found := containsWrongPassword(data, password)
		if found {
			data.wrong[pos].timer.Stop()
		} else {
			data.wrong = append(data.wrong, passtimer{})
			copy(data.wrong[pos+1:], data.wrong[pos:])
		}
		data.wrong[pos] = pass
	}
}

func removeCache(username, password string, ok bool) {
	defer mutex.Unlock()
	mutex.Lock()

	data, found := passwords[username]
	if !found {
		return
	}

	if ok {
		if data.correct != nil {
			data.correct.timer.Stop()
			data.correct = nil
		}
	} else {
		pos, found := containsWrongPassword(data, password)
		if found {
			data.wrong[pos].timer.Stop()
			data.wrong = data.wrong[:pos+copy(data.wrong[pos:], data.wrong[pos+1:])]
		}
	}

	if data.correct == nil && len(data.wrong) == 0 {
		delete(passwords, username)
	}
}
