CREATE SCHEMA IF NOT EXISTS neoway;

CREATE TABLE IF NOT EXISTS neoway.client (
    id                        SERIAL      NOT NULL PRIMARY KEY,
    cpf                       VARCHAR(14) NOT NULL,
    private                   BOOLEAN     NOT NULL DEFAULT false,
    incomplete                BOOLEAN     NOT NULL DEFAULT false,
    last_order_date           DATE,
    average_ticket            DECIMAL(10,2),
    last_order_ticket         DECIMAL(10,2),
    most_frequent_store       VARCHAR(14),
    last_order_store          VARCHAR(14),
    valid_cpf                 BOOLEAN,
    valid_most_frequent_store BOOLEAN,
    valid_last_order_store    BOOLEAN,
    created_at                TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    validated_at              TIMESTAMP
);