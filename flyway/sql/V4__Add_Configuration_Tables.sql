SET search_path TO template_service, public;
CREATE TABLE template_service.configuration
(
    id              SERIAL PRIMARY KEY,
    config_key      VARCHAR(100)                           NOT NULL UNIQUE,
    config_value    TEXT,
    description     TEXT,
    is_encrypted    BOOLEAN                  DEFAULT FALSE NOT NULL,
    last_updated_by VARCHAR(100),
    last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE template_service.template_config
(
    id           SERIAL PRIMARY KEY,
    template_id  UUID         NOT NULL REFERENCES template_service.template (id) ON DELETE CASCADE,
    config_key   VARCHAR(100) NOT NULL,
    config_value TEXT,
    description  TEXT,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE template_service.template_config
    ADD CONSTRAINT uk_template_id_config_key UNIQUE (template_id, config_key);
CREATE TABLE template_service.rendering_engine
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100)                          NOT NULL UNIQUE,
    description TEXT,
    engine_type VARCHAR(50)                           NOT NULL,
    config      JSONB,
    is_active   BOOLEAN                  DEFAULT TRUE NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE template_service.template_engine_mapping
(
    template_id UUID    NOT NULL REFERENCES template_service.template (id) ON DELETE CASCADE,
    engine_id   INTEGER NOT NULL REFERENCES template_service.rendering_engine (id) ON DELETE CASCADE,
    PRIMARY KEY (template_id, engine_id)
);
CREATE TRIGGER configuration_audit
    AFTER INSERT OR UPDATE OR DELETE
    ON template_service.configuration
    FOR EACH ROW
EXECUTE FUNCTION audit.log_change();

CREATE TRIGGER template_config_audit
    AFTER INSERT OR UPDATE OR DELETE
    ON template_service.template_config
    FOR EACH ROW
EXECUTE FUNCTION audit.log_change();
INSERT INTO template_service.configuration (config_key, config_value, description)
VALUES ('default_template_format', 'html', 'Default format for new templates'),
       ('max_template_size_kb', '512', 'Maximum template size in kilobytes'),
       ('enable_template_caching', 'true', 'Whether to cache rendered templates'),
       ('default_rendering_engine', 'freemarker', 'Default template rendering engine');
INSERT INTO template_service.rendering_engine (name, description, engine_type, config, is_active)
VALUES ('FreeMarker',
        'Apache FreeMarker template engine',
        'freemarker',
        '{
            "version": "2.3.31",
            "settings": {
                "locale": "en_US"
            }
        }',
        true),
       ('Velocity',
        'Apache Velocity template engine',
        'velocity',
        '{
            "version": "2.3",
            "settings": {
                "strict_mode": true
            }
        }',
        true),
       ('Thymeleaf',
        'Thymeleaf template engine',
        'thymeleaf',
        '{
            "version": "3.0.15",
            "settings": {
                "cache_ttl_ms": 3600000
            }
        }',
        true);
INSERT INTO template_service.system_info (version, description)
VALUES ('1.0.3', 'Added configuration tables and sample data');
