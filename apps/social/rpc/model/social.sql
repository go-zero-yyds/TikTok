-- Active: 1682338800757@@127.0.0.1@3306@tiktok_social

CREATE TABLE
    `follow` (
        `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字段ID',
        `user_id` bigint NOT NULL COMMENT '用户ID',
        `to_user_id` bigint NOT NULL COMMENT '关注者ID',
        `behavior` ENUM('1', '2') NOT NULL COMMENT '关注状态 2=>没关注 1=>关注',
        PRIMARY KEY (`id`) USING BTREE,
        UNIQUE INDEX idx_user_to_user (user_id, to_user_id) COMMENT '联合索引,联查查userid使用'
    ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

CREATE TABLE
    `message` (
        `id` bigint NOT NULL AUTO_INCREMENT COMMENT '消息ID',
        `from_user_id` bigint NOT NULL COMMENT '消息发生者ID',
        `to_user_id` bigint NOT NULL COMMENT '消息接收者ID',
        `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '消息内容',
        `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '消息创建时间',
        PRIMARY KEY (`id`) USING BTREE
    ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;