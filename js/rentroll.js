function populateRentRoll(email) {
    if (email === '') {
        console.log("populateRentRoll - empty email!");
        return;
    }
    $.ajax({
        url: '/tenants',
        data: { 'email': email },
        success: function( tenants ) {
            document.getElementById('tenants').innerHTML = tenants
        }
    });
}

function gRentRoll(resp) {
    populateRentRoll(gGetEmail(resp));
}

function fbRentRoll() {
    FB.api('/me', function(response) {
            populateRentRoll(response.email);
    });
}
 
