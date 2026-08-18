// HTTP server entry point for the Node.js stats API.
const app = require("./app");

const PORT = process.env.PORT || 3000;

app.listen(PORT, () => {
  // DECISION: Minimal startup log; in production use a structured logger.
});
