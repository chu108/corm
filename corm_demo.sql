/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50728
 Source Host           : localhost:3306
 Source Schema         : corm_demo

 Target Server Type    : MySQL
 Target Server Version : 50728
 File Encoding         : 65001

 Date: 14/11/2019 18:04:56
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for groups
-- ----------------------------
DROP TABLE IF EXISTS `groups`;
CREATE TABLE `groups` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '组名',
  `description` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=48 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户组';

-- ----------------------------
-- Records of groups
-- ----------------------------
BEGIN;
INSERT INTO `groups` VALUES (14, '用户组1', '用户组1', '2017-10-28 15:05:08', '2019-11-14 14:12:47');
INSERT INTO `groups` VALUES (15, '用户组2', '用户组2', '2017-10-28 17:13:29', '2019-11-14 14:12:51');
INSERT INTO `groups` VALUES (16, '用户组3', '用户组3', '2017-11-10 17:21:30', '2019-11-14 14:12:53');
INSERT INTO `groups` VALUES (17, '用户组4', '用户组4', '2017-11-20 16:14:05', '2019-11-14 14:12:56');
COMMIT;

-- ----------------------------
-- Table structure for user_groups
-- ----------------------------
DROP TABLE IF EXISTS `user_groups`;
CREATE TABLE `user_groups` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned NOT NULL COMMENT '用户',
  `group_id` int(11) unsigned NOT NULL COMMENT '组',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=211 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户与组关联表';

-- ----------------------------
-- Records of user_groups
-- ----------------------------
BEGIN;
INSERT INTO `user_groups` VALUES (59, 7, 15, '2018-03-13 15:40:19', '2018-03-13 15:44:52');
INSERT INTO `user_groups` VALUES (111, 9, 16, '2018-06-30 16:10:55', '2018-06-30 16:10:55');
INSERT INTO `user_groups` VALUES (115, 10, 28, '2018-07-06 10:26:36', '2018-07-06 10:26:36');
INSERT INTO `user_groups` VALUES (156, 5, 29, '2018-08-20 10:53:11', '2018-08-20 10:53:11');
COMMIT;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `nickname` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '昵称',
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '真实姓名',
  `age` int(4) NOT NULL DEFAULT '0',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '联系电话',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=216 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
INSERT INTO `users` VALUES (5, '张三', '张三', 22, '18735181111', '2016-10-25 23:28:29', '2019-11-14 14:14:11');
INSERT INTO `users` VALUES (8, '李四', '李四', 24, '18500082222', '2017-08-08 18:25:25', '2019-11-14 14:14:12');
INSERT INTO `users` VALUES (9, '大张伟', '大张伟', 20, '13888888888', '2017-08-08 22:14:01', '2019-11-14 17:41:49');
INSERT INTO `users` VALUES (10, '小明', '小明', 20, '13683624444', '2017-08-10 19:04:08', '2019-11-14 17:41:49');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
