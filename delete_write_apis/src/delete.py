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
def delete_root():
    return {"Hello": "Delete"}

@router.delete('/reset', response_model=None, status_code=200)
async def registry_reset(request: Request):
    # Parse request
    try:
        token = request.headers["X-Authorization"]
        userid = authentication.validate_jwt(token)
        assert userid != None

    except Exception:
        print(f"There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid: {traceback.print_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")

    print(f"Reset requested by {userid}")
    # TODO: Check if user is admin
    # raise HTTPException(status_code=401, detail="You do not have permission to reset the registry.")

    try:
        database.reset_database()
        bucket.empty_bucket()
    except Exception:
        print(f"Unable to empty database: {traceback.print_exc()}")
        raise HTTPException(status_code=500, detail="Internal Server Error")
   
    return Response(status_code=status.HTTP_200_OK)
   

@router.delete('/package/byName/{name}', response_model=None)
def package_by_name_delete(name: str, request: Request):
    """
    Delete all versions of this package.
    """
    # Parse request
    try:
        token = request.headers["X-Authorization"]
        userid = authentication.validate_jwt(token)
        assert userid != None
    except Exception:
        print(f"There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid: {traceback.print_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")

    print(f"Request to delete all versions of package: {name}")
    # Check if package exists
    try:
        pks = database.get_all_versions_of_package(name)
        print(f"Versions found: {pks}")
        assert len(pks) > 0
    except Exception:
        print(f"Error whe checking if package exists: {traceback.print_exc()}")
        raise HTTPException(status_code=404, detail="Package does not exist.")

    # Delete all versions
    try:
        for pk in pks:
            binary_pk = database.delete_package(pk)
            bucket.delete_blob(str(binary_pk))
            print(f"Deleted all data for package with id: {pk}")
    except Exception:
        print(f"Error when deleting package versions: {traceback.print_exc()}")
        raise HTTPException(status_code=500, detail="Internal Server Error")

    return Response(status_code=status.HTTP_200_OK)



@router.delete('/package/{id}', response_model=None)
def package_delete(id: str, request: Request):
    """
    Delete this version of the package.
    """
    # Parse request
    try:
        token = request.headers["X-Authorization"]
        userid = authentication.validate_jwt(token)
        assert userid != None
        packageId = int(id)
    except Exception:
        print(f"There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid: {traceback.print_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")

    print(f"Request to delete package with id: {packageId}")
    # Check if package exists
    try:
        packageExists = database.check_if_package_exists(packageId)
        assert packageExists == True
    except Exception:
        print(f"Error whe checking if package exists: {traceback.print_exc()}")
        raise HTTPException(status_code=404, detail="Package does not exist.")

    # Delete package
    try:
        binary_pk = database.delete_package(packageId)
        bucket.delete_blob(str(binary_pk))
    except Exception:
        print(f"Unable to empty database: {traceback.print_exc()}")
        raise HTTPException(status_code=500, detail="Internal Server Error")
   
    return Response(status_code=status.HTTP_200_OK)