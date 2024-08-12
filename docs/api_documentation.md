# Updated Task Management API Documentation

This document provides comprehensive information about the Task Management API, which uses MongoDB for persistent data storage. The API now includes user authentication and authorization features, enabling role-based access control.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Setup](#setup)
3. [API Endpoints](#api-endpoints)
   - [User Endpoints](#user-endpoints)
   - [Task Endpoints](#task-endpoints)
4. [Data Models](#data-models)
   - [User Model](#user-model)
   - [Task Model](#task-model)
5. [Authentication & Authorization](#authentication--authorization)
6. [Error Handling](#error-handling)
7. [Testing the API](#testing-the-api)
8. [MongoDB Inspection](#mongodb-inspection)
9. [API Versioning](#api-versioning)

---

## Prerequisites

- **Go**: Version 1.18 or later
- **MongoDB**: Version 4.4 or later
- **Gin Framework**: [Gin](https://github.com/gin-gonic/gin) for Go
- **MongoDB Go Driver**: [MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver)
- **Environment Variables**: Configure the following in a `.env` file
  - `MONGODB_URI`: MongoDB connection string
  - `JWT_SECRET`: Secret key for JWT token generation

## Setup

1. **Install Dependencies**: Ensure you have Go and MongoDB installed.

2. **Configure Environment Variables**: Create a `.env` file in the root directory of the project.

   ### Example `.env` File

   ```dotenv
   MONGODB_URI=mongodb://localhost:27017/taskDB
   JWT_SECRET=your_secret_key_here
   ```

3. **Run the Application**:

   Navigate to the root directory of the application and execute:

   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8000`.

## API Endpoints

### User Endpoints

1. **Register a New User**

   - **URL**: `/register`
   - **Method**: `POST`
   - **Description**: Creates a new user account.
   - **Request Body**: JSON object with user details

     ```json
     {
       "username": "string",
       "password": "string"
     }
     ```

   - **Response**:
     - **Status Code**: `201 Created` (on success), `400 Bad Request` (on validation errors), `500 Internal Server Error` (on server errors)
     - **Body**: JSON object with a success message or error details

2. **User Login**

   - **URL**: `/login`
   - **Method**: `POST`
   - **Description**: Authenticates a user and returns a JWT token.
   - **Request Body**: JSON object with login credentials

     ```json
     {
       "username": "string",
       "password": "string"
     }
     ```

   - **Response**:
     - **Status Code**: `200 OK` (on success), `400 Bad Request` (on validation errors), `401 Unauthorized` (on authentication failure)
     - **Body**: JSON object containing the JWT token or error details

     ```json
     {
       "token": "jwt_token_string"
     }
     ```

3. **Promote User to Admin** _(Admin Only)_

   - **URL**: `/promote/:username`
   - **Method**: `POST`
   - **Description**: Elevates a user's role to admin.
   - **Parameters**:
     - `username`: The username of the user to be promoted
   - **Headers**:
     - `Authorization`: `Bearer {jwt_token}`
   - **Response**:
     - **Status Code**: `200 OK` (on success), `403 Forbidden` (if not authorized), `500 Internal Server Error` (on server errors)
     - **Body**: JSON object with a success message or error details

### Task Endpoints

> **Note**: All task endpoints, except `GET /tasks` and `GET /tasks/:id`, require authentication. Creation, updating, and deletion of tasks are restricted to users with the **admin** role.

1. **Get All Tasks**

   - **URL**: `/tasks`
   - **Method**: `GET`
   - **Description**: Retrieves all tasks from the database.
   - **Headers**:
     - `Authorization`: `Bearer {jwt_token}`
   - **Response**:
     - **Status Code**: `200 OK`
     - **Body**: JSON array of task objects

2. **Get a Specific Task**

   - **URL**: `/tasks/:id`
   - **Method**: `GET`
   - **Description**: Retrieves a specific task by its ID.
   - **Parameters**:
     - `id`: The ID of the task (string)
   - **Headers**:
     - `Authorization`: `Bearer {jwt_token}`
   - **Response**:
     - **Status Code**: `200 OK` (if found), `404 Not Found` (if not found)
     - **Body**: JSON object of the task (if found)

3. **Create a New Task** _(Admin Only)_

   - **URL**: `/tasks`
   - **Method**: `POST`
   - **Description**: Creates a new task.
   - **Headers**:
     - `Authorization`: `Bearer {jwt_token}`
   - **Request Body**: JSON object with task details

     ```json
     {
       "id": "string",
       "title": "string",
       "description": "string",
       "due_date": "2023-08-09T00:00:00Z",
       "status": "string" // Allowed values: "pending", "in-progress", "completed"
     }
     ```

   - **Response**:
     - **Status Code**: `201 Created` (on success), `400 Bad Request` (on validation errors), `500 Internal Server Error` (on server errors)
     - **Body**: JSON object with a success message or error details

4. **Update a Task** _(Admin Only)_

   - **URL**: `/tasks/:id`
   - **Method**: `PUT`
   - **Description**: Updates an existing task.
   - **Parameters**:
     - `id`: The ID of the task to update (string)
   - **Headers**:
     - `Authorization`: `Bearer {jwt_token}`
   - **Request Body**: JSON object with updated task details

     ```json
     {
       "title": "string",
       "description": "string",
       "due_date": "2023-08-09T00:00:00Z",
       "status": "string" // Allowed values: "pending", "in-progress", "completed"
     }
     ```

   - **Response**:
     - **Status Code**: `200 OK` (if updated), `400 Bad Request` (on validation errors), `404 Not Found` (if not found), `500 Internal Server Error` (on server errors)
     - **Body**: JSON object with a success message or error details

5. **Delete a Task** _(Admin Only)_

   - **URL**: `/tasks/:id`
   - **Method**: `DELETE`
   - **Description**: Deletes a task by its ID.
   - **Parameters**:
     - `id`: The ID of the task to delete (string)
   - **Headers**:
     - `Authorization`: `Bearer {jwt_token}`
   - **Response**:
     - **Status Code**: `200 OK` (if deleted), `404 Not Found` (if not found), `500 Internal Server Error` (on server errors)
     - **Body**: JSON object with a success message or error details

## Data Models

### User Model

```go
type User struct {
    ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Username string             `json:"username" bson:"username"`
    Password string             `json:"password" bson:"password"`
    Role     string             `json:"role" bson:"role"` // Possible values: "user", "admin"
}
```

### Task Model

```go
type Task struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    DueDate     time.Time `json:"due_date"`
    Status      string    `json:"status"` // Allowed values: "pending", "in-progress", "completed"
}
```

## Authentication & Authorization

- **Authentication**: The API uses JSON Web Tokens (JWT) for authentication. Upon successful login, a JWT token is issued, which must be included in the `Authorization` header for protected routes.

  - **Header Format**:

    ```
    Authorization: Bearer {jwt_token}
    ```

- **Authorization**: Role-based access control is implemented. There are two roles:
  - **User**: Can view tasks.
  - **Admin**: Can create, update, delete tasks, and promote users.

> **Note**: The `JWT_SECRET` environment variable is crucial for token generation and validation. Ensure it's securely stored and consistent across all instances of the application.

## Error Handling

The API returns appropriate HTTP status codes and error messages in the response body when errors occur. Common error responses include:

- **400 Bad Request**: For invalid input data.
- **401 Unauthorized**: When authentication fails or the token is missing/invalid.
- **403 Forbidden**: When the user lacks necessary permissions.
- **404 Not Found**: When a requested resource doesn't exist.
- **500 Internal Server Error**: For server-side errors.

## Testing the API

Use [Postman](https://www.postman.com/) or [cURL](https://curl.se/) to test the API endpoints.

### Example: Testing `Get All Tasks` Endpoint Using cURL

```bash
curl -X GET http://localhost:8000/tasks \
     -H "Authorization: Bearer {jwt_token}"
```

Replace `{jwt_token}` with the actual token received upon successful login.

## MongoDB Inspection

You can use [MongoDB Compass](https://www.mongodb.com/products/compass) to inspect the data in your MongoDB instance.

## API Versioning

The current API version is `v1`. Future updates and changes may introduce new versions.