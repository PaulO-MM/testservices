// Integration tests for matrix routes using Supertest.
const request = require("supertest");
const jwt = require("jsonwebtoken");
const app = require("../src/app");

const JWT_SECRET = process.env.JWT_SECRET || "dev_secret_change_me";

function generateToken() {
  return jwt.sign({ sub: "test" }, JWT_SECRET, { expiresIn: "1h" });
}

describe("Matrix Routes", () => {
  describe("POST /api/v1/matrix/stats", () => {
    test("returns 401 without JWT", async () => {
      const res = await request(app)
        .post("/api/v1/matrix/stats")
        .send({ matrices: { q: [[1]], r: [[1]] } });

      expect(res.status).toBe(401);
      expect(res.body.success).toBe(false);
      expect(res.body.error.code).toBe("UNAUTHORIZED");
    });

    test("returns 403 with invalid JWT", async () => {
      const res = await request(app)
        .post("/api/v1/matrix/stats")
        .set("Authorization", "Bearer invalid_token_here")
        .send({ matrices: { q: [[1]], r: [[1]] } });

      expect(res.status).toBe(403);
      expect(res.body.success).toBe(false);
      expect(res.body.error.code).toBe("FORBIDDEN");
    });

    test("returns 400 for missing matrices", async () => {
      const token = generateToken();
      const res = await request(app)
        .post("/api/v1/matrix/stats")
        .set("Authorization", `Bearer ${token}`)
        .send({});

      expect(res.status).toBe(400);
      expect(res.body.success).toBe(false);
      expect(res.body.error.code).toBe("INVALID_MATRIX");
    });

    test("returns 200 with valid matrices", async () => {
      const token = generateToken();
      const res = await request(app)
        .post("/api/v1/matrix/stats")
        .set("Authorization", `Bearer ${token}`)
        .send({
          matrices: {
            q: [[1, 0], [0, 1]],
            r: [[2, 3], [0, 4]],
          },
        });

      expect(res.status).toBe(200);
      expect(res.body.success).toBe(true);
      expect(res.body.data).toHaveProperty("max");
      expect(res.body.data).toHaveProperty("min");
      expect(res.body.data).toHaveProperty("average");
      expect(res.body.data).toHaveProperty("sum");
      expect(res.body.data).toHaveProperty("diagonal");
      expect(res.body.data.max).toBe(4);
      expect(res.body.data.min).toBe(0);
      expect(res.body.data.diagonal.q).toBe(true);
      expect(res.body.data.diagonal.r).toBe(false);
    });

    test("returns diagonal.q=true for identity and diagonal.r=true for diagonal", async () => {
      const token = generateToken();
      const res = await request(app)
        .post("/api/v1/matrix/stats")
        .set("Authorization", `Bearer ${token}`)
        .send({
          matrices: {
            q: [[1, 0], [0, 1]],
            r: [[5, 0], [0, 3]],
          },
        });

      expect(res.status).toBe(200);
      expect(res.body.data.diagonal.q).toBe(true);
      expect(res.body.data.diagonal.r).toBe(true);
    });
  });
});
