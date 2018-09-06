package perm

import (
	"reflect"
	"testing"

	gorand "math/rand"

	"ajz.xyz/fastperm/rand"
)

func TestSlice(t *testing.T) {
	type args struct {
		r   rand.Rand
		dst []Item
	}
	tests := []struct {
		name string
		args args
		want []Item
	}{
		{"1..10", args{
			rand.Rand(5577006791947779410),
			[]Item{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			[]Item{1, 8, 2, 5, 9, 6, 7, 4, 3, 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := Rand(tt.args.r)
			rp.Slice(tt.args.dst)
			if ok := reflect.DeepEqual(tt.args.dst, tt.want); !ok {
				t.Errorf("Rand.Slice() = %v, want %v", tt.args.dst, tt.want)
			}
		})
	}
}

func BenchmarkSlice(b *testing.B) {
	tests := []struct {
		name string
		f    func(b *testing.B, t [][]Item)
	}{
		{"reference", func(b *testing.B, t [][]Item) {
			r := gorand.New(gorand.NewSource(5577006791947779410))
			for i := 0; i < len(t); i++ {
				r.Shuffle(len(t[i]), func(j, k int) { t[i][j], t[i][k] = t[i][k], t[i][j] })
			}
		}},
		{"package", func(b *testing.B, t [][]Item) {
			r := Rand(rand.Rand(5577006791947779410))
			for i := 0; i < len(t); i++ {
				r.Slice([]Item(t[i]))
			}
		}},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			r := gorand.NewSource(5577006791947779410)
			t := make([][]Item, 1000)
			for i := 0; i < len(t); i++ {
				t[i] = make([]Item, 1000)
				for j := 0; j < len(t[i]); j++ {
					t[i][j] = r.Int63()
				}
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tt.f(b, t)
			}
		})
	}
}
