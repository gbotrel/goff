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

// Code generated by goff (v0.2.2) DO NOT EDIT

// Package bn256 contains field arithmetic operations
package bn256

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular,
// there is no security guarantees such as constant time implementation
// or side-channel attack resistance
// /!\ WARNING /!\

//go:noescape
func mulAssignElement(res, y *Element)

//go:noescape
func fromMontElement(res *Element)

//go:noescape
func reduceElement(res *Element) // for test purposes

// Mul z = x * y mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *Element) Mul(x, y *Element) *Element {
	if z == x {
		mulAssignElement(z, y)
		return z
	} else if z == y {
		mulAssignElement(z, x)
		return z
	} else {
		z.Set(x)
		mulAssignElement(z, y)
		return z
	}
}

// MulAssign z = z * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *Element) MulAssign(x *Element) *Element {
	mulAssignElement(z, x)
	return z
}
