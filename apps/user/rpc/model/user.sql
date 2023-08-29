SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
                        `user_id` bigint(20) unsigned NOT NULL COMMENT '唯一用户ID，使用雪花算法生成',
                        `username` varchar(256) NOT NULL COMMENT '用户注册用户名，最长256个字符，唯一',
                        `password` varchar(256) NOT NULL COMMENT '密码，最长256字符',
                        `avatar` varchar(256) NOT NULL COMMENT '用户头像',
                        `background_image` varchar(256) NOT NULL COMMENT '用户个人页顶部大图',
                        `signature` varchar(32) NOT NULL COMMENT '个人简介',
                        PRIMARY KEY (`user_id`),
                        UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_as_ci COMMENT='TikTok用户表';

SET FOREIGN_KEY_CHECKS = 1;
