package jsonschema

import "testing"

type Benchmark struct {
	Name string `json:"name" jsonschema:"maxLength:5,pattern:^[0-8]+$"`
}

func BenchmarkValidator_Validate(b *testing.B) {
	bench := &Benchmark{
		Name: "1234567890",
	}

	b.ResetTimer()
	b.ReportAllocs()
	validator := NewValidator()
	for i := 0; i < b.N; i++ {
		validator.Validate(bench)
	}
}

func BenchmarkValidator_Validate_Parallel(b *testing.B) {
	bench := &Benchmark{
		Name: "1234567890",
	}

	b.ResetTimer()
	b.ReportAllocs()
	validator := NewValidator()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			validator.Validate(bench)
		}
	})
}
