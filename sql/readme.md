Having installed MySQL client v8.0(on mac)

# Create an alias for mysql (1) or go to directory where binary is located (2)
1. alias mysql=/usr/local/mysql/bin/mysql
2. cd /path/to/mysql/bin/mysql

# Run mysql cli
mysql --user=root -p

# Create and use test database
CREATE DATABASE test_db;
USE test_db;

# Run SQL files to create and fill tables with samples values
source /path/to/packit23/sql/test_env/create_tables.sql;
source /path/to/packit23/sql/test_env/fill_test_db.sql;

# Give permissions to users
CREATE USER 'db_user'@'localhost' IDENTIFIED BY 'oldpassword!!!';
GRANT ALL PRIVILEGES ON test_db.* TO 'db_user'@'localhost';