const state = {
  token: localStorage.getItem("token") || "",
  expiresAt: localStorage.getItem("expiresAt") || "",
  userId: localStorage.getItem("userId") || "",
  socket: null,
  desktopOnline: false,
  reconnectTimer: 0,
  reconnectAttempts: 0,
};

const $ = (id) => document.getElementById(id);
const isFilePreview = location.protocol === "file:";

function init() {
  $("loginButton").addEventListener("click", login);
  $("registerButton").addEventListener("click", registerAccount);
  $("logoutButton").addEventListener("click", logout);
  $("sendButton").addEventListener("click", sendText);
  $("clearButton").addEventListener("click", () => {
    $("textInput").value = "";
    $("sendHint").textContent = "";
  });
  if (isFilePreview) {
    showLogin();
    $("loginHint").textContent = "请通过 http://服务器IP:端口 或 https://域名 访问，不要直接打开 HTML 文件";
    return;
  }
  if (hasValidToken()) {
    showApp();
    connectSocket();
  } else {
    showLogin();
  }
}

function hasValidToken() {
  return state.token && state.expiresAt && new Date(state.expiresAt) > new Date();
}

async function login() {
  if (isFilePreview) {
    $("loginHint").textContent = "本地文件模式不能登录，请先启动服务端再访问网页地址";
    return;
  }
  $("loginHint").textContent = "登录中...";
  const res = await authRequest("/api/login", "loginHint");
  if (!res) return;
  if (!res.ok) {
    $("loginHint").textContent = "账号或密码不正确";
    return;
  }
  const session = await res.json();
  saveSession(session);
  showApp();
  connectSocket();
}

async function registerAccount() {
  if (isFilePreview) {
    $("loginHint").textContent = "本地文件模式不能注册，请先启动服务端再访问网页地址";
    return;
  }
  $("loginHint").textContent = "注册中...";
  const res = await authRequest("/api/register", "loginHint");
  if (!res) return;
  if (!res.ok) {
    const reason = await res.text();
    $("loginHint").textContent = `注册失败：${reason || "账号至少 3 位，密码至少 6 位，或账号已存在"}`;
    return;
  }
  const session = await res.json();
  saveSession(session);
  showApp();
  connectSocket();
}

async function authRequest(path, hintId) {
  try {
    return await fetch(path, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        username: $("usernameInput").value.trim(),
        password: $("passwordInput").value,
        remember: $("rememberInput").checked,
      }),
    });
  } catch (err) {
    $(hintId).textContent = "连接不上服务端，请先双击启动本地服务.cmd，或确认服务器正在运行";
    return null;
  }
}

async function logout() {
  await fetch("/api/logout", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ token: state.token }),
  }).catch(() => {});
  clearSession();
  showLogin();
}

function saveSession(session) {
  state.token = session.token;
  state.expiresAt = session.expiresAt;
  state.userId = session.userId;
  localStorage.setItem("token", state.token);
  localStorage.setItem("expiresAt", state.expiresAt);
  localStorage.setItem("userId", state.userId);
}

function clearSession() {
  if (state.socket) state.socket.close();
  ["token", "expiresAt", "userId"].forEach((key) => localStorage.removeItem(key));
  state.token = "";
  state.expiresAt = "";
  state.userId = "";
  state.desktopOnline = false;
}

function showLogin() {
  $("loginPanel").hidden = false;
  $("appPanel").hidden = true;
  $("loginHint").textContent = "未登录";
}

function showApp() {
  $("loginPanel").hidden = true;
  $("appPanel").hidden = false;
  setStatus("disconnected");
}

function connectSocket() {
  if (state.socket && state.socket.readyState === WebSocket.OPEN) return;
  if (state.reconnectTimer) {
    clearTimeout(state.reconnectTimer);
    state.reconnectTimer = 0;
  }
  const proto = location.protocol === "https:" ? "wss:" : "ws:";
  state.socket = new WebSocket(`${proto}//${location.host}/ws`);
  state.socket.onopen = () => {
    state.reconnectAttempts = 0;
    $("sendHint").textContent = "";
    sendSocket({ type: "web_join", token: state.token });
  };
  state.socket.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    if (msg.type === "presence") {
      state.desktopOnline = msg.desktop_online;
      setStatus(msg.status);
    }
    if (msg.type === "ack" && msg.ok && msg.delivered === false) {
      $("sendHint").textContent = "电脑端未连接，内容已保留，但不会自动粘贴";
    }
  };
  state.socket.onclose = () => {
    setStatus("server_down");
    scheduleReconnect();
  };
}

function scheduleReconnect() {
  if (!hasValidToken()) return;
  state.reconnectAttempts += 1;
  if (state.reconnectAttempts > 8) {
    $("sendHint").textContent = "服务器断开，请先双击启动本地服务.cmd，再刷新页面";
    return;
  }
  const delay = Math.min(1500 * state.reconnectAttempts, 10000);
  $("sendHint").textContent = `服务器断开，${Math.round(delay / 1000)} 秒后重连`;
  state.reconnectTimer = setTimeout(connectSocket, delay);
}

function setStatus(status) {
  const connected = status === "connected";
  $("statusDot").className = `dot ${connected ? "green" : "red"}`;
  $("statusText").textContent = connected ? "已连接" : "未连接";
  if (status === "server_down") $("statusText").textContent = "服务器断开";
}

function sendText() {
  const text = $("textInput").value;
  if (!text.trim()) {
    $("sendHint").textContent = "请输入要发送的文字";
    return;
  }
  sendSocket({ type: "send_text", token: state.token, text });
  $("sendHint").textContent = state.desktopOnline
    ? "已发送，有输入光标的电脑会自动粘贴"
    : "电脑端未连接，内容已保留";
}

function sendSocket(value) {
  if (!state.socket || state.socket.readyState !== WebSocket.OPEN) {
    $("sendHint").textContent = "服务器断开，正在重连";
    return;
  }
  state.socket.send(JSON.stringify(value));
}

init();
