const baseURL = "https://good-spec-d4rgapcc.uc.gateway.dev/";
/*
{
    "User": {
        "name": "ece30861defaultadminuser",
        "isAdmin": true
    },
    "Secret": {
        "password": "correcthorsebatterystaple123(!__+@**(A’”`;DROP TABLE packages;"
    }
}
*/
function authenticate(username, password) {
    
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