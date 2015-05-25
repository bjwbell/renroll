             function get(name){
        if(name=(new RegExp('[?&]'+encodeURIComponent(name)+'=([^&]*)')).exec(location.search))
        return decodeURIComponent(name[1]);
        }
     function signinCallback(authResult) {
         if (authResult['status']['signed_in']) {
             client_id = authResult['client_id']
             access_token = authResult['access_token']
             code = authResult['code']
             document.getElementById('signinButton').setAttribute('style', 'display: none');
             
             $.ajax({
                 url: "/oauth2callback",
                 data: {
                     client_id: client_id,
                     access_token: access_token,
                     code: code,                   
                 },
                 success: function( data ) {
                     document.getElementById('signoutButton').innerHTML = "" + data + "";
                 }
             });
         } else {
             // Update the app to reflect a signed out user
             // Possible error values:
             //   "user_signed_out" - User is signed-out
             //   "access_denied" - User denied access to your app
             //   "immediate_failed" - Could not automatically log in the user
             console.log('Sign-in state: ' + authResult['error']);
         }
     }
