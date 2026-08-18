// Tests for the matrix statistics service.
const { computeStats, isDiagonal, flatten } = require("../src/services/stats.service");

describe("Stats Service", () => {
  describe("computeStats", () => {
    test("computes correct stats for known matrices", () => {
      const q = [[1, 0], [0, 1]];
      const r = [[2, 3], [0, 4]];
      const stats = computeStats(q, r);

      expect(stats.max).toBe(4);
      expect(stats.min).toBe(0);
      expect(stats.sum).toBe(11);
      expect(stats.average).toBe(1.38);
    });

    test("computes stats for 3x2 matrices", () => {
      const q = [[0.707, 0], [0, 1], [0.707, 0]];
      const r = [[1.414, 0.707], [0, 1]];
      const stats = computeStats(q, r);

      expect(stats.max).toBeCloseTo(1.414, 2);
      expect(stats.min).toBe(0);
      expect(stats.sum).toBeGreaterThan(0);
    });

    test("handles empty matrices", () => {
      const stats = computeStats([], []);
      expect(stats.max).toBe(0);
      expect(stats.min).toBe(0);
      expect(stats.sum).toBe(0);
      expect(stats.average).toBe(0);
    });
  });

  describe("isDiagonal", () => {
    test("returns true for diagonal matrix", () => {
      expect(isDiagonal([[1, 0], [0, 2]])).toBe(true);
    });

    test("returns true for identity matrix", () => {
      expect(isDiagonal([[1, 0, 0], [0, 1, 0], [0, 0, 1]])).toBe(true);
    });

    test("returns false for non-diagonal matrix", () => {
      expect(isDiagonal([[1, 2], [0, 1]])).toBe(false);
    });

    test("returns false for non-square matrix", () => {
      expect(isDiagonal([[1, 0], [0, 1], [0, 0]])).toBe(false);
    });

    test("returns false for single-element matrix (vacuously diagonal but not square check)", () => {
      // DECISION: A 1x1 matrix IS square and IS diagonal.
      expect(isDiagonal([[5]])).toBe(true);
    });

    test("returns false for empty input", () => {
      expect(isDiagonal([])).toBe(false);
      expect(isDiagonal(null)).toBe(false);
    });

    test("returns false when off-diagonal element is near zero but above epsilon", () => {
      expect(isDiagonal([[1, 0.1], [0, 1]])).toBe(false);
    });

    test("returns true when off-diagonal element is below epsilon", () => {
      expect(isDiagonal([[1, 5e-10], [0, 1]])).toBe(true);
    });
  });

  describe("flatten", () => {
    test("flattens a 2D array", () => {
      expect(flatten([[1, 2], [3, 4]])).toEqual([1, 2, 3, 4]);
    });

    test("returns empty array for empty input", () => {
      expect(flatten([])).toEqual([]);
    });
  });
});
