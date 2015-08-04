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
function tenantActionsHtml(tenantId) {
    return "<a href=\"javascript:editTenant(" + tenantId + ")\">edit</a>, <a href=\"javascript:removeTenant(" + tenantId + ")\">remove</a>";
}

function addTenant() {
    var tr = $("#rentroll-table tr:last");
    var tenant = { };
    var lastTr = $("#rentroll-table tr:nth-last-child(2)");
    var newTr = document.createElement('tr');
    var dbName = $('#DbName').val();
    for (var i = 0; i < tr.children().length - 1; i++) {
        var td = document.createElement('td');
        td.className = 'tmplt-td';
        if (tr.children()[i].children.length == 0) {
            newTr.appendChild(td);
            continue;
        }
        child = tr.children()[i].children[0];
        if (child.tagName !== 'INPUT') {
            newTr.appendChild(td);
            continue;
        }
        tenant[child.name] = child.value;
        var value = child.value;
        if (value === parseFloat(value).toFixed(2)) {
            value = parseFloat(value).toFixed(2);
            value = formatMoney(value);
        }
        td.textContent = value;
        newTr.appendChild(td);
    }
    tenant['DbName'] = dbName;
    $.ajax({
        url: '/addtenant',
        data: tenant,
        success: function (tenantId) {
            newTr.id = 'tr-' + tenantId;
            var td = document.createElement('td');
            td.className = 'tmplt-td';
            td.innerHTML = tenantActionsHtml(tenantId);
            newTr.appendChild(td);
            lastTr.before(newTr);
            var editTr = editTRHtml('tr-edit-' + tenantId, 'edit', tenantId, tenant);
            lastTr.before(editTr);
            editTr = document.getElementById('tr-edit-' + tenantId);
            editTr.hidden = true;
        }
    });
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
                undo.innerHTML = '<a class="undo" href="javascript:undoSaveTenant(' + tenantId + ')">Undo</a>';
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

function printInvoices() {
    var dbName = $('#DbName').val();
}
