// Computes statistical properties (max, min, average, sum, diagonal) for Q and R matrices.

const EPSILON = 1e-9;

/**
 * Checks if a square matrix is diagonal (all off-diagonal elements ≈ 0).
 * Returns false automatically if the matrix is not square.
 */
function isDiagonal(matrix) {
  if (!Array.isArray(matrix) || matrix.length === 0) return false;
  const n = matrix.length;
  for (const row of matrix) {
    if (!Array.isArray(row) || row.length !== n) return false;
  }
  for (let i = 0; i < n; i++) {
    for (let j = 0; j < n; j++) {
      if (i !== j && Math.abs(matrix[i][j]) > EPSILON) {
        return false;
      }
    }
  }
  return true;
}

/**
 * Flattens a matrix into a 1D array of values.
 */
function flatten(matrix) {
  const values = [];
  for (const row of matrix) {
    for (const val of row) {
      values.push(val);
    }
  }
  return values;
}

/**
 * Computes stats for the combined Q and R matrices.
 * DECISION: Q and R values are combined into a single set for max/min/avg/sum.
 */
function computeStats(q, r) {
  const qFlat = flatten(q);
  const rFlat = flatten(r);
  const all = [...qFlat, ...rFlat];

  if (all.length === 0) {
    return { max: 0, min: 0, average: 0, sum: 0, diagonal: { q: false, r: false } };
  }

  let max = all[0];
  let min = all[0];
  let sum = 0;

  for (const val of all) {
    if (val > max) max = val;
    if (val < min) min = val;
    sum += val;
  }

  const average = Math.round((sum / all.length) * 100) / 100;

  return {
    max,
    min,
    average,
    sum,
    diagonal: {
      q: isDiagonal(q),
      r: isDiagonal(r),
    },
  };
}

module.exports = { computeStats, isDiagonal, flatten };
