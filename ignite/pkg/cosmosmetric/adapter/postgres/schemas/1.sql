BEGIN;

CREATE TABLE schema (
    version     SMALLINT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT schema_pk PRIMARY KEY (version)
);

INSERT INTO schema (version) VALUES (1);

CREATE TABLE tx (
    hash        CHAR(64) NOT NULL,
    height      BIGINT NOT NULL,
    index       BIGINT NOT NULL,
    block_time  TIMESTAMP NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT tx_pk PRIMARY KEY (hash)
);

CREATE INDEX tx_height_idx ON tx (height);

CREATE TABLE attribute (
    tx_hash     CHAR(64) NOT NULL,
    event_type  VARCHAR NOT NULL,
    event_index SMALLINT NOT NULL,
    name        VARCHAR NOT NULL,
    value       JSONB NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT attribute_tx_fk FOREIGN KEY (tx_hash) REFERENCES tx (hash) ON DELETE CASCADE
);

CREATE INDEX attribute_idx ON attribute (event_type, name);
CREATE INDEX attribute_event_idx ON attribute (event_type);

COMMIT;
