const getPackagesCall = "https://good-spec-d4rgapcc.uc.gateway.dev/packages";
const filterRegExInput = document.getElementById("filterRegExInput");
const filterIDInput = document.getElementById("filterIDInput");
const filterNameInput = document.getElementById("filterNameInput");

// should be called by body in respective page as an onload
function setupPage() {
    renderTable(getRegistry());
}

// Populate RegistryTable
// This should use /packages endpoint
function renderTable(data) {
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
}

function openPackage(id) {
    // get by package by ID and open 
    var url = 'package.html/package/' + id; // Replace with the URL for your package details page
    window.open(url, '_blank');
}

async function getRegistry() {
    const response = await fetch(getPackagesCall, {
        method: 'POST',
        headers: {
            'X-Authorization': "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODI2NDk0MDUsIm5iZiI6MTY4MjQ3NjYwNSwiaXNzIjoicGFja2l0MjMiLCJhdWQiOiJwYWNraXQyMyIsImlhdCI6MTY4MjQ3NjYwNSwic3ViIjoxfQ.mo04vigHZ9seVWUYbxNp_P5mMJZRQpeDRrd7gtwtwPg"
        }
    }).catch(error => console.log(error));
    console.log(response);
    console.log('end of submit by URL');
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