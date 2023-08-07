// SPDX-License-Identifier: MIT
//
// Copyright 2023 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

//go:build !safe

package fastrand

import (
	"unsafe"
)

//go:linkname u32 runtime.fastrand
func u32() uint32

//go:linkname u64 runtime.fastrand64
func u64() uint64

func putU64(p []byte, v uint64) {
	*(*uint64)(unsafe.Pointer(&p[0])) = v
}
