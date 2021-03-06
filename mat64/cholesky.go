// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Based on the CholeskyDecomposition class from Jama 1.0.3.

package mat64

import (
	"math"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/internal/asm"
)

const badTriangle = "mat64: invalid triangle"

// Cholesky calculates the Cholesky decomposition of the matrix A and returns
// whether the matrix is positive definite. The returned matrix is either a
// lower triangular matrix such that A = L * L^T or an upper triangular matrix
// such that A = U^T * U depending on the upper parameter.
func (t *TriDense) Cholesky(a *SymDense, upper bool) (ok bool) {
	n := a.Symmetric()
	if t.isZero() {
		t.mat = blas64.Triangular{
			N:      n,
			Stride: n,
			Diag:   blas.NonUnit,
			Data:   use(t.mat.Data, n*n),
		}
	} else if n != t.mat.N {
		panic(ErrShape)
	}
	mat := t.mat.Data
	stride := t.mat.Stride

	if upper {
		t.mat.Uplo = blas.Upper
		for j := 0; j < n; j++ {
			var d float64
			for k := 0; k < j; k++ {
				s := asm.DdotInc(
					mat, mat,
					uintptr(k),
					uintptr(stride), uintptr(stride),
					uintptr(k), uintptr(j),
				)
				s = (a.at(j, k) - s) / t.at(k, k)
				t.set(k, j, s)
				d += s * s
			}
			d = a.at(j, j) - d
			if d <= 0 {
				t.Reset()
				return false
			}
			t.set(j, j, math.Sqrt(math.Max(d, 0)))
		}
	} else {
		t.mat.Uplo = blas.Lower
		for j := 0; j < n; j++ {
			var d float64
			for k := 0; k < j; k++ {
				s := asm.DdotUnitary(mat[k*stride:k*stride+k], mat[j*stride:j*stride+k])
				s = (a.at(j, k) - s) / t.at(k, k)
				t.set(j, k, s)
				d += s * s
			}
			d = a.at(j, j) - d
			if d <= 0 {
				t.Reset()
				return false
			}
			t.set(j, j, math.Sqrt(math.Max(d, 0)))
		}
	}

	return true
}

// SolveCholesky finds the matrix x that solves A * X = B where A = L * L^T or
// A = U^T * U, and U or L are represented by t. The matrix A must be symmetric
// and positive definite.
func (m *Dense) SolveCholesky(t Triangular, b Matrix) {
	_, n := t.Dims()
	bm, bn := b.Dims()
	if n != bm {
		panic(ErrShape)
	}

	m.reuseAs(bm, bn)
	if b != m {
		m.Copy(b)
	}

	// TODO(btracey): Implement an algorithm that doesn't require a copy into
	// a blas64.Triangular.
	ta := getBlasTriangular(t)

	switch ta.Uplo {
	case blas.Upper:
		blas64.Trsm(blas.Left, blas.Trans, 1, ta, m.mat)
		blas64.Trsm(blas.Left, blas.NoTrans, 1, ta, m.mat)
	case blas.Lower:
		blas64.Trsm(blas.Left, blas.NoTrans, 1, ta, m.mat)
		blas64.Trsm(blas.Left, blas.Trans, 1, ta, m.mat)
	default:
		panic(badTriangle)
	}
}

// SolveTri finds the matrix x that solves op(A) * X = B where A is a triangular
// matrix and op is specified by trans.
func (m *Dense) SolveTri(a Triangular, trans bool, b Matrix) {
	n, _ := a.Triangle()
	bm, bn := b.Dims()
	if n != bm {
		panic(ErrShape)
	}

	m.reuseAs(bm, bn)
	if b != m {
		m.Copy(b)
	}

	// TODO(btracey): Implement an algorithm that doesn't require a copy into
	// a blas64.Triangular.
	ta := getBlasTriangular(a)

	t := blas.NoTrans
	if trans {
		t = blas.Trans
	}
	switch ta.Uplo {
	case blas.Upper, blas.Lower:
		blas64.Trsm(blas.Left, t, 1, ta, m.mat)
	default:
		panic(badTriangle)
	}
}
