function removeTenant(dbName, tenantId) {
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
        undo.innerHTML = '<a class="undo" href="javascript:undoRemoveTenant(\'' + dbName + '\', ' + tenantId + ', ' + idx + ')"' + '>Undo</a>';
    }
}


function editTenant(dbName, tenantId) {
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

function undoRemoveTenant(dbName, tenantId, trIdx) {
    tenantAction(dbName, tenantId, 'undoremovetenant');
    var tr = document.getElementById('tr-' + tenantId);
    tr.hidden = false;
    document.getElementById('undo').innerHTML = '';
}

function formatMoney(num, c, d, t) {
    var n = num, 
        c = isNaN(c = Math.abs(c)) ? 2 : c, 
        d = d == undefined ? "." : d, 
        t = t == undefined ? "," : t, 
        s = n < 0 ? "-" : "", 
        i = parseInt(n = Math.abs(+n || 0).toFixed(c)) + "", 
        j = (j = i.length) > 3 ? j % 3 : 0;
    return "$" + s + (j ? i.substr(0, j) + t : "") + i.substr(j).replace(/(\d{3})(?=\d)/g, "$1" + t) + (c ? d + Math.abs(n - i).toFixed(c).slice(2) : "");
};
       
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
            td.innerHTML = "<a href=\"\">edit</a>, <a href=\"javascript:removeTenant('" + $('#DbName').val() + "', " + tenantId + ")\">remove</a>";
            newTr.appendChild(td);
            lastTr.after(newTr);
        }
    });
}

function saveTenant(tenantId) {
    var tr = document.getElementById("tr-" + tenantId);
    var editTR = document.getElementById("tr-edit-" + tenantId);
    var tenant = { };
    var newTR = document.createElement('tr');
    var dbName = $('#DbName').val();
    for (var i = 0; i < editTR.children.length - 1; i++) {
        var td = document.createElement('td');
        td.className = 'tmplt-td';
        if (editTR.children[i].children.length == 0) {
            newTR.appendChild(td);
            continue;
        }
        var child = editTR.children[i].children[0];
        if (child.tagName !== 'INPUT') {
            newTR.appendChild(td);
            continue;
        }
        tenant[child.name] = child.value;
        var value = child.value;
        
        if (child.getAttribute('rent') !== null && child.getAttribute('rent') === 'true') {
            value = value.replace("$", "").replace(",", "");
            tenant[child.name] = value;
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
                td.innerHTML = "<a href=\"\">edit</a>, <a href=\"javascript:removeTenant('" + $('#DbName').val() + "', " + tenantId + ")\">remove</a>";
                newTR.appendChild(td);
                editTR.hidden = true;
                tr.parentNode.replaceChild(newTR, tr);
                newTR.id = 'tr-' + tenantId;
            } else {
                logError("Error updating tenant");
            }
        }
    });

}
