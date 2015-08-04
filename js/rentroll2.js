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

function printInvoices() {
    var dbName = $('#DbName').val();
}
