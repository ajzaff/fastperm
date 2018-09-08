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
			r := Rand(tt.fields.s)
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

// allow for Kolmogorov-Smirnov tolerance
const (
	n = 64
	k = 1e7

	tolerance = 5.4e-5
)

// BenchmarkFairness use Kolmogorov-Smirnov to test fairness
func BenchmarkFairness(b *testing.B) {
	seeds := gorand.NewSource(5577006791947779410)
	hist := make([]float64, n)
	for i := 0; i < b.N; i++ {
		// build the histogram
		r := Rand(seeds.Int63())
		for j := 0; j < k; j++ {
			x := r.Uint64()
			for i := 0; x > 0; x >>= 1 {
				hist[i] += float64(x & 1)
				i++
			}
		}
	}
	s := float64(0)
	for i := 0; i < n; i++ {
		// sum the hist
		s += hist[i]
	}
	for i := 0; i < n; i++ {
		// normalize the hist
		hist[i] /= s
	}
	cdf := make([]float64, n)
	edf := make([]float64, n)
	cdf[0] = 1 / float64(n)
	edf[0] = hist[0]
	for i := 1; i < n; i++ {
		// build the cdf & edf
		cdf[i] = 1/float64(n) + cdf[i-1]
		edf[i] = hist[i] + edf[i-1]
	}
	tol := float64(0)
	for i := 0; i < n; i++ {
		// find the max tolerance
		v := cdf[i] - edf[i]
		if v < 0 {
			v = -v
		}
		if v > tol {
			tol = v
		}
	}
	if tol > tolerance {
		b.Errorf("KS above acceptable threshold: %f > %f", tol, tolerance)
	}
}
