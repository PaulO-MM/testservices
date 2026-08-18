// Controller for the matrix statistics endpoint.
const { computeStats } = require("../services/stats.service");
const { successResponse, errorResponse } = require("../utils/errors");

function validateMatrix(matrix, name) {
  if (!Array.isArray(matrix) || matrix.length === 0) {
    return `${name} no es una matriz válida`;
  }
  const cols = matrix[0].length;
  if (cols === 0) return `${name} tiene filas vacías`;
  for (let i = 0; i < matrix.length; i++) {
    if (!Array.isArray(matrix[i]) || matrix[i].length !== cols) {
      return `${name} no es una matriz rectangular`;
    }
    for (let j = 0; j < matrix[i].length; j++) {
      if (typeof matrix[i][j] !== "number" || isNaN(matrix[i][j])) {
        return `${name} contiene elementos no numéricos en [${i}][${j}]`;
      }
    }
  }
  return null;
}

function computeMatrixStats(req, res) {
  try {
    const { matrices } = req.body;

    if (!matrices || !matrices.q || !matrices.r) {
      return res.status(400).json(
        errorResponse("INVALID_MATRIX", "Faltan las matrices q o r en el cuerpo de la solicitud")
      );
    }

    const errQ = validateMatrix(matrices.q, "q");
    if (errQ) {
      return res.status(400).json(errorResponse("INVALID_MATRIX", errQ));
    }

    const errR = validateMatrix(matrices.r, "r");
    if (errR) {
      return res.status(400).json(errorResponse("INVALID_MATRIX", errR));
    }

    const stats = computeStats(matrices.q, matrices.r);
    return res.status(200).json(successResponse(stats));
  } catch (err) {
    return res.status(500).json(
      errorResponse("INTERNAL_ERROR", "Error interno al calcular estadísticas")
    );
  }
}

module.exports = { computeMatrixStats };
