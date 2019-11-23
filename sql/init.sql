create table if not exists pbl_test.user (
    id         int          not null primary key,
    name       varchar(50)  not null,
    password   varchar(100) not null,
    created_at datetime,
    updated_at datetime,
    deleted_at datetime
)ENGINE=INNODB DEFAULT CHARSET=utf8;

set time_zone = '+09:00';

insert into pbl_test.user values (1, 'aaa', 'test', now(), now(), now());
insert into pbl_test.user values (2, 'bbb', 'pass', now(), now(), now());
