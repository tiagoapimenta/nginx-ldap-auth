package data

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

type Storage struct {
	passwords map[string]*userpass
	lock      sync.RWMutex
	success   time.Duration
	wrong     time.Duration
}

func NewStorage(success, wrong time.Duration) *Storage {
	return &Storage{
		passwords: map[string]*userpass{},
		lock:      sync.RWMutex{},
		success:   success,
		wrong:     wrong,
	}
}

func (p *Storage) Get(username, password string) (bool, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	data, found := p.passwords[username]
	if !found {
		return false, false
	}

	if data.correct != nil && (*data.correct).password == password {
		return true, true
	}

	_, found = containsWrongPassword(data, password)

	return false, found
}

func (p *Storage) Put(username, password string, ok bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	data, found := p.passwords[username]
	if !found {
		data = &userpass{}
		p.passwords[username] = data
	}

	timeout := p.wrong
	if ok {
		timeout = p.success
	}

	pass := passtimer{
		password: password,
		timer: time.AfterFunc(timeout, func() {
			p.remove(username, password, ok)
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

func (p *Storage) remove(username, password string, ok bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	data, found := p.passwords[username]
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
		delete(p.passwords, username)
	}
}

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
