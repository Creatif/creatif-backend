package main

import (
	"github.com/fatih/color"
)

const URL = "http://localhost:3002/api/v1"

const Cannot_Continue_Procedure = "cannot_continue"
const Can_Continue = "Everything is OK"

const Email = "mariofake@gmail.com"
const Password = "password"

var help = `
WARNING: THIS IS A DESTRUCTIVE COMMAND. IN CASE OF CERTAIN ERRORS, IT MIGHT DESTROY ALL THE DATA THAT YOU HAVE IN THE DATABASE. USE WITH CAUTION!!!

IMPORTANT:
This seed actually uploads images. Every clientVariable gets one image and every property gets 3 images. It would be wise to from time to time, just delete the 'var' and 'public' directories because they might get very large if you execute this function over and over again.

This program cannot start if you don't have the server up, so make sure that you open up a new terminal tab, hit 'docker compose up' on the main project
and only then execute this command.

This command seeds the initial application with seed data from real estate project. It has two structures: Clients and Properties. Clients is a map and Properties is a list. It generates five projects with those structure. Each project has 200 Client maps and 1000 (one thousand) Properties in 5 different locales. That means that this command will generate 5200 "entities" per project. There will be 5 projects so 26 thousand "entities" will be created in total.

This command will be used to test public SDKs. For now, there is only javascript SDK but hopefully, there will be more.

If you try to execute this command more than once, it not allow you to do that. Since the application can have a single admin (for now), you cannot create another admin, therefor the program will tell you that the app is already seeded.Calling this program without any flags will create five projects by default.

USAGE:

There is nothing special about this program. Just cd into this directory and run 'go run .' and that is it.

Flags:
--cleanup
    This flag will completely destroy all data in the database. USE WITH CAUTION!!! If you use this flag is the only thing that will be done even if you used other flags i.e. it will ignore all other flags.
--regenerate
    This will do what --cleanup does but will run other commands. Basically, you tell the program to wipe the database out start over
--projects={\d} 
    For how many projects should it seed the application. More will be slower.
--properties-per-type
    For every property type (House, Apartment, Land, Studio), generate that many properties. Every property is generated
    with 5 languages and 3 property statuses. That means for every language and every property status, generate that many
	properties per type. For example, if --properties-per-type=10, then 5 * 3 * 4 * 10 (600) properties will be generated.
--help
    If used, will output help. If used with any other flags, it will ignore them and just print help, i.e. it will ignore all other flags. 

Credentials:
Email: mariofake@gmail.com
Password: password

I know that password is weak, but there is a plan to put a password strength into effect in the application but until then, this is just fine.
`

var printers = map[string]*color.Color{
	"success": color.New(color.FgGreen).Add(color.Bold),
	"error":   color.New(color.FgRed).Add(color.Bold),
	"warning": color.New(color.FgYellow).Add(color.Bold),
	"info":    color.New(color.FgWhite).Add(color.Bold),
}
