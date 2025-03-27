use taskqueue_db;

SET @workerId = 2;

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



select trx_id, trx_state, trx_started, trx_mysql_thread_id, trx_tables_locked, trx_rows_locked, trx_isolation_level FROM information_schema.INNODB_TRX;
select ENGINE_TRANSACTION_ID, THREAD_ID, OBJECT_SCHEMA, OBJECT_NAME, INDEX_NAME, LOCK_TYPE, LOCK_MODE, LOCK_STATUS, LOCK_DATA FROM performance_schema.data_locks;

SELECT * FROM performance_schema.rwlock_instances;
SELECT * FROM performance_schema.table_lock_waits_summary_by_table;
SELECT * FROM performance_schema.metadata_locks;
SELECT * FROM performance_schema.data_lock_waits;
