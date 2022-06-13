CREATE TABLE tx (
    hash        CHAR(64) NOT NULL,
    "index"     BIGINT NOT NULL,
    height      BIGINT NOT NULL,
    block_time  TIMESTAMP NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT tx_pk PRIMARY KEY (hash)
);

CREATE INDEX tx_height_idx ON tx (height);

CREATE SEQUENCE event_id_seq AS INTEGER;

CREATE TABLE event (
    id          INTEGER NOT NULL DEFAULT nextval('event_id_seq'),
    tx_hash     CHAR(64) NOT NULL,
    "type"      VARCHAR NOT NULL,
    "index"     SMALLINT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT event_pk PRIMARY KEY (id),
    CONSTRAINT event_tx_fk FOREIGN KEY (tx_hash) REFERENCES tx (hash) ON DELETE CASCADE
);

ALTER SEQUENCE event_id_seq OWNED BY event.id;

CREATE INDEX event_type_idx ON event ("type");

CREATE TABLE attribute (
    event_id    INTEGER NOT NULL,
    name        VARCHAR NOT NULL,
    value       JSONB NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT attribute_pk PRIMARY KEY (event_id, name),
    CONSTRAINT attribute_event_fk FOREIGN KEY (event_id) REFERENCES event (id) ON DELETE CASCADE
);

CREATE TABLE raw_tx (
    hash        CHAR(64) NOT NULL,
    data        TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT raw_tx_pk PRIMARY KEY (hash)
);
