const authenticateCall = "https://good-spec-d4rgapcc.uc.gateway.dev/authenticate";
const baseURL = "https://good-spec-d4rgapcc.uc.gateway.dev/";
// const packageData = { 
//     id,
//     packagename,
//     // url,
//     version,
//     // author
// }

// const packageRatings = {
//     ratingid,
//     busfactor,
//     correctness,
//     rampup,
//     license,
//     pinningpractice,
//     pullrequest,
//     netscore
// }

// const packageDownload = {
//     downloadId,
//     downloadLink
// }

// global variables
{
var IDSlotEl = document.getElementById("IDSlot");
var NameSlotEl = document.getElementById("NameSlot");
var URLSlotEl = document.getElementById("URLSlot");
var VersionSlotEl = document.getElementById("VersionSlot");
var AuthorSlotEl = document.getElementById("AuthorSlot");

var DownloadSlotEl = document.getElementById("DownloadSlot");
var DownloadIdSlotEl = document.getElementById("DownloadIdSlotEl");

var RatingIDSlotEl = document.getElementById("RatingIDSlot");
var BusFactorSlotEl = document.getElementById("BusFactorSlot");
var CorrectnessSlotEl = document.getElementById("CorrectnessSlot");
var RampupSlotEl = document.getElementById("RampupSlot");
var LicenseSlotEl = document.getElementById("LicenseSlot");
var PPSlotEl = document.getElementById("PPSlot");
var PRSlotEl = document.getElementById("PRSlot");
var NetScoreSlotEl = document.getElementById("NetScoreSlot");
}

fetch('/package/:id', { // not 100% sure on this endpoint, might need more stuff
    method: 'POST',
    headers: {
    'Content-Type': 'application/json'
    },
    body: JSON.stringify(packageData)
    })
    .then(response => response.json())

// fetch('/users', { // not 100% sure on this endpoint, might need more stuff
//     method: 'POST',
//     headers: {
//     'Content-Type': 'application/json'
//     },
//     body: JSON.stringify(data)
//     })
//     .then(response => response.json())

function setupPage() {
    renderPackageInfo();
    renderRatingInfo();
    authenticate();
}

function renderPackageInfo() {
    // // HTML IDs: IDSlot, NameSlot, URLSlot, VersionSlot, AuthorSlot, DownloadSlot
    // IDSlotEl.value = packageData.id;
    // NameSlotEl.value = packageData.name;
    // URLSlotEl.value = packageData.url;
    // VersionSlotEl.value = packageData.version;
    // AuthorSlotEl.value = packageData.author;
}

function renderDownloadInfo() {
    // // DownloadIdSlotEl.value = packageDownload.downloadId; // wont be used
    // not quite sure how we will populate this
    // DownloadIdSlotEl.value = packageDownload.downloadLink;
}

function renderRatingInfo() {
    // // HTML IDs: RatingIDSlot, BusFactorSlot, CorrectnessSlot, RampupSlot, LicenseSlot, PPSlot, PRSlot, NetScoreSlot
    // RatingIDSlotEl.value = packageRatings.ratingid;
    // BusFactorSlotEl.value = packageRatings.busfactor;
    // CorrectnessSlotEl.value = packageRatings.correctness;
    // RampupSlotEl.value = packageRatings.rampup;
    // LicenseSlotEl.value = packageRatings.license;
    // PPSlotEl.value = packageRatings.pinningpractice;
    // PRSlotEl.value = packageRatings.pullrequest;
    // NetScoreSlotEl.value = packageRatings.netscore;
}

function openEditModal() {
    // // HTML IDs: packageEditModal
    checkPerms();
}

function closeEditModal() {
    // // HTML IDs: packageEditModal
    setupPage();
}

function savePackage() {
    // // HTML IDs: IDSlotModal, NameSlotModal, URLSlotModal, VersionSlotModal, AuthorSlotModal
    // // HTML IDs: packageSaveModal
    // IDSlotEl.value 
    // NameSlotEl.value 
    // URLSlotEl.value 
    // VersionSlotEl.value
    // AuthorSlotEl.value
}

async function deletePackage(id) {
    // HTML IDs: packageSaveModal, packageDeleteModal
    const response = await fetch(baseURL.concat("package/", id), {
        method: 'DELETE',
        headers: {
            'Access-Control-Request-Method': 'DELETE',
            'X-Authorization': "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODI2NDk0MDUsIm5iZiI6MTY4MjQ3NjYwNSwiaXNzIjoicGFja2l0MjMiLCJhdWQiOiJwYWNraXQyMyIsImlhdCI6MTY4MjQ3NjYwNSwic3ViIjoxfQ.mo04vigHZ9seVWUYbxNp_P5mMJZRQpeDRrd7gtwtwPg"
        },
    }).then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        console.log(response)
        console.log('Resource deleted successfully');
    }).catch(error => console.log(error));
    console.log(response);
    console.log("End of Package Deletion");
}

// ----------------
// MODAL FUNCTIONS
// ----------------

// const packageModal = document.getElementById('packageEditModal')
// const packageInput = document.getElementById('myInput')

// myModal.addEventListener('shown.bs.modal', () => {
//     myInput.focus()
// })