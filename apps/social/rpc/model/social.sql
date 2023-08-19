CREATE TABLE `follow`  (
                           `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字段ID',
                           `user_id` bigint NOT NULL COMMENT '用户ID',
                           `to_user_id` bigint NOT NULL COMMENT '关注者ID',
                           `behavior` ENUM('1', '2') NOT NULL  COMMENT '关注状态 2=>没关注 1=>关注',
                           PRIMARY KEY (`id`) USING BTREE,
                           UNIQUE INDEX idx_user_to_user (user_id, to_user_id) COMMENT '联合索引,联查查userid使用'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = uca1400_as_ci ROW_FORMAT = Dynamic;


CREATE TABLE `message`  (
                            `id` bigint NOT NULL AUTO_INCREMENT COMMENT '消息ID，由雪花算法生成',
                            `from_user_id` bigint NOT NULL COMMENT '消息发生者ID',
                            `to_user_id` bigint NOT NULL COMMENT '消息接收者ID',
                            `content` varchar(255) CHARACTER SET utf8mb4 COLLATE uca1400_as_ci NOT NULL COMMENT '消息内容',
                            `created_time` datetime NOT NULL COMMENT '消息创建时间',
                            PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = uca1400_as_ci ROW_FORMAT = Dynamic;
