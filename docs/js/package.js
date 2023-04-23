const packageData = { 
    id,
    packagename,
    url,
    version,
    author
}

const packageRatings = {
    ratingid,
    busfactor,
    correctness,
    rampup,
    license,
    pinningpractice,
    pullrequest,
    netscore
}

const packageDownload = {
    downloadid,
    downloadLink
}

// name: 'John Doe',

fetch('/package/{id}', { // not 100% sure on this endpoint, might need more stuff
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
}

function renderPackageInfo() {
    // HTML IDs: IDSlot, NameSlot, URLSlot, VersionSlot, AuthorSlot, DownloadSlot
}

function renderRatingInfo() {
    // HTML IDs: RatingIDSlot, BusFactorSlot, CorrectnessSlot, RampupSlot, LicenseSlot, PPSlot, PRSlot, NetScoreSlot
}

function checkPerms() {
    
}

function openEditModal() {
    // HTML IDs: packageEditModal
    checkPerms();
}

function closeEditModal() {
    // HTML IDs: packageEditModal
    setupPage();
}

function savePackage() {
    // HTML IDs: IDSlotModal, NameSlotModal, URLSlotModal, VersionSlotModal, AuthorSlotModal
    // HTML IDs: packageSaveModal
}

function deletePackage() {
    // HTML IDs: packageSaveModal
}