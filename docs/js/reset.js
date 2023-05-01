const basResetURL = "https://npm-registry-6dvk0w0m.uc.gateway.dev";
const resettoken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODI4MjU4MTgsIm5iZiI6MTY4MjY1MzAxOCwiaXNzIjoicGFja2l0MjMiLCJhdWQiOiJwYWNraXQyMyIsImlhdCI6MTY4MjY1MzAxOCwic3ViIjoxfQ.P1l3kX2lszRNdlVqC3H08NLfvuYmtG77U-yWAIlggjE";
const deletePkgByName = document.getElementById("deletePkgByName");
const deletePkgByID = document.getElementById("deletePkgByID");
const deleteSuccessMsg = document.getElementById("deleteSuccessMsg");
const resetSuccessMsg = document.getElementById("resetSuccessMsg");

async function ResetRegistry() {
    const resetResponse = await fetch(basResetURL, {
        method: 'DELETE',
        headers: {
            'Access-Control-Request-Method': 'DELETE',
            'X-Authorization': resettoken
        },
    }).then(resetResponse => {
        if (!resetResponse.ok) {
            throw new Error('Network response was not ok');
        } else {
            deleteSuccessMsg.style.display = "block";
        }
        console.log(resetResponse)
        console.log('Resource deleted successfully');
    }).catch(error => console.log(error));
    console.log(resetResponse);
}

async function DeletePackageByName() {
    const byNameResponse = await fetch(basResetURL + "/byName/" + deletePkgByName.value, {
        method: 'DELETE',
        headers: {
            'name': deletePkgByName.value,
            'X-Authorization': resettoken,
            'Access-Control-Request-Method': 'DELETE'
        },
    }).then(byNameResponse => {
        if (!byNameResponse.ok) {
            throw new Error('Network response was not ok');
        }
        console.log(byNameResponse)
        console.log('Package by Name deleted successfully');
    }).catch(error => console.log(error));
    console.log(byNameResponse);
}


async function DeletePackageByid() {
    const byNameResponse = await fetch(basResetURL + deletePkgByID.value, {
        method: 'DELETE',
        headers: {
            'id': deletePkgByID.value,
            'X-Authorization': resettoken,
            'Access-Control-Request-Method': 'DELETE'
        },
    }).then(byNameResponse => {
        if (!byNameResponse.ok) {
            throw new Error('Network response was not ok');
        }
        console.log(byNameResponse)
        console.log('Package by ID deleted successfully');
    }).catch(error => console.log(error));
    console.log(byNameResponse);
}
/*********************************************************************************
{
    "User": {
        "name": "ece30861defaultadminuser",
        "isAdmin": true
    },
    "Secret": {
        "password": "correcthorsebatterystaple123(!__+@**(A’”`;DROP TABLE packages;"
    }
}
*********************************************************************************/