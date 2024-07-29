document.documentElement.setAttribute('data-theme', 'mytheme');




document.addEventListener('DOMContentLoaded', function() {
  const titleInput = document.getElementById('title-input');
  const contentInput = document.querySelector("textarea")
  const titleWarning = document.getElementById('title-warning');
  const contentWarning = document.getElementById('content-warning');

  // Prevent form submission on Enter key for title input
  if (titleInput) {
    titleInput.addEventListener("keydown", e => {
      if (e.key === "Enter") {
        e.preventDefault();
        contentInput.focus();
      }
    })
  }

  if (titleWarning) {
    console.log(titleInput)
    titleInput.focus()
  }

  if (contentWarning) {
    contentInput.focus();
  }
});



