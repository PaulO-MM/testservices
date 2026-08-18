// Express application setup with middleware and routes.
const express = require("express");
const matrixRoutes = require("./routes/matrix.routes");
const { errorResponse } = require("./utils/errors");

const app = express();

app.use(express.json());

app.get("/health", (req, res) => {
  res.json({ status: "ok" });
});

app.use("/api/v1/matrix", matrixRoutes);

// Central error handler: translates uncaught exceptions to standard format.
// eslint-disable-next-line no-unused-vars
app.use((err, req, res, next) => {
  console.error("Unhandled error:", err);
  res.status(500).json(errorResponse("INTERNAL_ERROR", "Error interno del servidor"));
});

module.exports = app;
