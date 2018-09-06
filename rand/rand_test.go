/// xorshift64star Pseudo-Random Number Generator
/// This class is based on original code written and dedicated
/// to the public domain by Sebastiano Vigna (2014).
/// It has the following characteristics:
///
///  -  Outputs 64-bit numbers
///  -  Passes Dieharder and SmallCrush test batteries
///  -  Does not require warm-up, no zeroland to escape
///  -  Internal state is a single 64-bit integer
///  -  Period is 2^64 - 1
///  -  Speed: 1.60 ns/call (Core i7 @3.40GHz)
///
/// For further analysis see
///   <http://vigna.di.unimi.it/ftp/papers/xorshift.pdf>

package rand

import (
	gorand "math/rand"
	"testing"
)

func TestRandUint64(t *testing.T) {
	type fields struct {
		s uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{"seed 1337", fields{1337}, 4248603512886213755},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Rand{
				s: tt.fields.s,
			}
			if got := r.Uint64(); got != tt.want {
				t.Errorf("Rand.Uint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkRandom(b *testing.B) {
	type fields struct{ s int64 }
	tests := []struct {
		name string
		r    func(s int64) gorand.Source
	}{
		{"reference", func(s int64) gorand.Source { return gorand.NewSource(s) }},
		{"package", func(s int64) gorand.Source { r := Rand(uint64(s)); return &r }},
	}
	for _, tt := range tests {
		seeds := gorand.NewSource(5577006791947779410)
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := tt.r(seeds.Int63())
				for j := 0; j < 1000; j++ {
					r.Int63()
				}
			}
		})
	}
}
