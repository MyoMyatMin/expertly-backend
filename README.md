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
