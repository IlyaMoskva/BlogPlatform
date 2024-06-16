בס"ד

# BlogPlatform
### Design and implement blog platform REST API in Golang

[![Go](https://github.com/IlyaMoskva/BlogPlatform/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/IlyaMoskva/BlogPlatform/actions/workflows/go.yml)

## Problem Statement:
You are tasked with building a simple RESTful API for a blog platform. The API should allow creating, updating, deleting, and retrieving blog posts. Each blog post should have a title, content, and an author.

## Requirements:

* Implement CRUD (Create, Read, Update, Delete) operations for blog posts.
* Use a simple in-memory data store (e.g., a slice or a map) to store blog posts, using the attached JSON data sample.
* Design the API to follow RESTful principles.
* Include error handling for common scenarios (e.g., not found, validation errors).
* Write unit tests to ensure the reliability of your code.

## Endpoint Examples:

* GET /posts: Retrieve a list of all blog posts.
* GET /post/{id}: Retrieve details of a specific blog post.
* POST /post: Create a new blog post.
* PUT /post/{id}: Update an existing blog post.
* DELETE /post/{id}: Delete a blog post.

## Swagger documentation
Available here: https://localhost:8443/swagger/doc

# Initial Preparation

## Certificate Generation

Solution uses self-signed certificate to work over https. Pair of keys is prepared already and placed in the root folder.
To generate them again run
```sh
go run certgen/certgen.go
```
Generated files will be used by main program.
Use https://localhost:8443/ as a main API path.

## Run the Server
```sh
go buld main.go
go run main.go
```

## Import Initial Data (optional)
There is an Admin API to load data from json file with the given format.
Use Postman collection test "POST import from file"
```
curl --location 'https://localhost:8443/admin/import' \
--form 'file=@"/path_to_file/blog_data.json"'
```
Original file is placed in "files" folder.

Import can be done in any moment later, post Ids will be generated automatically. There is no duplication prevention.
