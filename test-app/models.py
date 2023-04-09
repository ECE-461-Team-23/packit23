from typing import List
from typing import Optional
from sqlalchemy import ForeignKey, String, DateTime
from sqlalchemy.orm import DeclarativeBase
from sqlalchemy.orm import Mapped
from sqlalchemy.orm import mapped_column
from sqlalchemy.orm import relationship

class Base(DeclarativeBase):  
    pass

class Users(Base):
    __tablename__ = "Users"
    id: Mapped[int] = mapped_column(primary_key = True)
    username: Mapped[str] = mapped_column(String(50))
    password: Mapped[str] = mapped_column(String(50))
    token: Mapped[str] = mapped_column(String(50))
    token_created: Mapped[DateTime] = mapped_column(DateTime)
    token_expiry: Mapped[DateTime] = mapped_column(DateTime)
    priviledge_level:Mapped[int] = mapped_column(int)
# CREATE TABLE Users ( 
# 	ID int, 
# 	USERNAME varchar(50),
# 	PASSWORD varchar(50),
# 	TOKEN varchar(50),
# 	TOKEN_CREATED datetime,
# 	TOKEN_EXPIRY datetime,
# 	PRIVILEDGE_LEVEL int,
# 	PRIMARY KEY (ID)
# );

