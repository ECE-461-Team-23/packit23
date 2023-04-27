const uploadAPICall = "https://good-spec-d4rgapcc.uc.gateway.dev/package";
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODI2NDk0MDUsIm5iZiI6MTY4MjQ3NjYwNSwiaXNzIjoicGFja2l0MjMiLCJhdWQiOiJwYWNraXQyMyIsImlhdCI6MTY4MjQ3NjYwNSwic3ViIjoxfQ.mo04vigHZ9seVWUYbxNp_P5mMJZRQpeDRrd7gtwtwPg";

const formPackageName = document.getElementById("formPackageName");
const formVersionNo = document.getElementById("formVersionNo");
const formZipUpload = document.getElementById("formZipUpload");
const formURL = document.getElementById("formURL");

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
            // 'Content-Type': "application/json",
            // 'Accept': "*/*",
            // 'Accept-Encoding': "gzip, deflate, br",
            // 'Connection': "keep-alive",
            'Fuck' : "6969696",
            'X-Authorization': "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODI2NDk0MDUsIm5iZiI6MTY4MjQ3NjYwNSwiaXNzIjoicGFja2l0MjMiLCJhdWQiOiJwYWNraXQyMyIsImlhdCI6MTY4MjQ3NjYwNSwic3ViIjoxfQ.mo04vigHZ9seVWUYbxNp_P5mMJZRQpeDRrd7gtwtwPg"
        },
        body: JSON.stringify({
            // "Content": "",
            'URL': inputUrl
            // ,"JSProgram": ""
        })
    }).catch(error => console.log(error));
    console.log(response);
    console.log('end of submit by URL');
}

function checkURLSubmission() {
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