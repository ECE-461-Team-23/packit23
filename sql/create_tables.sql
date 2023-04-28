CREATE TABLE packages (
	id INT,
	name VARCHAR (512) NOT NULL,
	rating_pk INT NOT NULL,
	author_pk INT NOT NULL, 
	url VARCHAR(512) NOT NULL,
	binary_pk INT NOT NULL,
	version VARCHAR (512) NOT NULL,
	uploaded_time DATETIME NOT NULL,
	is_external BOOLEAN  NOT NULL,
	PRIMARY KEY (ID),
	FOREIGN KEY(RATING_PK) REFERENCES ratings(ID),
	FOREIGN KEY(AUTHOR_PK) REFERENCES users(ID)
);

CREATE TABLE users ( 
	id INT, 
	username VARCHAR(512) NOT NULL,
	password BLOB NOT NULL,
	PRIMARY KEY (ID)
);

CREATE TABLE ratings (
	id INT,
	busFactor FLOAT NOT NULL,
	correctness FLOAT NOT NULL,
	rampUp FLOAT NOT NULL,
	responsiveMaintainer FLOAT NOT NULL,
	licenseScore FLOAT NOT NULL,
	goodPinningPractice FLOAT NOT NULL,
	pullRequest FLOAT NOT NULL,
	netScore FLOAT NOT NULL,
	PRIMARY KEY (ID)
);