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

function tenantAction(dbName, tenantId, action){
    $.ajax({
        url: '/' + action,
        data: { 'DbName': dbName, 'TenantId': tenantId },
        success: function (suc) {
            var tenantAction = document.getElementById('Tenant-' + tenantId);
            if (suc !== 'true') {
                logError('Error: tenantAction, action:' + action);
            }
        }
    });
}

function rentRollTableHeader() {
    return '<tr>\
        <th class="th-tenant">Name</th>\
        <th class="th-address">Address</th>\
        <th class="tmplt-th">Sq. Feet</th>\
        <th class="tmplt-th">Lease Start<br>(Date)</th>\
        <th class="tmplt-th">Lease End<br>(Date)</th>\
        <th class="tmplt-th">Base Rent</th>\
        <th class="tmplt-th">Electricity</th>\
        <th class="tmplt-th">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Gas&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</th>\
        <th class="tmplt-th">&nbsp;&nbsp;&nbsp;Water&nbsp;&nbsp;&nbsp;</th>\
        <th class="tmplt-th">Sewage/Trash/<br>Recyle</th>\
        <th class="tmplt-th">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Total&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</th>\
        <th class="th-comments">Comments</th>\
        <th class="th-th">Action</th>\
    </tr>'
}

function populateTenant(email, tenantId) {
    var html = '<table class="rentroll-table">';
    html += rentRollTableHeader();
    var dbName = document.getElementById("DbName");
    if (dbName !== null) {
        dbName.value = email;
    }
    getTenantTR(email, tenantId, function (tenant, TRhtml) {
        html += TRhtml;
        html += "</table>";
        html += "<br><br>";
        $("#tenant").append(html);
    });
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
function tenantActionsHtml(tenantId) {
    return "<a href=\"javascript:editTenant(" + tenantId + ")\">edit</a>, <a href=\"javascript:removeTenant(" + tenantId + ")\">remove</a>";
}

function tenantTRHtml(tenant) {
    var total = parseFloat(tenant.BaseRent) +
            parseFloat(tenant.Electricity) +
            parseFloat(tenant.Gas) +
            parseFloat(tenant.Water) +
            parseFloat(tenant.SewageTrashRecycle);
    var TRhtml = '<tr id="tr-' + tenant.Id + '">'
    TRhtml += '<td class="tmplt-td">' + tenant.Name + '</td>';
    TRhtml += '<td class="tmplt-td">' + tenant.Address + '</td>';
    TRhtml += '<td class="tmplt-td">' + tenant.SqFt + '</td>';
    TRhtml += '<td class="tmplt-td">' + tenant.LeaseStartDate + '</td>';
    TRhtml += '<td class="tmplt-td">' + tenant.LeaseEndDate + '</td>';
    TRhtml += '<td class="tmplt-td">' + formatMoney(parseFloat(tenant.BaseRent)) + '</td>';
    TRhtml += '<td class="tmplt-td">' + formatMoney(parseFloat(tenant.Electricity)) + '</td>';
    TRhtml += '<td class="tmplt-td">' + formatMoney(parseFloat(tenant.Gas)) + '</td>';
    TRhtml += '<td class="tmplt-td">' + formatMoney(parseFloat(tenant.Water)) + '</td>';
    TRhtml += '<td class="tmplt-td">' + formatMoney(parseFloat(tenant.SewageTrashRecycle)) + '</td>';
    TRhtml += '<td class="tmplt-td">' + formatMoney(total) + '</td>';
    TRhtml += '<td class="tmplt-td">' + tenant.Comments + '</td>';
    TRhtml += '<td class="tmplt-td">' + tenantActionsHtml(tenant.Id) + '</td>';
    TRhtml += '</tr>';
    return TRhtml;
}

function editTenant(tenantId) {
    var tr = document.getElementById("tr-" + tenantId);
    var editTR = document.getElementById("tr-edit-" + tenantId);
    if (tr === null || editTR === null) {
        logError("editTenant - no TR!");
        return;
    }
    tr.hidden = true;
    editTR.hidden = false;
}

var oldTR = null;
function saveTenant(tenantId) {
    var tr = document.getElementById("tr-" + tenantId);
    var editTR = document.getElementById("tr-edit-" + tenantId);
    var tenant = { };
    var newTR = document.createElement('tr');
    var dbName = $('#DbName').val();
    for (var i = 0; i < editTR.children.length - 1; i++) {
        var td = document.createElement('td');
        var child = null;
        td.className = 'tmplt-td';
        if (editTR.children[i].children.length == 0) {
            child = editTR.children[i];
            if (child.getAttribute('total') === 'true') {
                value = parseFloat(tenant['BaseRent']) +
                    parseFloat(tenant['Electricity']) +
                    parseFloat(tenant['Gas']) +
                    parseFloat(tenant['Water']) +
                    parseFloat(tenant['SewageTrashRecycle']);
                tenant[child.name] = value;
                value = formatMoney(value);
                td.textContent = value;
            }
            newTR.appendChild(td);
            continue;
        }
        child = editTR.children[i].children[0];
        if (child.tagName !== 'INPUT') {
            newTR.appendChild(td);
            continue;
        }
        tenant[child.name] = child.value;
        var value = child.value;
        if (child.getAttribute('rent') === 'true') {
            value = value.replace("$", "").replace(",", "");
            tenant[child.name] = value;
            value = formatMoney(value);
        }
        td.textContent = value;
        newTR.appendChild(td);
    }
    tenant['DbName'] = dbName;
    tenant['TenantId'] = tenantId;
    $.ajax({
        url: '/updatetenant',
        data: tenant,
        success: function (success) {
            if (success) {
                var td = document.createElement('td');
                td.className = 'tmplt-td';
                td.innerHTML = tenantActionsHtml(tenantId);
                newTR.appendChild(td);
                editTR.hidden = true;
                oldTR = tr.parentNode.replaceChild(newTR, tr);
                newTR.id = 'tr-' + tenantId;
                var undo = document.getElementById('undo');
                if (undo === null) {
                    logError("saveTenant - no undo element!");
                    return;
                } else {
                    undo.innerHTML = '<a class="undo" href="javascript:undoSaveTenant(' + tenantId + ')"' + '>Undo</a>';
                }
            } else {
                logError("Error updating tenant");
            }
        }
    });

}

function undoSaveTenant(tenantId) {
    var dbName = $('#DbName').val();
    tenantAction(dbName, tenantId, 'undoupdatetenant');
    var tr = document.getElementById('tr-' + tenantId);
    tr.parentNode.replaceChild(oldTR, tr);
    oldTR.hidden = false;
    document.getElementById('undo').innerHTML = '';
}

function tenantEditTRHtml(tenant) {
    var total = parseFloat(tenant.BaseRent) +
            parseFloat(tenant.Electricity) +
            parseFloat(tenant.Gas) +
            parseFloat(tenant.Water) +
            parseFloat(tenant.SewageTrashRecycle);
    var actionHtml = '<a href="javascript:saveTenant(' + tenant.Id + ')">save</a>';
    return "<tr hidden=\"true\" id=\"tr-edit-" + tenant.Id +"\">\
               <td class=\"tmplt-td\">\
               <input size=\"23\" type=\"text\" title=\"Name\" name=\"Name\" value=\"" + tenant.Name + "\" />\
               </td>\
               <td class=\"tmplt-td\">\
               <input size=\"27\" type=\"text\" title=\"Address\" name=\"Address\" value=\"" + tenant.Address + "\" />\
               </td>\
               <td class=\"tmplt-td\">\
               <input size=\"5\" type=\"text\" title=\"SqFt\" name=\"SqFt\" value=\"" + tenant.SqFt + "\" />\
               </td>\
               <td class=\"tmplt-td\">\
               <input size=\"10\" type=\"text\" title=\"LeaseStartDate\" name=\"LeaseStartDate\" value=\"" + tenant.LeaseStartDate + "\" />\
               </td>\
               <td class=\"tmplt-td\">\
               <input size=\"10\" type=\"text\" title=\"LeaseEndDate\" name=\"LeaseEndDate\" value=\"" + tenant.LeaseEndDate + "\"/></td>\
               <td class=\"tmplt-td\">\
               <input size=\"5\" type=\"text\" title=\"BaseRent\" rent=\"true\" name=\"BaseRent\" value=\"" + tenant.BaseRent + "\"/>\
               </td>\
               <td class=\"tmplt-td\">\
               <input size=\"5\" type=\"text\" title=\"Electricity\" rent=\"true\" name=\"Electricity\" value=\"" + tenant.Electricity + "\"/>\
               </td>\
               <td class=\"tmplt-td\">\
               <input size=\"5\" type=\"text\" title=\"Gas\" rent=\"true\" name=\"Gas\" value=\"" + tenant.Gas + "\"/>\
               </td>\
               <td class=\"tmplt-td\">\
               <input size=\"5\" type=\"text\" title=\"Water\" rent=\"true\" name=\"Water\" value=\"" + tenant.Water + "\"/>\
               </td>\
               <td class=\"tmplt-td\">\
               <input size=\"5\" type=\"text\" title=\"SewageTrashRecycle\" rent=\"true\" name=\"SewageTrashRecycle\" value=\"" + tenant.SewageTrashRecycle + "\"/>\
               </td>\
               <td class=\"tmplt-td\" total=\"true\">" + formatMoney(total) + "</td>\
               <td class=\"tmplt-td\">\
               <input size=\"14\" type=\"text\" title=\"Comments\" name=\"Comments\" value=\"" + tenant.Comments + "\"/>\
               </td>\
               <td class=\"tmplt-td\">\
               " + actionHtml + "\
               </td>\
               </tr>\
               ";
}

function getTenantTR(dbName, tenantId, callback) {
    getTenant(dbName, tenantId, function (tenant) {
        var  TRhtml = tenantTRHtml(tenant);
        TRhtml += tenantEditTRHtml(tenant);
        callback(tenant, TRhtml);
    });
}

function listTenants(email) {
    $.ajax({
        url: '/tenantsdata',
        dataType: 'json',
        data: { 'DbName': email },
        success: function (tenants) {
            var listHtml = "<ul>";
            for (var tenantId in tenants) {
                var tenant = tenants[tenantId];
                var tenantHtml = '<a class="tenantlist" href="/tenant?id=' + tenantId + '">' +
                    tenant.Name + ' (#' + tenantId + ')</a>';
                
                listHtml = listHtml + '<li>' + tenantHtml + '</li>';
            }
            listHtml = listHtml + "</ul>";
            $("#tenant-list").append(listHtml);
        }
    });
}
