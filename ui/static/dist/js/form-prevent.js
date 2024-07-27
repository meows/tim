document.addEventListener('DOMContentLoaded', function() {
  const titleInput = document.getElementById('title-input');
  titleInput.addEventListener("keydown", e => {
    if (e.key === "Enter") {
      e.preventDefault();
    }
  })
});
