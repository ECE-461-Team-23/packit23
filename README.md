# **ECE 461 Team 23**

RESTful API and web interface for package repository. 

Contributors:

- Ben Schwartz @benschwartz9
- Emile Baez @emilejbm
- Anna Shen @Ashassins
- Mimi Chon @19chonm

# **Components**

# delete_write_apis

- FastAPI python backend for deleting and writing functionality of API
- Includes authentication, uploading, and deleting packages
- Handles the following routes:

  /package **(POST)**

  /package/{id} **(PUT)**

  /authenticate **(PUT)**

  /reset **(DELETE)**

  /package/{id} **(DELETE)**

  /package/byName/{name} **(DELETE)**

# read_apis
- Golang mux backend for read related APIs
- Return contents and score of package requested as well as related metadata (id, name, versions stored in database)
- Query for packages based on semantic versioning
- Handles the following routes:

  /package/{id} **(GET)**

  /package/{id}/rate **(GET)**

  /package/byName/{name} **(GET)**

  /packages **(POST)**

  /package/byRegEx **(POST)**

# package_rater

- internal process that creates a score for github and npm packages based on:

  correctness, ramp up time, responsive maintainers, bus factor, valid licenses, dependencies, code review 

# docs
- Frontend

## sql

- setup tutorial for local test environment 

## terraform
- Continuous deployment
- Configures the following GCP services:

cloud run, IAM, cloud build, artifact registry, secret manager, cloud sql, and API gateway.

- Use the following commands to bring infrastructure up:
  - terraform init
  - terraform plan
  - terraform apply
