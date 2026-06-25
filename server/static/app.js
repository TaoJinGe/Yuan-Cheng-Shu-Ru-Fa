const state = {
  token: localStorage.getItem("token") || "",
  expiresAt: localStorage.getItem("expiresAt") || "",
  userId: localStorage.getItem("userId") || "",
  socket: null,
  desktopOnline: false,
  reconnectTimer: 0,
  reconnectAttempts: 0,
  pendingSendCount: 0,
  undoText: "",
};

const $ = (id) => document.getElementById(id);
const isFilePreview = location.protocol === "file:";

function init() {
  updateViewportHeight();
  window.addEventListener("resize", updateViewportHeight);
  if (window.visualViewport) {
    window.visualViewport.addEventListener("resize", updateViewportHeight);
    window.visualViewport.addEventListener("scroll", updateViewportHeight);
  }
  $("loginButton").addEventListener("click", login);
  $("registerButton").addEventListener("click", registerAccount);
  $("logoutButton").addEventListener("click", logout);
  bindActionButton($("undoButton"), undoText);
  bindActionButton($("sendButton"), sendText);
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

function updateViewportHeight() {
  const height = window.visualViewport ? window.visualViewport.height : window.innerHeight;
  document.documentElement.style.setProperty("--app-height", `${Math.round(height)}px`);
}

function bindActionButton(button, action) {
  let handledByPointer = false;
  let suppressClickUntil = 0;
  button.addEventListener("pointerdown", (event) => {
    button.classList.add("is-pressing");
    if (document.activeElement === $("textInput")) {
      handledByPointer = true;
      event.preventDefault();
    }
  });
  button.addEventListener("pointerup", (event) => {
    if (!handledByPointer) return;
    handledByPointer = false;
    suppressClickUntil = Date.now() + 700;
    event.preventDefault();
    button.classList.remove("is-pressing");
    flashButton(button);
    action();
  });
  button.addEventListener("pointercancel", () => {
    handledByPointer = false;
    button.classList.remove("is-pressing");
  });
  button.addEventListener("click", (event) => {
    button.classList.remove("is-pressing");
    if (Date.now() < suppressClickUntil) {
      event.preventDefault();
      return;
    }
    flashButton(button);
    action();
  });
}

function flashButton(button) {
  button.classList.remove("is-flashing");
  void button.offsetWidth;
  button.classList.add("is-flashing");
  setTimeout(() => button.classList.remove("is-flashing"), 420);
}

function setTextInput(text) {
  const input = $("textInput");
  input.value = text;
  input.setSelectionRange(text.length, text.length);
}

function focusTextInput() {
  requestAnimationFrame(() => {
    $("textInput").focus({ preventScroll: true });
  });
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
    if (msg.type === "ack") {
      handleAck(msg);
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

function handleAck(msg) {
  if (msg.message_type !== "send_text" || state.pendingSendCount < 1) return;
  state.pendingSendCount -= 1;
  if (!msg.ok) {
    if (!$("textInput").value && state.undoText) setTextInput(state.undoText);
    $("sendHint").textContent = msg.error ? `发送失败：${msg.error}` : "发送失败，请稍后重试";
    return;
  }
  if (msg.delivered === false) {
    $("sendHint").textContent = "电脑端未连接，可点撤销恢复";
  }
}

function undoText() {
  if (!state.undoText) {
    $("sendHint").textContent = "没有可撤销内容";
    focusTextInput();
    return;
  }
  setTextInput(state.undoText);
  $("sendHint").textContent = "已恢复上一条，可继续修改";
  focusTextInput();
}

function sendText() {
  const text = $("textInput").value;
  if (!text.trim()) {
    $("sendHint").textContent = "请输入要发送的文字";
    focusTextInput();
    return;
  }
  if (!sendSocket({ type: "send_text", token: state.token, text })) {
    focusTextInput();
    return;
  }
  state.pendingSendCount += 1;
  state.undoText = text;
  setTextInput("");
  $("sendHint").textContent = state.desktopOnline
    ? "已发送，输入框已清空"
    : "电脑端未连接，可点撤销恢复";
  focusTextInput();
}

function sendSocket(value) {
  if (!state.socket || state.socket.readyState !== WebSocket.OPEN) {
    $("sendHint").textContent = "服务器断开，正在重连";
    return false;
  }
  state.socket.send(JSON.stringify(value));
  return true;
}

init();
