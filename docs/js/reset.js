const deleteByIDCall = "https://good-spec-d4rgapcc.uc.gateway.dev/package/";
const resetCall = "https://good-spec-d4rgapcc.uc.gateway.dev/reset";

async function ResetRegistry() {
    const response = await fetch(resetCall, {
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