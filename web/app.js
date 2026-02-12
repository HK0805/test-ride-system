const output = document.getElementById("output");
const tokenState = document.getElementById("tokenState");
let token = localStorage.getItem("team_token") || "";
refreshTokenState();

function refreshTokenState() {
  tokenState.textContent = token ? "Token: logged in" : "Token: not logged in";
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

document.getElementById("signupForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const f = Object.fromEntries(new FormData(e.target).entries());
  logResult("Signup", await post("/signup", f));
});

document.getElementById("loginForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const f = Object.fromEntries(new FormData(e.target).entries());
  const res = await post("/login", f);
  if (res.ok && res.data.token) {
    token = res.data.token;
    localStorage.setItem("team_token", token);
    refreshTokenState();
  }
  logResult("Login", res);
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
