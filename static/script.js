function resetForm() {
  document.getElementById("form").reset();
}


document.addEventListener("DOMContentLoaded", () => {

  [...document.getElementsByClassName("edit-btn")].forEach((btn) => {
    btn.addEventListener("click", ({ currentTarget }) => {
      document.getElementById("form-id").value = currentTarget.getAttribute("data-id")
      document.getElementById("form-type").value = currentTarget.querySelector(".series-type").textContent
      document.getElementById("form-title").value = currentTarget.querySelector(".series-title").textContent
      document.getElementById("form-chapters").value = currentTarget.querySelector(".series-chapters").textContent
      document.getElementById("form-volumes").value = currentTarget.querySelector(".series-volumes").textContent
      document.getElementById("form-release-date").value = currentTarget.querySelector(".series-release-date").textContent
      document.getElementById("form-status").value = currentTarget.querySelector(".series-status").textContent
      document.getElementById("form-author").value = currentTarget.querySelector(".series-author").textContent
      document.getElementById("form-image").value = currentTarget.querySelector(".series-image").src
      document.getElementById("form-rating").value = currentTarget.querySelector(".series-rating").textContent
    })
  });

  [...document.getElementsByClassName("delete-btn")].forEach(btn => {
    btn.addEventListener("click", () => {
      const id = document.getElementById("form-id").value

      fetch(`/delete?id=${id}`, {
        method: 'DELETE',
      });
    })
  });

})