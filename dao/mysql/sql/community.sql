DROP TABLE IF EXISTS `community`;
CREATE TABLE `community` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `community_id`  BIGINT NOT NULL,
  `community_name` varchar(128) NOT NULL,
  `introduction` varchar(256) NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_community_id` (`community_id`),
  UNIQUE KEY `idx_community_name` (`community_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
