#!/bin/bash
set -e

# 1. Tambahkan izin replikasi ke pg_hba.conf
echo "host replication replica_user 0.0.0.0/0 md5" >> "$PGDATA/pg_hba.conf"

# 2. Buat user replikasi via SQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE ROLE replica_user WITH REPLICATION LOGIN PASSWORD 'rahasia';
EOSQL