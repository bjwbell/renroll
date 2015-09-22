/*jslint browser: true*/
/*globals $, FB, gapi, fbLogin, setSettings, gGetEmail*/

"use strict";

function populateRentRoll(email) {
    if (email === '') {
        console.log("populateRentRoll - empty email!");
        return;
    }
    $.ajax({
        url: '/tenants',
        data: { 'email': email },
        success: function (tenants) {
            document.getElementById('tenants').innerHTML = tenants;
            addTenantRow();
            var dbName = document.getElementById("DbName");
            if (dbName !== null) {
                dbName.value = email;
            }
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
}

function rentRollTemplateNotLoggedIn() { }
function gRentRollTemplate(resp) { }
function fbRentRollTemplate() { }


function removeTenant(tenantId) {
    var dbName = $('#DbName').val();
    tenantAction(dbName, tenantId, 'removetenant');
    var tr = null;
    var idx = null;
    var tbl = document.getElementById('rentroll-table');
    for (var i = 0; i < tbl.rows.length; i = i + 1) {
        var row = tbl.rows[i];
        if (row.id === 'tr-' + tenantId) {
            idx = i;
            tr = row;
            break;
        }
    }
    removedTr = tr;
    tr.hidden = true;
    var undo = document.getElementById('undo');
    if (undo === null) {
        logError("removeTenant - no undo element!");
        return;
    } else {
        undo.innerHTML = '<a class="undo" href="javascript:undoRemoveTenant(' + tenantId + ', ' + idx + ')"' + '>Undo</a>';
    }
}

function gSigninForm(resp) {
    window.location.href = '/rentroll';
}

function fbSigninForm() {
    window.location.href = '/rentroll';
}

function undoRemoveTenant(tenantId, trIdx) {
    var dbName = $('#DbName').val();
    tenantAction(dbName, tenantId, 'undoremovetenant');
    var tr = document.getElementById('tr-' + tenantId);
    tr.hidden = false;
    document.getElementById('undo').innerHTML = '';
}


function addTenant() {
    var tr = $("#rentroll-table tr:last");
    var tenant = { };
    var lastTr = $("#rentroll-table tr:nth-last-child(2)");
    var dbName = $('#DbName').val();
    for (var i = 0; i < tr.children().length - 1; i++) {
        if (tr.children()[i].children.length == 0) {
            continue;
        }
        child = tr.children()[i].children[0];
        if (child.tagName !== 'INPUT') {
            continue;
        }
        tenant[child.name] = child.value;
    }
    tenant['DbName'] = dbName;
    $.ajax({
        url: '/addtenant',
        data: tenant,
        success: function (tenantId) {
            tenant['Id'] = tenantId;
            tenant['Total'] = tenantTotalRent(tenant);
            var editTr = tenantEditTRHtml(tenant, true, true);
            var newTrHtml = tenantTRHtml(tenant, true, true);
            lastTr.before(newTrHtml);
            lastTr.before(editTr);
            editTr = document.getElementById('tr-edit-' + tenantId);
            editTr.hidden = true;
        }
    });
}

function printInvoices(month, year) {
    var dbName = $('#DbName').val();
    if (month === undefined || month == null || year === undefined || year === null) {
        console.log("printInvoices - no month/year set");
        logError("printInvoices - no month/year set");
    }
    window.location.href = '/printinvoices?DbName=' + encodeURIComponent(dbName) +
        '&month=' + encodeURIComponent(month) + '&year=' + encodeURIComponent(year);
}
