from typing import Union
from fastapi import FastAPI
from fastapi.exceptions import RequestValidationError
from fastapi.responses import PlainTextResponse

from . import write, delete

app = FastAPI()

app.include_router(delete.router)
app.include_router(write.router)


# Default validation error code is 422, change to 400 for every endpoint
# "There is missing field(s) in the AuthenticationRequest or it is formed improperly"
@app.exception_handler(RequestValidationError)
def validation_exception_handler(request, exc):
    return PlainTextResponse(str(exc), status_code=400)


@app.get("/")
def read_root():
    return {"Hello": "World"}
