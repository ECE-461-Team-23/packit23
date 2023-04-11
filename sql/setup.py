from google.cloud.sql.connector import Connector
import sqlalchemy

# initialize Connector object
connector = Connector()

#EMILE
# function to return the database connection object
def getconn():
    conn = connector.connect( 
        INSTANCE_CONNECTION_NAME,
        "pymysql",
        user=DB_USER,
        password=DB_PASS,
        db=DB_NAME
    )
    return conn

# create connection pool with 'creator' argument to our connection object function
pool = sqlalchemy.create_engine(
    "mysql+pymysql://",
    creator=getconn,
)

# connect to connection pool
with pool.connect() as db_conn:
    # create Registry table in our database
    db_conn.execute(
        sqlalchemy.text(
            "CREATE TABLE IF NOT EXISTS Registry (" \
            "ID int, NAME varchar (50), RATING_PK int, AUTHOR_PK int," \
            "URL varchar(255), BINARY_PK int, VERSION varchar (15)," \
            "UPLOADED datetime, IS_EXTERNAL boolean, PRIMARY KEY (ID)," \
            "FOREIGN KEY(RATING_PK) REFERENCES Rating(ID)," \
            "FOREIGN KEY(AUTHOR_PK) REFERENCES Users(ID)," \
            "FOREIGN KEY(BINARY_PK) REFERENCES Binaries(ID),);" 
        )
    )
    # commit transaction (SQLAlchemy v2.X.X is commit as you go)
    db_conn.commit()

    # create Users table in our database
    db_conn.execute(
        sqlalchemy.text(
            "CREATE TABLE IF NOT EXISTS Users (" \
            "ID int, USERNAME varchar(50), PASSWORD varchar(50)," \
            "TOKEN varchar(50), TOKEN_CREATED datetime, TOKEN_EXPIRY datetime," \
            "PRIVILEDGE_LEVEL int, PRIMARY KEY (ID));" 
        )
    )
    # commit transaction (SQLAlchemy v2.X.X is commit as you go)
    db_conn.commit()

    # create Ratings table in our database
    db_conn.execute(
        sqlalchemy.text(
            "CREATE TABLE IF NOT EXISTS Ratings (" \
            "ID int, BUS_FACTOR float, CORRECTNESS float," \
            "RAMP_UP float, RESPONSIVENESS float, LICENSE_SCORE float," \
            "PINNING_PRACTICE float, PULL_REQUEST float, NET_SCORE float, PRIMARY KEY (ID));" 
        ) 
    )
    # commit transaction (SQLAlchemy v2.X.X is commit as you go)
    db_conn.commit()

    # Unsure if this is needed
    # # create Binaries table in our database
    # db_conn.execute(
    #     sqlalchemy.text(
    #         "CREATE TABLE IF NOT EXISTS Binaries (" \
    #         "ID int, BINARY_FILE ????, JS_PROGRAM varchar(4),PRIMARY KEY (ID),);" \
    #     )
    # )
    # # commit transaction (SQLAlchemy v2.X.X is commit as you go)
    # db_conn.commit()