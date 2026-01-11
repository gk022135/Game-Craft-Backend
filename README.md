---

# Game Craft — Learn SQL Through Gameplay

## Overview

**Game Craft** is an interactive learning game designed to help users **learn SQL by playing**.
Instead of passively reading tutorials, users solve SQL-based challenges inside a game environment.

### Tech Stack

**Backend**

* **Golang**
* **PostgreSQL**
* **Prisma ORM**
* **Docker**
* **Air** (for live reload)

**Frontend**

* **Next.js**
* **TypeScript**
* **React.js**
* **Tailwind CSS**
* **Fetch & Axios** for API calls

---

## How the Backend Works

* The backend is written in **Go (Golang)**.
* It uses **two PostgreSQL database instances**:

  1. **Main Database (`maindb`)**

     * Stores user data
     * User stats
     * Question records
     * Table metadata
  2. **Game Database (`games_db`)**

     * Stores game-related tables
     * Executes SQL queries submitted by users for challenges

### Database Communication

* **Prisma ORM** is used to interact with PostgreSQL.
* Prisma Client is generated using `npx` and used within the Go backend for database operations.

---

## Running the Project

## Option 1: Run Databases Using Docker (Recommended)

### Prerequisites

* Go installed
* Docker installed
* Air installed globally

```bash
go install github.com/cosmtrek/air@latest
```

### Steps

1. Install Go dependencies:

   ```bash
   go mod tidy
   ```

2. Pull PostgreSQL Docker image:

   ```bash
   docker pull postgres
   ```

3. Create the **main database container**:

   * Map port `5432`
   * Name the container `maindb`
   * Create a database (e.g., `mydb`)
   * Update the database URL in `.env`

4. Create the **game database container**:

   * Name the container `games_db`
   * Repeat the same setup steps as above

5. Start the backend:

   ```bash
   air
   ```

⚠️ Make sure **Air** is installed globally before running the command.

---

## Option 2: Use an Online PostgreSQL Instance (No Docker)

* Create a free PostgreSQL instance using **Supabase**
* Copy the database connection URL
* Paste it into the `.env` file

![Supabase Setup](image.png)

Start the backend:

```bash
air
```

Backend runs on **[http://localhost:3001](http://localhost:3001)**

---

## API Routes

* All backend routes are defined in `routes.go`
* Use **Postman** or any API client to test endpoints
* Pick the corresponding route based on your requirement

---

## Running the Frontend (Next.js)

The frontend setup is simple.

### Steps

```bash
npm install
npm run dev
```

Frontend will start on **[http://localhost:3000](http://localhost:3000)**

---

## Authentication Flow

### Login / Signup Page

![Login Page](image-2.png)

* First-time users must **sign up**
* An **OTP is sent to your email**
* Use a valid email address

### OTP Verification

![OTP Email](image-3.png)

---

## Main Landing Page

![Landing Page](image-1.png)

---

## Game Screens

![Game Screen](image-4.png)
![Game Screen](image-5.png)
![Game Screen](image-6.png)
![Game Screen](image-7.png)
![Game Screen](image-8.png)

---












```markdown
# GameCraft-backend: Working on Collaboration
### Gaurav started working here


# Never Forget to change directory according to cd prisma_games or prisma(stay in root)

## Setup of a New Database Without Affecting the Previous One (runs in Docker)

This document explains how to add a **new Postgres database for games** inside Docker without affecting the existing **main Postgres + Prisma setup**.

---

## Project Structure

```
backend/
├── prisma/                 # existing main Prisma setup (main DB)
├── prisma\_games/           # NEW Prisma setup for games\_db
│    └── schema.prisma
├── docker-compose.yml      # docker setup for new Postgres
├── go.mod / go.sum         # Go modules

````

---
## Step 1: Docker Setup for Postgres
Create a `docker-compose.yml` file in the `backend/` directory:

```yaml
version: "3.9"

services:
  postgres_games:
    image: postgres:15
    container_name: postgres_games
    restart: always
    ports:
      - "5434:5432" # host:container
    environment:
      POSTGRES_USER: games_user
      POSTGRES_PASSWORD: games_pass
      POSTGRES_DB: games_db
    volumes:
      - postgres_games_data:/var/lib/postgresql/data

volumes:
  postgres_games_data:
````


Run the container:
```bash
docker-compose up -d
```

Verify the container is running:
```bash
docker ps
```

---

## Step 2: Prisma Schema for New Database

Create `backend/prisma_games/schema.prisma`:

```prisma
datasource db {
  provider = "postgresql"
  url      = "postgresql://games_user:games_pass@localhost:5434/games_db?schema=public"
}

generator client {
  provider = "prisma-client-go"
  output   = "./prisma_games_client" // separate client folder
}

model Game {
  id          Int      @id @default(autoincrement())
  title       String
  genre       String
  releaseDate DateTime
  developer   String
  createdAt   DateTime @default(now())
  updatedAt   DateTime @updatedAt
}
```

---

## Step 3: Apply the Migration

Run from the `backend/` directory:

```bash
npx prisma migrate dev --schema=prisma_games/schema.prisma --name init
```

---

## Step 4: Install Prisma Go Client

If Prisma Client for Go is not yet installed:

```bash
go get github.com/steebchen/prisma-client-go
```

---

## Step 5: Generate Prisma Client

Generate the Go client for the new database:

```bash
npx prisma generate --schema=prisma_games/schema.prisma
```

This creates the client inside:

```
backend/prisma_games/prisma_games_client
```

---

## Step 6: Import in Your Go Code

Example import inside controllers:

```go
import (
    db "gamecraft-backend/prisma_games/prisma_games_client"
)
```

You can now connect to and query the new `games_db` database without touching the main database.

---

## Troubleshooting

### Error: `file or directory not found`

Run Prisma commands from the `backend/` directory, not inside `prisma_games/`.

### Error: `port already in use`

Check if another Postgres instance is running on `5434`. Update the port in `docker-compose.yml` if needed.

### Error: `Could not find schema.prisma`

Make sure to pass the schema path explicitly:

```bash
npx prisma migrate dev --schema=prisma_games/schema.prisma --name init
npx prisma generate --schema=prisma_games/schema.prisma
```

### Error: `docker-compose not found`

Install it with:

```bash
sudo apt install docker-compose
```

or use:

```bash
docker compose up -d
```

---

## Summary

* **Existing database**: untouched (main Prisma setup in `backend/prisma/`).
* **New database**: runs in Docker on port `5434`, schema defined in `prisma_games/`.
* **Go client**: generated separately in `prisma_games_client`.
* Both databases can run independently.

```

Do you want me to also add a **Go code snippet** in this README that shows how to insert and fetch from the new `Game` model, so it’s directly usable?
```



Here’s a **clean README.md** file you can use for documenting how to connect and manage your Postgres container (based on the setup we discussed):

---

# PostgreSQL in Docker – Setup & Usage Guide

This guide explains how to connect to and manage a PostgreSQL instance running inside a Docker container.
We assume you already have the container up and running (container ID: `58d592fd188f` in this example).

---

## Prerequisites

* **Docker** must be installed and running on your system.
* A **PostgreSQL container** must already be running.
  Example from our setup:

  * `POSTGRES_USER=games_user`
  * `POSTGRES_PASSWORD=games_pass`
  * `POSTGRES_DB=games_db`

---

## Connecting to the Database

You cannot connect directly from your local machine with `psql` (because of Docker’s isolation).
Instead, you need to **enter the container shell** and use `psql` inside it.

### 1. Enter the container’s shell

Run the following command from your terminal:

```sh
docker exec -it 58d592fd188f /bin/bash
```

Your prompt should change to:

```
root@58d592fd188f:/#
```

---

### 2. Connect to PostgreSQL

Inside the container, connect to the database:

```sh
psql -U games_user -d postgres
```

If successful, your prompt will change to:

```
postgres=>
```

---

## Creating a New Database

1. **Create the database:**

```sql
CREATE DATABASE games;
```

Expected output:

```
CREATE DATABASE
```

2. **Verify the database:**

```sql
\l
```

You should now see the `games` database listed.

3. **Exit psql:**

```sql
\q
```

---

## Connecting Directly to the New Database

Reconnect to the new `games` database:

```sh
psql -U games_user -d games
```

Now your prompt should show:

```
games=>
```

---

## Stopping the Container

When you’re done:

1. Exit the container shell:

```sh
exit
```

2. Stop the running container:

```sh
docker stop 58d592fd188f
```

---

## Notes

* Replace `58d592fd188f` with your actual container ID or name (`games_postgres` if you used `docker-compose`).
* If you want to connect **from another container** (like Prisma or pgAdmin), make sure both are on the same Docker network.

---

Would you like me to also extend this README to include **how to hook this `games_db` into Prisma (Go client)**, so it’s all in one place for your team?




### run this command in root directory
---
``` bash
 npx prisma migrate dev --schema=prisma_games/schema.prisma --name init

 ```bash
 npx prisma generate --schema=prisma_games/schema.prisma


 npx prisma migrate dev --schema=prisma/schema.prisma --name init
  npx prisma generate --schema=prisma/schema.prisma


 


## psql shell open
sudo -u postgres psql
[sudo] password for gauravkrrr: 
psql (18.0 (Ubuntu 18.0-1.pgdg24.04+3), server 17.6 (Ubuntu 17.6-2.pgdg24.04+1))
Type "help" for help.

postgres=# psql
postgres-# /dt
postgres-# ;
ERROR:  syntax error at or near "psql"
LINE 1: psql
        ^
postgres=# /dt;
ERROR:  syntax error at or near "/"
LINE 1: /dt;
        ^
postgres=# dt/;
ERROR:  syntax error at or near "dt"
LINE 1: dt/;
        ^
postgres=# /c crm;
ERROR:  syntax error at or near "/"
LINE 1: /c crm;
        ^
postgres=# \dt
            List of tables
 Schema |   Name   | Type  |  Owner   
--------+----------+-------+----------
 public | students | table | postgres
(1 row)

postgres=# \c crm;
connection to server on socket "/var/run/postgresql/.s.PGSQL.5432" failed: FATAL:  database "crm" does not exist
Previous connection kept
postgres=# \db;
       List of tablespaces
    Name    |  Owner   | Location 
------------+----------+----------
 pg_default | postgres | 
 pg_global  | postgres | 
(2 rows)

postgres=# ALTER DATABASE mydb OWNER TO crm;
ALTER DATABASE
postgres=# /dt;
ERROR:  syntax error at or near "/"
LINE 1: /dt;
        ^
postgres=# \dt;
            List of tables
 Schema |   Name   | Type  |  Owner   
--------+----------+-------+----------
 public | students | table | postgres
(1 row)

postgres=# \db;
       List of tablespaces
    Name    |  Owner   | Location 
------------+----------+----------
 pg_default | postgres | 
 pg_global  | postgres | 
(2 rows)

postgres=# \c crm;
connection to server on socket "/var/run/postgresql/.s.PGSQL.5432" failed: FATAL:  database "crm" does not exist
Previous connection kept
postgres=# \c mydb;
psql (18.0 (Ubuntu 18.0-1.pgdg24.04+3), server 17.6 (Ubuntu 17.6-2.pgdg24.04+1))
You are now connected to database "mydb" as user "postgres".
mydb=# \dt;
               List of tables
 Schema |        Name        | Type  | Owner 
--------+--------------------+-------+-------
 public | Otp                | table | crm
 public | Question           | table | crm
 public | User               | table | crm
 public | _prisma_migrations | table | crm
 public | persons            | table | crm
(5 rows)

mydb=# select * from "Question";
 id | Title | Description | StarterSchema | StarterData | CorrectQuery 
----+-------+-------------+---------------+-------------+--------------
(0 rows)

mydb=# select * from "persons";
 personid | lastname | firstname | address | city 
----------+----------+-----------+---------+------
(0 rows)




DATABASE_URL="postgresql://crm:gk022135@localhost:5432/mydb?schema=public"
JWT_SECRET = "gauravkrrr"

DATABASE_URL_TESTING="postgresql://gagan:gagan@localhost:5436/games?schema=public"
DATABASE_URL_GAMES=""postgresql://games_user:games_pass@localhost:5434/games?schema=public"
