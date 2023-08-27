CREATE TABLE `follow`  (
                           `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字段ID',
                           `user_id` bigint NOT NULL COMMENT '用户ID',
                           `to_user_id` bigint NOT NULL COMMENT '关注者ID',
                           `behavior` ENUM('0', '1') NOT NULL  COMMENT '关注状态 0=>没关注 1=>关注',
                           `attribute` ENUM('0', '1', '2', '3') NOT NULL  COMMENT '关系 0=>陌生人 1=>关注 2=>粉丝 3=>好友',
                           PRIMARY KEY (`id`) USING BTREE,
                           UNIQUE INDEX idx_user_to_user (user_id, to_user_id) COMMENT '联合索引,联查查关系使用'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = uca1400_as_ci ROW_FORMAT = Dynamic;


CREATE TABLE `message`  (
                            `id` bigint NOT NULL AUTO_INCREMENT COMMENT '消息ID',
                            `from_user_id` bigint NOT NULL COMMENT '消息发生者ID',
                            `to_user_id` bigint NOT NULL COMMENT '消息接收者ID',
                            `content` varchar(255) CHARACTER SET utf8mb4 COLLATE uca1400_as_ci NOT NULL COMMENT '消息内容',
                            `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '消息创建时间',
                            INDEX idx_from_user_to_user (from_user_id, to_user_id) COMMENT '联合索引,联查聊天记录使用',
                            PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = uca1400_as_ci ROW_FORMAT = Dynamic;

CREATE TABLE `user_stats` (
                              `user_id` bigint NOT NULL COMMENT '用户ID',
                              `follow_count` int NOT NULL DEFAULT 0 COMMENT '关注数',
                              `follower_count` int NOT NULL DEFAULT 0 COMMENT '粉丝数',
                              PRIMARY KEY (`user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = uca1400_as_ci ROW_FORMAT = Dynamic;