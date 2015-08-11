function tableTenantColsHeader(includeAction) {
    var action = '<th class="th-th">Action</th>';
    if (includeAction === false) {
        action = '';
    }
    return '<th class="th-tenant">Name</th>\
        <th class="th-address">Address</th>\
        <th class="default-th">Sq. Feet</th>\
        <th class="default-th">Lease Start<br>(Date)</th>\
        <th class="default-th">Lease End<br>(Date)</th>\
        <th class="default-th">Base Rent</th>\
        <th class="default-th">Electricity</th>\
        <th class="default-th">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Gas&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</th>\
        <th class="default-th">&nbsp;&nbsp;&nbsp;Water&nbsp;&nbsp;&nbsp;</th>\
        <th class="default-th">Sewage/Trash/<br>Recyle</th>\
        <th class="default-th">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Total&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</th>\
        <th class="th-comments">Comments</th>' +
        action;
}

function tenantTotalRent(tenant) {
    var total = parseFloat(tenant.BaseRent) +
        parseFloat(tenant.Electricity) +
        parseFloat(tenant.Gas) +
        parseFloat(tenant.Water) +
        parseFloat(tenant.SewageTrashRecycle);
    return total;
}

function isMoneyField(name) {
    return tenantRentFields().indexOf(name) !== -1;
}

function tenantRentFields() {
    return ['BaseRent', 'Electricity', 'Gas', 'Water', 'SewageTrashRecycle', 'Total'];
}

function tenantFields() {
    var fields = ['Id', 'Name', 'Address', 'SqFt', 'LeaseStartDate', 'LeaseEndDate'];
    fields = fields.concat(tenantRentFields());
    fields = fields.concat(['Comments']);
    return fields;
}
    

function formatField(name, value) {
    if (isMoneyField(name)) {
        return formatMoney(parseFloat(value));
    } else {
        return value;
    }
}

function formatTenant(tenant) {
    var total = tenantTotalRent(tenant);    
    var t = {};
    var fields = tenantFields();
    tenant['Total'] = total;
    for (var i = 0; i < tenantFields().length; i += 1) {
        t[fields[i]] = formatField(fields[i], tenant[fields[i]]);
    }
    return t;
}


function tenantTDHtml(tenant, includeAction, includeRemove, includeLink) {
    var tdHtml = '';
    var name = tenant.Name;
    if (includeLink === true) {
        name = '<a class="tenant" href="/tenant?id=' + tenant.Id + '">' + tenant.Name + ' (#' + tenant.Id + ')</a>';
    }
    tdHtml += '<td class="default-td">' + name + '</td>';
    tdHtml += '<td class="default-td">' + tenant.Address + '</td>';
    tdHtml += '<td class="default-td">' + tenant.SqFt + '</td>';
    tdHtml += '<td class="default-td">' + tenant.LeaseStartDate + '</td>';
    tdHtml += '<td class="default-td">' + tenant.LeaseEndDate + '</td>';
    tdHtml += '<td class="default-td">' + tenant.BaseRent + '</td>';
    tdHtml += '<td class="default-td">' + tenant.Electricity + '</td>';
    tdHtml += '<td class="default-td">' + tenant.Gas + '</td>';
    tdHtml += '<td class="default-td">' + tenant.Water + '</td>';
    tdHtml += '<td class="default-td">' + tenant.SewageTrashRecycle + '</td>';
    tdHtml += '<td class="default-td">' + tenant.Total + '</td>';
    tdHtml += '<td class="default-td">' + tenant.Comments + '</td>';
    if (includeAction !== false) {
        tdHtml += '<td class="default-td">' + tenantActionsHtml(tenant.Id, includeRemove) + '</td>';
    }
    return tdHtml;    
}

function tenantActionsHtml(tenantId, includeRemove) {
    var remove = ", <a href=\"javascript:removeTenant(" + tenantId + ")\">remove</a>";
    if (includeRemove === false) {
        remove = '';
    }
    return "<a href=\"javascript:editTenant(" + tenantId + ")\">edit</a>" + remove;
}

function tenantSaveHtml(tenantId, canRemove, includeLink) {
    var remove = 'false';
    var link = 'false';
    if (canRemove === true) {
        remove = 'true';
    }
    if (includeLink === true) {
        link = 'true';
    }
    return '<a href="javascript:saveTenant(' + tenantId + ', ' + remove + ', ' + link + ')">save</a>';
}

function tenantTRHtml(tenant, includeRemove, includeLink) {
    return '<tr id="tr-' + tenant.Id + '">' + tenantTDHtml(formatTenant(tenant), true, includeRemove, includeLink) + '</tr>';
}

function rentRollTableHeader() {
    return '<tr>' + tableTenantColsHeader() + '</tr>';
}

function editSizes() {
    return {'Name': '23',
            'Address': '27',
            'SqFt': '5',
            'LeaseStartDate': '10',
            'LeaseEndDate': '10',
            'BaseRent': '5',
            'Electricity': '5',
            'Gas': '5',
            'Water': '5',
            'SewageTrashRecycle': '5',
            'Comments': '14',
    };
}

function tenantEditTRHtml(tenant, canRemove, includeLink) {
    var total = tenantTotalRent(tenant);
    var actionHtml = tenantSaveHtml(tenant.Id, canRemove, includeLink);
    var html = "<tr hidden=\"true\" id=\"tr-edit-" + tenant.Id +"\">";
    var fields = null;
    var fields = tenantFields();
    var size = editSizes();
    for (var i = 0; i < fields.length; i++) {
        field = fields[i];
        if (field === 'Id') {
            continue;
        }
        if (field === 'Total') {
            html += "<td class=\"default-td\" total=\"true\">" + formatMoney(total) + "</td>";
            continue;
        }
        
        html += "<td class=\"default-td\">";
        html += "<input size=\"" + size[field] + "\" type=\"text\" title=\"" + field + "\" name=\"" + field + "\" value=\"" + tenant[field] + "\" />";
        html += "</td>";
    }    
    html += "<td class=\"default-td\">\
               " + actionHtml + "\
               </td>\
               </tr>\
               ";
    return html;
}
