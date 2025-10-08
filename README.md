# Voucher Management System Backend

This document explains how to run the backend application for the Voucher Management System.

## Prerequisites

Make sure you have installed:

* Go (minimum version 1.20)
* Docker & Docker Compose
* Git (optional, if cloning from repository)

## Steps to run the application

1. **Install Go dependencies**

   Run the following command to download all required Go dependencies:

   ```bash
   go get
   ```

2. **Create the `.env` file**

   Copy the example environment file to `.env`:

   ```bash
   cp .env.example .env
   ```

   Adjust environment configurations as needed (database, port, etc.).

3. **Start the database using Docker**

   Use Docker Compose to run the database:

   ```bash
   docker-compose up -d
   ```

   Make sure the database is running before proceeding.

4. **Run database migrations**

   Run migrations to create the required tables:

   ```bash
   go run cmd/migrate.go
   ```

5. **Start the backend application**

   After migrations, run the backend server:

   ```bash
   go run cmd/main.go
   ```

   By default, the server will run at `http://localhost:8080`.

## Notes

* If the database schema changes, run `go run cmd/migrate.go` again.
* Make sure port `8080` is free.
* Use `docker-compose logs -f` to see database logs if needed.
