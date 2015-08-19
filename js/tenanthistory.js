function populateTenantHistory(tenantId) {
    var dbName = $("#DbName").val();
    getTenantHistory(dbName, tenantId, function (history) {
        var html = '<div id="tenant-history" class="tenant-content">';
        html += formatTenantHistory(history);
        html += '</div>';
        $("#tenant-history").replaceWith(html);
    });
}

function formatTenantHistory(history) {
    var html = '<table class="rentroll-table tenant-table">';
    html += '<tr>';
    html += '<th clss="default-th">Date/Time</th>';
    html += '<th clss="default-th">Action</th>';
    html += tableTenantColsHeader(false);
    html += '</tr>';
    var prevTenantValues = null;
    for (var i = 0; i < history.length; i++) {
        html += '<tr>';
        html += '<td class="default-td">';
        html += history[i]['DateTime'];
        html += '</td>'
        html += '<td class="default-td">';
        html += history[i]['Action'];
        html += '</td>';
        if (history[i]['HasValues'] === true) {
            if (prevTenantValues !== null) {
                html += tenantTDHtml(diffTenant(prevTenantValues, history[i]['TenantValues']), false);
            } else {
                html += tenantTDHtml(history[i]['TenantValues'], false);
            }
            prevTenantValues = history[i]['TenantValues'];
        }
        html += '</tr>';
    }
    html += '</table>';
    return html;
}

function diffValue(field, val1, val2) {
    if (val1 !== val2) {
        return '<b>' + formatField(field, val2) + '</b>';
    } 
    return '""';
}

function diffTenant(t1, t2) {
    t1.Total = tenantTotalRent(t1);
    t2.Total = tenantTotalRent(t2);
    var diff = {};
    for (var fieldName in t1) {
        var d = diffValue(fieldName, t1[fieldName], t2[fieldName]);
        if (isMoneyField(fieldName) && d === '""') {
            d = '<span class="small-diff-value">' + formatField(fieldName, t2[fieldName]) + '</span>';
        }
        diff[fieldName] = d;
    }
    return diff;
}



function getTenantHistory(dbName, tenantId, callback) {
    $.ajax({
        url: '/tenanthistory',
        dataType: 'json',
        data: { 'DbName': dbName, 'TenantId': tenantId },
        success: function (history) {
            callback(history);
        }
    });
}
