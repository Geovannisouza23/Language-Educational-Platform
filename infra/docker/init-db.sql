-- Create databases for each service
CREATE DATABASE auth_db;
CREATE DATABASE users_db;
CREATE DATABASE courses_db;
CREATE DATABASE tasks_db;
CREATE DATABASE progress_db;
CREATE DATABASE notifications_db;
CREATE DATABASE files_db;
CREATE DATABASE videos_db;
CREATE DATABASE tasks_db;
CREATE DATABASE progress_db;

-- Create users with appropriate permissions
CREATE USER auth_user WITH PASSWORD 'auth_password';
CREATE USER users_user WITH PASSWORD 'users_password';
CREATE USER courses_user WITH PASSWORD 'courses_password';
CREATE USER tasks_user WITH PASSWORD 'tasks_password';
CREATE USER progress_user WITH PASSWORD 'progress_password';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE auth_db TO auth_user;
GRANT ALL PRIVILEGES ON DATABASE users_db TO users_user;
GRANT ALL PRIVILEGES ON DATABASE courses_db TO courses_user;
GRANT ALL PRIVILEGES ON DATABASE tasks_db TO tasks_user;
GRANT ALL PRIVILEGES ON DATABASE progress_db TO progress_user;

-- Enable UUID extension for all databases
\c auth_db;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c users_db;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c courses_db;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c tasks_db;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c progress_db;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
