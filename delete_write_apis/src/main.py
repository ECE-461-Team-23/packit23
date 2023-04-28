from typing import Union
from fastapi import FastAPI
from fastapi.exceptions import RequestValidationError
from fastapi.responses import PlainTextResponse
from fastapi.middleware.cors import CORSMiddleware

from starlette.middleware import Middleware
from starlette_context import context, plugins
from starlette_context.middleware import RawContextMiddleware
from starlette.exceptions import HTTPException as StarletteHTTPException
from contextlib import asynccontextmanager

from . import write, delete, database

origins = [
    "https://growyourgrove.tech"
]

# Add middleware to access requestid using "context.data"
middleware = [
    Middleware(
        RawContextMiddleware,
        plugins=(
            plugins.RequestIdPlugin(),
            plugins.CorrelationIdPlugin()
        )
    )
]

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Setup
    database.setup_database()
    yield
    # Cleanup

app = FastAPI(middleware=middleware, lifespan=lifespan)
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.include_router(delete.router)
app.include_router(write.router)


# Default validation error code is 422, change to 400 for every endpoint
# "There is missing field(s) in the AuthenticationRequest or it is formed improperly"
@app.exception_handler(RequestValidationError)
def validation_exception_handler(request, exc):
    return PlainTextResponse(str(exc), status_code=400)

# Return errors as plain text
@app.exception_handler(StarletteHTTPException)
async def http_exception_handler(request, exc):
    return PlainTextResponse(str(exc.detail), status_code=exc.status_code)

@app.get("/")
def read_root():
    return {"Hello": "World"}
