# Database Setup dengan Docker

Dokumentasi ini menjelaskan cara menjalankan database menggunakan Docker beserta konfigurasi replikanya.

---

## Prasyarat

- [Docker](https://docs.docker.com/get-docker/) versi 20+
- [Docker Compose](https://docs.docker.com/compose/install/) versi 2+

---

## Struktur

```
.
├── docker-compose.yml        # Konfigurasi utama Docker Compose
├── docker-compose.replica.yml # Konfigurasi tambahan untuk replica
├── .env                      # Environment variables
└── init/
    └── init.sql              # Script inisialisasi database (opsional)
```

---

## Menjalankan Database

### 1. Salin file environment

```bash
cp .env.example .env
```

Sesuaikan variabel berikut di file `.env`:

```env
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=your_database

# Replica
REPLICA_USER=replica_user
REPLICA_PASSWORD=replica_password
```

### 2. Jalankan database utama (Primary)

```bash
docker compose up -d db
```

### 3. Jalankan beserta replica

```bash
docker compose -f docker-compose.yml -f docker-compose.replica.yml up -d
```

---

## Konfigurasi docker-compose.yml

```yaml
services:
  db:
    image: postgres:16
    container_name: db_primary
    restart: always
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    ports:
      - "5432:5432"
    volumes:
      - db_primary_data:/var/lib/postgresql/data
      - ./init:/docker-entrypoint-initdb.d
    command: >
      postgres
        -c wal_level=replica
        -c max_wal_senders=10
        -c wal_keep_size=512

volumes:
  db_primary_data:
```

---

## Konfigurasi Replica (docker-compose.replica.yml)

```yaml
services:
  db_replica:
    image: postgres:16
    container_name: db_replica
    restart: always
    environment:
      PGUSER: ${REPLICA_USER}
      PGPASSWORD: ${REPLICA_PASSWORD}
    ports:
      - "5433:5432"
    volumes:
      - db_replica_data:/var/lib/postgresql/data
    depends_on:
      - db
    command: >
      bash -c "
        until pg_basebackup -h db -D /var/lib/postgresql/data -U ${REPLICA_USER} -Fp -Xs -R -P; do
          echo 'Menunggu primary siap...';
          sleep 2;
        done
      "

volumes:
  db_replica_data:
```

> **Catatan:** Replica berjalan dalam mode _standby_ (read-only). Semua write tetap dilakukan ke primary.

---

## Menghubungkan Aplikasi

| Koneksi | Host        | Port   | Keterangan   |
| ------- | ----------- | ------ | ------------ |
| Primary | `localhost` | `5432` | Read & Write |
| Replica | `localhost` | `5433` | Read only    |

Contoh connection string:

```
# Primary
postgresql://your_user:your_password@localhost:5432/your_database

# Replica
postgresql://your_user:your_password@localhost:5433/your_database
```

---

## Perintah Umum

```bash
# Lihat status container
docker compose ps

# Lihat log primary
docker compose logs -f db

# Lihat log replica
docker compose logs -f db_replica

# Stop semua container
docker compose down

# Stop dan hapus volume (reset data)
docker compose down -v
```

---

## Cek Status Replikasi

Masuk ke container primary dan jalankan query berikut:

```bash
docker exec -it db_primary psql -U your_user -d your_database
```

```sql
SELECT client_addr, state, sent_lsn, write_lsn, flush_lsn, replay_lsn
FROM pg_stat_replication;
```

Jika replica terhubung dengan benar, akan muncul satu baris dengan `state = streaming`.

---

## Troubleshooting

**Replica tidak bisa terhubung ke primary**

- Pastikan `wal_level=replica` sudah di-set di primary.
- Cek apakah user replica punya hak `REPLICATION`:
  ```sql
  CREATE USER replica_user REPLICATION LOGIN ENCRYPTED PASSWORD 'replica_password';
  ```

**Port konflik**

- Ganti port di `docker-compose.yml` jika `5432` atau `5433` sudah dipakai proses lain.

**Data replica tidak sinkron**

- Stop replica, hapus volume-nya, lalu jalankan ulang agar `pg_basebackup` berjalan dari awal:
  ```bash
  docker compose stop db_replica
  docker volume rm <nama_volume_replica>
  docker compose up -d db_replica
  ```
