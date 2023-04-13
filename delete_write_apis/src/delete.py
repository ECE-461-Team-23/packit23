"""
APIs for DELETE methods:
/reset
/package/(id)
/package/byName/(name)
"""

import json

from fastapi import APIRouter, Request, HTTPException

from . import authentication, helper

router = APIRouter()

@router.get("/delete")
def delete_root():
    return {"Hello": "Delete"}

@router.delete('/reset', response_model=None)
async def registry_reset(request: Request) -> None:
    # Parse request
    try:
        token = request.headers["X-Authorization"]
        userid = authentication.validate_jwt(token)
        assert userid != None

    except Exception:
        print(f"There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid: {traceback.print_exc()}")
        raise HTTPException(status_code=400, detail="There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")

    print(f"reset requested by {userid}")

# @router.delete('/package/byName/{name}', response_model=None)
# def package_by_name_delete(
#     name: PackageName,
#     name: PackageName = ...,
#     x__authorization: AuthenticationToken = Header(..., alias='X-Authorization'),
# ) -> None:
#     """
#     Delete all versions of this package.
#     """
#     pass

# @router.delete('/package/{id}', response_model=None)
# def package_delete(
#     id: PackageID,
#     id: PackageID = ...,
#     x__authorization: AuthenticationToken = Header(..., alias='X-Authorization'),
# ) -> None:
#     """
#     Delete this version of the package.
#     """
#     pass