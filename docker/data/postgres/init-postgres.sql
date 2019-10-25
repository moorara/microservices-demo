-- Create database
CREATE DATABASE sensors;
\connect sensors;

-- Create sensor table
CREATE TABLE sensors (
  id varchar(256) PRIMARY KEY,
  site_id varchar(256) NOT NULL,
  name varchar(256) NOT NULL,
  unit varchar(256) NOT NULL,
  min_safe double precision NOT NULL,
  max_safe double precision NOT NULL
);
