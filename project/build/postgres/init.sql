CREATE TABLE jobs
(
    id varchar(36) PRIMARY KEY,
    status text,
    create_date bigint default (EXTRACT(epoch FROM now()))::bigint,
    finish_date bigint default (EXTRACT(epoch FROM now()))::bigint
);
