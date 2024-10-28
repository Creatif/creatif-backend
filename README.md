This repository is a part of Creatif but can be set up as a standalone project. 

In order to set it up, just hit `docker compose up` and it should all run.

To seed the project:

`cd cmd/creatif-sdk-seed && go run .`

The program itself has its own arguments and explanations how it works but if in doubt, hit `go run . --help` and
all will be explained. In order to set up [creatif-js-sdk](https://github.com/Creatif/creatif-js-sdk) or any future SDK,
you will need this command since testing in those projects will not work unless there is a live server (this application)
up and running with seed data. 

I do not believe in mocking. 