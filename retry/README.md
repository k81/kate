An interface for and collection of various implementations of retry behavior, designed to be easy to use and composable.

Some examples:

```go
// A very basic example.
strategy := &retry.CountStrategy{Tries: 3}
for done := false; !done && strategy.Next() {
	done = trySomething()
}
```

```go
// Compose a few strategies together
// Try not more than 3 times, with a 100ms delay between attempts
strategy := &retry.All{
	&retry.CountStrategy{Tries: 3},
	&retry.DelayStrategy{Delay: 100 * time.Millisecond},
}

for done := false; !done && strategy.Next() {
	done = trySomething()
}
```

```go
// More complex composition:
// Delay of 100ms between tries
// Keep trying for up to 10s
// At least 3 tries
strategy := &retry.All{
	&retry.Any{
		&retry.MaximumTimeStrategy{Duration: 10 * time.Second},
		&retry.CountStrategy{Tries: 3},
	},
	&retry.DelayStrategy{Delay: 100 * time.Millisecond},
}

for done := false; !done && strategy.Next() {
	done = trySomething()
}
```

```go
// Separate retry logic from a more complex function
func doComplexThings(strategy retry.Strategy)bool{
	for strategy.Next() {
		if success := trySomthingThatMightFail(); success {
			return true
		}
	}
	return false
}

doComplexThings(&retry.CountStrategy{Tries: 3})
doComplexThings(&retry.MaximumTimeStrategy{Duration: 10 * time.Second})
```

