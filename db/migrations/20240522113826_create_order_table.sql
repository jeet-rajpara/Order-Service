-- migrate:up

CREATE TABLE IF NOT EXISTS orders (
	id INT NOT NULL DEFAULT unique_rowid(),
	cus_name STRING NOT NULL,
	cus_email STRING NOT NULL,
    items JSONB NOT NULL,
    "status" STRING NOT NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC)
);

-- migrate:down

DROP TABLE IF EXISTS order;