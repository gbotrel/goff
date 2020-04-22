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

import (
	"crypto/rand"
	"math/big"
	"math/bits"
	mrand "math/rand"
	"testing"
)

func TestELEMENTCorrectnessAgainstBigInt(t *testing.T) {
	modulus, _ := new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)
	cmpEandB := func(e *Element, b *big.Int, name string) {
		var _e big.Int
		if e.FromMont().ToBigInt(&_e).Cmp(b) != 0 {
			t.Fatal(name, "failed")
		}
	}
	var modulusMinusOne, one big.Int
	one.SetUint64(1)

	modulusMinusOne.Sub(modulus, &one)

	var n int
	if testing.Short() {
		n = 20
	} else {
		n = 500
	}

	sAdx := supportAdx

	for i := 0; i < n; i++ {
		if i == n/2 && sAdx {
			supportAdx = false // testing without adx instruction
		}
		// sample 2 random big int
		b1, _ := rand.Int(rand.Reader, modulus)
		b2, _ := rand.Int(rand.Reader, modulus)
		rExp := mrand.Uint64()

		// adding edge cases
		// TODO need more edge cases
		switch i {
		case 0:
			rExp = 0
			b1.SetUint64(0)
		case 1:
			b2.SetUint64(0)
		case 2:
			b1.SetUint64(0)
			b2.SetUint64(0)
		case 3:
			rExp = 0
		case 4:
			rExp = 1
		case 5:
			rExp = ^uint64(0) // max uint
		case 6:
			rExp = 2
			b1.Set(&modulusMinusOne)
		case 7:
			b2.Set(&modulusMinusOne)
		case 8:
			b1.Set(&modulusMinusOne)
			b2.Set(&modulusMinusOne)
		}

		rbExp := new(big.Int).SetUint64(rExp)

		var bMul, bAdd, bSub, bDiv, bNeg, bLsh, bInv, bExp, bExp2, bSquare big.Int

		// e1 = mont(b1), e2 = mont(b2)
		var e1, e2, eMul, eAdd, eSub, eDiv, eNeg, eLsh, eInv, eExp, eSquare, eMulAssign, eSubAssign, eAddAssign Element
		e1.SetBigInt(b1)
		e2.SetBigInt(b2)

		// (e1*e2).FromMont() === b1*b2 mod q ... etc
		eSquare.Square(&e1)
		eMul.Mul(&e1, &e2)
		eMulAssign.Set(&e1)
		eMulAssign.MulAssign(&e2)
		eAdd.Add(&e1, &e2)
		eAddAssign.Set(&e1)
		eAddAssign.AddAssign(&e2)
		eSub.Sub(&e1, &e2)
		eSubAssign.Set(&e1)
		eSubAssign.SubAssign(&e2)
		eDiv.Div(&e1, &e2)
		eNeg.Neg(&e1)
		eInv.Inverse(&e1)
		eExp.Exp(e1, rExp)
		eLsh.Double(&e1)

		// same operations with big int
		bAdd.Add(b1, b2).Mod(&bAdd, modulus)
		bMul.Mul(b1, b2).Mod(&bMul, modulus)
		bSquare.Mul(b1, b1).Mod(&bSquare, modulus)
		bSub.Sub(b1, b2).Mod(&bSub, modulus)
		bDiv.ModInverse(b2, modulus)
		bDiv.Mul(&bDiv, b1).
			Mod(&bDiv, modulus)
		bNeg.Neg(b1).Mod(&bNeg, modulus)

		bInv.ModInverse(b1, modulus)
		bExp.Exp(b1, rbExp, modulus)
		bLsh.Lsh(b1, 1).Mod(&bLsh, modulus)

		cmpEandB(&eSquare, &bSquare, "Square")
		cmpEandB(&eMul, &bMul, "Mul")
		cmpEandB(&eMulAssign, &bMul, "MulAssign")
		cmpEandB(&eAdd, &bAdd, "Add")
		cmpEandB(&eAddAssign, &bAdd, "AddAssign")
		cmpEandB(&eSub, &bSub, "Sub")
		cmpEandB(&eSubAssign, &bSub, "SubAssign")
		cmpEandB(&eDiv, &bDiv, "Div")
		cmpEandB(&eNeg, &bNeg, "Neg")
		cmpEandB(&eInv, &bInv, "Inv")
		cmpEandB(&eExp, &bExp, "Exp")

		cmpEandB(&eLsh, &bLsh, "Lsh")

		// legendre symbol
		if e1.Legendre() != big.Jacobi(b1, modulus) {
			t.Fatal("legendre symbol computation failed")
		}
		if e2.Legendre() != big.Jacobi(b2, modulus) {
			t.Fatal("legendre symbol computation failed")
		}

		// these are slow, killing circle ci
		if n <= 5 {
			// sqrt
			var eSqrt, eExp2 Element
			var bSqrt big.Int
			bSqrt.ModSqrt(b1, modulus)
			eSqrt.Sqrt(&e1)
			cmpEandB(&eSqrt, &bSqrt, "Sqrt")

			bits := b2.Bits()
			exponent := make([]uint64, len(bits))
			for k := 0; k < len(bits); k++ {
				exponent[k] = uint64(bits[k])
			}
			eExp2.Exp(e1, exponent...)
			bExp2.Exp(b1, b2, modulus)
			cmpEandB(&eExp2, &bExp2, "Exp multi words")
		}
	}
	supportAdx = sAdx
}

func TestELEMENTIsRandom(t *testing.T) {
	for i := 0; i < 50; i++ {
		var x, y Element
		x.SetRandom()
		y.SetRandom()
		if x.Equal(&y) {
			t.Fatal("2 random numbers are unlikely to be equal")
		}
	}
}

func TestByte(t *testing.T) {
	modulus := ElementModulus()
	sample, _ := rand.Int(rand.Reader, modulus)
	var witness Element

	witness.SetBigInt(sample)

	b := witness.ToBytes()

	// check consistency conversion
	var test Element
	test.SetBytes(b)
	if !test.Equal(&witness) {
		t.Fatal("Inconsistancy during conversion ToBytes/SetBytes")
	}

}

// -------------------------------------------------------------------------------------------------
// benchmarks
// most benchmarks are rudimentary and should sample a large number of random inputs
// or be run multiple times to ensure it didn't measure the fastest path of the function

var benchResElement Element

func BenchmarkInverseELEMENT(b *testing.B) {
	var x Element
	x.SetRandom()
	benchResElement.SetRandom()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchResElement.Inverse(&x)
	}

}
func BenchmarkExpELEMENT(b *testing.B) {
	var x Element
	x.SetRandom()
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Exp(x, mrand.Uint64())
	}
}

func BenchmarkDoubleELEMENT(b *testing.B) {
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Double(&benchResElement)
	}
}

func BenchmarkAddELEMENT(b *testing.B) {
	var x Element
	x.SetRandom()
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Add(&x, &benchResElement)
	}
}

func BenchmarkSubELEMENT(b *testing.B) {
	var x Element
	x.SetRandom()
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Sub(&x, &benchResElement)
	}
}

func BenchmarkNegELEMENT(b *testing.B) {
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Neg(&benchResElement)
	}
}

func BenchmarkDivELEMENT(b *testing.B) {
	var x Element
	x.SetRandom()
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Div(&x, &benchResElement)
	}
}

func BenchmarkFromMontELEMENT(b *testing.B) {
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.FromMont()
	}
}

func BenchmarkToMontELEMENT(b *testing.B) {
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.ToMont()
	}
}
func BenchmarkSquareELEMENT(b *testing.B) {
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Square(&benchResElement)
	}
}

func BenchmarkSqrtELEMENT(b *testing.B) {
	var a Element
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Sqrt(&a)
	}
}

func BenchmarkMulAssignELEMENT(b *testing.B) {
	x := Element{
		17522657719365597833,
		13107472804851548667,
		5164255478447964150,
		493319470278259999,
	}
	benchResElement.SetOne()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.MulAssign(&x)
	}
}

func TestELEMENTAsm(t *testing.T) {
	// ensure ASM implementations matches the ones using math/bits
	modulus, _ := new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)
	sadx := supportAdx
	for i := 0; i < 500; i++ {
		// sample 2 random big int
		if i == 250 && sadx {
			// going the no_adx path
			supportAdx = false
		}
		b1, _ := rand.Int(rand.Reader, modulus)
		b2, _ := rand.Int(rand.Reader, modulus)

		// e1 = mont(b1), e2 = mont(b2)
		var e1, e2, eTestMul, eMulAssign, eSquare, eTestSquare Element
		e1.SetBigInt(b1)
		e2.SetBigInt(b2)

		eTestMul = e1
		eTestMul.testMulAssign(&e2)
		eMulAssign = e1
		eMulAssign.MulAssign(&e2)

		if !eTestMul.Equal(&eMulAssign) {
			if supportAdx {
				t.Fatal("mul assembly implementation WITH adx instructions doesn't match non-assembly one")
			} else {
				t.Fatal("mul assembly implementation WITHOUT adx instructions doesn't match non-assembly one")
			}
		}

		// square
		eSquare.Square(&e1)
		eTestSquare.testSquare(&e1)

		if !eTestSquare.Equal(&eSquare) {
			if supportAdx {
				t.Fatal("square assembly implementation WITH adx instructions doesn't match non-assembly one")
			} else {
				t.Fatal("square assembly implementation WITHOUT adx instructions doesn't match non-assembly one")
			}
		}
	}
	supportAdx = sadx
}

func TestELEMENTreduce(t *testing.T) {
	q := Element{
		4332616871279656263,
		10917124144477883021,
		13281191951274694749,
		3486998266802970665,
	}

	var testData []Element
	{
		a := q
		a[3] -= 1
		testData = append(testData, a)
	}
	{
		a := q
		a[0] -= 1
		testData = append(testData, a)
	}
	{
		a := q
		a[3] += 1
		testData = append(testData, a)
	}
	{
		a := q
		a[0] += 1
		testData = append(testData, a)
	}
	{
		a := q
		testData = append(testData, a)
	}

	for _, s := range testData {
		expected := s
		reduceElement(&s)
		expected.testReduce()
		if !s.Equal(&expected) {
			t.Fatal("reduce failed")
		}
	}

}

// this is here for consistency purposes, to ensure MulAssign on AMD64 using asm implementation gives consistent results
func (z *Element) testMulAssign(x *Element) *Element {

	var t [4]uint64
	var c [3]uint64
	{
		// round 0
		v := z[0]
		c[1], c[0] = bits.Mul64(v, x[0])
		m := c[0] * 9786893198990664585
		c[2] = madd0(m, 4332616871279656263, c[0])
		c[1], c[0] = madd1(v, x[1], c[1])
		c[2], t[0] = madd2(m, 10917124144477883021, c[2], c[0])
		c[1], c[0] = madd1(v, x[2], c[1])
		c[2], t[1] = madd2(m, 13281191951274694749, c[2], c[0])
		c[1], c[0] = madd1(v, x[3], c[1])
		t[3], t[2] = madd3(m, 3486998266802970665, c[0], c[2], c[1])
	}
	{
		// round 1
		v := z[1]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * 9786893198990664585
		c[2] = madd0(m, 4332616871279656263, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, 10917124144477883021, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, 13281191951274694749, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		t[3], t[2] = madd3(m, 3486998266802970665, c[0], c[2], c[1])
	}
	{
		// round 2
		v := z[2]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * 9786893198990664585
		c[2] = madd0(m, 4332616871279656263, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, 10917124144477883021, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, 13281191951274694749, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		t[3], t[2] = madd3(m, 3486998266802970665, c[0], c[2], c[1])
	}
	{
		// round 3
		v := z[3]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * 9786893198990664585
		c[2] = madd0(m, 4332616871279656263, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], z[0] = madd2(m, 10917124144477883021, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], z[1] = madd2(m, 13281191951274694749, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		z[3], z[2] = madd3(m, 3486998266802970665, c[0], c[2], c[1])
	}

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

func (z *Element) testReduce() *Element {

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

// this is here for consistency purposes, to ensure Square on AMD64 using asm implementation gives consistent results
func (z *Element) testSquare(x *Element) *Element {

	var p [4]uint64

	var u, v uint64
	{
		// round 0
		u, p[0] = bits.Mul64(x[0], x[0])
		m := p[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, p[0])
		var t uint64
		t, u, v = madd1sb(x[0], x[1], u)
		C, p[0] = madd2(m, 10917124144477883021, v, C)
		t, u, v = madd1s(x[0], x[2], t, u)
		C, p[1] = madd2(m, 13281191951274694749, v, C)
		_, u, v = madd1s(x[0], x[3], t, u)
		p[3], p[2] = madd3(m, 3486998266802970665, v, C, u)
	}
	{
		// round 1
		m := p[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, p[0])
		u, v = madd1(x[1], x[1], p[1])
		C, p[0] = madd2(m, 10917124144477883021, v, C)
		var t uint64
		t, u, v = madd2sb(x[1], x[2], p[2], u)
		C, p[1] = madd2(m, 13281191951274694749, v, C)
		_, u, v = madd2s(x[1], x[3], p[3], t, u)
		p[3], p[2] = madd3(m, 3486998266802970665, v, C, u)
	}
	{
		// round 2
		m := p[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, p[0])
		C, p[0] = madd2(m, 10917124144477883021, p[1], C)
		u, v = madd1(x[2], x[2], p[2])
		C, p[1] = madd2(m, 13281191951274694749, v, C)
		_, u, v = madd2sb(x[2], x[3], p[3], u)
		p[3], p[2] = madd3(m, 3486998266802970665, v, C, u)
	}
	{
		// round 3
		m := p[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, p[0])
		C, z[0] = madd2(m, 10917124144477883021, p[1], C)
		C, z[1] = madd2(m, 13281191951274694749, p[2], C)
		u, v = madd1(x[3], x[3], p[3])
		z[3], z[2] = madd3(m, 3486998266802970665, v, C, u)
	}

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
