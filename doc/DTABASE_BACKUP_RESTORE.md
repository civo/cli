
# Database Backup and Restore

> This repository provides a solution for backing up and restoring databases. Whether you need scheduled backups for automated routine tasks or manual backups for specific events, this guide will walk you through the process.

## Backup Types

### Scheduled Backups

Scheduled backups are automated processes that run at predefined intervals. They help ensure that your data is regularly backed up without manual intervention.

To set up scheduled backups, follow the steps outlined in the Configuration section.

* Support: PostgreSQL

### Manual Backups

Manual backups are initiated by a user on-demand. These are useful for creating backups before system updates, major changes, or any event that requires a snapshot of the database at a specific point in time.

To manually create a backup, follow the steps outlined in the Usage section.

* Support: PostgreSQL

### Restore Process

The restore process allows you to recover your database from a previously created backup.

* Support: PostgreSQL

# CLI Operations

## Create a PostgreSQL database

```bash
❯ civo database create postgres-demo --size g3.db.medium --software PostgreSQL --version 14
Database (postgres-demo) with ID 65dd8173-f754-4c6c-b50a-7ddb6d5446c5 has been created
```

```bash
❯ civo database ls
+--------+---------------+--------------+-------+------------+------------------+--------------+------+--------+
| ID     | Name          | Size         | Nodes | Software   | Software Version | Host         | Port | Status |
+--------+---------------+--------------+-------+------------+------------------+--------------+------+--------+
| 65dd81 | postgres-demo | g3.db.medium |     1 | PostgreSQL |               14 | 31.28.88.149 | 5432 | Ready  |
+--------+---------------+--------------+-------+------------+------------------+--------------+------+--------+
To get the credentials for a database, use `civo db credential <name/ID>`
```

### List database backups

```bash
❯ civo database backups ls postgres-demo
```

## PostgreSQL

> **Note:** PostgreSQL 14 is deprecated and will be removed in a future release. Please migrate to PostgreSQL 17. For migration guidance, see: https://www.civo.com/docs/database/postgresql/migrate-from-14-to-17

### Create Scheduled Backup

```bash
❯ civo database backups create postgres-demo --name every10minutes --schedule "*/10 * * * *"
Database backup (every10minutes) for database postgr-b697-c429d7 has been created
```

### Create Manual Backup

```bash
❯ civo  database backups create  postgres-dem --name firstbackup --type manual
Database backup (firstbackup) for database postgres-demo has been created
```

### List backup

```bash
Scheduled backup
+-------------+---------------+------------+--------------+----------------+------------------+
| Database ID | Database Name | Software   | Schedule     | Backup Name    | Backup           |
+-------------+---------------+------------+--------------+----------------+------------------+
| 65dd81      | postgres-demo | PostgreSQL | */10 * * * * | every10minutes | 20240131-100009F |
+-------------+---------------+------------+--------------+----------------+------------------+
```

### Create Manual Backup

```bash
❯ civo database backups create postgres-demo
```

### List

```bash
❯ civo database backups ls postgres-demo
Scheduled backup
+-------------+---------------+------------+--------------+----------------+------------------+
| Database ID | Database Name | Software   | Schedule     | Backup Name    | Backup           |
+-------------+---------------+------------+--------------+----------------+------------------+
| 65dd81      | postgres-demo | PostgreSQL | */10 * * * * | every10minutes | 20240131-100009F |
+-------------+---------------+------------+--------------+----------------+------------------+
Manual backups
+-------------+---------------+------------+------------------+
| Database ID | Database Name | Software   | Backup           |
+-------------+---------------+------------+------------------+
| 65dd81      | postgres-demo | PostgreSQL | 20240131-095615F |
+-------------+---------------+------------+------------------+
```

### Restore from scheduled

```bash
❯ civo database restore postgres-demo --name restorefromscheduledbackup --backup 20240131-102006F
Warning: Are you sure you want to restore db postgres-demo from 20240131-102006F backup (y/N) ? y
Restoring database postgres-demo from from backup 20240131-102006F
```

### Restore from manual

```bash
❯ civo database restore postgres-demo --name restorefromscheduledbackup --backup 20240131-095615F
Warning: Are you sure you want to restore db postgres-demo from 20240131-095615F backup (y/N) ? y
Restoring database postgres-demo from from backup 20240131-095615F
```

## Point-in-time restore

Besides restoring a backup, a PostgreSQL database can be restored in place to any moment inside its restore window. The window starts at the oldest full backup kept for the database and ends now; `civo database show` prints its start as "Restore Window From". Data written after the chosen moment is lost, so the command asks for confirmation.

```bash
❯ civo database show postgres-demo
...
Storage             : 120 GB (40 GB plan + 80 GB add-on)
Restore Window From : 2026-10-21T00:00:00Z

❯ civo database restore postgres-demo --name restore-before-cleanup --target-time 2026-10-21T10:00:00Z
Warning: Are you sure you want to restore db postgres-demo in place to its state at 2026-10-21T10:00:00Z (data written after that moment is lost) (y/N) ? y
Restoring database postgres-demo to its state at 2026-10-21T10:00:00Z. Data written after that moment is lost.
```

The time is RFC 3339; an offset such as `2026-10-21T12:00:00+02:00` is accepted and converted to UTC.

## Fork a database

A fork is a new, independent database built from the backups of an existing one. By default it reflects the source at its latest backup; `--target-time` forks the source as it was at that moment, within the same restore window as above. Size, node count and network default to the source's.

```bash
❯ civo database fork postgres-demo --name postgres-demo-copy --wait
Database postgres-demo-copy has been forked from postgres-demo (its latest backup) in 4m12s
To get the credentials, run : civo db credential postgres-demo-copy

❯ civo database fork postgres-demo --name postgres-demo-yesterday --target-time 2026-10-21T10:00:00Z --size g3.db.medium --nodes 3
Database postgres-demo-yesterday is being forked from postgres-demo (its state at 2026-10-21T10:00:00Z)
```

## Backup destinations

Scheduled backups go to the Civo Object Store of the database's region unless a destination is given. `--destination-region` sends them to the Civo Object Store of another region, and the `--s3-*` flags send them to any S3-compatible bucket. Both work on `backup create` and `backup update` for the scheduled backup; manual backups always stay in the database's region.

```bash
❯ civo database backup create postgres-demo --name offsite --schedule "0 3 * * *" --destination-region FRA1
❯ civo database backup create postgres-demo --name external --schedule "0 3 * * *" \
    --s3-endpoint https://s3.example.com --s3-bucket my-backups --s3-region eu-1 \
    --s3-access-key KEY --s3-secret-key SECRET --s3-path-prefix postgres-demo/
```

## Storage, extensions and parameters

The storage of every node can be grown beyond the plan (never shrunk; the extra GB are billed as an add-on), PostgreSQL extensions can be enabled, and server parameters can be set. All three go through `civo database update`.

```bash
❯ civo database extension list
+----------+---------------------------+----------+------------------+
| Name     | Description               | Versions | Requires Restart |
+----------+---------------------------+----------+------------------+
| pg_cron  | Job scheduler for Postgres| 1.6      | Yes              |
| postgis  | Geospatial types          | 3.4      | No               |
+----------+---------------------------+----------+------------------+

❯ civo database update postgres-demo --storage 120
❯ civo database update postgres-demo --extensions pg_cron,postgis
❯ civo database update postgres-demo --parameter max_connections=200 --parameter shared_buffers=1GB
```

`--extensions` and `--parameter` replace the current set; pass `--extensions ""` to disable every extension. Extensions that require a restart trigger a brief rolling restart.

## MySQL Backup (Deprecated)

### Create

```bash
❯ civo database backups create mysql-demo --name firstbackup --type manual
Database backup (firstbackup) for database mysql-demo has been created
```

### List backup

```bash
❯ civo database backups ls mysql-demo
Manual backups
+-------------+---------------+-----------+-------------+----------+--------+
| Database ID | Database Name | Backup ID | Backup Name | Software | Status |
+-------------+---------------+-----------+-------------+----------+--------+
| 0d328d      | mysql-demo    | ba0466    | firstbackup | MySQL    | READY  |
+-------------+---------------+-----------+-------------+----------+--------+
```

### Restore

```bash
❯ civo database restore mysql-demo --name restorefirstbackup --backup firstbackup
Warning: Are you sure you want to restore db mysql-demo from firstbackup backup (y/N) ? y
Restoring database mysql-demo from from backup firstbackup
```
