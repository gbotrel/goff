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

// Code generated by goff (v0.3.2) DO NOT EDIT

// Package fp contains field arithmetic operations
package fp

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular,
// there is no security guarantees such as constant time implementation
// or side-channel attack resistance
// /!\ WARNING /!\

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"math/big"
	"math/bits"
	"strconv"
	"sync"
)

// Element represents a field element stored on 4 words (uint64)
// Element are assumed to be in Montgomery form in all methods
// field modulus q =
//
// 21888242871839275222246405745257275088696311157297823662689037894645226208583
type Element [4]uint64

// Limbs number of 64 bits words needed to represent Element
const Limbs = 4

// Bits number bits needed to represent Element
const Bits = 254

// field modulus stored as big.Int
var _modulus big.Int
var onceModulus sync.Once

// Modulus returns q as a big.Int
// q =
//
// 21888242871839275222246405745257275088696311157297823662689037894645226208583
func Modulus() *big.Int {
	onceModulus.Do(func() {
		_modulus.SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)
	})
	return new(big.Int).Set(&_modulus)
}

// q (modulus)
var qElement = Element{
	4332616871279656263,
	10917124144477883021,
	13281191951274694749,
	3486998266802970665,
}

// q'[0], see montgommery multiplication algorithm
var qElementInv0 uint64 = 9786893198990664585

// rSquare
var rSquare = Element{
	17522657719365597833,
	13107472804851548667,
	5164255478447964150,
	493319470278259999,
}

// Bytes returns the regular (non montgomery) value
// of z as a big-endian byte slice.
func (z *Element) Bytes() []byte {
	_z := z.ToRegular()
	var res [Limbs * 8]byte
	binary.BigEndian.PutUint64(res[24:32], _z[0])
	binary.BigEndian.PutUint64(res[16:24], _z[1])
	binary.BigEndian.PutUint64(res[8:16], _z[2])
	binary.BigEndian.PutUint64(res[0:8], _z[3])

	return res[:]
}

// SetBytes interprets e as the bytes of a big-endian unsigned integer,
// sets z to that value (in Montgomery form), and returns z.
func (z *Element) SetBytes(e []byte) *Element {
	var tmp big.Int
	tmp.SetBytes(e)
	z.SetBigInt(&tmp)
	return z
}

// SetUint64 z = v, sets z LSB to v (non-Montgomery form) and convert z to Montgomery form
func (z *Element) SetUint64(v uint64) *Element {
	*z = Element{v}
	return z.Mul(z, &rSquare) // z.ToMont()
}

// Set z = x
func (z *Element) Set(x *Element) *Element {
	z[0] = x[0]
	z[1] = x[1]
	z[2] = x[2]
	z[3] = x[3]
	return z
}

// SetInterface converts i1 from uint64, int, string, or Element, big.Int into Element
// panic if provided type is not supported
func (z *Element) SetInterface(i1 interface{}) *Element {
	switch c1 := i1.(type) {
	case Element:
		return z.Set(&c1)
	case *Element:
		return z.Set(c1)
	case uint64:
		return z.SetUint64(c1)
	case int:
		return z.SetString(strconv.Itoa(c1))
	case string:
		return z.SetString(c1)
	case *big.Int:
		return z.SetBigInt(c1)
	case big.Int:
		return z.SetBigInt(&c1)
	case []byte:
		return z.SetBytes(c1)
	default:
		panic("invalid type")
	}
}

// SetZero z = 0
func (z *Element) SetZero() *Element {
	z[0] = 0
	z[1] = 0
	z[2] = 0
	z[3] = 0
	return z
}

// SetOne z = 1 (in Montgomery form)
func (z *Element) SetOne() *Element {
	z[0] = 15230403791020821917
	z[1] = 754611498739239741
	z[2] = 7381016538464732716
	z[3] = 1011752739694698287
	return z
}

// Div z = x*y^-1 mod q
func (z *Element) Div(x, y *Element) *Element {
	var yInv Element
	yInv.Inverse(y)
	z.Mul(x, &yInv)
	return z
}

// Equal returns z == x
func (z *Element) Equal(x *Element) bool {
	return (z[3] == x[3]) && (z[2] == x[2]) && (z[1] == x[1]) && (z[0] == x[0])
}

// IsZero returns z == 0
func (z *Element) IsZero() bool {
	return (z[3] | z[2] | z[1] | z[0]) == 0
}

// SetRandom sets z to a random element < q
func (z *Element) SetRandom() *Element {
	bytes := make([]byte, 32)
	io.ReadFull(rand.Reader, bytes)
	z[0] = binary.BigEndian.Uint64(bytes[0:8])
	z[1] = binary.BigEndian.Uint64(bytes[8:16])
	z[2] = binary.BigEndian.Uint64(bytes[16:24])
	z[3] = binary.BigEndian.Uint64(bytes[24:32])
	z[3] %= 3486998266802970665

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[3] < 3486998266802970665 || (z[3] == 3486998266802970665 && (z[2] < 13281191951274694749 || (z[2] == 13281191951274694749 && (z[1] < 10917124144477883021 || (z[1] == 10917124144477883021 && (z[0] < 4332616871279656263))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 4332616871279656263, 0)
		z[1], b = bits.Sub64(z[1], 10917124144477883021, b)
		z[2], b = bits.Sub64(z[2], 13281191951274694749, b)
		z[3], _ = bits.Sub64(z[3], 3486998266802970665, b)
	}

	return z
}

// One returns 1 (in montgommery form)
func One() Element {
	var one Element
	one.SetOne()
	return one
}

// MulAssign is deprecated
// Deprecated: use Mul instead
func (z *Element) MulAssign(x *Element) *Element {
	return z.Mul(z, x)
}

// AddAssign is deprecated
// Deprecated: use Add instead
func (z *Element) AddAssign(x *Element) *Element {
	return z.Add(z, x)
}

// SubAssign is deprecated
// Deprecated: use Sub instead
func (z *Element) SubAssign(x *Element) *Element {
	return z.Sub(z, x)
}

// API with assembly impl

// Mul z = x * y mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *Element) Mul(x, y *Element) *Element {
	mul(z, x, y)
	return z
}

// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *Element) Square(x *Element) *Element {
	square(z, x)
	return z
}

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func (z *Element) FromMont() *Element {
	fromMont(z)
	return z
}

// Add z = x + y mod q
func (z *Element) Add(x, y *Element) *Element {
	add(z, x, y)
	return z
}

// Double z = x + x mod q, aka Lsh 1
func (z *Element) Double(x *Element) *Element {
	double(z, x)
	return z
}

// Sub  z = x - y mod q
func (z *Element) Sub(x, y *Element) *Element {
	sub(z, x, y)
	return z
}

// Neg z = q - x
func (z *Element) Neg(x *Element) *Element {
	neg(z, x)
	return z
}

// Exp z = x^exponent mod q
func (z *Element) Exp(x Element, exponent *big.Int) *Element {
	var bZero big.Int
	if exponent.Cmp(&bZero) == 0 {
		return z.SetOne()
	}

	z.Set(&x)

	for i := exponent.BitLen() - 2; i >= 0; i-- {
		z.Square(z)
		if exponent.Bit(i) == 1 {
			z.Mul(z, &x)
		}
	}

	return z
}

// ToMont converts z to Montgomery form
// sets and returns z = z * r^2
func (z *Element) ToMont() *Element {
	return z.Mul(z, &rSquare)
}

// ToRegular returns z in regular form (doesn't mutate z)
func (z Element) ToRegular() Element {
	return *z.FromMont()
}

// String returns the string form of an Element in Montgomery form
func (z *Element) String() string {
	var _z big.Int
	return z.ToBigIntRegular(&_z).String()
}

// ToBigInt returns z as a big.Int in Montgomery form
func (z *Element) ToBigInt(res *big.Int) *big.Int {
	var b [Limbs * 8]byte
	binary.BigEndian.PutUint64(b[24:32], z[0])
	binary.BigEndian.PutUint64(b[16:24], z[1])
	binary.BigEndian.PutUint64(b[8:16], z[2])
	binary.BigEndian.PutUint64(b[0:8], z[3])

	return res.SetBytes(b[:])
}

// ToBigIntRegular returns z as a big.Int in regular form
func (z *Element) ToBigIntRegular(res *big.Int) *big.Int {
	zRegular := z.ToRegular()
	var b [Limbs * 8]byte
	binary.BigEndian.PutUint64(b[24:32], zRegular[0])
	binary.BigEndian.PutUint64(b[16:24], zRegular[1])
	binary.BigEndian.PutUint64(b[8:16], zRegular[2])
	binary.BigEndian.PutUint64(b[0:8], zRegular[3])
	return res.SetBytes(b[:])
}

// SetBigInt sets z to v (regular form) and returns z in Montgomery form
func (z *Element) SetBigInt(v *big.Int) *Element {
	z.SetZero()

	zero := big.NewInt(0)
	q := Modulus()

	// fast path
	c := v.Cmp(q)
	if c == 0 {
		return z
	} else if c != 1 && v.Cmp(zero) != -1 {
		// v should
		vBits := v.Bits()
		for i := 0; i < len(vBits); i++ {
			z[i] = uint64(vBits[i])
		}
		return z.ToMont()
	}

	// copy input
	vv := new(big.Int).Set(v)
	vv.Mod(v, q)

	// v should
	vBits := vv.Bits()
	for i := 0; i < len(vBits); i++ {
		z[i] = uint64(vBits[i])
	}
	return z.ToMont()
}

// SetString creates a big.Int with s (in base 10) and calls SetBigInt on z
func (z *Element) SetString(s string) *Element {
	x, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("Element.SetString failed -> can't parse number in base10 into a big.Int")
	}
	return z.SetBigInt(x)
}

var (
	_bLegendreExponentElement *big.Int
	_bSqrtExponentElement     *big.Int
)

func init() {
	_bLegendreExponentElement, _ = new(big.Int).SetString("183227397098d014dc2822db40c0ac2ecbc0b548b438e5469e10460b6c3e7ea3", 16)
	const sqrtExponentElement = "c19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f52"
	_bSqrtExponentElement, _ = new(big.Int).SetString(sqrtExponentElement, 16)
}

// Legendre returns the Legendre symbol of z (either +1, -1, or 0.)
func (z *Element) Legendre() int {
	var l Element
	// z^((q-1)/2)
	l.Exp(*z, _bLegendreExponentElement)

	if l.IsZero() {
		return 0
	}

	// if l == 1
	if (l[3] == 1011752739694698287) && (l[2] == 7381016538464732716) && (l[1] == 754611498739239741) && (l[0] == 15230403791020821917) {
		return 1
	}
	return -1
}

// Sqrt z = √x mod q
// if the square root doesn't exist (x is not a square mod q)
// Sqrt leaves z unchanged and returns nil
func (z *Element) Sqrt(x *Element) *Element {
	// q ≡ 3 (mod 4)
	// using  z ≡ ± x^((p+1)/4) (mod q)
	var y, square Element
	y.Exp(*x, _bSqrtExponentElement)
	// as we didn't compute the legendre symbol, ensure we found y such that y * y = x
	square.Square(&y)
	if square.Equal(x) {
		return z.Set(&y)
	}
	return nil
}

// Inverse z = x^-1 mod q
// Algorithm 16 in "Efficient Software-Implementation of Finite Fields with Applications to Cryptography"
// if x == 0, sets and returns z = x
func (z *Element) Inverse(x *Element) *Element {
	if x.IsZero() {
		return z.Set(x)
	}

	// initialize u = q
	var u = Element{
		4332616871279656263,
		10917124144477883021,
		13281191951274694749,
		3486998266802970665,
	}

	// initialize s = r^2
	var s = Element{
		17522657719365597833,
		13107472804851548667,
		5164255478447964150,
		493319470278259999,
	}

	// r = 0
	r := Element{}

	v := *x

	var carry, borrow, t, t2 uint64
	var bigger, uIsOne, vIsOne bool

	for !uIsOne && !vIsOne {
		for v[0]&1 == 0 {

			// v = v >> 1
			t2 = v[3] << 63
			v[3] >>= 1
			t = t2
			t2 = v[2] << 63
			v[2] = (v[2] >> 1) | t
			t = t2
			t2 = v[1] << 63
			v[1] = (v[1] >> 1) | t
			t = t2
			v[0] = (v[0] >> 1) | t

			if s[0]&1 == 1 {

				// s = s + q
				s[0], carry = bits.Add64(s[0], 4332616871279656263, 0)
				s[1], carry = bits.Add64(s[1], 10917124144477883021, carry)
				s[2], carry = bits.Add64(s[2], 13281191951274694749, carry)
				s[3], _ = bits.Add64(s[3], 3486998266802970665, carry)

			}

			// s = s >> 1
			t2 = s[3] << 63
			s[3] >>= 1
			t = t2
			t2 = s[2] << 63
			s[2] = (s[2] >> 1) | t
			t = t2
			t2 = s[1] << 63
			s[1] = (s[1] >> 1) | t
			t = t2
			s[0] = (s[0] >> 1) | t

		}
		for u[0]&1 == 0 {

			// u = u >> 1
			t2 = u[3] << 63
			u[3] >>= 1
			t = t2
			t2 = u[2] << 63
			u[2] = (u[2] >> 1) | t
			t = t2
			t2 = u[1] << 63
			u[1] = (u[1] >> 1) | t
			t = t2
			u[0] = (u[0] >> 1) | t

			if r[0]&1 == 1 {

				// r = r + q
				r[0], carry = bits.Add64(r[0], 4332616871279656263, 0)
				r[1], carry = bits.Add64(r[1], 10917124144477883021, carry)
				r[2], carry = bits.Add64(r[2], 13281191951274694749, carry)
				r[3], _ = bits.Add64(r[3], 3486998266802970665, carry)

			}

			// r = r >> 1
			t2 = r[3] << 63
			r[3] >>= 1
			t = t2
			t2 = r[2] << 63
			r[2] = (r[2] >> 1) | t
			t = t2
			t2 = r[1] << 63
			r[1] = (r[1] >> 1) | t
			t = t2
			r[0] = (r[0] >> 1) | t

		}

		// v >= u
		bigger = !(v[3] < u[3] || (v[3] == u[3] && (v[2] < u[2] || (v[2] == u[2] && (v[1] < u[1] || (v[1] == u[1] && (v[0] < u[0])))))))

		if bigger {

			// v = v - u
			v[0], borrow = bits.Sub64(v[0], u[0], 0)
			v[1], borrow = bits.Sub64(v[1], u[1], borrow)
			v[2], borrow = bits.Sub64(v[2], u[2], borrow)
			v[3], _ = bits.Sub64(v[3], u[3], borrow)

			// r >= s
			bigger = !(r[3] < s[3] || (r[3] == s[3] && (r[2] < s[2] || (r[2] == s[2] && (r[1] < s[1] || (r[1] == s[1] && (r[0] < s[0])))))))

			if bigger {

				// s = s + q
				s[0], carry = bits.Add64(s[0], 4332616871279656263, 0)
				s[1], carry = bits.Add64(s[1], 10917124144477883021, carry)
				s[2], carry = bits.Add64(s[2], 13281191951274694749, carry)
				s[3], _ = bits.Add64(s[3], 3486998266802970665, carry)

			}

			// s = s - r
			s[0], borrow = bits.Sub64(s[0], r[0], 0)
			s[1], borrow = bits.Sub64(s[1], r[1], borrow)
			s[2], borrow = bits.Sub64(s[2], r[2], borrow)
			s[3], _ = bits.Sub64(s[3], r[3], borrow)

		} else {

			// u = u - v
			u[0], borrow = bits.Sub64(u[0], v[0], 0)
			u[1], borrow = bits.Sub64(u[1], v[1], borrow)
			u[2], borrow = bits.Sub64(u[2], v[2], borrow)
			u[3], _ = bits.Sub64(u[3], v[3], borrow)

			// s >= r
			bigger = !(s[3] < r[3] || (s[3] == r[3] && (s[2] < r[2] || (s[2] == r[2] && (s[1] < r[1] || (s[1] == r[1] && (s[0] < r[0])))))))

			if bigger {

				// r = r + q
				r[0], carry = bits.Add64(r[0], 4332616871279656263, 0)
				r[1], carry = bits.Add64(r[1], 10917124144477883021, carry)
				r[2], carry = bits.Add64(r[2], 13281191951274694749, carry)
				r[3], _ = bits.Add64(r[3], 3486998266802970665, carry)

			}

			// r = r - s
			r[0], borrow = bits.Sub64(r[0], s[0], 0)
			r[1], borrow = bits.Sub64(r[1], s[1], borrow)
			r[2], borrow = bits.Sub64(r[2], s[2], borrow)
			r[3], _ = bits.Sub64(r[3], s[3], borrow)

		}
		uIsOne = (u[0] == 1) && (u[3]|u[2]|u[1]) == 0
		vIsOne = (v[0] == 1) && (v[3]|v[2]|v[1]) == 0
	}

	if uIsOne {
		z.Set(&r)
	} else {
		z.Set(&s)
	}

	return z
}
