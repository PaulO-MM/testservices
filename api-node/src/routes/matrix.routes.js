// Route definitions for matrix-related endpoints.
const express = require("express");
const authMiddleware = require("../middleware/auth.middleware");
const { computeMatrixStats } = require("../controllers/matrix.controller");

const router = express.Router();

router.post("/stats", authMiddleware, computeMatrixStats);

module.exports = router;
