import datetime
import json
import os

from google.cloud.sql.connector import Connector, IPTypes
import pymysql

import sqlalchemy
from sqlalchemy import inspect
from sqlalchemy import Table, Column, Integer, String, MetaData, LargeBinary, DateTime, Float, Boolean

from . import bucket, helper

# Table representations
metadata = MetaData()
packages = Table(
    'packages', metadata, 
    Column('id', Integer, primary_key = True),
    Column('name', String(512), nullable=False),
    Column('rating_pk', Integer, nullable=False),
    Column('author_pk', Integer, nullable=False),
    Column('url', String(512), nullable=False),
    Column('binary_pk', Integer, nullable=False),
    Column('version', String(512), nullable=False),
    Column('upload_time', DateTime, nullable=False),
    Column('is_external', Boolean, nullable=False),
)

users = Table(
    'users', metadata, 
    Column('id', Integer, primary_key = True),
    Column('username', String(512), nullable=False), 
    Column('password', LargeBinary, nullable=False), 
)

ratings = Table(
    'ratings', metadata, 
    Column('id', Integer, primary_key = True),
    Column('busFactor', Float, nullable=False), 
    Column('correctness', Float, nullable=False), 
    Column('rampUp', Float, nullable=False), 
    Column('responsiveMaintainer', Float, nullable=False), 
    Column('licenseScore', Float, nullable=False), 
    Column('goodPinningPractice', Float, nullable=False), 
    Column('pullRequest', Float, nullable=False), 
    Column('netScore', Float, nullable=False),     
)


def connect_with_connector() -> sqlalchemy.engine.base.Engine:
    """
    Initializes a connection pool for a Cloud SQL instance of MySQL.

    Uses the Cloud SQL Python Connector package.
    """
    # Note: Saving credentials in environment variables is convenient, but not
    # secure - consider a more secure solution such as
    # Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
    # keep secrets safe.

    instance_connection_name = os.environ["INSTANCE_CONNECTION_NAME"]  # e.g. 'project:region:instance'
    db_user = os.environ.get("DB_USER", "")  # e.g. 'my-db-user'
    db_pass = os.environ["DB_PASS"]  # e.g. 'my-db-password'
    db_name = os.environ["DB_NAME"]  # e.g. 'my-database'

    ip_type = IPTypes.PRIVATE if os.environ.get("PRIVATE_IP") else IPTypes.PUBLIC

    connector = Connector(ip_type)

    def getconn() -> pymysql.connections.Connection:
        conn: pymysql.connections.Connection = connector.connect(
            instance_connection_name,
            "pymysql",
            user=db_user,
            password=db_pass,
            db=db_name,
        )
        return conn

    pool = sqlalchemy.create_engine(
        "mysql+pymysql://",
        creator=getconn,
        # ...
    )
    return pool



## External functions
def check_if_default_exists():
    inspector = inspect(engine)
    return "users" in inspector.get_table_names()

def create_default():
    if not check_if_default_exists():
        helper.log("Creating tables")
        metadata.create_all(engine) # Create tables

        helper.log("Creating users")
        logins_str = os.environ["USER_LOGINS"]
        logins_json = json.loads(logins_str)
        logins_list = [{"username": u, "password": bytes(p, 'utf-8')} for u, p in logins_json.items()]

        with engine.begin() as conn:
            result = conn.execute(users.insert(), logins_list)
        helper.log("Finished creating tables & users")
    else:
        helper.log("Tables already detected, not creating")

def get_data_for_user(username: str) -> str:
    # Return (userid, username, password) for a user
    with engine.begin() as conn:
        s = users.select().where(users.c.username==username)
        result = conn.execute(s)
        if result.rowcount > 1:
            raise Exception(f"Error: {result.rowcount} users found for {username}")
        row = result.fetchone()
        helper.log(f"get_password_for_user({username}): found {row}")

        if row:
            return row
        else:
            return None, None, None

def check_if_package_exists(packageId: int):
    # Check if a given package name & version already exists
    with engine.begin() as conn:
        s = packages.select().where(packages.c.id==packageId)
        result = conn.execute(s)
        return result.rowcount > 0

def get_package_id(packageName: str, packageVersion: str):
    # Check if a given package name & version already exists
    with engine.begin() as conn:
        s = packages.select().where(packages.c.name==packageName, packages.c.version==packageVersion)
        row = conn.execute(s).fetchone()
        if row:
            return row.id
        else:
            return None

def get_all_versions_of_package(packageName: str) -> list:
    # Return the pk's of all versions of a package name
    with engine.begin() as conn:
        s = packages.select().where(packages.c.name==packageName)
        return [row.id for row in conn.execute(s)]

def upload_package(name: str, version: str, author_pk: str, rating, url: str, content, isExternal):
    # Upload to ratings table
    helper.log("Uploading to rating table..")
    with engine.begin() as conn:
        ins = ratings.insert().values(
            busFactor=rating["BUS_FACTOR_SCORE"],
            correctness=rating["CORRECTNESS_SCORE"],
            rampUp=rating["RAMP_UP_SCORE"],
            responsiveMaintainer=rating["RESPONSIVENESS_MAINTAINER_SCORE"],
            licenseScore=rating["LICENSE_SCORE"],
            goodPinningPractice=0, #rating["GOOD_PINNING_PRACTICE_SCORE"], #TODO:
            pullRequest=0, #rating["PULL_REQUEST"], TODO:
            # VERSION_SCORE?
            netScore=rating["NET_SCORE"]
        )
        result = conn.execute(ins)
        rating_pk = result.inserted_primary_key[0]
        helper.log(f"Rating inserted, rating_pk: {rating_pk}")

    # Upload to cloud bucket
    helper.log("Uploading binary to cloud bucket..")
    binary_pk = rating_pk # binary_pk is same as rating_pk
    bucket.upload_b64_blob(content, str(binary_pk)) 
    helper.log(f"Uploaded binary, binary_pk:{binary_pk}")

    # Upload to packages table
    helper.log("Uploading to packages..")
    with engine.begin() as conn:
        currentTime = datetime.datetime.now(tz=datetime.timezone.utc)

        ins = packages.insert().values(
            name=name,
            rating_pk=rating_pk,
            author_pk=author_pk,
            url=url,
            binary_pk=binary_pk,
            version=version,
            upload_time=currentTime,
            is_external=isExternal,
        )

        result = conn.execute(ins)
        package_pk = result.inserted_primary_key[0]
        helper.log(f"Package inserted, package_pk: {package_pk}")
        return package_pk

def delete_package(packageId: int) -> int:
    # Packages table
    with engine.begin() as conn:
        # Get other pk's from package table
        s = packages.select().where(packages.c.id==packageId)
        row = conn.execute(s).first()
        _, _, rating_pk, _, _, binary_pk, _, _, _ = row

        # Delete package
        d1 = packages.delete().where(packages.c.id==packageId)
        conn.execute(d1)

    # Delete from ratings table
    with engine.begin() as conn:
        d2 = ratings.delete().where(ratings.c.id==rating_pk)
        conn.execute(d2)

    return binary_pk # to be deleted

def reset_database():
    with engine.begin() as conn:
        d1 = packages.delete()
        conn.execute(d1)
        d2 = ratings.delete()
        conn.execute(d2)

def update_package(id: int, author_pk: str, rating, url: str, content, isExternal):
    # Upload to ratings table
    helper.log("Uploading to rating table..")
    with engine.begin() as conn:
        ins = ratings.insert().values(
            busFactor=rating["BUS_FACTOR_SCORE"],
            correctness=rating["CORRECTNESS_SCORE"],
            rampUp=rating["RAMP_UP_SCORE"],
            responsiveMaintainer=rating["RESPONSIVENESS_MAINTAINER_SCORE"],
            licenseScore=rating["LICENSE_SCORE"],
            goodPinningPractice=0, #rating["GOOD_PINNING_PRACTICE_SCORE"], #TODO:
            pullRequest=0, #rating["PULL_REQUEST"], TODO:
            # VERSION_SCORE?
            netScore=rating["NET_SCORE"]
        )
        result = conn.execute(ins)
        rating_pk = result.inserted_primary_key[0]
        helper.log(f"Rating inserted, rating_pk: {rating_pk}")

    # Upload to cloud bucket
    helper.log("Uploading binary to cloud bucket..")
    binary_pk = rating_pk # binary_pk is same as rating_pk
    bucket.upload_b64_blob(content, str(binary_pk)) 
    helper.log(f"Uploaded binary, binary_pk:{binary_pk}")

    # Upload to packages table
    helper.log("Updating packages table..")
    with engine.begin() as conn:
        currentTime = datetime.datetime.now(tz=datetime.timezone.utc)

        # Get old pks
        s = packages.select().where(packages.c.id==id)
        row = conn.execute(s).first()
        _, _, old_rating_pk, _, _, old_binary_pk, _, _, _ = row

        # Update with new data
        ins = packages.update().where(packages.c.id==id).values(
            rating_pk=rating_pk,
            author_pk=author_pk,
            url=url,
            binary_pk=binary_pk,
            upload_time=currentTime,
            is_external=isExternal,
        )

        result = conn.execute(ins)
    
    # Delete old rating & binary
    with engine.begin() as conn:
        d1 = ratings.delete().where(ratings.c.id==old_rating_pk)
        conn.execute(d1)
    bucket.delete_blob(str(old_binary_pk))

    helper.log(f"Package updated, package_pk: {id}")


# Functions on startup
def setup_database():
    global engine
    engine = connect_with_connector()
    create_default()