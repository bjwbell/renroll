/*jslint browser: true*/
/*globals $, FB, gapi, fbLogin, setSettings, gGetEmail*/

"use strict";

function populateRentRoll(email) {
    if (email === '') {
        console.log("populateRentRoll - empty email!");
        return;
    }
    var dbName = document.getElementById("DbName");
    if (dbName !== undefined) {
        dbName.value = email;
    }
    $.ajax({
        url: '/tenants',
        data: { 'email': email },
        success: function (tenants) {
            document.getElementById('tenants').innerHTML = tenants;
        }
    });
}

function gRentRoll(resp) {
    populateRentRoll(gGetEmail(resp));
}

function fbRentRoll() {
    FB.api('/me', function (response) {
        populateRentRoll(response.email);
    });
}

function rentRollNotLoggedIn() {
    window.location.href = "/submit";
    /*var signinForm = document.getElementById("signinform");
    if (signinForm == null) {
        console.log("rentRollNotLoggedIn - ERROR NO SIGNIN FORM ELEMENT");
        return;
    }
    $.ajax({
        url: '/signinform',
        success: function( form ) {
            signinForm.innerHTML = form;
        }
    });*/
}

function rentRollTemplateNotLoggedIn() { }
function gRentRollTemplate(resp) { }
function fbRentRollTemplate() { }

