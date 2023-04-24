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
from fastapi import APIRouter, Request, HTTPException, Response, status
from fastapi.responses import PlainTextResponse

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

MINIMUM_ACCEPTABLE_NET_SCORE = 0.1 #TODO: figure out requirement for this
PACKAGE_RATER_RETRY_MAX = 2

@router.get("/write")
async def write_root(request: Request):
    await helper.log_request(request)
    return {"Hello": "Write"}

@router.put('/authenticate')
async def create_auth_token(request: Request):
    await helper.log_request(request)
    # Parsing to make sure valid request (Need to manually decode request to allow unicode characters)
    try:
        payload = await request.body()
        payloadDecoded = payload.decode("UTF-8")
        helper.log("payloadDecoded: ", payloadDecoded)
        parsed_body = json.loads(payloadDecoded, strict=False)
        helper.log("parsed_body: ", parsed_body)
        username = parsed_body["User"]["name"]
        password = parsed_body["Secret"]["password"]

        # get_hashed_password is only ran by an administrator to add users by hand
        # helper.log("New password:")
        # helper.log(authentication.get_hashed_password(password)) 
    except Exception:
        helper.log(f"Unable to get/parse request body: {traceback.format_exc()}")
        raise HTTPException(status_code=400, detail="Unable to get/parse request body")

    # See if username and password are valid
    try:
        userid, _, stored_password = database.get_data_for_user(username)
        helper.log("userid, stored_password:", userid, stored_password)
        assert stored_password != None
        assert authentication.check_password(password, stored_password)
    except Exception:
        helper.log(f"The user or password is invalid: {traceback.format_exc()}")
        raise HTTPException(status_code=401, detail="The user or password is invalid")

    # Generate and return token
    try:
        assert userid != None
        token = authentication.generate_jwt(userid)
        helper.log("token: ", token)
        assert token != None
        return PlainTextResponse(token, status_code=200)
    except Exception:
        helper.log(f"Error when trying to generate token: {traceback.format_exc()}")
        raise HTTPException(status_code=501, detail="Internal server error")



@router.post('/package', response_model=None, status_code=201)
async def package_create(request: Request) -> Union[None, Package]:
    await helper.log_request(request)
    # Parse request
    try:
        token = request.headers["X-Authorization"]
        userid = authentication.validate_jwt(token)
        assert userid != None

        payload = await request.body()
        payloadDecoded = payload.decode("UTF-8")
        helper.log("payloadDecoded: ", payloadDecoded)
        parsed_body = json.loads(payloadDecoded, strict=False)
        helper.log("parsed_body: ", parsed_body)

        # On package upload, either Content or URL should be set.
        assert ("Content" in parsed_body) or ("URL" in parsed_body) # At least one should be set
        assert not ( ("Content" in parsed_body) and ("URL" in parsed_body) ) # Both shouldn't be set
    except Exception:
        helper.log(f"There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid: {traceback.format_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid.")

    # Validate URL from request body
    helper.log(f"HTTP Request body: {parsed_body}")
    try:
        packageName, packageVersion, packageUrl = helper.grabPackageDataFromRequest(parsed_body)
        helper.log("packageName, packageVersion, packageUrl: ", packageName, packageVersion, packageUrl)
        assert packageName != None and packageName != ""
        assert packageVersion != None and packageVersion != ""
        assert packageUrl != None and packageUrl != ""
        assert helper.checkGithubUrl(packageUrl)
    except Exception:
        helper.log(f"Unable to get/validate package data from request body: {traceback.format_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid")
    
    # Check if the package already exists
    if database.get_package_id(packageName, packageVersion) != None:
        helper.log(f"Package already exists: {packageName}, {packageVersion}")
        raise HTTPException(status_code=409, detail="Package exists already.")   

    # Send url to package_rater container, read response
    # Error if the package has a disqualified rating
    try:
        retry = 0
        while retry < PACKAGE_RATER_RETRY_MAX:
            helper.log("Sending request to Package Rater")
            response = requests.post(url=os.environ["PACKAGE_RATER_URL"], data=packageUrl, timeout=90)
            responseBody = response.text
            helper.log(f"Response body from Package Rater: {responseBody}")
            if response.status_code == 500:
                retry += 1
            else:
                break
        assert responseBody != None and responseBody != ""
        rating = json.loads(responseBody)
        helper.log("rating: ", rating)
        netscore = rating["NET_SCORE"]
        if netscore < MINIMUM_ACCEPTABLE_NET_SCORE:
            helper.log(f"Package has a disqualifying rating: {rating}")
            raise HTTPException(status_code=424, detail="Package is not uploaded due to the disqualified rating.")           
    except Exception:
        helper.log(f"Unable to get data from package_rater: {traceback.format_exc()}")
        raise HTTPException(status_code=500, detail="Internal server error")
    
    # Upload package
    helper.log("Uploading package...")
    try:
        # Ingest external packages if required
        if "Content" in parsed_body:
            content = parsed_body["Content"]
            isExternal = False
        else:
            helper.log("Ingesting Package..")
            content = helper.downloadGithubRepo(packageUrl)
            isExternal = True

        packageId = database.upload_package(name=packageName, 
                                            version=packageVersion, 
                                            author_pk=userid,
                                            rating=rating,
                                            url=packageUrl,
                                            content=content,
                                            isExternal=isExternal)           
    except Exception:
        helper.log(f"Unable to upload package: {traceback.format_exc()}")
        raise HTTPException(status_code=500, detail="Internal server error")
    helper.log("Upload complete")

    # Build response
    packageMetadata = PackageMetadata(
        Name=PackageName(__root__=packageName),
        Version=packageVersion,
        ID=PackageID(__root__=packageId),
    )
    packageData = PackageData(Content=content)

    return Package(metadata=packageMetadata, data=packageData)



@router.put('/package/{id}', response_model=None)
async def package_update(id: str, request: Request) -> Union[None, Package]:
    """
    Update this content of the package.
    """
    await helper.log_request(request)
    # Parse request
    try:
        token = request.headers["X-Authorization"]
        userid = authentication.validate_jwt(token)
        assert userid != None

        payload = await request.body()
        helper.log("payloadDecoded: ", payloadDecoded)
        payloadDecoded = payload.decode("UTF-8")
        helper.log("parsed_body: ", parsed_body)
        parsed_body = json.loads(payloadDecoded, strict=False)

        # On package upload, either Content or URL should be set.
        assert ("Content" in parsed_body["data"]) or ("URL" in parsed_body["data"]) # At least one should be set
        assert not ( ("Content" in parsed_body["data"]) and ("URL" in parsed_body["data"]) ) # Both shouldn't be set
    except Exception:
        helper.log(f"There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid: {traceback.format_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")

    # Validate URL from request body
    helper.log(f"HTTP Request body: {parsed_body}")
    try:
        packageName, packageVersion, packageId = parsed_body["metadata"]["Name"], parsed_body["metadata"]["Version"], parsed_body["metadata"]["ID"]
        helper.log("packageName, packageVersion, packageUrl: ", packageName, packageVersion, packageUrl)
        packageId = int(packageId)
        _, _, packageUrl = helper.grabPackageDataFromRequest(parsed_body["data"]) # Ignore data from package.json, use metadata given
        helper.log("packageId, packageUrl", packageId, packageUrl)
        assert packageName != None and packageName != ""
        assert packageVersion != None and packageVersion != ""
        assert packageUrl != None and packageUrl != ""
        assert packageId != None
        assert int(id) == packageId
        assert helper.checkGithubUrl(packageUrl)
    except Exception:
        helper.log(f"Unable to get/validate package data from request body: {traceback.format_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
    
    # Check that name, version, and id match
    try:
        db_id = database.get_package_id(packageName, packageVersion)
        helper.log("db_id", db_id)
        assert db_id == packageId
    except Exception:
        helper.log(f"Unable to find match in database: {traceback.format_exc()}")
        raise HTTPException(status_code=404, detail="Package does not exist.")

    # Send url to package_rater container, read response
    # Check if the package already exists
    try:
        retry = 0
        while retry < PACKAGE_RATER_RETRY_MAX:
            helper.log("Sending request to Package Rater")
            response = requests.post(url=os.environ["PACKAGE_RATER_URL"], data=packageUrl, timeout=90)
            responseBody = response.text
            helper.log(f"Response body from Package Rater: {responseBody}")
            if response.status_code == 500:
                retry += 1
            else:
                break
        assert responseBody != None and responseBody != ""
        rating = json.loads(responseBody)   
    except Exception:
        helper.log(f"Unable to get data from package_rater: {traceback.format_exc()}")
        raise HTTPException(status_code=500, detail="Internal server error")

    # Update package
    helper.log("Update package...")
    try:
        # Ingest external packages if required
        if "Content" in parsed_body:
            content = parsed_body["Content"]
            isExternal = False
        else:
            helper.log("Ingesting Package..")
            content = helper.downloadGithubRepo(packageUrl)
            isExternal = True

        database.update_package(id=packageId, 
                                author_pk=userid,
                                rating=rating,
                                url=packageUrl,
                                content=content,
                                isExternal=isExternal)           
    except Exception:
        helper.log(f"Unable to update package: {traceback.format_exc()}")
        raise HTTPException(status_code=500, detail="Internal server error")
    helper.log("Update complete")

    # Return response
    return Response(status_code=status.HTTP_200_OK)
