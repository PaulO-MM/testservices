// JWT authentication middleware for Express.
const jwt = require("jsonwebtoken");

const JWT_SECRET = process.env.JWT_SECRET || "dev_secret_change_me";

function authMiddleware(req, res, next) {
  const authHeader = req.headers.authorization;
  if (!authHeader) {
    return res.status(401).json({
      success: false,
      error: { code: "UNAUTHORIZED", message: "Falta el header de autorización" },
    });
  }

  const parts = authHeader.split(" ");
  if (parts.length !== 2 || parts[0] !== "Bearer") {
    return res.status(401).json({
      success: false,
      error: {
        code: "UNAUTHORIZED",
        message: "Formato de autorización inválido. Use: Bearer <token>",
      },
    });
  }

  try {
    const decoded = jwt.verify(parts[1], JWT_SECRET);
    req.user = decoded;
    req.token = parts[1];
    next();
  } catch (err) {
    return res.status(403).json({
      success: false,
      error: { code: "FORBIDDEN", message: "Token inválido o expirado" },
    });
  }
}

module.exports = authMiddleware;
