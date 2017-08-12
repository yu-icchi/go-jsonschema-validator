# go-jsonschema-validator

supported draft 4

```go
type Sample struct {
	Name string   `json:"name" jsonschema:"required,pattern:[a-zA-Z0-9],maxLength:20"`
	Age  int      `json:"age" jsonschema:"required,minimum:0,maximum:20"`
	Arr  []string `json:"arr" jsonschema:"enum:[a,b,c]"`
}

sample := Sample{
	Name: "a-test",
	Age:  10,
}

validator := jsonschema.NewValidator()
validator.AddFormat("my-format", func(data interface{}) (valid bool) {
	...
	return
})
err := validator.Validate(sample)
```

### Number

- int [8,16,32,64]
- float [32,64]
- uint [8,16,32,64]

#### maximum

```go
Num int `json:"num" jsonschema:"maximum:100"`
Num int `json:"num" jsonschema:"maximum:100,exclusiveMaximum:true"` // draft 4
```

#### minimum

```go
Num int `json:"num" jsonschema:"minimum:0"`
Num int `json:"num" jsonschema:"minimum:0,exclusiveMinimum:true"` // draft 4
```

#### multipleOf

```go
Num int `json:"num" jsonschema:"multipleOf:5"`
```

### String

- string

#### maxLength

```go
Str string `json:"str" jsonschema:"maxLength:10"`
```

#### minLength

```go
Str string `json:"str" jsonschema:"minLength:100"`
```

#### pattern

```go
Str string `json:"str" jsonschema:"pattern:[a-zA-Z0-9]+"`
```

#### format

```go
Str string `json:"str" jsonschema:"format:my-format"`
```

### Array

- Array
- slice

#### maxItems

```go
Arr []string `json:"arr" jsonschema:"maxItems:5"`
```

#### minItems

```go
Arr []string `json:"arr" jsonschema:"minItems:1"`
```

#### uniqueItems

```go
Arr []string `json:"arr" jsonschema:"uniqueItems:true"`
```

### Object

- map
- struct

#### maxProperties

map only

#### minProperties

map only

#### required

### All

#### enum
