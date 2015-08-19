/*jslint browser: true*/
/*globals $, FB, gapi, fbLogin, setSettings, gGetEmail*/

"use strict";

function populate(email) {
    var tenantId = null;
    if (email === '') {
        console.log("populateTenant - empty email!");
        logError('populateTenant - empty email!');
        return;
    }
    tenantId = get('id');
    if (tenantId !== '') {
        $("#tenant-heading").show();
        $("#history-heading").show();
        populateTenant(email, tenantId);
    } else {
        listTenants(email);
    }
}

function gTenant(resp) {
    populate(gGetEmail(resp));
}

function fbTenant() {
    FB.api('/me', function (response) {
        populate(response.email);
    });
}

function tenantNotLoggedIn() {
    logError("Error logging in on tenant url");
    alert("Please login!");
}

function tenantAction(dbName, tenantId, action, sucFunc){
    $.ajax({
        url: '/' + action,
        data: { 'DbName': dbName, 'TenantId': tenantId },
        success: function (suc) {
            var tenantAction = document.getElementById('Tenant-' + tenantId);
            if (suc !== 'true') {
                logError('Error: tenantAction, action:' + action);
            } else {
                if (sucFunc !== null && sucFunc !== undefined) {
                    sucFunc(tenantId);
                }
            }
        }
    });
}

function populateTenant(email, tenantId) {
    var dbNameEl = document.getElementById("DbName");
    if (dbNameEl !== null) {
        dbNameEl.value = email;
    }
    var dbName = $("#DbName").val();
    getTenantTR(dbName, tenantId, false);
    populateTenantHistory(tenantId);
}

function getTenant(dbName, tenantId, callback) {
    $.ajax({
        url: '/tenantdata',
        dataType: 'json',
        data: { 'DbName': dbName, 'TenantId': tenantId },
        success: function (tenant) {
            callback(tenant);
        }
    });
}

function getTenantTrAndEditTr(tenantId) {
    var tr = document.getElementById("tr-" + tenantId);
    var editTR = document.getElementById("tr-edit-" + tenantId);
    if (tr === null || editTR === null) {
        logError("editTenant - no TR!");
        return null;
    }
    return {'tr': tr, 'editTr': editTR};
}

function editTenant(tenantId) {
    var trs = getTenantTrAndEditTr(tenantId);
    if (trs === null) {
        logError("editTenant - no TR!");
        return;
    }
    trs['tr'].hidden = true;
    trs['editTr'].hidden = false;
}

function cancelEditTenant(tenantId) {
    var trs = getTenantTrAndEditTr(tenantId);
    if (trs === null) {
        logError("editTenant - no TR!");
        return;
    }
    trs['editTr'].hidden = true;
    trs['tr'].hidden = false;
}

var oldTR = null;
function saveTenant(tenantId, canRemove, includeLink) {
    var tr = document.getElementById("tr-" + tenantId);
    var editTR = document.getElementById("tr-edit-" + tenantId);
    var tenant = { };
    var dbName = $('#DbName').val();
    for (var i = 0; i < editTR.children.length - 1; i++) {
        var child = null;
        if (editTR.children[i].children.length == 0) {
            continue;
        }
        child = editTR.children[i].children[0];
        if (child.tagName !== 'INPUT') {
            continue;
        }
        tenant[child.name] = child.value;
        var value = child.value;
        if (isMoneyField(child.name)) {
            value = value.replace("$", "").replace(",", "");
            tenant[child.name] = value;
        }
    }
    tenant['DbName'] = dbName;
    tenant['TenantId'] = tenantId;
    tenant['Id'] = tenantId;
    tenant['Total'] = tenantTotalRent(tenant);
    $.ajax({
        url: '/updatetenant',
        data: tenant,
        success: function (success) {
            if (success) {
                editTR.hidden = true;
                var newTrHtml = tenantTRHtml(tenant, canRemove, includeLink);
                oldTR = $('#tr-' + tenantId).replaceWith(newTrHtml);
                $('#tr-edit-' + tenantId).replaceWith(tenantEditTRHtml(tenant, canRemove, includeLink));
                var undo = document.getElementById('undo');
                if (undo === null) {
                    logError("saveTenant - no undo element!");
                    return;
                } else {
                    undo.innerHTML = '<a class="undo" href="javascript:undoSaveTenant(' + tenantId + ')"' + '>Undo</a>';
                }
                if (typeof(populateTenantHistory) !== "undefined") {
                    populateTenantHistory(tenantId);
                }
            } else {
                logError("Error updating tenant");
            }
        }
    });

}

function undoSaveTenant(tenantId) {
    var dbName = $('#DbName').val();
    if (typeof(populateTenantHistory) !== 'undefined') {
        tenantAction(dbName, tenantId, 'undoupdatetenant', populateTenantHistory);
    } else {
        tenantAction(dbName, tenantId, 'undoupdatetenant');
    }
    $('#tr-' + tenantId).replaceWith(oldTR);
    oldTR[0].hidden = false;
    document.getElementById('undo').innerHTML = '';
}

function getTenantTR(dbName, tenantId, canRemove) {
    var html = '<table class="rentroll-table tenant-table">';
    html += rentRollTableHeader();
    getTenant(dbName, tenantId, function (tenant) {
        html += tenantTRHtml(tenant, canRemove);
        html += tenantEditTRHtml(tenant);
        html += "</table>";
        $("#tenant").append(html);
    });
}

function listTenants(email) {
    $.ajax({
        url: '/tenantsdata',
        dataType: 'json',
        data: { 'DbName': email },
        success: function (tenants) {
            var dbName = document.getElementById("DbName");
            if (dbName !== null) {
                dbName.value = email;
            }
            var listHtml = "<ul>";
            for (var tenantId in tenants) {
                var tenant = tenants[tenantId];
                var tenantHtml = '<a class="tenantlist" href="/tenant?id=' + tenantId + '">' +
                    formattedTenantName(tenant) + '</a>';
                
                listHtml = listHtml + '<li>' + tenantHtml + '</li>';
            }
            listHtml = listHtml + "</ul>";
            $("#tenant-list").append(listHtml);
        }
    });
}
