(() => {
  const key = "quizzy-theme";
  const themes = ["graphite", "crimson"];
  const normalize = value => themes.includes(value) ? value : "graphite";

  function apply(value) {
    const theme = normalize(value);
    document.documentElement.dataset.theme = theme;
    document.querySelectorAll("[data-theme-choice]").forEach(button => {
      button.setAttribute("aria-pressed", String(button.dataset.themeChoice === theme));
    });
    const meta = document.querySelector('meta[name="theme-color"]');
    if (meta) meta.content = theme === "crimson" ? "#230305" : "#292a2c";
  }

  // Apply before the first paint to avoid flashing the other theme on navigation.
  let saved = "graphite";
  try { saved = localStorage.getItem(key); } catch (_) { /* Storage may be disabled. */ }
  apply(saved);

  document.addEventListener("DOMContentLoaded", () => {
    apply(document.documentElement.dataset.theme);
    document.querySelectorAll("[data-theme-choice]").forEach(button => {
      button.addEventListener("click", () => {
        const theme = normalize(button.dataset.themeChoice);
        apply(theme);
        try { localStorage.setItem(key, theme); } catch (_) { /* Keep the session usable. */ }
      });
    });
  });

  // Keep host and player tabs consistent when the theme changes in either one.
  window.addEventListener("storage", event => {
    if (event.key === key || event.key === null) apply(event.newValue);
  });
})();
