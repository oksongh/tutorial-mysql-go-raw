create database if not exists test;

create table if not exists user(
    `id` int not null auto_increment,
    `name` varchar(50) not null,
    `stmt` varchar(255) not null,
    `otaku_id` int not null,
    primary key (id)
);
show columns from user;

insert into user (name,stmt,otaku_id) values ('kaniman','yolo','2'), ('kaiman','yolo','1');
select * from user;
