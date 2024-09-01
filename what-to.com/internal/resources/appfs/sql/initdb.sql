CREATE TABLE IF NOT EXISTS entities_reference (
    entity_reference_id SERIAL PRIMARY KEY,
    entity_reference_name VARCHAR(40) NOT NULL,   -- Name of the entity
    entity_reference_comment VARCHAR(256)         -- Comment or description of the entity
);
COMMENT ON COLUMN entities_reference.entity_reference_name IS 'entity''s name';
COMMENT ON COLUMN entities_reference.entity_reference_comment IS 'entity''s comment';

CREATE TABLE IF NOT EXISTS entities (
    entity_id SERIAL PRIMARY KEY,
    entity_reference INT NOT NULL REFERENCES entities_reference(entity_reference_id)
);
COMMENT ON COLUMN entities.entity_id IS 'entity''s id';
COMMENT ON COLUMN entities.entity_reference IS 'entity''s reference';

CREATE TABLE IF NOT EXISTS entities_data (
  entity_data_entity INT NOT NULL,
  entity_data_order INT NOT NULL,
  entity_data_value INT NOT NULL,
  PRIMARY KEY (entity_data_entity, entity_data_order)
);
COMMENT ON COLUMN entities_data.entity_data_entity IS 'object''s id';
COMMENT ON COLUMN entities_data.entity_data_order IS 'data order';
COMMENT ON COLUMN entities_data.entity_data_value IS 'data''s value';

CREATE TABLE IF NOT EXISTS entities_data_reference (
  entity_data_reference_entity_reference INT NOT NULL,
  entity_data_reference_order INT NOT NULL,
  entity_data_reference_name VARCHAR(40) NOT NULL,
  entity_data_reference_type VARCHAR(10) NOT NULL,
  entity_data_reference_comment VARCHAR(256) NOT NULL,
  PRIMARY KEY (entity_data_reference_entity_reference, entity_data_reference_order)
);
COMMENT ON COLUMN entities_data_reference.entity_data_reference_entity_reference IS 'entity reference''s id';
COMMENT ON COLUMN entities_data_reference.entity_data_reference_order IS 'data order in entity';
COMMENT ON COLUMN entities_data_reference.entity_data_reference_name IS 'data name';
COMMENT ON COLUMN entities_data_reference.entity_data_reference_type IS 'data type';
COMMENT ON COLUMN entities_data_reference.entity_data_reference_comment IS 'comment for data';

CREATE TABLE IF NOT EXISTS entities_data_val_char (
  entity_data_val_char_id SERIAL PRIMARY KEY,
  entity_data_val_char_value varchar(2048) NOT NULL
);
COMMENT ON COLUMN entities_data_val_char.entity_data_val_char_id IS 'value''s id';
COMMENT ON COLUMN entities_data_val_char.entity_data_val_char_value IS 'value';

CREATE TABLE IF NOT EXISTS entities_data_val_float (
  entity_data_val_float_id SERIAL PRIMARY KEY,
  entity_data_val_float_value float NOT NULL
);
COMMENT ON COLUMN entities_data_val_float.entity_data_val_float_id IS 'value''s id';
COMMENT ON COLUMN entities_data_val_float.entity_data_val_float_value IS 'data''s value';

CREATE TABLE IF NOT EXISTS entities_data_val_time (
  entity_data_val_time_id SERIAL PRIMARY KEY,
  entity_data_val_time_value timestamp NOT NULL
);
COMMENT ON COLUMN entities_data_val_time.entity_data_val_time_id IS 'value''s id';
COMMENT ON COLUMN entities_data_val_time.entity_data_val_time_value IS 'value';

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  user_name VARCHAR(40) NOT NULL,
  user_email VARCHAR(256) NOT NULL,
  user_password VARCHAR(256) NOT NULL,
  user_role VARCHAR(10) NOT NULL,
  user_person_id int,
  UNIQUE(user_email)
);

-- CREATE TABLE IF NOT EXISTS addresses
-- CREATE TABLE addresses(
--     id SERIAL NOT NULL,
--     created_at timestamp with time zone,
--     updated_at timestamp with time zone,
--     deleted_at timestamp with time zone,
--     PRIMARY KEY(id)
-- );
-- ALTER TABLE addresses
-- ADD COLUMN IF NOT EXISTS address_unit VARCHAR(255) NOT NULL;
-- ALTER TABLE addresses
-- ADD COLUMN IF NOT EXISTS address_building VARCHAR(255) NOT NULL;
-- ALTER TABLE addresses
-- ADD COLUMN IF NOT EXISTS address_street VARCHAR(255) NOT NULL;
-- ALTER TABLE addresses
-- ADD COLUMN IF NOT EXISTS address_city VARCHAR(255) NOT NULL;
-- ALTER TABLE addresses
-- ADD COLUMN IF NOT EXISTS address_state VARCHAR(255) NOT NULL;
-- ALTER TABLE addresses
-- ADD COLUMN IF NOT EXISTS address_zip_code VARCHAR(20) NOT NULL;
-- CREATE INDEX idx_addresses_deleted_at ON addresses USING btree ("deleted_at");


-- -- CREATE TABLE IF NOT EXISTS emails
-- CREATE TABLE emails(
--     id SERIAL NOT NULL,
--     created_at timestamp with time zone,
--     updated_at timestamp with time zone,
--     deleted_at timestamp with time zone,
--     email_name varchar(255),
--     email_address varchar(255),
--     PRIMARY KEY(id)
-- );
-- ALTER TABLE emails
-- ADD COLUMN IF NOT EXISTS email_name VARCHAR(255) NOT NULL;
-- ALTER TABLE emails
-- ADD COLUMN IF NOT EXISTS email_address VARCHAR(255) NOT NULL;
-- CREATE UNIQUE INDEX idx_emails_email ON emails USING btree ("email");
-- CREATE INDEX idx_emails_deleted_at ON emails USING btree ("deleted_at");

-- -- CREATE TABLE IF NOT EXISTS people
-- CREATE TABLE IF NOT EXISTS people (
--   id SERIAL PRIMARY KEY NOT NULL,
-- --   id SERIAL NOT NULL,
--   created_at timestamp with time zone,
--   updated_at timestamp with time zone,
--   deleted_at timestamp with time zone,
--   PRIMARY KEY(id)
-- );
-- ALTER TABLE people
-- ADD COLUMN IF NOT EXISTS person_first_name VARCHAR(40) NOT NULL;
-- ALTER TABLE people
-- ADD COLUMN IF NOT EXISTS person_middle_name VARCHAR(40) NOT NULL;
-- ALTER TABLE people
-- ADD COLUMN IF NOT EXISTS person_last_name VARCHAR(40) NOT NULL;
-- ALTER TABLE people
-- ADD COLUMN IF NOT EXISTS person_birthday VARCHAR(40) NOT NULL;
-- ALTER TABLE people
-- ADD COLUMN IF NOT EXISTS person_id VARCHAR(40) NOT NULL;

-- -- CREATE TABLE IF NOT EXISTS users
-- CREATE TABLE users(
--   id SERIAL NOT NULL,
--   created_at timestamp with time zone,
--   updated_at timestamp with time zone,
--   deleted_at timestamp with time zone,
-- );
-- ALTER TABLE users
-- ADD COLUMN IF NOT EXISTS user_name varchar(255);
-- ALTER TABLE users
-- ADD COLUMN IF NOT EXISTS user_password varchar(255);
-- -- CREATE INDEX IF NOT EXISTS PRIMARY KEY(id);
-- CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users USING btree ("deleted_at");


/*
-- Insert into entities_reference
INSERT INTO entities_reference (entity_reference_name, entity_reference_comment)
VALUES
    ('Person', 'Contains personal data'),
    ('Score', 'Stores float values of scores'),
    ('Event Time', 'Records timestamp of events'),
    ('Identifier', 'Contains integer IDs');

-- Insert into entities_data_reference for 'Person'
INSERT INTO entities_data_reference (entity_data_reference_entity_reference, entity_data_reference_order, entity_data_reference_name, entity_data_reference_type, entity_data_reference_comment)
VALUES
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Person'), 1, 'First Name', 'char', 'First name of the person'),
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Person'), 2, 'Last Name', 'char', 'Last name of the person'),
-- Insert into entities_data_reference for 'Score'
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Score'), 1, 'Test Score', 'float', 'Score achieved in a test'),
-- Insert into entities_data_reference for 'Event Time'
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Event Time'), 1, 'Event Start Time', 'timestamp', 'Start time of the event'),
-- Insert into entities_data_reference for 'Identifier'
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Identifier'), 1, 'ID Number', 'int', 'Identifier number for the entity');

-- Insert into entities
INSERT INTO entities (entity_reference)
VALUES 
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Person')),
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Score')),
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Event Time')),
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Identifier')),
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Person')),
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Score')),
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Event Time')),
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Identifier')),
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Person')),
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Score')),
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Event Time')),
    ((SELECT entity_reference_id FROM entities_reference WHERE entity_reference_name = 'Identifier'));

INSERT INTO entities_data_val_char (entity_data_val_char_value)
VALUES ('John'), ('Doe'), ('John1'), ('Doe1'), ('John2'), ('Doe2');
-- Insert into entities_data_val_float for 'Score'
INSERT INTO entities_data_val_float (entity_data_val_float_value)
VALUES (95.75), (85.45), (75.30);


INSERT INTO entities_data_val_time (entity_data_val_time_value)
VALUES ('2024-08-09 14:00:00'), ('2024-08-09 13:00:00'), ('2024-08-09 12:00:00');

INSERT INTO entities_data (entity_data_entity, entity_data_order, entity_data_value)
VALUES 
-- Link these values to the 'Person' entity in entities_data
    (1, 1, (SELECT entity_data_val_char_id FROM entities_data_val_char WHERE entity_data_val_char_value = 'John')),
    (1, 2, (SELECT entity_data_val_char_id FROM entities_data_val_char WHERE entity_data_val_char_value = 'Doe')),
-- Link this value to the 'Score' entity in entities_data
    (2, 1, (SELECT entity_data_val_float_id FROM entities_data_val_float WHERE entity_data_val_float_value = 95.75)),
-- Link this value to the 'Event Time' entity in entities_data
    (3, 1, (SELECT entity_data_val_time_id FROM entities_data_val_time WHERE entity_data_val_time_value = '2024-08-09 14:00:00')),
-- Insert integer ID directly into entities_data for 'Identifier'
    (4, 1, 123456);

INSERT INTO entities_data (entity_data_entity, entity_data_order, entity_data_value)
VALUES 
-- Link these values to the 'Person' entity in entities_data
    (5, 1, (SELECT entity_data_val_char_id FROM entities_data_val_char WHERE entity_data_val_char_value = 'John1')),
    (5, 2, (SELECT entity_data_val_char_id FROM entities_data_val_char WHERE entity_data_val_char_value = 'Doe1')),
-- Link this value to the 'Score' entity in entities_data
    (6, 1, (SELECT entity_data_val_float_id FROM entities_data_val_float WHERE entity_data_val_float_value = 85.45)),
-- Link this value to the 'Event Time' entity in entities_data
    (7, 1, (SELECT entity_data_val_time_id FROM entities_data_val_time WHERE entity_data_val_time_value = '2024-08-09 13:00:00')),
-- Insert integer ID directly into entities_data for 'Identifier'
    (8, 1, 123457);

INSERT INTO entities_data (entity_data_entity, entity_data_order, entity_data_value)
VALUES 
-- Link these values to the 'Person' entity in entities_data
    (9, 1, (SELECT entity_data_val_char_id FROM entities_data_val_char WHERE entity_data_val_char_value = 'John2')),
    (9, 2, (SELECT entity_data_val_char_id FROM entities_data_val_char WHERE entity_data_val_char_value = 'Doe2')),
-- Link this value to the 'Score' entity in entities_data
    (10, 1, (SELECT entity_data_val_float_id FROM entities_data_val_float WHERE entity_data_val_float_value = 75.30)),
-- Link this value to the 'Event Time' entity in entities_data
    (11, 1, (SELECT entity_data_val_time_id FROM entities_data_val_time WHERE entity_data_val_time_value = '2024-08-09 12:00:00')),
-- Insert integer ID directly into entities_data for 'Identifier'
    (12, 1, 123458);
*/
