// @license
// Copyright (C) 2019  Dinko Korunic
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// +build darwin

package smc

import (
	"encoding/binary"
	"fmt"
	"github.com/panotza/gosmc"
	"math"
	"strings"
)

// FPConv type used for AppleFPConv map
type FPConv struct {
	Div    float32
	Signed bool
}

// AppleFPConv maps floating point type conversion constants and signedness property
var AppleFPConv = map[string]FPConv{
	"fp1f": FPConv{Div: 32768.0},
	"fp2e": FPConv{Div: 16384.0},
	"fp3d": FPConv{Div: 8192.0},
	"fp4c": FPConv{Div: 4096.0},
	"fp5b": FPConv{Div: 2048.0},
	"fp6a": FPConv{Div: 1024.0},
	"fp79": FPConv{Div: 512.0},
	"fp88": FPConv{Div: 256.0},
	"fpa6": FPConv{Div: 64.0},
	"fpc4": FPConv{Div: 16.0},
	"fpe2": FPConv{Div: 4.0},
	"sp1e": FPConv{Div: 16384.0, Signed: true},
	"sp2d": FPConv{Div: 8192.0, Signed: true},
	"sp3c": FPConv{Div: 4096.0, Signed: true},
	"sp4b": FPConv{Div: 2048.0, Signed: true},
	"sp5a": FPConv{Div: 1024.0, Signed: true},
	"sp69": FPConv{Div: 512.0, Signed: true},
	"sp78": FPConv{Div: 256.0, Signed: true},
	"sp87": FPConv{Div: 128.0, Signed: true},
	"sp96": FPConv{Div: 64.0, Signed: true},
	"spa5": FPConv{Div: 32.0, Signed: true},
	"spb4": FPConv{Div: 16.0, Signed: true},
	"spf0": FPConv{Div: 1.0, Signed: true},
}

// fpToFloat32 converts fp* SMC types to float32
func fpToFloat32(t string, x gosmc.SMCBytes, size uint32) (float32, error) {
	if v, ok := AppleFPConv[t]; ok {
		res := binary.BigEndian.Uint16(x[:size])
		if v.Signed {
			return float32(int16(res)) / v.Div, nil
		} else {
			return float32(res) / v.Div, nil
		}
	}

	return 0.0, fmt.Errorf("unable to convert to float32 type %q, bytes %v to float32", t, x)
}

// fltToFloat32 converts flt SMC type to float32
func fltToFloat32(k string, x gosmc.SMCBytes, size uint32) (float32, error) {
	return math.Float32frombits(binary.LittleEndian.Uint32(x[:size])), nil
}

// smcTypeToString converts UInt32Char array to regular Go string removing trailing null and whitespace
func smcTypeToString(x gosmc.UInt32Char) string {
	return strings.TrimRight(x.ToString(), "\x00 ")
}

// smcBytesToUint32 converts ui8/ui16/ui32 SMC types to uint32
func smcBytesToUint32(x gosmc.SMCBytes, size uint32) uint32 {
	var total uint32
	for i := uint32(0); i < size; i++ {
		total += uint32(x[i]) << ((size - 1 - i) * 8)
	}
	return total
}

// smcBytesToFloat32 converts ui8/ui16/ui32 SMC types to float32
func smcBytesToFloat32(x gosmc.SMCBytes, size uint32) float32 {
	return float32(smcBytesToUint32(x, size))
}
