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

 Date: 18/08/2023 20:27:41
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '消息ID，由雪花算法生成',
  `from_user_id` bigint NOT NULL COMMENT '消息发生者ID',
  `to_user_id` bigint NOT NULL COMMENT '消息接收者ID',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '消息内容',
  `created_time` datetime NOT NULL COMMENT '消息创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4316 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
