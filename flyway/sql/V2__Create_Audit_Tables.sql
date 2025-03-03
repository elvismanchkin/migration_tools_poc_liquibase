SET search_path TO template_service, public;
CREATE SCHEMA IF NOT EXISTS audit;
CREATE TABLE audit.audit_log
(
    id          SERIAL PRIMARY KEY,
    entity_type VARCHAR(100) NOT NULL,
    entity_id   VARCHAR(100) NOT NULL,
    action      VARCHAR(50)  NOT NULL,
    user_id     VARCHAR(100) NOT NULL,
    change_data JSONB,
    timestamp   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    client_ip   VARCHAR(50),
    user_agent  TEXT
);
CREATE INDEX idx_audit_entity ON audit.audit_log (entity_type, entity_id);
CREATE INDEX idx_audit_timestamp ON audit.audit_log (timestamp);
DROP FUNCTION IF EXISTS audit.log_change();

CREATE OR REPLACE FUNCTION audit.log_change()
    RETURNS TRIGGER
    LANGUAGE plpgsql AS
'BEGIN
    INSERT INTO audit.audit_log (entity_type,
                                 entity_id,
                                 action,
                                 user_id,
                                 change_data)
    VALUES (TG_TABLE_NAME,
            NEW.id::text,
            TG_OP,
            COALESCE(current_setting(''app.current_user_id'', true), ''system''),
            jsonb_build_object(
                    ''old'', to_jsonb(OLD),
                    ''new'', to_jsonb(NEW)
            ));
    RETURN NEW;
END;';
INSERT INTO template_service.system_info (version, description)
VALUES ('1.0.1', 'Added audit tables and functions');
