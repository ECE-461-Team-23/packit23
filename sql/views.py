from google.cloud.sql.connector import Connector
import sqlalchemy
import setup
from sqlalchemy import select

# /packages - POST
# This is probably fine, but only for lots at once
# I dont really get what this returns and how we'd pass it to JS
packagesResults = setup.db_conn.execute(
    sqlalchemy.text(
        "BEGIN" \
        "SELECT * FROM Registry AS A" \
            "INNER JOIN Binaries AS B" \
                "ON A.BINARY_PK = B.ID" \
            "INNER JOIN Users AS C" \
                "ON A.USER_PK = C.ID" \
            "INNER JOIN Ratings AS D" \
                "ON A.RATING_PK = D.ID" \
        "END;"
    )
)

requestedID = 1 #EMILE: This needs to be changed in the future tbh
session = Session(engine) 
stmt = (
    select(Registry)
    .where(Registry.id == requestedID)
    .join(Users)
)
# /package/(id) - GET
# This is probably not fine, you'd want som specific values
singlePackageResult = setup.db_conn.execute(
    sqlalchemy.text(
        "BEGIN" \
        "SELECT * FROM Registry AS A" \
            "WHERE" \
            "A.ID = @@Id;" \
            "INNER JOIN Binaries AS B" \
                "ON A.BINARY_PK = B.ID" \
            "INNER JOIN Users AS C" \
                "ON A.USER_PK = C.ID" \
            "INNER JOIN Ratings AS D" \
                "ON A.RATING_PK = D.ID" \
        "END;"
    )
)
# print(singlePackageResult)
