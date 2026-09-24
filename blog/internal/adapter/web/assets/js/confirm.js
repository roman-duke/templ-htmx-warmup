document.addEventListener("htmx:confirm", function (e) {
  if (!e.detail.question) return; // element had no hx-confirm, do nothing
  e.preventDefault(); // prevent the native dialog from showing

  const dialog = document.getElementById("confirm-dialog");
  const confirmMessage = document.getElementById("confirm-message");

  confirmMessage.textContent = e.detail.question;
  dialog.showModal();

  document.getElementById("confirm-yes").onclick = () => { dialog.close(); e.detail.issueRequest(true) }
  document.getElementById("confirm-no").onclick = () => { dialog.close(); }
});
