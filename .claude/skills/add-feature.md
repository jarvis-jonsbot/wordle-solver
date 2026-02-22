# Add Feature

## Steps

1. Create or modify files in `internal/<package>/`
2. Write tests first or alongside — every exported function needs a test
3. Run `make test` and `make lint` before committing
4. Use `fmt.Errorf("context: %w", err)` for error wrapping
5. Commit with conventional commit: `feat(<scope>): <description>`

## Testing Pattern

```go
func TestFeatureName(t *testing.T) {
    tests := []struct {
        name string
        // inputs
        // expected
    }{
        {"case 1", ...},
        {"case 2", ...},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // act + assert
        })
    }
}
```
