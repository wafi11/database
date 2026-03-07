SELECT
    application_name,
    state,
    sent_lsn,
    write_lsn,
    flush_lsn
FROM pg_stat_replication;