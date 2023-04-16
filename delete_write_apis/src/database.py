import os

from google.cloud.sql.connector import Connector, IPTypes
import pymysql

import sqlalchemy

# https://github.com/GoogleCloudPlatform/cloud-sql-python-connector/blob/main/README.md#iam-authentication
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
    db_name = os.environ["DB_NAME"]  # e.g. 'my-database'

    ip_type = IPTypes.PRIVATE if os.environ.get("PRIVATE_IP") else IPTypes.PUBLIC

    connector = Connector(ip_type)

    def getconn() -> pymysql.connections.Connection:
        conn: pymysql.connections.Connection = connector.connect(
            instance_connection_name,
            "pymysql",
            user=db_user,
            db=db_name,
            enable_iam_auth=True,
        )
        return conn

    pool = sqlalchemy.create_engine(
        "mysql+pymysql://",
        creator=getconn,
        # ...
    )
    return pool

def read_rows():
    engine = connect_with_connector()
    connection = engine.connect()
    metadata = sqlalchemy.db.MetaData()
    print(metadata)

    census = sqlalchemy.db.Table('census', metadata, autoload=True, autoload_with=engine)
    query = sqlalchemy.db.select([census]) 
    ResultProxy = connection.execute(query)
    ResultSet = ResultProxy.fetchall()

    print(ResultSet)

