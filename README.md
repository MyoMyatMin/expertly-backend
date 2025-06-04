# Expertly Backend

A robust backend API for a social platform with features for content sharing, user interactions, and moderation.

## Setup Instructions

To run this project, follow these steps:

1. **Install Go**: Follow the instructions [here](https://go.dev/doc/install) to install Go on your laptop.

2. **Install sqlc and goose**:
    ```sh
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```

    For more information, refer to the documentation:
    - [goose documentation](https://pkg.go.dev/github.com/pressly/goose/v3)
    - [sqlc documentation](https://sqlc.dev/)

3. **Run database migrations**:
    ```sh
    cd sql/schema
    goose postgres "postgresql://username:password@databasename?sslmode=require" up
    ```

4. **Generate SQL code**:
    ```sh
    cd ../..  
    sqlc generate
    ```

5. **Create `sqlc.yaml`**:
    ```yaml
    version: "2"
    sql:
      - schema: "sql/schema"
        queries: "sql/queries"
        engine: "postgresql"
        gen:
          go:
            out: "pkg/database"
    ```

6. **Build and run the project**:
    ```sh
    go build && ./expertly-backend
    ```

7. **Set environment variables**:
    Create a `.env` file with the following content:
    ```env
    PORT=
    DB_URL=
    SECRET_KEY=
    CLOUDINARY_CLOUD_NAME=
    CLOUDINARY_API_KEY=
    CLOUDINARY_API_SECRET=
    ```

## Features

### User Authentication
- User registration and login
- JWT-based authentication
- Middleware for protected routes

### Posts
- Create, read, update, and delete posts
- Get posts by user
- Get post details
- Feed generation

### Comments
- Create, update, and delete comments
- Nested comments support (replies to comments)
- Get all comments for a post

### Interactions
- Upvote/downvote posts
- Follow/unfollow users
- Save posts for later viewing

### User Profiles
- View user profiles
- Update profile information
- View user's posts

### Moderation System
- Report users and contributors
- Admin login and moderator creation
- Review and update report status

### Appeals System
- Create appeals for reported content
- Review and update appeal status
- View appeals by type (user or contributor)

### Contributor Applications
- Apply to become a contributor
- Review and update application status

### Search
- Search for posts
- Search for users

## API Endpoints

### Authentication
- `POST /v1/auth/register` - Register a new user
- `POST /v1/auth/login` - Login a user

### Posts
- `POST /v1/posts` - Create a new post
- `GET /v1/posts/{id}` - Get post details
- `GET /v1/users/{id}/posts` - Get posts by user
- `GET /v1/feed` - Get feed posts

### Comments
- `POST /v1/posts/{id}/comments` - Create a comment
- `PUT /v1/comments/{id}` - Update a comment
- `DELETE /v1/comments/{id}` - Delete a comment
- `GET /v1/posts/{id}/comments` - Get all comments for a post

### Interactions
- `POST /v1/posts/{id}/upvote` - Upvote a post
- `DELETE /v1/posts/{id}/upvote` - Remove upvote
- `POST /v1/users/{id}/follow` - Follow a user
- `DELETE /v1/users/{id}/follow` - Unfollow a user
- `POST /v1/posts/{id}/save` - Save a post
- `DELETE /v1/posts/{id}/save` - Unsave a post
- `GET /v1/users/{id}/saved` - Get saved posts

### User Profiles
- `GET /v1/users/{id}` - Get user profile
- `PUT /v1/users/{id}` - Update user profile

### Moderation
- `POST /v1/admin/login` - Admin login
- `POST /v1/admin/moderators` - Create a moderator
- `GET /v1/admin/moderators` - Get all moderators

### Reports
- `POST /v1/reports` - Create a report
- `GET /v1/reports/contributors` - Get reported contributors
- `GET /v1/reports/users` - Get reported users
- `PUT /v1/reports/{id}` - Update report status

### Appeals
- `POST /v1/appeals` - Create an appeal
- `GET /v1/appeals` - Get all appeals
- `GET /v1/appeals/contributors` - Get contributor appeals
- `GET /v1/appeals/users` - Get user appeals
- `GET /v1/appeals/{id}` - Get appeal by ID
- `PUT /v1/appeals/{id}` - Update appeal status

### Contributor Applications
- `POST /v1/contributor-applications` - Apply to be a contributor
- `GET /v1/contributor-applications` - Get all applications
- `GET /v1/contributor-applications/{id}` - Get application by ID
- `PUT /v1/contributor-applications/{id}` - Update application status

### Search
- `GET /v1/search/posts` - Search posts
- `GET /v1/search/users` - Search users

## Architecture

- **Language**: Go
- **Database**: PostgreSQL
- **ORM**: sqlc for type-safe SQL
- **Migration**: goose for database migrations
- **Router**: chi for HTTP routing
- **Authentication**: JWT for secure authentication
- **File Storage**: Cloudinary for media storage

## Security

- JWT-based authentication
- Protected routes with middleware
- Environment variables for sensitive information
- Database connection security with SSL
