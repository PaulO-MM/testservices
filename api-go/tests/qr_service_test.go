// Tests for the QR decomposition service using Modified Gram-Schmidt.
package tests

import (
	"math"
	"testing"

	"challenge-api-go/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testEpsilon = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < testEpsilon
}

func matAlmostEqual(a, b [][]float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if !almostEqual(a[i][j], b[i][j]) {
				return false
			}
		}
	}
	return true
}

func transpose(m [][]float64) [][]float64 {
	if len(m) == 0 {
		return nil
	}
	rows := len(m)
	cols := len(m[0])
	t := make([][]float64, cols)
	for j := 0; j < cols; j++ {
		t[j] = make([]float64, rows)
		for i := 0; i < rows; i++ {
			t[j][i] = m[i][j]
		}
	}
	return t
}

func matMul(a, b [][]float64) [][]float64 {
	m := len(a)
	n := len(b[0])
	k := len(b)
	result := make([][]float64, m)
	for i := 0; i < m; i++ {
		result[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			sum := 0.0
			for l := 0; l < k; l++ {
				sum += a[i][l] * b[l][j]
			}
			result[i][j] = sum
		}
	}
	return result
}

func isUpperTriangular(r [][]float64) bool {
	n := len(r)
	for i := 1; i < n; i++ {
		for j := 0; j < i && j < len(r[i]); j++ {
			if !almostEqual(r[i][j], 0) {
				return false
			}
		}
	}
	return true
}

func TestComputeQR_Identity2x2(t *testing.T) {
	input := [][]float64{
		{1, 0},
		{0, 1},
	}

	q, r, err := services.ComputeQR(input)
	require.NoError(t, err)
	assert.Equal(t, 2, len(q))
	assert.Equal(t, 2, len(q[0]))
	assert.Equal(t, 2, len(r))

	expectedQ := [][]float64{{1, 0}, {0, 1}}
	assert.True(t, matAlmostEqual(q, expectedQ), "Q should be identity")

	expectedR := [][]float64{{1, 0}, {0, 1}}
	assert.True(t, matAlmostEqual(r, expectedR), "R should be identity")

	qt := transpose(q)
	qtq := matMul(qt, q)
	expectedI := [][]float64{{1, 0}, {0, 1}}
	assert.True(t, matAlmostEqual(qtq, expectedI), "Q^T * Q should be identity")

	qr := matMul(q, r)
	assert.True(t, matAlmostEqual(qr, input), "Q * R should equal original matrix")
}

func TestComputeQR_Rectangular3x2(t *testing.T) {
	input := [][]float64{
		{1, 2},
		{0, 1},
		{1, 0},
	}

	q, r, err := services.ComputeQR(input)
	require.NoError(t, err)
	assert.Equal(t, 3, len(q), "Q should have 3 rows")
	assert.Equal(t, 2, len(q[0]), "Q should have 2 columns")
	assert.Equal(t, 2, len(r), "R should have 2 rows")
	assert.Equal(t, 2, len(r[0]), "R should have 2 columns")

	qt := transpose(q)
	qtq := matMul(qt, q)
	expectedI := [][]float64{{1, 0}, {0, 1}}
	assert.True(t, matAlmostEqual(qtq, expectedI), "Q^T * Q should be identity")

	assert.True(t, isUpperTriangular(r), "R should be upper triangular")

	qr := matMul(q, r)
	assert.True(t, matAlmostEqual(qr, input), "Q * R should equal original matrix")
}

func TestComputeQR_SingularMatrix(t *testing.T) {
	input := [][]float64{
		{1, 1},
		{1, 1},
	}

	_, _, err := services.ComputeQR(input)
	assert.ErrorIs(t, err, services.ErrSingularMatrix)
}

func TestComputeQR_UnsupportedDimensions(t *testing.T) {
	input := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}

	_, _, err := services.ComputeQR(input)
	assert.ErrorIs(t, err, services.ErrUnsupportedDimensions)
}
