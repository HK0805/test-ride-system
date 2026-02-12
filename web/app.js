const output = document.getElementById("output");
const loginView = document.getElementById("loginView");
const appView = document.getElementById("appView");
const memberName = document.getElementById("memberName");

let token = localStorage.getItem("team_token") || "";
let savedName = localStorage.getItem("team_name") || "Team Member";

function showLogin() {
  loginView.classList.remove("hidden");
  appView.classList.add("hidden");
}

function showApp() {
  loginView.classList.add("hidden");
  appView.classList.remove("hidden");
  memberName.textContent = savedName;
}

function logResult(title, data) {
  output.textContent = `${title}\n${JSON.stringify(data, null, 2)}`;
}

async function post(url, body, auth = false) {
  const headers = { "Content-Type": "application/json" };
  if (auth && token) headers.Authorization = `Bearer ${token}`;
  const res = await fetch(url, { method: "POST", headers, body: JSON.stringify(body) });
  const text = await res.text();
  let parsed;
  try { parsed = JSON.parse(text); } catch { parsed = { message: text }; }
  return { ok: res.ok, status: res.status, data: parsed };
}

if (token) {
  showApp();
} else {
  showLogin();
}

document.getElementById("loginForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const f = Object.fromEntries(new FormData(e.target).entries());
  const res = await post("/login", f);
  if (res.ok && res.data.token) {
    token = res.data.token;
    savedName = res.data.name || "harikeerthan";
    localStorage.setItem("team_token", token);
    localStorage.setItem("team_name", savedName);
    showApp();
  }
  logResult("Login", res);
});

document.getElementById("logoutBtn").addEventListener("click", () => {
  token = "";
  savedName = "Team Member";
  localStorage.removeItem("team_token");
  localStorage.removeItem("team_name");
  showLogin();
  output.textContent = "Logged out.";
});

document.getElementById("registerForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const f = Object.fromEntries(new FormData(e.target).entries());
  logResult("Register", await post("/register", f, true));
});

document.getElementById("verifyForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const f = Object.fromEntries(new FormData(e.target).entries());
  logResult("Verify OTP", await post("/verify-otp", f, true));
});
