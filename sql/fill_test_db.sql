-- Fill Users table
-- what does privilege level look like?

-- user with valid token
INSERT INTO users (id, username, password) VALUES (
    1, 
    'valid-user',
    'strongpassword'
);

-- user with expired token
INSERT INTO users (id, username, password) VALUES (
    2, 
    'invalid-user',
    'strongpassword'
);

-- Fill Ratings table

-- Good rating
INSERT INTO ratings (id, busFactor, correctness, rampUp, responsiveMaintainer, licenseScore, goodPinningPractice, pullRequest, netScore) VALUES (
	1,
	1.0,
    1.0,
	1.0,
	1.0,
	1.0,
	1.0,
	1.0,
	1.0
);

-- Bad rating
INSERT INTO ratings (id, busFactor, correctness, rampUp, responsiveMaintainer, licenseScore, goodPinningPractice, pullRequest, netScore) VALUES (
	2,
	0.0,
    0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0
);

-- Fill Registry table

INSERT INTO packages (id, name, rating_pk, author_pk, url, binary_pk, version, upload_time, is_external) VALUES (
	1,
	'package1-name',
	1,
	1, 
	'github.com/packit461/packit23',
	1,
	'2.0',
	'2023-04-14 02:20:20',
	False
);
INSERT INTO packages (id, name, rating_pk, author_pk, url, binary_pk, version, upload_time, is_external) VALUES (
	2,
	'package2-name',
	2,
	2, 
	'github.com/19chonm/461_1_23',
	2,
	'1.5',
	'2023-04-14 02:20:20',
	False
);
INSERT INTO packages (id, name, rating_pk, author_pk, url, binary_pk, version, upload_time, is_external) VALUES (
	3,
	'package2-name',
	2,
	2, 
	'github.com/19chonm/461_1_23',
	2,
	'1.5.0',
	'2023-04-14 02:20:20',
	False
);
INSERT INTO packages (id, name, rating_pk, author_pk, url, binary_pk, version, upload_time, is_external) VALUES (
	3,
	'package2-name',
	2,
	2, 
	'github.com/19chonm/461_1_23',
	2,
	'1.5.5',
	'2023-04-14 02:20:20',
	False
);
INSERT INTO packages (id, name, rating_pk, author_pk, url, binary_pk, version, upload_time, is_external) VALUES (
	4,
	'package2-name',
	2,
	2, 
	'github.com/19chonm/461_1_23',
	2,
	'1.6.0',
	'2023-04-14 02:20:20',
	False
);
INSERT INTO packages (id, name, rating_pk, author_pk, url, binary_pk, version, upload_time, is_external) VALUES (
	6,
	'newpackage',
	2,
	2, 
	'github.com/19chonm/461_1_23',
	2,
	'1.2.3',
	'2023-04-14 02:20:20',
	False
);

INSERT INTO packages (id, name, rating_pk, author_pk, url, binary_pk, version, upload_time, is_external) VALUES (
	7,
	'newpackage',
	2,
	2, 
	'github.com/19chonm/461_1_23',
	2,
	'1.2.6',
	'2023-04-14 02:20:20',
	False
);
INSERT INTO packages (id, name, rating_pk, author_pk, url, binary_pk, version, upload_time, is_external) VALUES (
	8,
	'newpackage',
	2,
	2, 
	'github.com/19chonm/461_1_23',
	2,
	'1.5.0',
	'2023-04-14 02:20:20',
	False
);
INSERT INTO packages (id, name, rating_pk, author_pk, url, binary_pk, version, upload_time, is_external) VALUES (
	9,
	'newpackage',
	2,
	2, 
	'github.com/19chonm/461_1_23',
	2,
	'2.1.0',
	'2023-04-14 02:20:20',
	False
);