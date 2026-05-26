# edu-planner
This is a platform designed for educational planning and administration.

## Description
Edu-planner is an application that handles scheduling, curriculum management, and user administration for educational environments. It features a main service to manage educational data, an administrative dashboard for user management, and incorporates a robust database structure to handle institutional planning efficiently.

## Setup

### 1. Clone the Repository
First, download the project files by running the following command:

```bash
git clone https://github.com/Lapakin/edu-planner.git

```

### 2. Configure Environment Variables

Before launching the application, create a `.env` file in the project root to store your database credentials and default administrator details. Use the exact format below:

```env
DB_USER=dev
DB_PASSWORD=12345

ADMIN_EMAIL=admin@admin.com
ADMIN_PASSWORD=admin123

ADMIN_FIRST_NAME=Admin
ADMIN_LAST_NAME=User

```

### 3. Launch the Application

Start the services using Docker via the provided Makefile:

```bash
make up

```

## Services

| Service | Endpoint |
| --- | --- |
| Web Application / API | `http://localhost:8080` |

## Testing

To execute the full suite of unit tests, use the following command:

```bash
make go-test

```
