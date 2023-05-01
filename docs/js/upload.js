const uploadtoken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODMwODE0MjgsIm5iZiI6MTY4MjkwODYyOCwiaXNzIjoicGFja2l0MjMiLCJhdWQiOiJwYWNraXQyMyIsImlhdCI6MTY4MjkwODYyOCwic3ViIjoxfQ.A1NUmac8P5IfyAuSida8DHfTf-5y6rl5JnQLoF9T_g8";
const uploadAPICall = "https://npm-registry-6dvk0w0m.uc.gateway.dev/package";

const formPackageName = document.getElementById("formPackageName");
const formVersionNo = document.getElementById("formVersionNo");
const formZipUpload = document.getElementById("formZipUpload");
const formURL = document.getElementById("formURL");

const successMsg = document.getElementById("successMsg");
const errPermsMsg = document.getElementById("errPermsMsg");
const errMsg = document.getElementById("errMsg");

// async function submitPackagebyZip(author, encodedFile, name, versionNo, zip, url) {
//     try {
//         const response = await fetch(uploadAPICall, {
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

/*********************************************************************************
{
"Content": "string",
"URL": "string",
"JSProgram": "string"
}
*********************************************************************************/
async function submitPackageByURL(inputUrl) {
    const response = await fetch(uploadAPICall, {
        // mode: 'no-cors',
        method: 'POST',
        headers: {
            // Access-Control-Expose-Headers: 
            // 'Content-Type': "application/json",
            // 'Accept': "*/*",
            // 'Accept-Encoding': "gzip, deflate, br",
            // 'Connection': "keep-alive",
            // 'Fuck' : "6969696",
            'X-Authorization': uploadtoken
        },
        body: JSON.stringify({
            // "Content": "",
            'URL': inputUrl
            // ,"JSProgram": ""
        })
    }).catch(error => console.log(error));
    console.log(response);
    console.log('end of submit by URL');
    successMsg.style.display = "block";
}

function checkURL() {
    if (formURL.value == "") {
        console.log("empty URL upload");
        errURLMsg.style.display = "block";
    } else {
        submitPackageByURL(formURL.value);
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
    // if (formZipUpload.value == "" ) {
    //     errMsg.style.display = "block";
    //     console.log("empty zip upload");
    // } 

    errMsg.style.display = "none";
    successMsg.style.display = "none";
    if (formURL.value != "") {
        submitPackageByURL(formURL.value);
    } else if (formZipUpload.value != "") {
        // submitPackageByZip();
    } else {
        console.log("empty upload");
        errMsg.style.display = "block";
    }

    console.log(formZipUpload.value);
}

function encodeFile() {

}

function getAuthor() {

}