function get(name){
    if(name=(new RegExp('[?&]'+encodeURIComponent(name)+'=([^&]*)')).exec(location.search))
        return decodeURIComponent(name[1]);
}

function gGetEmail(resp){
    if (resp.result.emails.length == 0) {
        console.log("gRentRollCallback - error no email!");
        return ''
    }
    return resp.result.emails[0].value;
}

function gGetName(resp){
    if (resp.result.emails.length == 0) {
        console.log("gRentRollCallback - error no email!");
        return ''
    }
    return resp.result.displayName;    
}

var gSigninFailed = 0;

function gSignin(authResult) {
    if (authResult['status']['signed_in']) {
        client_id = authResult['client_id'];
        access_token = authResult['access_token'];
        code = authResult['code'];
        gSigninButton = document.getElementById('gSigninButton');
        if (gSigninButton !== null) {
            gSigninButton.setAttribute('style', 'display: none');
        }
        fbSigninButton = document.getElementById('fbSigninButton');
        if (fbSigninButton !== null) {
            fbSigninButton.setAttribute('style', 'display: none');
        }
        document.getElementById('gsettings').setAttribute('style', 'display: inline');
        gapi.client.load('plus', 'v1').then(function() {
            var request = gapi.client.plus.people.get({
                'userId': 'me'
            });
            request.then(function(resp) {
                email = "";
                if (resp.result.emails.length == 0) {
                    console.log("Error no email!");
                    email = "dummy@dummy.com";
                } else {
                    email = resp.result.emails[0].value;
                }
                gSignoutButton = document.getElementById('gSignoutButton');
                if (gSignoutButton !== null) {
                    gSignoutButton.innerHTML = "" + email + "";
                }
                executeCallback('g-callback', resp);
            }, function(reason) {
                console.log('Error: ' + reason.result.error.message);
            });
        });
    } else {
        // Update the app to reflect a signed out user
        // Possible error values:
        //   "user_signed_out" - User is signed-out
        //   "access_denied" - User denied access to your app
        //   "immediate_failed" - Could not automatically log in the user
        if (gSigninFailed == 0){
            console.log('Sign-in state: ' + authResult['error']);
            startFBLogin();
            gSigninFailed += 1;
        }
    }
}

function fbLogin(response){
    // The response object is returned with a status field that lets the
    // app know the current login status of the person.
    // Full docs on the response object can be found in the documentation
    // for FB.getLoginStatus().
    if (response.status === 'connected') {
        // Logged into your app and Facebook.
        document.getElementById('gSigninButton').setAttribute('style', 'display: none');
        document.getElementById('fbsettings').setAttribute('style', 'display: inline');
        FB.api('/me', function(response) {
            document.getElementById('fbEmail').innerHTML = response.email;            
        });
        executeCallback('fb-callback', response);
    } else if (response.status === 'not_authorized') {
        // The person is logged into Facebook, but not your app.
        document.getElementById('fbstatus').innerHTML = 'Please log into this app.';
        executeCallback('notloggedin-callback', '');        
    } else {
        // The person is not logged into Facebook
        executeCallback('notloggedin-callback', '');
        
    }
}

function executeCallback(name, response) {
    var callback = document.getElementById(name);
    if (callback === null) {
        console.log('executeCallback - ERROR NO CALLBACK: ' + name);
        return;
    }
    var callbackAttr = callback.attributes['callback'];
    if (callbackAttr === null) {
        console.log('executeCallback - ERROR NO CALLBACK: ' + name);
        return;
    } 
    var callbackName = callbackAttr.value;
    if (callbackName == null || callbackName === '') {
        console.log('executeCallback - empty callback name');
        return;
    }
    // call the callback.
    window[callbackName](response);
}

function createAccount(email, loginMethod) {
    if (email === '') {
        console.log("populateRentRoll - empty email!");
        return;
    }
    $.ajax({url: '/createaccount',
            data: { 'email': email, 'loginmethod': loginMethod },
            success: function( data ) {
                window.location.href = '/rentroll';
            }
        });
}
function gSignup(resp) {
    createAccount(gGetEmail(resp), 'Google');
}

function fbSignup() {
    FB.api('/me', function(response) {
        createAccount(response.email, 'Facebook');
    });
}

function gSettings(resp){
    setSettings(gGetName(resp), gGetEmail(resp));
}

function fbSettings(){
    FB.api('/me', function(response) {
        setSettings(response.name, response.email);
    });

}

function setSettings(name, email) {
    var nameEle = document.getElementById("name");
    var emailEle = document.getElementById("email");
    nameEle.innerHTML = name;
    emailEle.innerHTML = email;
}
