CREATE TABLE IF NOT EXISTS entities (
    id SERIAL PRIMARY KEY COMMENT 'entity''s id',
    reference INT NOT NULL COMMENT 'entity''s reference',
    comment VARCHAR(256) NOT NULL COMMENT 'entity''s reference'
);

CREATE TABLE IF NOT EXISTS entities_data (
    entity INT NOT NULL COMMENT 'object''s id',
    order INT NOT NULL COMMENT 'data order',
    value INT NOT NULL COMMENT 'data''s value',
    PRIMARY KEY (entity, order)
);

CREATE TABLE IF NOT EXISTS entities_data_ref (
    entity_ref INT NOT NULL COMMENT 'entity reference''s id',
    order INT NOT NULL COMMENT 'data order in entity',
    name CHAR(40) NOT NULL COMMENT 'data name',
    type VARCHAR(10) NOT NULL COMMENT 'data type',
    comment VARCHAR(256) NOT NULL COMMENT 'comment for data',
    PRIMARY KEY (entity_ref, order)
) COMMENT 'Entity''s data reference';

CREATE TABLE `objects_data_val_char` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'value''s id',
  `value` varchar(2048) NOT NULL COMMENT 'value',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COMMENT='Value for char object''s data';
CREATE TABLE `objects_data_val_float` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'value''s id',
  `value` double NOT NULL COMMENT 'data''s value',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Values for float object''s data';
CREATE TABLE `objects_data_val_time` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'value''s id',
  `value` datetime NOT NULL COMMENT 'value',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Time object''s data values';

CREATE TABLE `objects_ref` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'reference''s id',
  `name` char(40) NOT NULL COMMENT 'reference''s name',
  `comment` varchar(256) NOT NULL COMMENT 'comment',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
