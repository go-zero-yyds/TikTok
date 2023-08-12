SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `social`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字段ID',
  `user_id` bigint NOT NULL COMMENT '用户ID，由雪花算法生成',
  `follow_count` bigint NOT NULL COMMENT '关注总数',
  `follower_count` bigint NOT NULL COMMENT '粉丝总数',
  `total_favorited` bigint NOT NULL COMMENT '获赞数量',
  `favorite_count` bigint NOT NULL COMMENT '喜欢数',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;