CREATE TABLE `message`  (
                          `id` bigint NOT NULL COMMENT '消息ID，由雪花算法生成',
                          `to_user_id` bigint NOT NULL COMMENT '消息接收者ID',
                          `from_user_id` bigint NOT NULL COMMENT '消息发生着ID',
                          `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '消息内容',
                          `create_time` datetime NOT NULL COMMENT '消息创建时间',
                          PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
