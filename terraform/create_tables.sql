CREATE TABLE IF NOT EXISTS packages (
    id int PRIMARY KEY,
    name VAR(255) NOT NULL,
    rating INT,
    author INT,
    url VARCHAR(255),
    binary INT,
    version INT,
    upload_time DATETIME,
    is_external BOOL
);

CREATE TABLE IF NOT EXISTS users (
    id int PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    token VARCHAR(255),
    token_creation_time DATETIME,
    token_expiry DATETIME,
    privilege_level INT
);

CREATE TABLE IF NOT EXISTS ratings (
    id int PRIMARY KEY,
    bus_factor_score FLOAT,
    ramp_up_score FLOAT,
    responsive_maintainers_score FLOAT,
    license_score FLOAT,
    pinning_score FLOAT,
    pull_request_score FLOAT,
    net_score FLOAT
);
