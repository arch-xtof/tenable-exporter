CREATE DATABASE tenable;

CREATE USER 'tenable'@'%' IDENTIFIED BY 'tenable';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP ON tenable.* TO 'tenable'@'%';
FLUSH PRIVILEGES;

USE tenable;

CREATE TABLE Asset (
    uuid varchar(63) NOT NULL,
    hostname varchar(63) NOT NULL,
    asset_group varchar(31),
    PRIMARY KEY (uuid)
);

CREATE TABLE Vuln (
    id int NOT NULL AUTO_INCREMENT,
    name varchar(255),
    asset_uuid varchar(63) NOT NULL,
    severity varchar(15) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (asset_uuid) REFERENCES Asset(uuid) ON DELETE CASCADE
);
