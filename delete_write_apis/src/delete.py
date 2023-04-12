"""
APIs for DELETE methods:
/reset
/package/(id)
/package/byName/(name)
"""

from fastapi import APIRouter

router = APIRouter()

@router.get("/delete")
def delete_root():
    return {"Hello": "Delete"}

# @router.delete('/reset', response_model=None)
# def registry_reset(
#     x__authorization: AuthenticationToken = Header(..., alias='X-Authorization')
# ) -> None:
#     """
#     Reset the registry
#     """
#     pass

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