CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SCHEMA IF NOT EXISTS template_service;
SET search_path TO template_service, public;
CREATE TABLE template_service.system_info
(
    id          SERIAL PRIMARY KEY,
    version     VARCHAR(50) NOT NULL,
    description TEXT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO template_service.system_info (version, description)
VALUES ('1.0.0', 'Initial schema');
