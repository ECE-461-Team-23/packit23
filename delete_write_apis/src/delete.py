"""
APIs for DELETE methods:
/reset
/package/(id)
/package/byName/(name)
"""

import traceback

from fastapi import APIRouter, Request, HTTPException, Response, status

from . import authentication, helper, database, bucket

router = APIRouter()

@router.get("/delete")
async def delete_root(request: Request):
    await helper.log_request(request)
    return {"Hello": "Delete"}

@router.delete('/reset', response_model=None, status_code=200)
async def registry_reset(request: Request):
    # Parse request
    await helper.log_request(request)
    try:
        token = request.headers["X-Authorization"]
        userid = authentication.validate_jwt(token)
        assert userid != None

    except Exception:
        helper.log(f"There is missing field(s) in the AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid: {traceback.format_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")

    helper.log(f"Reset requested by {userid}")
    try:
        database.reset_database()
        bucket.empty_bucket()
    except Exception:
        helper.log(f"Unable to empty database: {traceback.format_exc()}")
        raise HTTPException(status_code=500, detail="Internal Server Error")
   
    return Response(status_code=status.HTTP_200_OK)
   

@router.delete('/package/byName/{name}', response_model=None)
async def package_by_name_delete(name: str, request: Request):
    """
    Delete all versions of this package.
    """
    await helper.log_request(request)
    # Parse request
    try:
        token = request.headers["X-Authorization"]
        userid = authentication.validate_jwt(token)
        assert userid != None
    except Exception:
        helper.log(f"There is missing field(s) in the PackageName/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid: {traceback.format_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageName/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")

    helper.log(f"Request to delete all versions of package: {name}")
    # Check if package exists
    try:
        pks = database.get_all_versions_of_package(name)
        helper.log(f"Versions found: {pks}")
        assert len(pks) > 0
    except Exception:
        helper.log(f"Error when checking if package exists: {traceback.format_exc()}")
        raise HTTPException(status_code=404, detail="Package does not exist.")

    # Delete all versions
    try:
        for pk in pks:
            binary_pk = database.delete_package(pk)
            bucket.delete_blob(str(binary_pk))
            helper.log(f"Deleted all data for package with id: {pk}")
    except Exception:
        helper.log(f"Error when deleting package versions: {traceback.format_exc()}")
        raise HTTPException(status_code=500, detail="Internal Server Error")

    return Response(status_code=status.HTTP_200_OK)



@router.delete('/package/{id}', response_model=None)
async def package_delete(id: str, request: Request):
    """
    Delete this version of the package.
    """
    await helper.log_request(request)
    # Parse request
    try:
        token = request.headers["X-Authorization"]
        userid = authentication.validate_jwt(token)
        assert userid != None
        packageId = int(id)
    except Exception:
        helper.log(f"There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid: {traceback.format_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")

    helper.log(f"Request to delete package with id: {packageId}")
    # Check if package exists
    try:
        packageExists = database.check_if_package_exists(packageId)
        assert packageExists == True
    except Exception:
        helper.log(f"Error whe checking if package exists: {traceback.format_exc()}")
        raise HTTPException(status_code=404, detail="Package does not exist.")

    # Delete package
    try:
        binary_pk = database.delete_package(packageId)
        bucket.delete_blob(str(binary_pk))
    except Exception:
        helper.log(f"Unable to empty database: {traceback.format_exc()}")
        raise HTTPException(status_code=500, detail="Internal Server Error")
   
    return Response(status_code=status.HTTP_200_OK)
