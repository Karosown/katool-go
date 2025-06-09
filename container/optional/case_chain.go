package optional

type OptSwitch struct {
	isRun      bool
	isBreak    bool
	lastResult bool
}

func (t *OptSwitch) Case(flag bool, run ...func()) *OptSwitch {
	if t.isBreak {
		return t
	}
	if flag {
		t.lastResult = true
		run[0]()
	} else {
		t.lastResult = false
		run[1]()
	}
	return t
}

func (t *OptSwitch) CaseFunc(flag func() bool, run ...func()) *OptSwitch {
	return t.Case(flag(), run...)
}
func (t *OptSwitch) Break() *OptSwitch {
	if t.lastResult {
		t.isBreak = true
	}
	return t
}
func (t *OptSwitch) Default(run ...func()) *OptSwitch {
	if !t.isRun {
		for _, fn := range run {
			fn()
		}
	}
	return t
}
