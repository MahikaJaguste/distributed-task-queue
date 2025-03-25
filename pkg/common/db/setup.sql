use taskqueue_db;

create table tasks (
	id mediumint not null auto_increment primary key ,
    name varchar(20),
    pickedAt timestamp,
    processedAt timestamp,
    completedAt timestamp,
    workerId mediuminint,
    status int
);

create table if not exists workerCount (
    id mediumint not null auto_increment primary key ,
    count int
);

insert into workerCount (count) values (1);