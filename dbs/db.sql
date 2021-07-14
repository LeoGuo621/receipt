drop database if exists receipt;
create database receipt character set utf8;
use receipt;

create table t_user (
    user_id int not null primary key auto_increment,
    email varchar(50),
    name varchar(30),
    password varchar(100),
    address varchar(256)
);

create table t_content (
    content_id int not null primary key auto_increment,
    title varchar(100),
    content varchar(256),
    content_hash varchar(100),
    address varchar(100),
    token_id bigint,
    time_stamp timestamp
);