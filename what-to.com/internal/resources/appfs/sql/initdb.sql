CREATE TABLE IF NOT EXISTS entities (
    id SERIAL PRIMARY KEY,
    reference INT NOT NULL,
    comment VARCHAR(256) NOT NULL
);

COMMENT ON COLUMN entities.id IS 'entity''s id';
COMMENT ON COLUMN entities.reference IS 'entity''s reference';
COMMENT ON COLUMN entities.comment IS 'entity''s comment';

CREATE TABLE IF NOT EXISTS entities_data (
  entity INT NOT NULL,
  entity_order INT NOT NULL,
  value INT NOT NULL,
  PRIMARY KEY (entity, entity_order)
);

COMMENT ON COLUMN entities_data.entity IS 'object''s id';
COMMENT ON COLUMN entities_data.entity_order IS 'data order';
COMMENT ON COLUMN entities_data.value IS 'data''s value';

CREATE TABLE IF NOT EXISTS entities_data_ref (
  entity_ref INT NOT NULL,
  entity_data_ref_order INT NOT NULL,
  name CHAR(40) NOT NULL,
  type VARCHAR(10) NOT NULL,
  comment VARCHAR(256) NOT NULL,
  PRIMARY KEY (entity_ref, entity_data_ref_order)
);

COMMENT ON COLUMN entities_data_ref.entity_ref IS 'entity reference''s id';
COMMENT ON COLUMN entities_data_ref.entity_data_ref_order IS 'data order in entity';
COMMENT ON COLUMN entities_data_ref.name IS 'data name';
COMMENT ON COLUMN entities_data_ref.type IS 'data type';
COMMENT ON COLUMN entities_data_ref.comment IS 'comment for data';

CREATE TABLE IF NOT EXISTS entities_data_val_char (
  id SERIAL PRIMARY KEY,
  value varchar(2048) NOT NULL
);

COMMENT ON COLUMN entities_data_val_char.id IS 'value''s id';
COMMENT ON COLUMN entities_data_val_char.value IS 'value';

CREATE TABLE IF NOT EXISTS entities_data_val_float (
  id SERIAL PRIMARY KEY,
  value float NOT NULL
);


COMMENT ON COLUMN entities_data_val_float.id IS 'value''s id';
COMMENT ON COLUMN entities_data_val_float.value IS 'data''s value';

CREATE TABLE IF NOT EXISTS entities_data_val_float (
  id SERIAL PRIMARY KEY,
  value double NOT NULL
);

COMMENT ON COLUMN entities_data_val_float.id IS 'value''s id';
COMMENT ON COLUMN entities_data_val_float.value IS 'data''s value';

CREATE TABLE IF NOT EXISTS entities_data_val_time (
  id SERIAL PRIMARY KEY,
  value timestamp NOT NULL
);

COMMENT ON COLUMN entities_data_val_time.id IS 'value''s id';
COMMENT ON COLUMN entities_data_val_time.value IS 'value';

CREATE TABLE IF NOT EXISTS entities_ref (
  id SERIAL PRIMARY KEY,
  entities_ref_name char(40) NOT NULL,
  entities_ref_comment varchar(256) NOT NULL,
  UNIQUE(entities_ref_name)
);

COMMENT ON COLUMN entities_ref.id IS 'reference''s id';
COMMENT ON COLUMN entities_ref.entities_ref_name IS 'reference''s name';
COMMENT ON COLUMN entities_ref.entities_ref_comment IS 'comment';
