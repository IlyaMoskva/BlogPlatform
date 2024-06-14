בס"ד

# BlogPlatform
### Design and implement blog platform REST API in Golang

[![Go](https://github.com/IlyaMoskva/BlogPlatform/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/IlyaMoskva/BlogPlatform/actions/workflows/go.yml)

# Problem Statement:
You are tasked with building a simple RESTful API for a blog platform. The API should allow creating, updating, deleting, and retrieving blog posts. Each blog post should have a title, content, and an author.

# Requirements:

* Implement CRUD (Create, Read, Update, Delete) operations for blog posts.
* Use a simple in-memory data store (e.g., a slice or a map) to store blog posts, using the attached JSON data sample.
* Design the API to follow RESTful principles.
* Include error handling for common scenarios (e.g., not found, validation errors).
* Write unit tests to ensure the reliability of your code.

# Endpoint Examples:

* GET /posts: Retrieve a list of all blog posts.
* GET /posts/{id}: Retrieve details of a specific blog post.
* POST /posts: Create a new blog post.
* PUT /posts/{id}: Update an existing blog post.
* DELETE /posts/{id}: Delete a blog post.
