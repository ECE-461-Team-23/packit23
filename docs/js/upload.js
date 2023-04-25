import { endpointPackage, xAuth } from 'overarching.js';
const packageData = {
    id,
    packagename,
    url,
    version,
    author,
    file,
    encodedFile
}
const formPackageName = document.getElementById("formPackageName");
const formVersionNo = document.getElementById("formVersionNo");
const formZipUpload = document.getElementById("formZipUpload");
const formURL = document.getElementById("formURL");

const errPermsMsg = document.getElementById("errPermsMsg");
const errMsg = document.getElementById("errMsg");

var xhr = new XMLHttpRequest();
// HTML IDs
// formPackageName
// formVersionNo
// formZipUpload
// formURL
// formSubmit

// async function submitPackage(author, encodedFile, name, versionNo, zip, url) {
//     try {
//         const response = await fetch(endpointPackage, {
//             method: 'POST',
//             headers: {
//                 xAuth: 'application/json'
//             },
//             body: JSON.stringify(data)
//         });
//         const json = await response.json();
//         console.log(json);
//     } catch (error) {
//         console.error('Error:', error);
//     }
// }

async function submitPackageByURL(inputUrl) {
    // var data = 
    try {
        const response = await fetch(overarching.endpointPackage, {
            method: 'POST',
            headers: {
                'X-Authorization': token
            },
            body: {
                'URL': inputUrl
            }
        });
        const json = await response.json();
        console.log(json);
    } catch (error) {
        console.error('Error:', error);
    }
}
function checkPackage() {
    // if (authenticate() == false) {
    //     errPermsMsg.style.display = "block";
    // } else 
    // if (formPackageName.value == "" || formVersionNo.value == "" || (formZipUpload.value == "" && formURL.value == "")) {
    //     errMsg.style.display = "block";
    // } else {
    //     var author = getAuthor();
    //     var encodedFile = encodeFile();
    //     submitPackage(author, encodedFile, formPackageName.value, formVersionNo.value, formZipUpload.value, formURL.value);
    // }
    if (formZipUpload.value == "" ) {
        errMsg.style.display = "block";
    } else {
        submitPackageByURL(url); 
    }
}

function encodeFile() {

}

function getAuthor() {

}