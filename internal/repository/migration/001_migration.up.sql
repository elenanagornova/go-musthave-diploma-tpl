CREATE SCHEMA IF NOT EXISTS gophermart;
-- DROP SCHEMA gophermart CASCADE ;
--CREATE SCHEMA gophermart;
SET SEARCH_PATH TO gophermart;
-- добавить дату создания и удаления?
CREATE TABLE IF NOT EXISTS users
(
    id              serial primary key,
    user_uid        varchar,
    login           varchar,
    password        varchar
);
CREATE UNIQUE INDEX IF NOT EXISTS users_idx ON users USING btree (id);
CREATE UNIQUE INDEX IF NOT EXISTS users_login_uniq_idx ON users USING btree (login);

CREATE TABLE IF NOT EXISTS orders
(
    id              serial primary key,
    user_uid        varchar,
    order_num       varchar,
    uploaded_at     TIMESTAMP,
    status          varchar,
    accrual         decimal (15,2)
);
CREATE UNIQUE INDEX IF NOT EXISTS orders_idx ON orders USING btree (id);
CREATE UNIQUE INDEX IF NOT EXISTS order_num_idx ON orders USING btree (order_num);


CREATE TABLE IF NOT EXISTS withdrawals
(
    id              serial primary key,
    order_num       varchar,
    sum             decimal (15,2),
    processed_at    TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS withdrawals_idx ON withdrawals USING btree (id);

CREATE TABLE IF NOT EXISTS account_info
(
    id              serial primary key,
    user_uid        varchar,
    balance         decimal (15,2),
    withdrawal      decimal (15,2)
);
CREATE UNIQUE INDEX IF NOT EXISTS account_info_idx ON account_info USING btree (id);
CREATE UNIQUE INDEX IF NOT EXISTS account_uid_uniq_idx ON users USING btree (user_uid);
