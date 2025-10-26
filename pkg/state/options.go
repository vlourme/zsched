package state

type StateParameter struct {
	value any
}

// Get returns the parameter by name
func (c *State) Get(name string) StateParameter {
	p, ok := c.Parameters[name]
	if !ok {
		return StateParameter{nil}
	}
	return StateParameter{p}
}

// String returns the string value of the parameter
func (p StateParameter) String(defaultValue ...string) string {
	s, ok := p.value.(string)
	if !ok {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return s
}

// Int returns the int value of the parameter
func (p StateParameter) Int(defaultValue ...int) int {
	i, ok := p.value.(float64)
	if !ok {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return int(i)
}

// Float returns the float value of the parameter
func (p StateParameter) Float(defaultValue ...float64) float64 {
	f, ok := p.value.(float64)
	if !ok {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return f
}

// Bool returns the bool value of the parameter
func (p StateParameter) Bool(defaultValue ...bool) bool {
	b, ok := p.value.(bool)
	if !ok {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	return b
}
