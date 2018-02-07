-- Create database
CREATE DATABASE microservices;
\connect microservices;

-- Create sensor table
CREATE TABLE votes (
  id character varying(256) PRIMARY KEY,
  link_id character varying(256) NOT NULL,
  stars int
);
