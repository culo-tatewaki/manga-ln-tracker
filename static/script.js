function resetForm() {
    document.getElementById("form").reset(); 
}


document.addEventListener("DOMContentLoaded", () => {

    [...document.getElementsByClassName("edit-btn")].forEach((btn) => {
        btn.addEventListener("click", ({target}) => {
            const card = target.parentElement.parentElement.parentElement.parentElement.parentElement

            document.getElementById("form-id").value = card.getAttribute("data-id")
            document.getElementById("form-type").value = card.querySelector(".series-type").textContent 
            document.getElementById("form-title").value = card.querySelector(".series-title").textContent
            document.getElementById("form-chapters").value = card.querySelector(".series-chapters").textContent
            document.getElementById("form-volumes").value = card.querySelector(".series-volumes").textContent
            document.getElementById("form-author").value = card.querySelector(".series-author").textContent
            document.getElementById("form-image").value = card.querySelector(".series-image").src 
            document.getElementById("form-rating").value = card.querySelector(".series-rating").textContent
        })
    })
    
})