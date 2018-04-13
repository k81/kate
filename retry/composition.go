package retry

// A composite strategy.  In order for Next or HasNext to
// succeed, all of the included strategies must succeed
type All []Strategy

func (s All) Next() bool {
	for _, ss := range s {
		if !ss.Next() {
			return false
		}
	}
	return true
}

func (s All) HasNext() bool {
	for _, ss := range s {
		if !ss.HasNext() {
			return false
		}
	}
	return true
}

// A composite strategy.  In order for Next or HasNext to
// succeed, any one of the included strategies must succeed
type Any []Strategy

func (s Any) Next() bool {
	// Call all strategies even one returns true
	// otherwise, they might lose count
	var succ = false
	for _, ss := range s {
		if ss.Next() {
			succ = true
		}
	}
	return succ
}

func (s Any) HasNext() bool {
	for _, ss := range s {
		if ss.HasNext() {
			return true
		}
	}
	return false
}

type AllResettable []ResettableStrategy

func (s AllResettable) Next() bool {
	for _, ss := range s {
		if !ss.Next() {
			return false
		}
	}
	return true
}

func (s AllResettable) HasNext() bool {
	for _, ss := range s {
		if !ss.HasNext() {
			return false
		}
	}
	return true
}

func (s AllResettable) Reset() {
	for _, ss := range s {
		ss.Reset()
	}
}

type AnyResettable []ResettableStrategy

func (s AnyResettable) Next() bool {
	// Call all strategies even one returns true
	// otherwise, they might lose count
	var succ = false
	for _, ss := range s {
		if ss.Next() {
			succ = true
		}
	}
	return succ
}

func (s AnyResettable) HasNext() bool {
	for _, ss := range s {
		if ss.HasNext() {
			return true
		}
	}
	return false
}

func (s AnyResettable) Reset() {
	for _, ss := range s {
		ss.Reset()
	}
}
