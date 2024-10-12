function resetForm() {
    document.getElementById("form").reset(); 
}


document.addEventListener("DOMContentLoaded", () => {

    [...document.getElementsByClassName("edit-btn")].forEach((btn) => {
        btn.addEventListener("click", ({target}) => {
            const card = target.parentElement.parentElement.parentElement.parentElement

            document.getElementById("form-id").value = card.getAttribute("data-id")
            document.getElementById("form-type").value = card.querySelector(".book-type").textContent 
            document.getElementById("form-series").value = card.querySelector(".book-series").textContent
            document.getElementById("form-volume").value = card.querySelector(".book-volume").textContent
            document.getElementById("form-author").value = card.querySelector(".book-author").textContent
            document.getElementById("form-image").value = card.querySelector(".book-image").src 
            document.getElementById("form-rating").value = card.querySelector(".book-rating").textContent
        })
    })
    
})