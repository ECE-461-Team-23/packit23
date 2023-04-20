import bcrypt
import datetime
import hashlib
import jwt
import os
import traceback

from . import helper

JWT_LIFESPAN = datetime.timedelta(minutes=2880)


### JWT ###
def generate_jwt(userid: str):
    # Generate a JWT token for a given userid. Returns None on error
    currentTime = datetime.datetime.now(tz=datetime.timezone.utc)

    payload = {
        "exp": currentTime + JWT_LIFESPAN, # Expiration Time
        "nbf": currentTime, # Not Before Time Claim
        "iss": "packit23", # Issuer Claim
        "aud": "packit23", # Audience Claim
        "iat": currentTime, # Issued At Claim
        "sub": userid, # Subject
        # "jti": uniqueJWTID, # Unique JWT ID
    }

    # Get secret
    key = os.environ["JWT_SECRET"]

    return jwt.encode(payload, key, algorithm="HS256")


def validate_jwt(token: str) -> str:
    # Validates a JWT and returns the given userid. Returns None on error
    
    # Get secret
    try:
        key = os.environ["JWT_SECRET"]
    except Exception:
        helper.log(f"Unable to find environment variable: {traceback.format_exc()}")
        return
    
    # Decode JWT
    try:
        options = {
            "verify_signature": True,
            "require": ["exp", "nbf", "aud", "sub"],
            "verify_aud": True,
            "verify_exp": True,
            "verify_nbf": True,
        }
        decoded = jwt.decode(token, key, audience="packit23", algorithms=["HS256"], options=options)
        return decoded["sub"]
    except jwt.ExpiredSignatureError:
        helper.log("JWT token has expired")
    except Exception:
        helper.log(f"Unable to decode JWT: {traceback.format_exc()}")
    
    return 


### Password Hashing ###
# Still needs work:
# https://stackoverflow.com/questions/9594125/salt-and-hash-a-password-in-python/
# https://stackabuse.com/hashing-passwords-in-python-with-bcrypt/
# https://security.stackexchange.com/questions/39849/does-bcrypt-have-a-maximum-password-length
def get_hashed_password(plain_text_password: str):
    # Hash a password for the first time (Using bcrypt, the salt is saved into the hash itself
    salt = bcrypt.gensalt()
    compressed_password = hashlib.sha256(plain_text_password.encode("UTF-8")).digest()
    return bcrypt.hashpw(compressed_password, salt)


def check_password(plain_text_password: str, hashed_password: bytes):
    # Check hashed password. Using bcrypt, the salt is saved into the hash itself
    compressed_password = hashlib.sha256(plain_text_password.encode("UTF-8")).digest()
    return bcrypt.checkpw(compressed_password, hashed_password)