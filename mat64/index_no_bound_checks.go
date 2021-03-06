// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file must be kept in sync with index_bound_checks.go.

//+build !bounds

package mat64

import "github.com/gonum/blas"

// At returns the element at row r, column c.
func (m *Dense) At(r, c int) float64 {
	if r >= m.mat.Rows || r < 0 {
		panic(ErrRowAccess)
	}
	if c >= m.mat.Cols || c < 0 {
		panic(ErrColAccess)
	}
	return m.at(r, c)
}

func (m *Dense) at(r, c int) float64 {
	return m.mat.Data[r*m.mat.Stride+c]
}

// Set sets the element at row r, column c to the value v.
func (m *Dense) Set(r, c int, v float64) {
	if r >= m.mat.Rows || r < 0 {
		panic(ErrRowAccess)
	}
	if c >= m.mat.Cols || c < 0 {
		panic(ErrColAccess)
	}
	m.set(r, c, v)
}

func (m *Dense) set(r, c int, v float64) {
	m.mat.Data[r*m.mat.Stride+c] = v
}

// At returns the element at row r, column c. It panics if c is not zero.
func (m *Vector) At(r, c int) float64 {
	if r < 0 || r >= m.n {
		panic(ErrRowAccess)
	}
	if c != 0 {
		panic(ErrColAccess)
	}
	return m.at(r)
}

func (m *Vector) at(r int) float64 {
	return m.mat.Data[r*m.mat.Inc]
}

// Set sets the element at (r,c) to the value v. It panics if c is not zero.
func (m *Vector) Set(r, c int, v float64) {
	if r < 0 || r >= m.n {
		panic(ErrRowAccess)
	}
	if c != 0 {
		panic(ErrColAccess)
	}
	m.set(r, v)
}

func (m *Vector) set(r int, v float64) {
	m.mat.Data[r*m.mat.Inc] = v
}

// At returns the element at row r and column c.
func (s *SymDense) At(r, c int) float64 {
	if r >= s.mat.N || r < 0 {
		panic(ErrRowAccess)
	}
	if c >= s.mat.N || c < 0 {
		panic(ErrColAccess)
	}
	return s.at(r, c)
}

func (s *SymDense) at(r, c int) float64 {
	if r > c {
		r, c = c, r
	}
	return s.mat.Data[r*s.mat.Stride+c]
}

// SetSym sets the elements at (r,c) and (c,r) to the value v.
func (s *SymDense) SetSym(r, c int, v float64) {
	if r >= s.mat.N || r < 0 {
		panic(ErrRowAccess)
	}
	if c >= s.mat.N || c < 0 {
		panic(ErrColAccess)
	}
	s.set(r, c, v)
}

func (s *SymDense) set(r, c int, v float64) {
	if r > c {
		r, c = c, r
	}
	s.mat.Data[r*s.mat.Stride+c] = v
}

// At returns the element at row r, column c.
func (t *TriDense) At(r, c int) float64 {
	if r >= t.mat.N || r < 0 {
		panic(ErrRowAccess)
	}
	if c >= t.mat.N || c < 0 {
		panic(ErrColAccess)
	}
	return t.at(r, c)
}

func (t *TriDense) at(r, c int) float64 {
	if t.mat.Uplo == blas.Upper {
		if r > c {
			return 0
		}
		return t.mat.Data[r*t.mat.Stride+c]
	}
	if r < c {
		return 0
	}
	return t.mat.Data[r*t.mat.Stride+c]
}

// SetTri sets the element at row r, column c to the value v.
// It panics if the location is outside the appropriate half of the matrix.
func (t *TriDense) SetTri(r, c int, v float64) {
	if r >= t.mat.N || r < 0 {
		panic(ErrRowAccess)
	}
	if c >= t.mat.N || c < 0 {
		panic(ErrColAccess)
	}
	if t.mat.Uplo == blas.Upper && r > c {
		panic("mat64: triangular set out of bounds")
	}
	if t.mat.Uplo == blas.Lower && r < c {
		panic("mat64: triangular set out of bounds")
	}
	t.set(r, c, v)
}

func (t *TriDense) set(r, c int, v float64) {
	t.mat.Data[r*t.mat.Stride+c] = v
}
