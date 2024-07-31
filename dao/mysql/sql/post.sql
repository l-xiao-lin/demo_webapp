CREATE TABLE `post` (
                        `id` bigint(20) NOT NULL AUTO_INCREMENT,
                        `post_id` bigint(20) NOT NULL,
                        `title` varchar(128) NOT NULL,
                        `content` varchar(8192) NOT NULL,
                        `author_id` bigint(20) NOT NULL,
                        `community_id` bigint(20) NOT NULL,
                        `status` tinyint(4) NOT NULL DEFAULT '1',
                        `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                        `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `idx_post_id` (`post_id`),
                        KEY `idx_author_id` (`author_id`),
                        KEY `idx_community_id` (`community_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;