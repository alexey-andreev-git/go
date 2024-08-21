-- CREATE TABLE IF NOT EXISTS entities_reference (
--     entity_reference_id SERIAL PRIMARY KEY,
--     entity_reference_name VARCHAR(40) NOT NULL,   -- Name of the entity
--     entity_reference_comment VARCHAR(256)         -- Comment or description of the entity
-- );
-- COMMENT ON COLUMN entities_reference.entity_reference_name IS 'entity''s name';
-- COMMENT ON COLUMN entities_reference.entity_reference_comment IS 'entity''s comment';

-- CREATE TABLE IF NOT EXISTS entities (
--     entity_id SERIAL PRIMARY KEY,
--     entity_reference INT NOT NULL REFERENCES entities_reference(entity_reference_id)
-- );
-- COMMENT ON COLUMN entities.entity_id IS 'entity''s id';
-- COMMENT ON COLUMN entities.entity_reference IS 'entity''s reference';

-- CREATE TABLE IF NOT EXISTS entities_data_records (
--   entity_data_record_id SERIAL PRIMARY KEY,
--   entity_data_record_entity INT NOT NULL);
-- COMMENT ON COLUMN entities_data_records.entity_data_record_id IS 'data record''s id';
-- COMMENT ON COLUMN entities_data_records.entity_data_record_entity IS 'data record''s entity';

-- CREATE TABLE IF NOT EXISTS entities_data (
--   entity_data_record INT NOT NULL,
--   entity_data_order INT NOT NULL,
--   entity_data_value INT NOT NULL,
--   PRIMARY KEY (entity_data_record, entity_data_order)
-- );
-- COMMENT ON COLUMN entities_data.entity_data_record IS 'object''s id';
-- COMMENT ON COLUMN entities_data.entity_data_order IS 'data order';
-- COMMENT ON COLUMN entities_data.entity_data_value IS 'data''s value';

-- CREATE TABLE IF NOT EXISTS entities_data_reference (
--   entity_data_reference_entity INT NOT NULL,
--   entity_data_reference_order INT NOT NULL,
--   entity_data_reference_name CHAR(40) NOT NULL,
--   entity_data_reference_type VARCHAR(10) NOT NULL,
--   entity_data_reference_comment VARCHAR(256) NOT NULL,
--   PRIMARY KEY (entity_data_reference_entity, entity_data_reference_order)
-- );
-- COMMENT ON COLUMN entities_data_reference.entity_data_reference_entity IS 'entity reference''s id';
-- COMMENT ON COLUMN entities_data_reference.entity_data_reference_entity IS 'data order in entity';
-- COMMENT ON COLUMN entities_data_reference.entity_data_reference_name IS 'data name';
-- COMMENT ON COLUMN entities_data_reference.entity_data_reference_type IS 'data type';
-- COMMENT ON COLUMN entities_data_reference.entity_data_reference_comment IS 'comment for data';

-- CREATE TABLE IF NOT EXISTS entities_data_val_char (
--   entity_data_val_char_id SERIAL PRIMARY KEY,
--   entity_data_val_char_value varchar(2048) NOT NULL
-- );
-- COMMENT ON COLUMN entities_data_val_char.entity_data_val_char_id IS 'value''s id';
-- COMMENT ON COLUMN entities_data_val_char.entity_data_val_char_value IS 'value';

-- CREATE TABLE IF NOT EXISTS entities_data_val_float (
--   entity_data_val_float_id SERIAL PRIMARY KEY,
--   entity_data_val_float_value float NOT NULL
-- );
-- COMMENT ON COLUMN entities_data_val_float.entity_data_val_float_id IS 'value''s id';
-- COMMENT ON COLUMN entities_data_val_float.entity_data_val_float_value IS 'data''s value';

-- CREATE TABLE IF NOT EXISTS entities_data_val_time (
--   entity_data_val_time_id SERIAL PRIMARY KEY,
--   entity_data_val_time_value timestamp NOT NULL
-- );
-- COMMENT ON COLUMN entities_data_val_time.entity_data_val_time_id IS 'value''s id';
-- COMMENT ON COLUMN entities_data_val_time.entity_data_val_time_value IS 'value';

-- CREATE TABLE IF NOT EXISTS users (
--   id SERIAL PRIMARY KEY,
--   user_name VARCHAR(40) NOT NULL,
--   user_email VARCHAR(256) NOT NULL,
--   user_password VARCHAR(256) NOT NULL,
--   user_role VARCHAR(10) NOT NULL,
--   user_person_id int,
--   UNIQUE(user_email)
-- );

-- CREATE TABLE IF NOT EXISTS addresses
CREATE TABLE addresses(
    id SERIAL NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    PRIMARY KEY(id)
);
ALTER TABLE addresses
ADD COLUMN IF NOT EXISTS address_unit VARCHAR(255) NOT NULL;
ALTER TABLE addresses
ADD COLUMN IF NOT EXISTS address_building VARCHAR(255) NOT NULL;
ALTER TABLE addresses
ADD COLUMN IF NOT EXISTS address_street VARCHAR(255) NOT NULL;
ALTER TABLE addresses
ADD COLUMN IF NOT EXISTS address_city VARCHAR(255) NOT NULL;
ALTER TABLE addresses
ADD COLUMN IF NOT EXISTS address_state VARCHAR(255) NOT NULL;
ALTER TABLE addresses
ADD COLUMN IF NOT EXISTS address_zip_code VARCHAR(20) NOT NULL;
CREATE INDEX idx_addresses_deleted_at ON addresses USING btree ("deleted_at");


-- CREATE TABLE IF NOT EXISTS emails
CREATE TABLE emails(
    id SERIAL NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    email_name varchar(255),
    email_address varchar(255),
    PRIMARY KEY(id)
);
ALTER TABLE emails
ADD COLUMN IF NOT EXISTS email_name VARCHAR(255) NOT NULL;
ALTER TABLE emails
ADD COLUMN IF NOT EXISTS email_address VARCHAR(255) NOT NULL;
CREATE UNIQUE INDEX idx_emails_email ON emails USING btree ("email");
CREATE INDEX idx_emails_deleted_at ON emails USING btree ("deleted_at");

-- CREATE TABLE IF NOT EXISTS people
CREATE TABLE IF NOT EXISTS people (
  id SERIAL PRIMARY KEY NOT NULL,
--   id SERIAL NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
  PRIMARY KEY(id)
);
ALTER TABLE people
ADD COLUMN IF NOT EXISTS person_first_name VARCHAR(40) NOT NULL;
ALTER TABLE people
ADD COLUMN IF NOT EXISTS person_middle_name VARCHAR(40) NOT NULL;
ALTER TABLE people
ADD COLUMN IF NOT EXISTS person_last_name VARCHAR(40) NOT NULL;
ALTER TABLE people
ADD COLUMN IF NOT EXISTS person_birthday VARCHAR(40) NOT NULL;
ALTER TABLE people
ADD COLUMN IF NOT EXISTS person_id VARCHAR(40) NOT NULL;

-- CREATE TABLE IF NOT EXISTS users
CREATE TABLE users(
  id SERIAL NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
);
ALTER TABLE users
ADD COLUMN IF NOT EXISTS user_name varchar(255);
ALTER TABLE users
ADD COLUMN IF NOT EXISTS user_password varchar(255);
-- CREATE INDEX IF NOT EXISTS PRIMARY KEY(id);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users USING btree ("deleted_at");


-- Insert into entities_reference
-- INSERT INTO entities_reference (entity_reference_name, entity_reference_comment)
-- VALUES
--     ('Person', 'Contains personal data'),
--     ('Score', 'Stores float values of scores'),
--     ('Event Time', 'Records timestamp of events'),
--     ('Identifier', 'Contains integer IDs');

-- Insert into entities
-- INSERT INTO entities (entity_reference)
-- VALUES 
--     ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Person')),
--     ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Score')),
--     ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Event Time')),
--     ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Identifier'));

-- Insert into entities_data_reference for 'Person'
-- INSERT INTO entities_data_reference (entity_data_reference_entity, entity_data_reference_order, entity_data_reference_name, entity_data_reference_type, entity_data_reference_comment)
-- VALUES
--     ((SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Person')), 1, 'First Name', 'char', 'First name of the person'),
--     ((SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Person')), 2, 'Last Name', 'char', 'Last name of the person');

-- Insert into entities_data_reference for 'Score'
-- INSERT INTO entities_data_reference (entity_data_reference_entity, entity_data_reference_order, entity_data_reference_name, entity_data_reference_type, entity_data_reference_comment)
-- VALUES
--     ((SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Score')), 1, 'Test Score', 'float', 'Score achieved in a test');

-- Insert into entities_data_reference for 'Event Time'
-- INSERT INTO entities_data_reference (entity_data_reference_entity, entity_data_reference_order, entity_data_reference_name, entity_data_reference_type, entity_data_reference_comment)
-- VALUES
--     ((SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Event Time')), 1, 'Event Start Time', 'timestamp', 'Start time of the event');

-- Insert into entities_data_reference for 'Identifier'
-- INSERT INTO entities_data_reference (entity_data_reference_entity, entity_data_reference_order, entity_data_reference_name, entity_data_reference_type, entity_data_reference_comment)
-- VALUES
--     ((SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Identifier')), 1, 'ID Number', 'int', 'Identifier number for the entity');

-- Insert into entities_data_records
-- INSERT INTO entities_data_records (entity_data_record_entity)
-- VALUES 
--     ((SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Person'))),
--     ((SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Score'))),
--     ((SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Event Time'))),
--     ((SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Identifier')));

-- Insert into entities_data_val_char for 'Person'
-- INSERT INTO entities_data_val_char (entity_data_val_char_value)
-- VALUES ('John'), ('Doe');

-- Link these values to the 'Person' entity in entities_data
-- INSERT INTO entities_data (entity_data_record, entity_data_order, entity_data_value)
-- VALUES 
--     ((SELECT entity_data_record_id FROM entities_data_records WHERE entity_data_record_entity = (SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Person'))), 1, (SELECT entity_data_val_char_id FROM entities_data_val_char WHERE entity_data_val_char_value = 'John')),
--     ((SELECT entity_data_record_id FROM entities_data_records WHERE entity_data_record_entity = (SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Person'))), 2, (SELECT entity_data_val_char_id FROM entities_data_val_char WHERE entity_data_val_char_value = 'Doe'));

-- Insert into entities_data_val_float for 'Score'
-- INSERT INTO entities_data_val_float (entity_data_val_float_value)
-- VALUES (95.75);

-- Link this value to the 'Score' entity in entities_data
-- INSERT INTO entities_data (entity_data_record, entity_data_order, entity_data_value)
-- VALUES 
--     ((SELECT entity_data_record_id FROM entities_data_records WHERE entity_data_record_entity = (SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Score'))), 1, (SELECT entity_data_val_float_id FROM entities_data_val_float WHERE entity_data_val_float_value = 95.75));

-- Insert into entities_data_val_time for 'Event Time'
-- INSERT INTO entities_data_val_time (entity_data_val_time_value)
-- VALUES ('2024-08-09 14:00:00');

-- Link this value to the 'Event Time' entity in entities_data
-- INSERT INTO entities_data (entity_data_record, entity_data_order, entity_data_value)
-- VALUES 
--     ((SELECT entity_data_record_id FROM entities_data_records WHERE entity_data_record_entity = (SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Event Time'))), 1, (SELECT entity_data_val_time_id FROM entities_data_val_time WHERE entity_data_val_time_value = '2024-08-09 14:00:00'));

-- Insert integer ID directly into entities_data for 'Identifier'
-- INSERT INTO entities_data (entity_data_record, entity_data_order, entity_data_value)
-- VALUES 
--     ((SELECT entity_data_record_id FROM entities_data_records WHERE entity_data_record_entity = (SELECT entity_id FROM entities WHERE entity_reference = (SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Identifier'))), 1, 123456);
