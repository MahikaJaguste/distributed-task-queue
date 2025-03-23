use taskqueue_db;

create table tasks (
	id mediumint not null auto_increment primary key ,
    name varchar(20),
    pickedAt timestamp,
    processedAt timestamp,
    completedAt timestamp
);