SELECT 
    application_name, 
    state,            -- checking streaming or not
    sent_lsn,         -- posisition logs latest retreived
    write_lsn,        -- Posisi log terakhir yang ditulis di Replica
    flush_lsn         -- Posisi log terakhir yang sudah aman di disk Replica
FROM pg_stat_replication;