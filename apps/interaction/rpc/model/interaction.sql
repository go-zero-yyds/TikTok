-- Active: 1682338800757@@127.0.0.1@3306@tiktok_interaction

CREATE TABLE
    comment (
        comment_id BIGINT NOT NULL AUTO_INCREMENT COMMENT "自增id",
        user_id BIGINT NOT NULL COMMENT "用户id",
        video_id BIGINT NOT NULL COMMENT "视频id",
        create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "创建日期",
        content TEXT NOT NULL COMMENT "用户评论内容",
        is_deleted ENUM('0', '1') NOT NULL DEFAULT '0' COMMENT "0:未删除 1:已删除",
        ip_address TEXT NOT NULL COMMENT "用户IP地址",
        location TEXT NOT NULL COMMENT "IP地址归属地",
        INDEX idx_user_video (video_id, user_id) COMMENT "联合索引，video用的较多",
        PRIMARY KEY (comment_id)
    );

CREATE TABLE
    favorite (
        favorite_id BIGINT NOT NULL AUTO_INCREMENT COMMENT "自增id，优化插入效率",
        user_id BIGINT NOT NULL COMMENT "用户id",
        video_id BIGINT NOT NULL COMMENT "视频id",
        behavior ENUM('0', '1') NOT NULL COMMENT "1:点赞 2:未点赞",
        PRIMARY KEY (favorite_id),
        UNIQUE INDEX idx_user_video (user_id, video_id) COMMENT "联合索引，联查或查userid使用",
        INDEX idx_video (video_id) COMMENT "查videoId使用"
    );

-- New table to record user likes

CREATE TABLE
    user_likes (
        user_id BIGINT NOT NULL COMMENT "用户id",
        like_count BIGINT NOT NULL DEFAULT 0 COMMENT "用户点赞数",
        PRIMARY KEY (user_id)
    );

-- New table to record video likes and comments

CREATE TABLE
    video_stats (
        video_id BIGINT NOT NULL COMMENT "视频id",
        like_count BIGINT NOT NULL DEFAULT 0 COMMENT "视频获赞数",
        comment_count BIGINT NOT NULL DEFAULT 0 COMMENT "视频评论数",
        PRIMARY KEY (video_id)
    );