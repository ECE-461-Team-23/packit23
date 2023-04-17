import os

from google.cloud.sql.connector import Connector, IPTypes
import pymysql

import sqlalchemy


def connect_with_connector_auto_iam_authn() -> sqlalchemy.engine.base.Engine:
    """
    Initializes a connection pool for a Cloud SQL instance of MySQL.

    Uses the Cloud SQL Python Connector with Automatic IAM Database Authentication.
    """
    # Note: Saving credentials in environment variables is convenient, but not
    # secure - consider a more secure solution such as
    # Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
    # keep secrets safe.
    instance_connection_name = os.environ["INSTANCE_CONNECTION_NAME"]  # e.g. 'project:region:instance'
    db_iam_user = os.environ["DB_IAM_USER"]  # e.g. 'sa-name@project-id.iam'
    db_name = os.environ["DB_NAME"]  # e.g. 'my-database'

    ip_type = IPTypes.PRIVATE if os.environ.get("PRIVATE_IP") else IPTypes.PUBLIC

    # initialize Cloud SQL Python Connector object
    connector = Connector()

    def getconn() -> pymysql.connections.Connection:
        conn: pymysql.connections.Connection = connector.connect(
            instance_connection_name,
            "pymysql",
            user=db_iam_user,
            db=db_name,
            enable_iam_auth=True,
            ip_type=ip_type,
        )
        return conn

    # The Cloud SQL Python Connector can be used with SQLAlchemy
    # using the 'creator' argument to 'create_engine'
    pool = sqlalchemy.create_engine(
        "mysql+pymysql://",
        creator=getconn,
        # ...
    )
    return pool


# import os

# from google.cloud.sql.connector import Connector, IPTypes
# import pymysql

# import sqlalchemy

# # https://github.com/GoogleCloudPlatform/cloud-sql-python-connector/blob/main/README.md#iam-authentication
# def connect_with_connector() -> sqlalchemy.engine.base.Engine:
#     """
#     Initializes a connection pool for a Cloud SQL instance of MySQL.

#     Uses the Cloud SQL Python Connector package.
#     """
#     # Note: Saving credentials in environment variables is convenient, but not
#     # secure - consider a more secure solution such as
#     # Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
#     # keep secrets safe.

#     instance_connection_name = os.environ["INSTANCE_CONNECTION_NAME"]  # e.g. 'project:region:instance'
#     db_user = os.environ.get("DB_USER", "")  # e.g. 'my-db-user'
#     db_name = os.environ["DB_NAME"]  # e.g. 'my-database'

#     ip_type = IPTypes.PRIVATE if os.environ.get("PRIVATE_IP") else IPTypes.PUBLIC

#     connector = Connector(ip_type)

#     def getconn() -> pymysql.connections.Connection:
#         conn: pymysql.connections.Connection = connector.connect(
#             instance_connection_name,
#             "pymysql",
#             user=db_user,
#             db=db_name,
#             enable_iam_auth=True,
#         )
#         return conn

#     pool = sqlalchemy.create_engine(
#         "mysql+pymysql://",
#         creator=getconn,
#         # ...
#     )
#     return pool

def read_rows():
    engine = connect_with_connector_auto_iam_authn()
    connection = engine.connect()
    print(f"made connection! {connection}")
    # metadata = sqlalchemy.db.MetaData()
    # print(metadata)

    # census = sqlalchemy.db.Table('census', metadata, autoload=True, autoload_with=engine)
    # query = sqlalchemy.db.select([census]) 
    # ResultProxy = connection.execute(query)
    # ResultSet = ResultProxy.fetchall()

    # print(ResultSet)

