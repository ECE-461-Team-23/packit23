import datetime
import json
import os

from google.cloud.sql.connector import Connector, IPTypes
import pymysql

import sqlalchemy
from sqlalchemy import inspect
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from sqlalchemy import Table, Column, Integer, String, MetaData, LargeBinary, DateTime, Float, Boolean

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
        print("Creating tables")
        metadata.create_all(engine) # Create tables

        print("Creating users")
        logins_str = os.environ("USER_LOGINS")
        logins_json = json.loads(logins_str)
        logins_list = [{"username": u, "password": bytes(p, 'utf-8')} for u, p in logins_json.items()]

        with engine.begin() as conn:
            result = conn.execute(users.insert(), logins_list)
        print("Finished creating tables & users")
    else:
        print("Tables already detected, not creating")

def read_rows():
    ### FOR TESTING, CAN DELETE ###
    # engine = connect_with_connector()

    with engine.begin() as conn:
        s = users.select().where(users.c.username=="sampleruser")
        result = conn.execute(s)
        row = result.fetchone()
        print(row)
        if row:
            return {"id": row[0], "username": row[1], "password": row[2]}
        else:
            return None
    # connection = engine.connect()
    # print(f"made connection! {connection}")

    # inspector = inspect(engine)
    # print(inspector.get_table_names())

def get_data_for_user(username: str) -> str:
    # Return (userid, username, password) for a user
    with engine.begin() as conn:
        s = users.select().where(users.c.username==username)
        result = conn.execute(s)
        if result.rowcount > 1:
            raise Exception(f"Error: {result.rowcount} users found for {username}")
        row = result.fetchone()
        print(f"get_password_for_user({username}): found {row}")

        if row:
            return row
        else:
            return None

def check_if_package_exists(packageName: str, packageVersion: str):
    # Check if a given package name & version already exists
    with engine.begin() as conn:
        s = packages.select().where(packages.c.name==packageName, packages.c.version==packageVersion)
        result = conn.execute(s)
        return result.rowcount > 0

def upload_package(name: str, version: str, author_pk: str, rating, url: str, content):
    # Upload to ratings table
    print("Uploading to rating table..")
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
        print(f"Rating inserted at id: {rating_pk}")

    # TODO: upload to cloud bucket
    print("Uploading binary to cloud bucket..")
    binary_pk = 0

    # Upload to packages table
    print("Uploading to packages..")
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
            is_external=(content == None),
        )

        result = conn.execute(ins)
        package_pk = result.inserted_primary_key[0]
        print(f"Package inserted at id: {package_pk}")
        return package_pk

# Functions on startup
engine = connect_with_connector()
create_default()