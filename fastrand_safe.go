// SPDX-License-Identifier: MIT
//
// Copyright 2023 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

//go:build safe

package fastrand

import "hash/maphash"

func u32() uint32 {
	return uint32(maphash.Bytes(maphash.MakeSeed(), nil) >> 32)
}

func u64() uint64 {
	return maphash.Bytes(maphash.MakeSeed(), nil)
}

func putU64(p []byte, v uint64) {
	_ = p[7] // Early bounds check to guarantee safety of writes below.
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}
