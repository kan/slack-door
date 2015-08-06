CREATE TABLE channels (
    id varchar(10) NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    latest varchar(20),
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL
);

CREATE TABLE users (
    id varchar(10) NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    icon_url TEXT NOT NULL,
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL
);

CREATE TABLE messages (
    id int unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    channel_id varchar(10) NOT NULL,
    user_id varchar(10) NOT NULL,
    text TEXT NOT NULL,
    timestamp datetime NOT NULL,
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL,
    INDEX channel_id (channel_id)
    -- FULLTEXT INDEX (text)
); -- ENGINE=Mroonga;
