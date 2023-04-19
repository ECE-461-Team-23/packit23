// should be called by body in respective page as an onload
function setupPage() {
    renderTable();
}

// Populate RegistryTable
// This should use /packages endpoint
function renderTable() {
    $.ajax({
        'url': "/packages", // the go server URL I think? or /registry.html im not sure
        'method': "POST",
        'contentType': 'application/json'
    }).done( function(data) {
        $('#registryTable').dataTable( {
            "aaData": data,
            "columns": [
                { "data": "ID"},
                { 
                    "data": "Name", 
                    "render": function(data, type, row, meta) {
                        return '<a href="#" onclick="openEntry(' + data + ')">' + data + '</a>';
                    } 
                },
                { "data": "Version" },
            ]
        })
    })
        
    // "<i href=openEntry(" Data.ID") class='bi bi-folder2-open'> </i>"
    
    // button onClick=""
    // onRowClick("my-table-id", function (row){
    //     var value = row.getElementsByTagName("td")[0].innerHTML;
    //     document.getElementById('click-response').innerHTML = value + " clicked!";
    //     console.log("value>>", value);
    // });
}

function deleteEntry() {

}

function editEntry() {

}

function openEntry(id) {
    // get by package by ID and open 
}

function searchRegistry() {

}