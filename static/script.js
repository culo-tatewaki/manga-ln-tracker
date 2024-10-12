document.addEventListener("DOMContentLoaded", () => {

    document.getElementById("add-form").addEventListener("submit", function(event) {
        event.preventDefault();
    
        // Capture form data
        const email = document.getElementById("exampleInputEmail1").value;
        const message = document.getElementById("exampleInputPassword1").value;
    
        console.log(email, message);

        document.getElementById("add-form").reset();
    })

})