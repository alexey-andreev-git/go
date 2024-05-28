CREATE TABLE IF NOT EXISTS entities (
    entity_id SERIAL PRIMARY KEY,
    entity_name VARCHAR(40) NOT NULL,
    entity_comment VARCHAR(256) NOT NULL
);
COMMENT ON COLUMN entities.entity_id IS 'entity''s id';
COMMENT ON COLUMN entities.entity_name IS 'entity''s reference';
COMMENT ON COLUMN entities.entity_comment IS 'entity''s comment';

CREATE TABLE IF NOT EXISTS entities_data (
  entities_data_entity INT NOT NULL,
  entities_data_order INT NOT NULL,
  entities_data_value INT NOT NULL,
  PRIMARY KEY (entities_data_entity, entities_data_order)
);
COMMENT ON COLUMN entities_data.entities_data_entity IS 'object''s id';
COMMENT ON COLUMN entities_data.entities_data_order IS 'data order';
COMMENT ON COLUMN entities_data.entities_data_value IS 'data''s value';

CREATE TABLE IF NOT EXISTS entities_data_reference (
  entities_data_reference_entity INT NOT NULL,
  entities_data_reference_order INT NOT NULL,
  entities_data_reference_name CHAR(40) NOT NULL,
  entities_data_reference_type VARCHAR(10) NOT NULL,
  entities_data_reference_comment VARCHAR(256) NOT NULL,
  PRIMARY KEY (entities_data_reference_entity, entities_data_reference_order)
);
COMMENT ON COLUMN entities_data_reference.entities_data_reference_entity IS 'entity reference''s id';
COMMENT ON COLUMN entities_data_reference.entities_data_reference_entity IS 'data order in entity';
COMMENT ON COLUMN entities_data_reference.entities_data_reference_name IS 'data name';
COMMENT ON COLUMN entities_data_reference.entities_data_reference_type IS 'data type';
COMMENT ON COLUMN entities_data_reference.entities_data_reference_comment IS 'comment for data';

CREATE TABLE IF NOT EXISTS entities_data_val_char (
  entities_data_val_char_id SERIAL PRIMARY KEY,
  entities_data_val_char_value varchar(2048) NOT NULL
);
COMMENT ON COLUMN entities_data_val_char.entities_data_val_char_id IS 'value''s id';
COMMENT ON COLUMN entities_data_val_char.entities_data_val_char_value IS 'value';

CREATE TABLE IF NOT EXISTS entities_data_val_float (
  entities_data_val_float_id SERIAL PRIMARY KEY,
  entities_data_val_float_value float NOT NULL
);
COMMENT ON COLUMN entities_data_val_float.entities_data_val_float_id IS 'value''s id';
COMMENT ON COLUMN entities_data_val_float.entities_data_val_float_value IS 'data''s value';

CREATE TABLE IF NOT EXISTS entities_data_val_float (
  entities_data_val_float_id SERIAL PRIMARY KEY,
  entities_data_val_float_value double NOT NULL
);
COMMENT ON COLUMN entities_data_val_float.entities_data_val_float_id IS 'value''s id';
COMMENT ON COLUMN entities_data_val_float.entities_data_val_float_value IS 'data''s value';

CREATE TABLE IF NOT EXISTS entities_data_val_time (
  entities_data_val_time_id SERIAL PRIMARY KEY,
  entities_data_val_time_value timestamp NOT NULL
);
COMMENT ON COLUMN entities_data_val_time.entities_data_val_time_id IS 'value''s id';
COMMENT ON COLUMN entities_data_val_time.entities_data_val_time_value IS 'value';

-- CREATE TABLE IF NOT EXISTS entities_ref (
--   id SERIAL PRIMARY KEY,
--   entities_ref_name char(40) NOT NULL,
--   entities_ref_comment varchar(256) NOT NULL,
--   UNIQUE(entities_ref_name)
-- );

-- COMMENT ON COLUMN entities_ref.id IS 'reference''s id';
-- COMMENT ON COLUMN entities_ref.entities_ref_name IS 'reference''s name';
-- COMMENT ON COLUMN entities_ref.entities_ref_comment IS 'comment';
