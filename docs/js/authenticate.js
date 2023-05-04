const baseURL = "https://npm-registry-6dvk0w0m.uc.gateway.dev/";
const usernameInput = document.getElementById("usernameInput");
const passwordInput = document.getElementById("passwordInput");
const loginSuccess = document.getElementById("loginSuccess");

async function authenticate() {
    const authResponse = await fetch(baseURL + "/authenticate", {
        method: 'PUT',
        // body: JSON.stringify([
        //     {
        //         "User": {
        //             "name": "ece30861defaultadminuser",
        //             "isAdmin": true
        //         },
        //         "Secret": {
        //             "password": "correcthorsebatterystaple123(!__+@**(A’”`;DROP TABLE packages;"
        //         }
        //     }
        // ])
    }).then(authResponse => {
        if (!authResponse.ok) {
            throw new Error('Network response was not ok');
        } else {
            // IF SUCCESSFUL:
            loginSuccess.style.display = "block";
            // settoken
        }
        console.log(authResponse)
        console.log('Resource deleted successfully');
    }).catch(error => console.log(error));
    console.log(authResponse);
}

function setToken(token) {
    document.cookie = "token=" + token + "; path=/";
}

function getToken() {
    const tokenName = "token=";
    const decodedCookie = decodeURIComponent(document.cookie);
    const cookieArray = decodedCookie.split(';');
    for (let i = 0; i < cookieArray.length; i++) {
        let cookie = cookieArray[i];
        while (cookie.charAt(0) === ' ') {
            cookie = cookie.substring(1);
        }
        if (cookie.indexOf(tokenName) === 0) {
            return cookie.substring(tokenName.length, cookie.length);
        }
    }
    return "";
}

function setAdmin(admin) {
    document.cookie = "isAdmin=" + admin + "; path=/";
}

function getAdmin() {
    const adminName = "admin=";
    const decodedCookie = decodeURIComponent(document.cookie);
    const cookieArray = decodedCookie.split(';');
    for (let i = 0; i < cookieArray.length; i++) {
        let cookie = cookieArray[i];
        while (cookie.charAt(0) === ' ') {
            cookie = cookie.substring(1);
        }
        if (cookie.indexOf(adminName) === 0) {
            return cookie.substring(adminName.length, cookie.length);
        }
    }
    return "";
}