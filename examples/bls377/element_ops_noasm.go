// +build !amd64

// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by goff (v0.3.1) DO NOT EDIT

// Package fp contains field arithmetic operations
package fp

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular,
// there is no security guarantees such as constant time implementation
// or side-channel attack resistance
// /!\ WARNING /!\

import "math/bits"

func Mul(z, x, y *Element) {
	_mulGenericElement(z, x, y)
}

func Square(z, x *Element) {
	_squareGenericElement(z, x)
}

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func FromMont(z *Element) {
	_fromMontGenericElement(z)
}

// Add z = x + y mod q
func Add(z, x, y *Element) {
	var carry uint64

	z[0], carry = bits.Add64(x[0], y[0], 0)
	z[1], carry = bits.Add64(x[1], y[1], carry)
	z[2], carry = bits.Add64(x[2], y[2], carry)
	z[3], carry = bits.Add64(x[3], y[3], carry)
	z[4], carry = bits.Add64(x[4], y[4], carry)
	z[5], _ = bits.Add64(x[5], y[5], carry)

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[5] < 121098312706494698 || (z[5] == 121098312706494698 && (z[4] < 14284016967150029115 || (z[4] == 14284016967150029115 && (z[3] < 1883307231910630287 || (z[3] == 1883307231910630287 && (z[2] < 2230234197602682880 || (z[2] == 2230234197602682880 && (z[1] < 1660523435060625408 || (z[1] == 1660523435060625408 && (z[0] < 9586122913090633729))))))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 9586122913090633729, 0)
		z[1], b = bits.Sub64(z[1], 1660523435060625408, b)
		z[2], b = bits.Sub64(z[2], 2230234197602682880, b)
		z[3], b = bits.Sub64(z[3], 1883307231910630287, b)
		z[4], b = bits.Sub64(z[4], 14284016967150029115, b)
		z[5], _ = bits.Sub64(z[5], 121098312706494698, b)
	}
}

// Double z = x + x mod q, aka Lsh 1
func Double(z, x *Element) {
	var carry uint64

	z[0], carry = bits.Add64(x[0], x[0], 0)
	z[1], carry = bits.Add64(x[1], x[1], carry)
	z[2], carry = bits.Add64(x[2], x[2], carry)
	z[3], carry = bits.Add64(x[3], x[3], carry)
	z[4], carry = bits.Add64(x[4], x[4], carry)
	z[5], _ = bits.Add64(x[5], x[5], carry)

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[5] < 121098312706494698 || (z[5] == 121098312706494698 && (z[4] < 14284016967150029115 || (z[4] == 14284016967150029115 && (z[3] < 1883307231910630287 || (z[3] == 1883307231910630287 && (z[2] < 2230234197602682880 || (z[2] == 2230234197602682880 && (z[1] < 1660523435060625408 || (z[1] == 1660523435060625408 && (z[0] < 9586122913090633729))))))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 9586122913090633729, 0)
		z[1], b = bits.Sub64(z[1], 1660523435060625408, b)
		z[2], b = bits.Sub64(z[2], 2230234197602682880, b)
		z[3], b = bits.Sub64(z[3], 1883307231910630287, b)
		z[4], b = bits.Sub64(z[4], 14284016967150029115, b)
		z[5], _ = bits.Sub64(z[5], 121098312706494698, b)
	}
}

// Sub  z = x - y mod q
func Sub(z, x, y *Element) {
	var b uint64
	z[0], b = bits.Sub64(x[0], y[0], 0)
	z[1], b = bits.Sub64(x[1], y[1], b)
	z[2], b = bits.Sub64(x[2], y[2], b)
	z[3], b = bits.Sub64(x[3], y[3], b)
	z[4], b = bits.Sub64(x[4], y[4], b)
	z[5], b = bits.Sub64(x[5], y[5], b)
	if b != 0 {
		var c uint64
		z[0], c = bits.Add64(z[0], 9586122913090633729, 0)
		z[1], c = bits.Add64(z[1], 1660523435060625408, c)
		z[2], c = bits.Add64(z[2], 2230234197602682880, c)
		z[3], c = bits.Add64(z[3], 1883307231910630287, c)
		z[4], c = bits.Add64(z[4], 14284016967150029115, c)
		z[5], _ = bits.Add64(z[5], 121098312706494698, c)
	}
}

// Neg z = q - x
func Neg(z, x *Element) {
	if x.IsZero() {
		z.SetZero()
		return
	}
	var borrow uint64
	z[0], borrow = bits.Sub64(9586122913090633729, x[0], 0)
	z[1], borrow = bits.Sub64(1660523435060625408, x[1], borrow)
	z[2], borrow = bits.Sub64(2230234197602682880, x[2], borrow)
	z[3], borrow = bits.Sub64(1883307231910630287, x[3], borrow)
	z[4], borrow = bits.Sub64(14284016967150029115, x[4], borrow)
	z[5], _ = bits.Sub64(121098312706494698, x[5], borrow)
}
