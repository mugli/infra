PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE `identities` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`name` text,`last_seen_at` datetime,`created_by` integer,PRIMARY KEY (`id`));
INSERT INTO identities VALUES(64404683855568897,'2022-06-27 17:21:13.212609471+00:00','2022-06-27 17:21:13.212609471+00:00',NULL,'connector','0001-01-01 00:00:00+00:00',1);
CREATE TABLE `providers` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`name` text,`url` text,`client_id` text,`client_secret` text,`created_by` integer, `auth_url` text, `scopes` text, PRIMARY KEY (`id`));
INSERT INTO providers VALUES(64398317355081729,'2022-06-27 16:55:55.320743754+00:00','2022-06-27 16:55:55.320743754+00:00',NULL,'infra','','','AAAAENAV3AEOwupzNbFRKjqW+B8GYWVzZ2NtBGW5mRYbL3Jvb3QvLmluZnJhL3NxbGl0ZTMuZGIua2V5DIswJidILk0q/uLQhw',1,'','');
CREATE TABLE `provider_users` (`identity_id` integer,`provider_id` integer, `id` integer, `created_at` datetime, `updated_at` datetime, `deleted_at` datetime, `email` text, `groups` text, `last_update` datetime, `redirect_url` text, `access_token` text, `refresh_token` text, `expires_at` datetime,PRIMARY KEY (`identity_id`,`provider_id`),CONSTRAINT `fk_provider_users_identity` FOREIGN KEY (`identity_id`) REFERENCES `identities`(`id`),CONSTRAINT `fk_provider_users_provider` FOREIGN KEY (`provider_id`) REFERENCES `providers`(`id`));
CREATE TABLE `groups` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`name` text,`created_by` integer,`created_by_provider` integer,PRIMARY KEY (`id`));
CREATE TABLE `identities_groups` (`group_id` integer,`identity_id` integer,PRIMARY KEY (`group_id`,`identity_id`),CONSTRAINT `fk_identities_groups_group` FOREIGN KEY (`group_id`) REFERENCES `groups`(`id`),CONSTRAINT `fk_identities_groups_identity` FOREIGN KEY (`identity_id`) REFERENCES `identities`(`id`));
CREATE TABLE `grants` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`subject` text,`privilege` text,`resource` text,`created_by` integer,PRIMARY KEY (`id`));
INSERT INTO grants VALUES(64404683855568898,'2022-06-27 17:21:13.212907554+00:00','2022-06-27 17:21:13.212907554+00:00',NULL,'i:9EUZQmQAM8','connector','infra',1);
CREATE TABLE `destinations` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`name` text,`unique_id` text,`connection_url` text,`connection_ca` text,`resources` text,`roles` text,PRIMARY KEY (`id`));
INSERT INTO destinations VALUES(67067378731917312,'2022-07-05 01:41:49.143574+00.00','2022-07-05 01:41:49.143574+00.00',NULL,'docker-desktop','unique-id','localhost:123','','','');
CREATE TABLE `access_keys` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`name` text,`issued_for` integer,`provider_id` integer,`scopes` text,`expires_at` datetime,`extension` integer,`extension_deadline` datetime,`key_id` text,`secret_checksum` blob,PRIMARY KEY (`id`),CONSTRAINT `fk_access_keys_issued_for_identity` FOREIGN KEY (`issued_for`) REFERENCES `identities`(`id`));
CREATE TABLE `settings` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`private_jwk` blob,`public_jwk` blob,PRIMARY KEY (`id`));
INSERT INTO settings VALUES(64404683851374592,'2022-06-27 17:21:13.211997596+00:00','2022-06-27 17:21:13.211997596+00:00',NULL,X'41414141346833344766467153774c584971432b37776851466b3051333743667135367474657442613741795932727476735655596355416565774f6e777048593357654a75774845537a7a48326f6e2b7243376c75667846736e3154456f32666f6b66795a467a357456516a487a6b735063414e6c4855687068785857713531417553414b76546978694a71534a6f4c3868446a667a57354b3268535051304c3732765138486679637a78412f655336452b5142796e77505969736a38336e3671704c6c4a3047694666354b2f784a593362723659773665397a4877456266482b6e59487668594f53724d3151486c6d37375856447958556e34746f4e334d565052744e754d336c30447772446535755070674341306a686f634a454f4c4f714e723158794e59563672324d316e70457873475957567a5a324e74424e5873554e77624c334a76623351764c6d6c755a6e4a684c334e7862476c305a544d755a47497561325635444c4a637558534f4e6f713468764c6f7041',X'7b22757365223a22736967222c226b7479223a224f4b50222c226b6964223a22395f617a517837667236467865645864786761742d385f617a62436469544f762d377270387a70355550303d222c22637276223a2245643235353139222c22616c67223a2245443235353139222c2278223a22716e45396b784a746167785559745430556e6c5133524467686b3277674b7576487564505f343465784334227d');
CREATE TABLE `encryption_keys` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`key_id` integer,`name` text,`encrypted` blob,`algorithm` text,`root_key_id` text,PRIMARY KEY (`id`));
INSERT INTO encryption_keys VALUES(64404683822014464,'2022-06-27 17:21:13.204509137+00:00','2022-06-27 17:21:13.204509137+00:00',NULL,171020364,'dbkey',X'414141414d417248304b6b56515a46527743666c567251542b4e35777a34597959622b75632b4d347a7666764d7061347132794e553559566f5330784464786f7a58434147415a685a584e6e59323045343744455167414d4a6f6772704f6772574533506e476d39','aesgcm','/root/.infra/sqlite3.db.key');
CREATE TABLE `credentials` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`identity_id` integer,`password_hash` blob,`one_time_password` numeric,PRIMARY KEY (`id`));
CREATE TABLE migrations (id VARCHAR(255) PRIMARY KEY);
INSERT INTO migrations VALUES('SCHEMA_INIT');
INSERT INTO migrations VALUES('202203231621');
INSERT INTO migrations VALUES('202203241643');
INSERT INTO migrations VALUES('202203301642');
INSERT INTO migrations VALUES('202203301652');
INSERT INTO migrations VALUES('202203301643');
INSERT INTO migrations VALUES('202203301645');
INSERT INTO migrations VALUES('202203301646');
INSERT INTO migrations VALUES('202203301647');
INSERT INTO migrations VALUES('202203301648');
INSERT INTO migrations VALUES('202204061643');
INSERT INTO migrations VALUES('202204111503');
INSERT INTO migrations VALUES('202204181613');
INSERT INTO migrations VALUES('202204211705');
INSERT INTO migrations VALUES('202204281130');
INSERT INTO migrations VALUES('202204291613');
INSERT INTO migrations VALUES('202206081027');
INSERT INTO migrations VALUES('202206161733');
CREATE UNIQUE INDEX `idx_identities_name` ON `identities`(`name`) WHERE deleted_at is NULL;
CREATE UNIQUE INDEX `idx_providers_name` ON `providers`(`name`) WHERE deleted_at is NULL;
CREATE UNIQUE INDEX `idx_groups_name` ON `groups`(`name`) WHERE deleted_at is NULL;
CREATE UNIQUE INDEX `idx_destinations_unique_id` ON `destinations`(`unique_id`) WHERE deleted_at is NULL;
CREATE UNIQUE INDEX `idx_access_keys_key_id` ON `access_keys`(`key_id`) WHERE deleted_at is NULL;
CREATE UNIQUE INDEX `idx_access_keys_name` ON `access_keys`(`name`) WHERE deleted_at is NULL;
CREATE UNIQUE INDEX `idx_encryption_keys_key_id` ON `encryption_keys`(`key_id`);
CREATE UNIQUE INDEX `idx_credentials_identity_id` ON `credentials`(`identity_id`) WHERE deleted_at is NULL;
COMMIT;
