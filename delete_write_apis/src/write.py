"""
APIs for POST methods:
/package
/package/(id)
/authenticate

function names will have method specified if already
in use (other apis) 
"""

import json
import requests
import traceback


from typing import List, Optional, Union
from fastapi import APIRouter, Request, HTTPException

from . import authentication, helper

from .models import (
    AuthenticationRequest,
    AuthenticationToken,
    EnumerateOffset,
    Error,
    Package,
    PackageData,
    PackageHistoryEntry,
    PackageID,
    PackageMetadata,
    PackageName,
    PackageQuery,
    PackageRating,
)

router = APIRouter()

@router.get("/write")
def write_root():
    return {"Hello": "Write"}

test_db = {
    "sampleuser": b'$2b$12$iUwbqGXqpky94Xul/l0RMe.nIe0HGFha9eQ0M9.82MIeSNvxoXNje'
}

@router.put('/authenticate', response_model=AuthenticationToken)
async def create_auth_token(request: Request) -> AuthenticationToken:
    # Parsing to make sure valid request (Need to manually decode request to allow unicode characters)
    try:
        payload = await request.body()
        payloadDecoded = payload.decode("UTF-8")
        parsed_body = json.loads(payloadDecoded, strict=False)

        username = parsed_body["User"]["name"]
        password = parsed_body["Secret"]["password"]

        # get_hashed_password is only ran by an administrator to add users by hand
        # print("New password:")
        # print(authentication.get_hashed_password(password)) 
    except Exception:
        print(f"Unable to get/parse request body: {traceback.print_exc()}")
        raise HTTPException(status_code=400, detail="Unable to get/parse request body")

    # See if username and password are valid
    try:
        stored_password = test_db[username]
        assert stored_password != None
        assert authentication.check_password(password, stored_password)
    except Exception:
        print(f"The user or password is invalid: {traceback.print_exc()}")
        raise HTTPException(status_code=401, detail="The user or password is invalid")

    # Generate and return token
    try:
        token = authentication.generate_jwt(username)
        assert token != None
        return AuthenticationToken(__root__=token)
    except Exception:
        print(f"Error when trying to generate token: {traceback.print_exc()}")
        raise HTTPException(status_code=501, detail="Internal server error")


@router.post('/package', response_model=None, status_code=201)
async def package_create(request: Request) -> Union[None, Package]:
    # Parse request
    try:
        token = request.headers["X-Authorization"]
        userid = authentication.validate_jwt(token)
        assert userid != None

        payload = await request.body()
        payloadDecoded = payload.decode("UTF-8")
        parsed_body = json.loads(payloadDecoded, strict=False)

        # On package upload, either Content or URL should be set.
        assert ("Content" in parsed_body) or ("URL" in parsed_body) # At least one should be set
        assert not ( ("Content" in parsed_body) and ("URL" in parsed_body) ) # Both shouldn't be set

        # TODO: Is this required?
        # assert "JSProgram" in parsed_body
    except Exception:
        print(f"There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid: {traceback.print_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")


    # Handle Request
    print(parsed_body)
    url = None
    if "Content" in parsed_body:       
        # Package contents. This is the zip file uploaded by the user. (Encoded as text using a Base64 encoding).
        # This will be a zipped version of an npm package's GitHub repository, minus the ".git/" directory." It will, for example, include the "package.json" file that can be used to retrieve the project homepage.
        # See https://docs.npmjs.com/cli/v7/configuring-npm/package-json#homepage.
        try:
            url = helper.grabUrl(parsed_body["Content"])
        except Exception:
            print(f"Unable to get url from Content: {traceback.print_exc()}")
            raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
    elif "URL" in parsed_body:
        # Ingest package from public URL
        url = parsed_body["URL"]

    try:
        assert helper.checkGithubUrl(url)
    except Exception:
        print(f"Unable to validate url ({url}): {traceback.print_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")

    # TODO: Send url to package_rater container, read response
    #           Error if the package has a disqualified rating
    headers = {}
    timeout = "??"
    try:
        response = requests.post(url="API_ENDPOINT", headers=headers, data=url, timeout=timeout)
        rating = response.json()
    except Exception:
        print(f"Unable to get data from package_rater: {traceback.print_exc()}")
        raise HTTPException(status_code=500, detail="Internal server error")
    
    # TODO: Check if package exists, error if it already does
    print(rating)
    # TODO: Upload package

    # Build response
    packageMetadata = PackageMetadata(
        Name=PackageName(__root__="ExampleName"),
        Version="ExampleVersion",
        ID=PackageID(__root__="ExampleName"),
    )
    packageData = PackageData(Content="Example Content")
    # packageData = PackageData(URL="Example URL") # TODO: is this an option?

    return Package(metadata=packageMetadata, data=packageData)



# @router.put('/package/{id}', response_model=None)
# def package_update(
#     id: PackageID,
#     id: PackageID = ...,
#     x__authorization: AuthenticationToken = Header(..., alias='X-Authorization'),
#     body: Package = ...,
# ) -> None:
#     """
#     Update this content of the package.
#     """
#     pass


