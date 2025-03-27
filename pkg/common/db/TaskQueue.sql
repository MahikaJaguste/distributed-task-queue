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

update tasks set pickedAt = NULL, processedAt=NULL, completedAt=null, workerId =null, status=0 where id>-1;

CREATE INDEX idx_status ON tasks (status);

SHOW index in tasks;

-- update tasks set pickedAt = now() WHERE id  in (1,2);
-- update tasks set pickedAt = now() WHERE id in (1,2)

show variables like 'innodb_lock_wait_timeout';

create table if not exists workerCount (
    id mediumint not null auto_increment primary key ,
    count int
);

insert into workerCount (count) values (1);

truncate table tasks;

insert into tasks (name, status) values ("Task 1", 0), ("Task 2", 0), ("Task 3", 0), ("Task 4", 0);


START TRANSACTION;
SELECT @id := id FROM tasks 
WHERE status = 0 
LIMIT 1
FOR UPDATE OF tasks 
SKIP LOCKED;
do sleep(5);
COMMIT;

--

SELECT * FROM information_schema.INNODB_TRX;
SELECT * FROM performance_schema.data_locks;

select trx_id, trx_state, trx_started, trx_mysql_thread_id, trx_tables_locked, trx_rows_locked, trx_isolation_level FROM information_schema.INNODB_TRX;
select ENGINE_TRANSACTION_ID, THREAD_ID, OBJECT_SCHEMA, OBJECT_NAME, INDEX_NAME, LOCK_TYPE, LOCK_MODE, LOCK_STATUS, LOCK_DATA FROM performance_schema.data_locks;

-- SELECT * FROM performance_schema.rwlock_instances;
-- SELECT * FROM performance_schema.table_lock_waits_summary_by_table;
-- SELECT * FROM performance_schema.metadata_locks;
-- SELECT * FROM performance_schema.data_lock_waits;

--

SET @workerId = 1;

START TRANSACTION;
SELECT @id := id FROM tasks 
WHERE status = 0 
LIMIT 1
FOR UPDATE OF tasks 
SKIP LOCKED;
UPDATE tasks 
SET pickedAt = now(), processedAt = now(), workerId = @workerId, status = 1 
WHERE id = @id;
COMMIT;


--