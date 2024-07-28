--name: insert-multi
INSERT INTO neoway.client (cpf, private, incomplete, last_order_date, average_ticket, last_order_ticket, most_frequent_store, last_order_store, valid_cpf, valid_most_frequent_store, valid_last_order_store)
VALUES %s;

--name: get-paginated
SELECT id, cpf, private, incomplete, last_order_date, average_ticket, last_order_ticket, most_frequent_store, last_order_store
FROM neoway.client
ORDER BY id
LIMIT $1 OFFSET $2;

--name: update-multi-documents-validity
INSERT INTO neoway.client (id, cpf, valid_cpf, valid_most_frequent_store, valid_last_order_store, validated_at)
VALUES %s
ON CONFLICT (id) DO UPDATE
    SET valid_cpf = EXCLUDED.valid_cpf,
    valid_most_frequent_store = EXCLUDED.valid_most_frequent_store,
    valid_last_order_store = EXCLUDED.valid_last_order_store,
    validated_at = EXCLUDED.validated_at;
