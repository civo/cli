
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

## Create and PostgreSQL database

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
