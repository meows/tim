document.addEventListener('DOMContentLoaded', function() {
  const titleInput = document.getElementById('title-input');
  const contentInput = document.querySelector("textarea")
  titleInput.addEventListener("keydown", e => {
    if (e.key === "Enter") {
      e.preventDefault();
      contentInput.focus();
    }
  })
});
