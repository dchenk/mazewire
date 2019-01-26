# Database Schema

This file describes the database schema used.

Some data types, such as `BYTES`, are named differently in different DBMSs. The types and functions
used here use the CockroachDB syntax.

The `SEQUENCE` schemas are necessary only for database systems that require such a counter outside
of the definition of a table; other systems, such as MySQL, have increment features built into the
tables.

None of the columns in any of the tables can have a `NULL` value. Columns that don't have a default
value defined require a value upon insertion of a row.

All foreign key references use the `ON DELETE CASCADE` feature so that there are never child records
remaining after a parent record is deleted.

## Sequence sites_id

The sequence that sets the primary key of rows in the `sites` table.

## Table sites

| Column     | Type       | Default             |
| :--------- | :--------- | :------------------ |
| id         | INT        | nextval('sites_id') |
| domain     | STRING     |                     |
| name       | STRING     |                     |
| logo       | STRING     | (empty string)      |
| favicon    | STRING     | (empty string)      |
| tls        | INT        | 0                   |
| created    | TIMESTAMP  | now()               |

### Indexes for sites

| Name      | Type      | Columns       |
| :-------- | :-------- | :------------ |
| pk_id     | PRIMARY   | id            |
| un_domain | UNIQUE    | domain        |

## Sequence users_id

The sequence that sets the primary key of rows in the `users` table.

## Table users 

| Column     | Type       | Default             |
| :--------- | :--------- | :------------------ |
| id         | INT        | nextval('users_id') |
| username   | STRING     |                     |
| email      | STRING     |                     |
| pass       | BYTES      |                     |
| fname      | STRING     |                     |
| lname      | STRING     |                     |
| updated    | TIMESTAMP  | now()               |
| registered | TIMESTAMP  | now()               |

### Indexes for users

| Name      | Type      | Columns       |
| :-------- | :-------- | :------------ |
| pk_id     | PRIMARY   | id            |
| un_uname  | UNIQUE    | username      |
| un_email  | UNIQUE    | email         |

## Table usermeta

Various settings for and details about users.

| Column     | Type       | Default    |
| :--------- | :--------- | :--------- |
| user_id    | INT        |            |
| k          | STRING     |            |
| v          | BYTES      |            |
| updated    | TIMESTAMP  | now()      |

### Indexes for usermeta

| Name       | Type      | Columns      | References |
| :--------- | :-------- | :----------- | :--------- |
| pk_id_k    | PRIMARY   | user_id, k   |            |
| fk_user_id | FOREIGN   | user_id      | users (id) |

## Sequence content_id

The sequence that sets the primary key of rows in the `content` table.

## Table content

The records in the `content` table represent things that are viewable from a URL; that is, they
have a URL slug that will be part of some URL on the site. A row in this table may represent an
item that only functions as a parent to other items, where the parent slug is used in the URL
along with the child slug.

| Column     | Type       | Default               | Constraints              |
| :--------- | :--------- | :-------------------- | ------------------------ |
| id         | INT        | nextval('content_id') |                          |
| site       | INT        |                       |                          |
| slug       | STRING     |                       | CHECK (length(slug) > 0) |
| author     | INT        |                       |                          |
| type       | STRING     |                       | CHECK (length(type) > 0) |
| parent     | INT        | 0                     |                          |
| title      | STRING     |                       |                          |
| meta_title | STRING     | (empty string)        |                          |
| meta_desc  | STRING     | (empty string)        |                          |
| body       | BYTES      | (empty string)        |                          |
| status     | STRING     | 'draft'               | CHECK (status IN ('draft', 'published', 'unsaved', 'trashed')) |
| updated    | TIMESTAMP  | now()                 |                          |

### Indexes for content

| Name             | Type      | Columns            | References |
| :--------------- | :-------- | :----------------- | :--------- |
| pk_id            | PRIMARY   | id                 |            |
| fk_site_id       | FOREIGN   | site               | sites (id) |
| fk_author_id     | FOREIGN   | author             | users (id) |
| site_slug        | UNIQUE    | site, slug         |            |
| site_type_status | INDEX     | site, type, status |            |

## Sequence blobs_id

The sequence that sets the primary key of rows in the `blobs` table.
Such a sequence is replicated for each website, named "blobs_id##" where "##" is the site's ID.

## Table blobs 

Arbitrary data keyed by 'role' string and/or 'k' integer.
Such a table is replicated for each website, named "blobs##" where "##" is the site's ID.

| Column     | Type       | Default             |
| :--------- | :--------- | :------------------ |
| id         | INT        | nextval('blobs_id') |
| role       | STRING     |                     |
| k          | INT        | 0                   |
| v          | BYTES      |                     |
| updated    | TIMESTAMP  | now()               |

### Indexes for blobs

| Name        | Type      | Columns       |
| :---------- | :-------- | :------------ |
| pk_id       | PRIMARY   | id            |
| indx_role_k | INDEX     | role, k       |

## Table options 

Basic settings that can be set on any site or globally for all sites (with site = 1).

| Column     | Type       | Default  |
| :--------- | :--------- | :------- |
| site       | INT        |          |
| k          | STRING     |          |
| v          | BYTES      |          |

### Indexes for options

| Name        | Type      | Columns   | References |
| :---------- | :-------- | :-------- | :--------- |
| pk_id       | PRIMARY   | site, id  |            |
| fk_site_id  | FOREIGN   | site      | sites (id) |

## Sequence site_messages_id

The sequence that sets the primary key of rows in the `site_messages` table.

## Table site_messages

Alert messages pertaining to sites.

| Column     | Type       | Default                     |
| :--------- | :--------- | :-------------------------- |
| id         | INT        | nextval('site_messages_id') |
| site       | INT        |                             |
| role       | STRING     |                             |
| k          | STRING     |                             |
| v          | STRING     |                             |
| created    | TIMESTAMP  | now()                       |

### Indexes for site_messages

| Name        | Type     | Columns  | References |
| :---------- | :------- | :------- | :--------- |
| pk_id       | PRIMARY  | id       |            |
| fk_site_id  | FOREIGN  | site     | sites (id) |
| indx_role_k | INDEX    | role, k  |            |

## Sequence user_messages_id

The sequence that sets the primary key of rows in the `user_messages` table.

## Table user_messages 

Alert messages pertaining to users.

| Column     | Type       | Default                     |
| :--------- | :--------- | :-------------------------- |
| id         | INT        | nextval('user_messages_id') |
| user_id    | INT        |                             |
| k          | STRING     |                             |
| v          | STRING     |                             |
| created    | TIMESTAMP  | now()                       |

### Indexes for user_messages

| Name        | Type     | Columns  | References |
| :---------- | :------- | :------- | :--------- |
| pk_id       | PRIMARY  | id       |            |
| fk_user_id  | FOREIGN  | user     | users (id) |

## Table media 

Static media files uploaded through and managed by the system.

| Column      | Type       | Default                     |
| :---------- | :--------- | :-------------------------- |
| id          | STRING     |                             |
| ext         | STRING     | (empty string)              |
| site        | INT        |                             |
| name        | STRING     | (empty string)              |
| alt         | STRING     | (empty string)              |
| description | STRING     | (empty string)              |
| uploaded    | TIMESTAMP  | now()                       |

### Indexes for media

| Name        | Type     | Columns  | References |
| :---------- | :------- | :------- | :--------- |
| pk_id       | PRIMARY  | id       |            |
| fk_site_id  | FOREIGN  | site     | sites (id) |
| indx_name   | INDEX    | name     |            |
