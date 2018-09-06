# fastperm

## `tl;dr`

Approaching the limits of shuffling performance in Go.

## Benchmark

```
goos: darwin
goarch: amd64
pkg: ajz.xyz/fastperm/perm
BenchmarkSlice/reference-4                  1000          29908776 ns/op            5376 B/op          1 allocs/op
BenchmarkSlice/package-4                    5000           5505743 ns/op               0 B/op          0 allocs/op
PASS
ok      ajz.xyz/fastperm/perm   61.331s
```

## Background

The Go library provides a shuffle implementation in `math/rand` suitable for generating good permutations.

```golang
arr := []int{1, 2, 3, 4, 5}
rand.Shuffle(len(arr), func(i, j int) {
    arr[i], arr[j] = arr[j], arr[i]
})
```

Of course, it just works, even edge cases like permutations of length `1<<32` and greater are handled. Wow! The implementation sources its random numbers from the default random source.

But it's possible to beat the library implementation with few tradeoffs. Here are some things I did to improve performance:

### New PRNG

I used a XOR shift PRNG with a 64 bit state. More details are available [here](http://vigna.di.unimi.it/ftp/papers/xorshift.pdf). This PRNG is favorable on machines with 64 bit word size. It is faster because the internal state is smaller than the library implementation. 

### Eliminate divisions

Divisions are among the slowest arithmetic operations out there. Computer scientists use modulo division to fit a number to a fixed range.

```golang
r.Uint64() % n // numbers in the range [0, n)
```

But it is possible to replace this call with a multiplication and a shift.

```golang
(r.Uint32() * n) >> 32
```

You can read more about this trick [here](https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/).

### Reduce stack footprint

I was able a cut time by changing the function signature to accept the slice directly.

```golang
// before
Shuffle(n int, swap func(i, j int))

// after
Slice(dst []Item)
```

The `Rand` type definition was changed from a `struct` with one field to a decorated `uint64`. This works since the state is a single number.

```golang
// before
type Rand struct {
    state uint64
}

// after
type Rand uint64
```