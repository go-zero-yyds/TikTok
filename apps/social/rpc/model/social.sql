SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user_stats
-- ----------------------------
DROP TABLE IF EXISTS `user_stats`;
CREATE TABLE `user_stats` (
                              `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
                              `follow_count` int(11) NOT NULL DEFAULT 0 COMMENT '关注数',
                              `follower_count` int(11) NOT NULL DEFAULT 0 COMMENT '粉丝数',
                              PRIMARY KEY (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_as_ci;

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
                           `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '消息ID',
                           `from_user_id` bigint(20) unsigned NOT NULL COMMENT '消息发生者ID',
                           `to_user_id` bigint(20) unsigned NOT NULL COMMENT '消息接收者ID',
                           `content` varchar(255) NOT NULL COMMENT '消息内容',
                           `create_time` datetime NOT NULL DEFAULT current_timestamp() COMMENT '消息创建时间',
                           PRIMARY KEY (`id`) USING BTREE,
                           KEY `idx_from_user_to_user` (`from_user_id`,`to_user_id`) COMMENT '联合索引,联查聊天记录使用'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_as_ci;

-- ----------------------------
-- Table structure for follow
-- ----------------------------
DROP TABLE IF EXISTS `follow`;
CREATE TABLE `follow` (
                          `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '字段ID',
                          `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
                          `to_user_id` bigint(20) unsigned NOT NULL COMMENT '关注者ID',
                          `behavior` enum('0','1') NOT NULL COMMENT '关注状态 0=>没关注 1=>关注',
                          `attribute` enum('0','1','2','3') NOT NULL COMMENT '关系 0=>陌生人 1=>关注 2=>粉丝 3=>好友',
                          PRIMARY KEY (`id`) USING BTREE,
                          UNIQUE KEY `idx_user_to_user` (`user_id`,`to_user_id`) COMMENT '联合索引,联查查关系使用'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_as_ci;

SET FOREIGN_KEY_CHECKS = 1;
