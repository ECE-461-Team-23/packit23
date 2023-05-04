const baseURLCall = "https://npm-registry-6dvk0w0m.uc.gateway.dev/";
const filterRegExInput = document.getElementById("filterRegExInput");
const filterIDInput = document.getElementById("filterIDInput");
const filterNameInput = document.getElementById("filterNameInput");
const registrytoken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODMwODE0MjgsIm5iZiI6MTY4MjkwODYyOCwiaXNzIjoicGFja2l0MjMiLCJhdWQiOiJwYWNraXQyMyIsImlhdCI6MTY4MjkwODYyOCwic3ViIjoxfQ.A1NUmac8P5IfyAuSida8DHfTf-5y6rl5JnQLoF9T_g8";

const bodyData = [
    { Version: "1.6.0-10.0.0", Name: "*" }
];
// should be called by body in respective page as an onload
function setupPage() {
    renderTable(getRegistry());
}

// Populate RegistryTable
// This should use /packages endpoint
function renderTable(data) {
    // [{"ID":"51","Name":"superagent","Version":"8.0.9"},{"ID":"52","Name":"superagent","Version":"8.0.9"}]
    $('#registryTable').dataTable( {
        "aaData": data,
        "columns": [
            { "data": "ID"},
            { 
                "data": "Name", 
                "render": function(data, type, row, meta) {
                    return '<a href="#" onclick="openPackage(' + data.ID + ')">' + data.Name + '</a>';
                } 
            },
            { "data": "Version" },
        ]
    });
    console.log(data);
}

function openPackage(id) {
    // get by package by ID and open 
    var url = 'package.html/package/' + id; // Replace with the URL for your package details page
    window.open(url, '_blank');
}

async function getRegistry() {
    const response = await fetch(baseURLCall + "packages", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-Authorization': registrytoken
        },
        body: JSON.stringify([
            { 'Version': "1.6.0-10.0.0", 'Name': "*" }
        ])
    }).catch(error => console.log(error));
    console.log(response);
    console.log('end of getRegistry');
    return response;
}

async function searchRegistryByRegex() {
    // filterRegExInput.value;
    // Add Fetch blocks here to return reponse. 
    // renderTable(response);
}

async function searchRegistryByID() {
    // filterIDInput.value;
    // renderTable(response);
}

async function searchRegistryByName() {
    // filterName.value;
    // renderTable(response);
}