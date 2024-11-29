package lock

type LockSupport struct {
	wt    chan bool
	state chan bool
}

func NewLockSupport() *LockSupport {
	return &LockSupport{
		wt:    make(chan bool),
		state: make(chan bool),
	}
}

func (l *LockSupport) Park() bool {
	if len(l.state) >= 1 {
		return false
	}
	l.state <- true
	return <-l.wt
}

func (l *LockSupport) Unpark() (err error) {
	defer func() error {
		e := recover()
		if e != nil {
			err = e.(error)
		}
		if err != nil {
			return err
		}
		return nil
	}()
	<-l.state
	if len(l.wt) >= 1 {
		return nil
	}
	l.wt <- true
	return nil
}
