// SPDX-License-Identifier: MIT
//
// Copyright 2023 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

// Package fastrand provides quickly generated pseudo-random numbers
// with no repeatability guarantees on the stream of values.
package fastrand

import (
	"io"

	"golang.org/x/exp/constraints"
)

const (
	maxUint32 = (1 << 32) - 1
	maxUint64 = (1 << 64) - 1
	maxInt32  = maxUint32 >> 1
	maxInt64  = maxUint64 >> 1
)

// Float32 returns a pseudo-random float32 in the half-open interval [0,n).
func Float32() float32 {
	const (
		mask = 1<<24 - 1
		mult = 0x1.0p-24
	)
	return float32(u32()&mask) * mult
}

// Float64 returns a pseudo-random float64 in the half-open interval [0,n).
func Float64() float64 {
	const (
		mask = 1<<53 - 1
		mult = 0x1.0p-53
	)
	return float64(u64()&mask) * mult
}

// Int31 returns a non-negative pseudo-random int32.
func Int31() int32 {
	return int32(u32() >> 1)
}

// Int31n returns a non-negative pseudo-random int32 in the half-open interval [0,n).
// It panics if n <= 0.
func Int31n(n int32) int32 {
	if n <= 0 {
		panic("fastrand.Int31n: invalid argument")
	}
	if n&(n-1) == 0 { // n is power of two, can mask
		return Int31() & (n - 1)
	}
	max := maxInt32 - maxInt32%n
	v := Int31()
	for v >= max {
		v = Int31()
	}
	return v % n
}

// Int63 returns a non-negative pseudo-random int64.
func Int63() int64 {
	return int64(u64() >> 1)
}

// Int63n returns a non-negative pseudo-random int64 in the half-open interval [0,n).
// It panics if n <= 0.
func Int63n(n int64) int64 {
	if n <= 0 {
		panic("fastrand.Int63n: invalid argument")
	}
	if n&(n-1) == 0 { // n is power of two, can mask
		return Int63() & (n - 1)
	}
	max := maxInt64 - maxInt64%n
	v := Int63()
	for v >= max {
		v = Int63()
	}
	return v % n
}

// Uint32 returns a pseudo-random uint32.
func Uint32() uint32 {
	return u32()
}

// Uint32n returns a pseudo-random uint32 in the half-open interval [0,n).
func Uint64nUint32n(n uint32) uint32 {
	if n&(n-1) == 0 { // n is power of two, can mask
		return u32() & (n - 1)
	}
	max := maxUint32 - maxUint32%uint32(n)
	v := u32()
	for v >= max {
		v = u32()
	}
	return v % n
}

// Uint64 returns a pseudo-random uint64.
func Uint64() uint64 {
	return u64()
}

// Uint64n returns a pseudo-random uint64 in the half-open interval [0,n).
func Uint64n(n uint64) uint64 {
	if n&(n-1) == 0 { // n is power of two, can mask
		return u64() & (n - 1)
	}
	max := maxUint64 - maxUint64%uint64(n)
	v := u64()
	for v >= max {
		v = u64()
	}
	return v % n
}

// A Real is a real number.
type Real interface {
	constraints.Signed | constraints.Unsigned
}

// Jitter returns a pseudo-random value in the interval [v - factor*v, v + factor*v].
func Jitter[T Real](v T, factor float64) T {
	r := Float64()
	// r = [0, 1)
	// 2*r = [0, 2)
	// 2*r - 1 = [-1, 1)
	// j*(2*r - 1) = [-j, j)
	// 1 + j*(2*r - 1) = [1 - j, 1 + j)
	// b*(1 + j*(2*r - 1)) = [b - j*b, b + j*b)
	return T(float64(v) * (1 + (factor * (2*r - 1))))
}

// Shuffle pseudo-randomizes the order of elements in s.
func Shuffle[E any](s []E) {
	// Fisher-Yates shuffle: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	i := len(s) - 1
	// Shuffle really ought not be called with slice indices that requires more than 31 bits.
	// Nevertheless, handle it as best we can.
	for ; i >= maxInt32-1; i-- {
		j := Int63n(int64(i + 1))
		s[i], s[j] = s[j], s[i]
	}
	// Switch to 31-bit indices.
	for ; i > 0; i-- {
		j := Int31n(int32(i + 1))
		s[i], s[j] = s[j], s[i]
	}
}

var ioReader io.Reader = &reader{}

// Reader returns an io.Reader that fills the read buffer with
// pseudo-random bytes and never returns an error.
func Reader() io.Reader {
	return ioReader
}

type reader struct{}

func (*reader) Read(p []byte) (int, error) {
	Fill(p)
	return len(p), nil
}

// Fill fills b with pseudo-random bytes.
func Fill(p []byte) {
	for len(p) >= 8 {
		putU64(p, u64())
		p = p[8:]
	}
	switch {
	case len(p) > 4:
		fill(p, u64())
	case len(p) > 0:
		fill(p, u32())
	}
}

func fill[T interface{ uint32 | uint64 }](p []byte, v T) {
	for i := range p {
		p[i] = byte(v >> (i * 8))
	}
}
