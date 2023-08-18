/*
 Navicat Premium Data Transfer

 Source Server         : TG
 Source Server Type    : MySQL
 Source Server Version : 80028
 Source Host           : localhost:3306
 Source Schema         : tiktok_social

 Target Server Type    : MySQL
 Target Server Version : 80028
 File Encoding         : 65001

 Date: 18/08/2023 20:27:46
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for social
-- ----------------------------
DROP TABLE IF EXISTS `social`;
CREATE TABLE `social`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字段ID',
  `user_id` bigint NOT NULL COMMENT '用户ID，由雪花算法生成',
  `follow_count` bigint NOT NULL COMMENT '关注总数',
  `follower_count` bigint NOT NULL COMMENT '粉丝总数',
  `version` bigint NOT NULL COMMENT 'for update版本号',
  PRIMARY KEY (`id`, `user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1473 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
