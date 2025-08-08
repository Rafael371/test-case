# test-case
test case for Roketin

This project contains example HTTP requests for use with [Insomnia](https://insomnia.rest/). You can import these requests to test the available API endpoints.

---

## How to Run the Project

This project uses **PostgreSQL** as the database that can be run with Docker.

---


### Option 1: Run with Docker

1. Open your terminal.

2. Navigate to the project directory.

3. Run the following command and wait for the process to complete:

    ```
    docker compose up --build
    ```

4. Open another terminal and run the following command to access the PostgreSQL container:

    ```
    docker exec -it test_case_db psql -U user -d testdb
    ```

5. Create the `movies` table by running this SQL command:

    ```sql
    CREATE TABLE movies (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        duration INT,
        description TEXT,
        artists TEXT[],
        genres TEXT[]
    );
    ```

6. The project is now ready to use.

---

### Option 2: Run without Docker

1. Open the `main.go` file.

2. Comment out the following lines:

    ```go
    import "test-case/models" // line 8
    models.InitDB()           // line 12
    ```

3. Open your terminal.

4. Navigate to the project directory.

5. Run the project using the following command:

    ```
    go run main.go
    ```

6. The project is now ready to use.
