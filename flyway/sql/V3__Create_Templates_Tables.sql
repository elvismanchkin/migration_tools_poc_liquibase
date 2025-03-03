SET search_path TO template_service, public;
CREATE TABLE template_service.template_category
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE template_service.template
(
    id          UUID                     DEFAULT uuid_generate_v4() PRIMARY KEY,
    name        VARCHAR(200)                            NOT NULL,
    category_id INTEGER REFERENCES template_service.template_category (id),
    content     TEXT                                    NOT NULL,
    format      VARCHAR(50)              DEFAULT 'html' NOT NULL,
    version     INTEGER                  DEFAULT 1      NOT NULL,
    is_active   BOOLEAN                  DEFAULT TRUE   NOT NULL,
    created_by  VARCHAR(100)                            NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(100),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE template_service.template_version
(
    id           SERIAL PRIMARY KEY,
    template_id  UUID         NOT NULL REFERENCES template_service.template (id) ON DELETE CASCADE,
    version      INTEGER      NOT NULL,
    content      TEXT         NOT NULL,
    format       VARCHAR(50)  NOT NULL,
    created_by   VARCHAR(100) NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    change_notes TEXT
);
ALTER TABLE template_service.template_version
    ADD CONSTRAINT uk_template_id_version UNIQUE (template_id, version);
CREATE TABLE template_service.template_variable
(
    id            SERIAL PRIMARY KEY,
    template_id   UUID                                      NOT NULL REFERENCES template_service.template (id) ON DELETE CASCADE,
    variable_name VARCHAR(100)                              NOT NULL,
    description   TEXT,
    default_value TEXT,
    is_required   BOOLEAN                  DEFAULT FALSE    NOT NULL,
    variable_type VARCHAR(50)              DEFAULT 'string' NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE template_service.template_variable
    ADD CONSTRAINT uk_template_id_variable_name UNIQUE (template_id, variable_name);
CREATE TRIGGER template_audit
    AFTER INSERT OR UPDATE OR DELETE
    ON template_service.template
    FOR EACH ROW
EXECUTE FUNCTION audit.log_change();

CREATE TRIGGER template_version_audit
    AFTER INSERT OR UPDATE OR DELETE
    ON template_service.template_version
    FOR EACH ROW
EXECUTE FUNCTION audit.log_change();
INSERT INTO template_service.template_category (name, description)
VALUES ('Email', 'Email templates'),
       ('Notification', 'System notification templates'),
       ('Report', 'Report templates');

INSERT INTO template_service.system_info (version, description)
VALUES ('1.0.2', 'Added template tables and sample data');
