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

function populateTenant(email, tenantId) {
    $.ajax({
        url: '/tenantdata',
        dataType: 'json',
        data: { 'DbName': email, 'TenantId': tenantId },
        success: function (tenant) {
            var html = '<table class="rentroll-table">';
            html += '<tr>\
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
    </tr>';
            var total = parseFloat(tenant.BaseRent) +
                parseFloat(tenant.Electricity) +
                parseFloat(tenant.Gas) +
                parseFloat(tenant.Water) +
                parseFloat(tenant.SewageTrashRecycle);
                
            html += '<tr id="tr-' + tenantId + '">';
            html += '<td class="tmplt-td">' + tenant.Name + '</td>';
            html += '<td class="tmplt-td">' + tenant.Address + '</td>';
            html += '<td class="tmplt-td">' + tenant.SqFt + '</td>';
            html += '<td class="tmplt-td">' + tenant.LeaseStartDate + '</td>';
            html += '<td class="tmplt-td">' + tenant.LeaseEndDate + '</td>';
            html += '<td class="tmplt-td">' + formatMoney(parseFloat(tenant.BaseRent)) + '</td>';
            html += '<td class="tmplt-td">' + formatMoney(parseFloat(tenant.Electricity)) + '</td>';
            html += '<td class="tmplt-td">' + formatMoney(parseFloat(tenant.Gas)) + '</td>';
            html += '<td class="tmplt-td">' + formatMoney(parseFloat(tenant.Water)) + '</td>';
            html += '<td class="tmplt-td">' + formatMoney(parseFloat(tenant.SewageTrashRecycle)) + '</td>';
            html += '<td class="tmplt-td">' + formatMoney(total) + '</td>';
            html += '<td class="tmplt-td">' + tenant.Comments + '</td>';
            html += '</tr>';
            html += "</table>";
            html += "<br><br>";
            $("#tenant").append(html);
        }
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
