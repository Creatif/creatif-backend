This repository is a part of Creatif but can be set up as a standalone project. 

In order to set it up, just hit `docker compose up` and it should all run.

To seed the project:

`cd cmd/creatif-sdk-seed && go run .`

The program itself has its own arguments and explanations how it works but if in doubt, hit `go run . --help` and
all will be explained. In order to set up [creatif-js-sdk](https://github.com/Creatif/creatif-js-sdk) or any future SDK,
you will need this command since testing in those projects will not work unless there is a live server (this application)
up and running with seed data. 

**IMPORTANT**

Since the seed program also uploads a lot of images, if you use this program a lot, there could be a lot of images uploaded
and created in your `var` and `public` directories. These directories are safe to delete and are recreated if they do not exist.
Note that you will have to use `sudo` to do so since these are docker volumes:

`sudo rm -rf var && sudo rm -rf public`