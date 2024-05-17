CREATE TABLE emails
(
    id    int IDENTITY(1,1) PRIMARY KEY,
    email varchar(30) NOT NULL UNIQUE
);