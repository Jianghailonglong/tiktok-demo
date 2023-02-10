SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '用户id，自增',
    `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名，唯一',
    `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码，加密过的',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `username_idx`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for relations
-- ----------------------------
DROP TABLE IF EXISTS `relations`;
CREATE TABLE `relations` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '自增关系id',
    `user_id` int NOT NULL COMMENT '关注者id',
    `to_user_id` int NOT NULL COMMENT '被关注者id',
    `subscribed` int NOT NULL COMMENT '1为正在关注，0为取关', PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos`  (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '视频id，自增',
    `author_id` int NOT NULL COMMENT '作者id',
    `play_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '视频url',
    `cover_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '封面url',
    `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '标题',
    `publish_time` bigint(20) NULL DEFAULT NULL COMMENT '发布时间',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `publish_time_idx`(`publish_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`  (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '评论自增id',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评论发布日期',
    `video_id` int NOT NULL COMMENT '视频id',
    `user_id` int NOT NULL COMMENT '用户id',
    `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '评论内容',
    `is_comment` int NOT NULL COMMENT '是否评论，0取消评论，1评论',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for favorites
-- ----------------------------
DROP TABLE IF EXISTS `favorites`;
CREATE TABLE `favorites`  (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '点赞自增id',
    `user_id` int NOT NULL COMMENT '用户id',
    `video_id` int NOT NULL COMMENT '视频id',
    `is_favorite` int NOT NULL COMMENT '是否点赞，0取消点赞，1点赞',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;

