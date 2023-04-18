import os

from google.cloud.sql.connector import Connector, IPTypes
import pymysql

import sqlalchemy
from sqlalchemy import inspect
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from sqlalchemy import Table, Column, Integer, String, MetaData, LargeBinary, DateTime, Float

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
    Column('version', Integer, nullable=False),
    Column('upload_time', Integer, nullable=False),
    Column('is_external', DateTime, nullable=False),
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
engine = connect_with_connector()

def create_default():
    # engine = connect_with_connector()
    metadata.create_all(engine) # Create tables

    with engine.begin() as conn:
        ins = users.insert().values(username="sampleuser", password=b'$2b$12$iUwbqGXqpky94Xul/l0RMe.nIe0HGFha9eQ0M9.82MIeSNvxoXNje')
        result = conn.execute(ins)

def read_rows():
    # create_default()
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

def get_password_for_user(username: str) -> str:
    with engine.begin() as conn:
        s = users.select().where(users.c.username==username)
        result = conn.execute(s)
        row = result.fetchone()
        print(f"get_password_for_user({username}): found {row}")

        if row:
            return row[2] #password
        else:
            return None

def check_if_package_exists(packageName: str, packageVersion: str):
    # maybe some others?
    pass

def upload_package():
    pass