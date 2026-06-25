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
