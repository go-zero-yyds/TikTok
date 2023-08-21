CREATE TABLE comment(
    commentId BIGINT NOT NULL COMMENT "评论id 雪花算法生成",
    userId BIGINT NOT NULL COMMENT "用户id" ,
    videoId BIGINT NOT NULL COMMENT "视频id" ,
    createDate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP  COMMENT "创建日期 ",
    content TEXT NOT NULL COMMENT "用户评论内容",
    INDEX idx_user_video (videoId, userId) COMMENT "联合索引,video用的较多",
    PRIMARY KEY (commentId)
);

CREATE TABLE favorite(
    favoriteId BIGINT NOT NULL AUTO_INCREMENT COMMENT "自增id,优化插入效率",
    userId BIGINT NOT NULL COMMENT "用户id",
    videoId BIGINT NOT NULL COMMENT "视频id",
    behavior ENUM('1', '2') NOT NULL COMMENT "1:点赞 2:未点赞",
    PRIMARY KEY (favoriteId),
    UNIQUE INDEX idx_user_video (userId, videoId) COMMENT "联合索引,联查或查userid使用",
    INDEX idx_video (videoId) COMMENT "查videoId使用"
);
