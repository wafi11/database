Ini adalah draft dokumentasi profesional dalam format Markdown (`.md`) untuk repositori kamu. Dokumentasi ini dirancang dengan standar industri agar siapa pun (termasuk interviewer kamu) bisa melihat betapa rapinya cara kamu bekerja.

Kamu bisa menyimpan file ini dengan nama `README.md` di root folder proyek kamu.

---

````markdown
# PostgreSQL Streaming Replication (Primary-Replica) with Docker

Repositori ini berisi implementasi **PostgreSQL 14 High Availability** menggunakan arsitektur _Physical Streaming Replication_ (Asynchronous). Proyek ini dirancang untuk mensimulasikan lingkungan produksi di mana beban baca (_Read_) dipisahkan dari beban tulis (_Write_).

## 🏗️ Arsitektur Sistem

- **Primary (db_main):** Bertugas menangani semua transaksi tulis (INSERT, UPDATE, DELETE). Mencatat perubahan ke dalam Write Ahead Log (WAL).
- **Replica (db_replica):** Bertugas menangani beban baca (SELECT). Melakukan sinkronisasi data secara real-time dari Primary.

---

## 🚀 Cara Menjalankan

### 1. Persiapan

Pastikan folder data database dalam keadaan bersih (jika ini instalasi baru):

```bash
rm -rf ./db ./db_replica_data
```
````

### 2. Jalankan Container

Gunakan Docker Compose untuk membangun seluruh infrastruktur:

```bash
docker-compose up -d

```

### 3. Otomatisasi Inisialisasi

Sistem menggunakan script `init-replication.sh` yang secara otomatis:

1. Membuat `replica_user` dengan izin khusus replikasi.
2. Mendaftarkan alamat IP Replica ke dalam `pg_hba.conf` (Host-Based Authentication).
3. Melakukan `pg_basebackup` pada container Replica saat pertama kali dijalankan.

---

## 📊 Monitoring Replikasi

Gunakan query berikut untuk memastikan status replikasi berjalan dengan baik:

### Di Sisi Primary

Jalankan di container `db_main` untuk melihat status streaming:

```sql
SELECT
    application_name,
    state,
    sent_lsn,
    write_lsn,
    flush_lsn
FROM pg_stat_replication;

```

_Pastikan `state` bernilai `streaming`._

### Di Sisi Replica

Cek jeda waktu (_replication lag_) terhadap data asli:

```sql
SELECT now() - pg_last_xact_replay_timestamp() AS replication_lag;

```

---

## 🛠️ Troubleshoot & DBA Knowledge

### Konfigurasi Penting (Primary)

- `wal_level = replica`: Memastikan log cukup detail untuk replikasi.
- `max_wal_senders`: Jumlah maksimum koneksi dari replica.
- `hot_standby = on`: Memungkinkan replica menerima query baca meskipun dalam mode recovery.

### Mekanisme Failover (Manual)

Jika Primary mati, promosikan Replica menjadi Primary baru dengan perintah:

```bash
docker exec -it db_replica pg_ctl promote

```

---

## 👨‍💻 Tech Stack

- **Database:** PostgreSQL 14
- **Orchestration:** Docker Compose
- **Scripting:** Bash (Auto-provisioning)

---

```

### Tips Tambahan buat Kamu:
Kalau kamu mau upload ke GitHub, pastikan file `init-replication.sh` tadi juga di-upload dan diberi izin eksekusi (`chmod +x init-replication.sh`).

Dokumentasi ini menunjukkan kalau kamu bukan cuma "bisa bikin jalan", tapi kamu paham **apa** yang kamu bikin. Ini nilai plus banget buat posisi DBA atau Backend Engineer.

**Apakah ada bagian teknis lain yang ingin kamu tambahkan ke dokumentasinya? Mungkin bagian cara koneksi dari C#?**

```
