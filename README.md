# MYIHX GOLANG CODEBASE

A Golang restful API boilerplate based on Echo framework v4. Includes tools for module generation, db migration, authorization, authentication and more.
Any feedback and pull requests are welcome and highly appreciated. Feel free to open issues just for comments and discussions.

<!--toc-->
- [MYIHX GOLANG CODEBASE]
  - [HOW TO USE THIS TEMPLATE](#how-to-use-this-template)
  - [Features](#features)
  - [Running the project](#running-the-project)
  - [Environment variables](#environment-variables)
  - [Commands](#commands)
  - [Folder structure](#folder-structure)
  - [Open source refs](#open-source-refs)
  - [Contributing](#contributing)
  - [TODOs](#todos)

<!-- tocstop -->

## HOW TO USE THIS TEMPLATE

> **DO NOT FORK** this is meant to be used from **[Use this template](https://github.com/dzungtran/echo-rest-api/generate)** feature.

1. Give a name to your project at *go.mod*
   (e.g. `my_awesome_project` recommendation is to use all lowercase and underscores separation for repo names.)
2. Rename all projectname in directory **bin** with new project name (e.g. `my_awesome_project`) 
   
<!--
## Overview
 
![Request processing flow - Sequence Diagram](out/docs/diagrams/overview/request_flow.svg) -->


## Running the project

- Make sure you have docker installed.
- Copy `.env.example` to `.env`
- Run `make run`.
- Go to `localhost:8088` to verify if the API server works.

## Swagger Docs

To create swagger api documentation you just have to run `make docs`.    
after the command executed successfully, run the app and open `http://localhost:8088/docs/index.html` to see swagger documentation page.



## Commands

| Command                                  | Description                                                 |
|------------------------------------------|-------------------------------------------------------------|
| `make run`                               | Start DEV REST API application                              |
| `make start`                             | Start REST API application                                  |

## Folder Structure

```
.
├── bin         # Thirdparty configs
   ├── app
   ├── config            # configuration of env
   ├── middlewares        # Auth middleware
   └── modules          # Core module, includes apis: users, orgs        # Thirdparty configs
   └── pkg          # helpers
├── key
│   └── private.key             # key for jwt key
│   └── public.key             # key for jwt key


```

## Open Source Refs
- https://cuelang.org/docs/about/
- https://www.openpolicyagent.org/docs/latest/
- https://echo.labstack.com/guide/
- https://firebase.google.com/docs/auth/admin/
- https://pkg.go.dev/firebase.google.com/go/auth


## Contributing

Please open issues if you want the template to add some features that is not in todos.

Create a PR with relevant information if you want to contribute in this template.

## TODOs

- [x] Update docker compose for ory/kratos.
- [x] Update README.md.
- [ ] Update API docs.
- [ ] Write more tests.

## AUTHOR
- Farid Tri Wicaksono [https://github.com/farid-alfernass]