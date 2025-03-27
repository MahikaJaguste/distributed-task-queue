use taskqueue_db;

create table tasks (
	id mediumint not null auto_increment primary key ,
    name varchar(20),
    pickedAt timestamp,
    processedAt timestamp,
    completedAt timestamp,
    workerId int,
    status int
);

create table if not exists workerCount (
    id mediumint not null auto_increment primary key ,
    count int
);

create index idx_status ON tasks (status);

insert into workerCount (count) values (1);