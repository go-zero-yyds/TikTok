SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for video_stats
-- ----------------------------
DROP TABLE IF EXISTS `video_stats`;
CREATE TABLE `video_stats` (
                               `video_id` bigint(20) unsigned NOT NULL COMMENT '视频id',
                               `like_count` bigint(20) NOT NULL DEFAULT 0 COMMENT '视频获赞数',
                               `comment_count` bigint(20) NOT NULL DEFAULT 0 COMMENT '视频评论数',
                               PRIMARY KEY (`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_as_ci;

-- ----------------------------
-- Table structure for user_likes
-- ----------------------------
DROP TABLE IF EXISTS `user_likes`;
CREATE TABLE `user_likes` (
                              `user_id` bigint(20) unsigned NOT NULL COMMENT '用户id',
                              `like_count` bigint(20) NOT NULL DEFAULT 0 COMMENT '用户点赞数',
                              PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_as_ci;

-- ----------------------------
-- Table structure for favorite
-- ----------------------------
DROP TABLE IF EXISTS `favorite`;
CREATE TABLE `favorite` (
                            `favorite_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id，优化插入效率',
                            `user_id` bigint(20) unsigned NOT NULL COMMENT '用户id',
                            `video_id` bigint(20) unsigned NOT NULL COMMENT '视频id',
                            `behavior` enum('0','1') NOT NULL COMMENT '0:未点赞 1:点赞',
                            PRIMARY KEY (`favorite_id`),
                            UNIQUE KEY `idx_user_video` (`user_id`,`video_id`) COMMENT '联合索引，联查或查userid使用'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_as_ci;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
                           `comment_id` bigint(20) unsigned NOT NULL COMMENT '雪花id',
                           `user_id` bigint(20) unsigned NOT NULL COMMENT '用户id',
                           `video_id` bigint(20) unsigned NOT NULL COMMENT '视频id',
                           `content` text NOT NULL COMMENT '用户评论内容',
                           `ip_address` varchar(16) NOT NULL COMMENT '用户IP地址',
                           `location` varchar(64) NOT NULL COMMENT 'IP地址归属地',
                           `create_time` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建日期',
                           `is_deleted` enum('0','1') NOT NULL DEFAULT '0' COMMENT '0:未删除 1:已删除',
                           PRIMARY KEY (`comment_id`),
                           KEY `idx_video` (`video_id`) COMMENT '查询视频评论列表'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_as_ci;

SET FOREIGN_KEY_CHECKS = 1;
