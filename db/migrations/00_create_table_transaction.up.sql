CREATE TABLE "transaction" (
	txHash bytea,
	fromIndex NUMERIC(78),
	toIndex NUMERIC(78),
	amount NUMERIC(78),
	fee NUMERIC(78),
	nonce NUMERIC(78),
	signature bytea
);
