const output = document.getElementById("output");
const loginView = document.getElementById("loginView");
const appView = document.getElementById("appView");
const memberName = document.getElementById("memberName");

let token = localStorage.getItem("team_token") || "";
let savedName = localStorage.getItem("team_name") || "harikeerthan";

function showLogin() {
  loginView.classList.remove("hidden");
  appView.classList.add("hidden");
}

function showApp() {
  loginView.classList.add("hidden");
  appView.classList.remove("hidden");
  memberName.textContent = savedName;
}

function printResult(title, result) {
  output.textContent = `${title}\n${JSON.stringify(result, null, 2)}`;
}

async function postJSON(url, body, withAuth = false) {
  const headers = { "Content-Type": "application/json" };
  if (withAuth && token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const res = await fetch(url, {
    method: "POST",
    headers,
    body: JSON.stringify(body),
  });

  const raw = await res.text();
  let data;
  try {
    data = JSON.parse(raw);
  } catch {
    data = { message: raw };
  }

  return { ok: res.ok, status: res.status, data };
}

if (token) {
  showApp();
} else {
  showLogin();
}

document.getElementById("loginForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  const payload = Object.fromEntries(new FormData(event.target).entries());
  const result = await postJSON("/login", payload);

  if (result.ok && result.data.token) {
    token = result.data.token;
    savedName = result.data.name || "harikeerthan";
    localStorage.setItem("team_token", token);
    localStorage.setItem("team_name", savedName);
    showApp();
  }

  printResult("Login", result);
});

document.getElementById("logoutBtn").addEventListener("click", () => {
  token = "";
  savedName = "harikeerthan";
  localStorage.removeItem("team_token");
  localStorage.removeItem("team_name");
  showLogin();
  output.textContent = "Logged out.";
});

document.getElementById("registerForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  const payload = Object.fromEntries(new FormData(event.target).entries());
  const result = await postJSON("/register", payload, true);
  printResult("Register", result);
});

document.getElementById("verifyForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  const payload = Object.fromEntries(new FormData(event.target).entries());
  const result = await postJSON("/verify-otp", payload, true);
  printResult("Verify OTP", result);
});
