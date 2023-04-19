"""
APIs for POST methods:
/package
/package/(id)
/authenticate

function names will have method specified if already
in use (other apis) 
"""

import json
import os
import requests
import traceback


from typing import List, Optional, Union
from fastapi import APIRouter, Request, HTTPException

from . import authentication, helper, database

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

MINIMUM_ACCEPTABLE_NET_SCORE = 0.1
PACKAGE_RATER_RETRY_MAX = 2

@router.get("/write")
def write_root():
    return {"Hello": "Write"}

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
        userid, _, stored_password = database.get_data_for_user(username)
        assert stored_password != None
        assert authentication.check_password(password, stored_password)
    except Exception:
        print(f"The user or password is invalid: {traceback.print_exc()}")
        raise HTTPException(status_code=401, detail="The user or password is invalid")

    # Generate and return token
    try:
        assert userid != None
        token = authentication.generate_jwt(userid)
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

        # TODO: Is this the following line required?
        # assert "JSProgram" in parsed_body
    except Exception:
        print(f"There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid: {traceback.print_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")

    # Validate URL from request body
    print(f"HTTP Request body: {parsed_body}")
    try:
        packageName, packageVersion, packageUrl = helper.grabPackageDataFromRequest(parsed_body)
        assert packageName != None and packageName != ""
        assert packageVersion != None and packageVersion != ""
        assert packageUrl != None and packageUrl != ""
        # packageUrl = "https://github.com/axios/axios" #TODO: undo
        assert helper.checkGithubUrl(packageUrl)
    except Exception:
        print(f"Unable to get/validate package data from request body: {traceback.print_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
    
    if database.check_if_package_exists(packageName, packageVersion):
        print(f"Package already exists: {packageName}, {packageVersion}")
        raise HTTPException(status_code=409, detail="Package exists already.")   

    # Send url to package_rater container, read response
    # Error if the package has a disqualified rating
    # Check if the package already exists
    try:
        retry = 0
        while retry < PACKAGE_RATER_RETRY_MAX:
            print("Sending request to Package Rater")
            response = requests.post(url=os.environ["PACKAGE_RATER_URL"], data=packageUrl, timeout=90)
            responseBody = response.text
            print(f"Response body from Package Rater: {responseBody}")
            if response.status_code == 500:
                retry += 1
            else:
                break
        assert responseBody != None and responseBody != ""
        rating = json.loads(responseBody)
        netscore = rating["NET_SCORE"]
        if netscore < MINIMUM_ACCEPTABLE_NET_SCORE:
            print(f"Package has a disqualifying rating: {rating}")
            raise HTTPException(status_code=424, detail="Package is not uploaded due to the disqualified rating.")           
    except Exception:
        print(f"Unable to get data from package_rater: {traceback.print_exc()}")
        raise HTTPException(status_code=500, detail="Internal server error")
    
    # Upload package
    print("Uploading package...")
    try:
        # Ingest external packages if required
        if "Content" in parsed_body:
            content = parsed_body["Content"]
        else:
            print("Ingesting Package..")
            content = helper.downloadGithubRepo(packageUrl)

        packageId = database.upload_package(name=packageName, 
                                            version=packageVersion, 
                                            author_pk=userid,
                                            rating=rating,
                                            url=packageUrl,
                                            content=content)           
    except Exception:
        print(f"Unable to upload package: {traceback.print_exc()}")
        raise HTTPException(status_code=500, detail="Internal server error")
    print("Upload complete")

    # Build response
    packageMetadata = PackageMetadata(
        Name=PackageName(__root__=packageName),
        Version=packageVersion,
        ID=PackageID(__root__=packageId),
    )

    if "Content" in parsed_body: 
        packageData = PackageData(Content=parsed_body["Content"])
    elif "URL" in parsed_body:
        packageData = PackageData(URL=parsed_body["URL"])

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


