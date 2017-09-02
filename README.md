# go-jsonschema-validator

[![Build Status](https://travis-ci.org/yu-ichiko/go-jsonschema-validator.svg?branch=master)](https://travis-ci.org/yu-ichiko/go-jsonschema-validator)

supported draft 4 and draft 6

```go
type Sample struct {
	Name string   `json:"name" jsonschema:"pattern:[a-zA-Z0-9],maxLength:20"`
	Age  int      `json:"age" jsonschema:"minimum:0,maximum:20"`
}

sample := Sample{
	Name: "test",
	Age:  10,
}

validator := jsonschema.NewValidator()
validator.AddFormat("my-format", func(value *reflect.Value, field *reflect.StructField) (err error) {
	...
	return
})
err := validator.Validate(sample)
```
