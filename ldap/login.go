package ldap

func (p *Pool) Validate(username, password string) (bool, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	err := p.auth()
	if err != nil {
		return false, err
	}

	p.admin = false
	err = p.networkJail(func() error {
		return p.conn.Bind(username, password)
	})
	if err != nil {
		return false, err
	}

	err = p.auth()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (p *Pool) auth() error {
	if p.admin || p.username == "" && p.password == "" {
		return nil
	}

	err := p.networkJail(func() error {
		return p.conn.Bind(p.username, p.password)
	})
	if err == nil {
		p.admin = true
	}
	return err
}
