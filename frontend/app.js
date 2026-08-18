// Frontend application for QR Matrix Decomposition.
// DECISION: API_BASE_URL points to the Go API gateway.
const API_BASE_URL = "http://localhost:8080";

let authToken = null;

// Renders a 2D matrix as an HTML table.
function renderMatrix(containerId, matrix) {
  const container = document.getElementById(containerId);
  let html = '<div class="matrix-wrapper"><table>';
  for (const row of matrix) {
    html += "<tr>";
    for (const val of row) {
      html += `<td>${Number.isFinite(val) ? parseFloat(val.toFixed(4)) : val}</td>`;
    }
    html += "</tr>";
  }
  html += "</table></div>";
  container.innerHTML = html;
}

// Renders the stats table.
function renderStats(stats) {
  const tbody = document.querySelector("#statsTable tbody");
  tbody.innerHTML = "";
  const rows = [
    ["Máximo", stats.max],
    ["Mínimo", stats.min],
    ["Promedio", stats.average],
    ["Suma total", stats.sum],
  ];
  for (const [label, value] of rows) {
    const tr = document.createElement("tr");
    tr.innerHTML = `<td>${label}</td><td>${parseFloat(value.toFixed(4))}</td>`;
    tbody.appendChild(tr);
  }
}

// Renders the diagonal flags table.
function renderDiagonal(diagonal) {
  const tbody = document.querySelector("#diagonalTable tbody");
  tbody.innerHTML = "";
  const rows = [
    ["Q", diagonal.q],
    ["R", diagonal.r],
  ];
  for (const [label, value] of rows) {
    const tr = document.createElement("tr");
    tr.innerHTML = `<td>${label}</td><td>${value ? "Sí" : "No"}</td>`;
    tbody.appendChild(tr);
  }
}

// Shows an error message to the user.
function showError(message) {
  const el = document.getElementById("error");
  el.textContent = message;
  el.classList.remove("hidden");
  document.getElementById("results").classList.add("hidden");
  document.getElementById("warning").classList.add("hidden");
}

// Shows a warning message.
function showWarning(message) {
  const el = document.getElementById("warning");
  el.textContent = message;
  el.classList.remove("hidden");
}

// Hides all messages.
function hideMessages() {
  document.getElementById("error").classList.add("hidden");
  document.getElementById("warning").classList.add("hidden");
}

// Performs the mock login to obtain a JWT token.
async function login() {
  const res = await fetch(`${API_BASE_URL}/api/v1/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username: "candidate", password: "challenge2024" }),
  });
  const data = await res.json();
  if (!data.success) {
    throw new Error(data.error?.message || "Error al autenticar");
  }
  return data.data.token;
}

// Main calculation flow: login → fetch QR → display.
async function calculate() {
  hideMessages();
  const btn = document.getElementById("calculateBtn");
  btn.disabled = true;
  btn.textContent = "Calculando...";

  try {
    const matrixText = document.getElementById("matrixInput").value.trim();
    let matrix;
    try {
      matrix = JSON.parse(matrixText);
    } catch {
      showError("El formato JSON de la matriz es inválido");
      return;
    }

    if (!authToken) {
      authToken = await login();
    }

    const res = await fetch(`${API_BASE_URL}/api/v1/matrix/qr`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
      body: JSON.stringify({ matrix }),
    });

    const data = await res.json();

    if (!data.success) {
      showError(data.error?.message || "Error desconocido");
      return;
    }

    if (data.warning) {
      showWarning(data.warning);
    }

    renderMatrix("qTable", data.data.qr.q);
    renderMatrix("rTable", data.data.qr.r);

    if (data.data.stats) {
      renderStats(data.data.stats);
      renderDiagonal(data.data.stats.diagonal);
    } else {
      document.querySelector("#statsTable tbody").innerHTML =
        '<tr><td colspan="2">Estadísticas no disponibles</td></tr>';
      document.querySelector("#diagonalTable tbody").innerHTML =
        '<tr><td colspan="2">Estadísticas no disponibles</td></tr>';
    }

    document.getElementById("results").classList.remove("hidden");
  } catch (err) {
    showError(err.message || "Error de conexión con la API");
  } finally {
    btn.disabled = false;
    btn.textContent = "Calcular";
  }
}
