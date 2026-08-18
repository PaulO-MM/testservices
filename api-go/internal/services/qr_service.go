// Package services implements the QR decomposition using Modified Gram-Schmidt.
package services

import (
	"errors"
	"math"
)

var (
	// ErrUnsupportedDimensions is returned when m < n.
	ErrUnsupportedDimensions = errors.New("matrices with more columns than rows are not supported")
	// ErrSingularMatrix is returned when columns are linearly dependent.
	ErrSingularMatrix = errors.New("matrix columns are linearly dependent (singular)")
)

const epsilon = 1e-9

// ComputeQR performs QR decomposition on a rectangular matrix (m x n, m >= n)
// using the Modified Gram-Schmidt algorithm for numerical stability.
// Returns Q (m x n orthonormal columns) and R (n x n upper triangular).
func ComputeQR(matrix [][]float64) (q [][]float64, r [][]float64, err error) {
	if len(matrix) == 0 {
		return nil, nil, ErrUnsupportedDimensions
	}

	m := len(matrix)
	n := len(matrix[0])

	if m < n {
		return nil, nil, ErrUnsupportedDimensions
	}

	// Validate rectangular: all rows must have the same length.
	for i := 1; i < m; i++ {
		if len(matrix[i]) != n {
			return nil, nil, errors.New("matrix rows have inconsistent lengths")
		}
	}

	// Initialize Q as copy of input, R as n x n zeros.
	q = make([][]float64, m)
	for i := range q {
		q[i] = make([]float64, n)
		copy(q[i], matrix[i])
	}

	r = make([][]float64, n)
	for i := range r {
		r[i] = make([]float64, n)
	}

	// Modified Gram-Schmidt process.
	for j := 0; j < n; j++ {
		// Compute norm of q[:,j].
		norm := 0.0
		for i := 0; i < m; i++ {
			norm += q[i][j] * q[i][j]
		}
		norm = math.Sqrt(norm)

		if norm < epsilon {
			return nil, nil, ErrSingularMatrix
		}

		r[j][j] = norm

		// Normalize column j of Q.
		for i := 0; i < m; i++ {
			q[i][j] /= norm
		}

		// Subtract projection of remaining columns onto q[:,j].
		for k := j + 1; k < n; k++ {
			dot := 0.0
			for i := 0; i < m; i++ {
				dot += q[i][j] * q[i][k]
			}
			r[j][k] = dot
			for i := 0; i < m; i++ {
				q[i][k] -= dot * q[i][j]
			}
		}
	}

	return q, r, nil
}
