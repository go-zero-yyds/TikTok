create table Favorite(
    userId bigint COMMENT "用户id",
    videoId bigint COMMENT "视频id",
    primary key(userId , videoId) COMMENT "联合主键"
);
create table Comment(
    commentId bigint primary key COMMENT "评论id 雪花算法生成",
    userId bigint not null COMMENT "用户id" ,
    videoId bigint not null COMMENT "视频id" ,
    createDate DATE DEFAULT CURDATE() COMMENT "创建日期 mm-dd格式",
    content text not null COMMENT "用户评论内容",
    INDEX idx_user_video (videoId, userId) COMMENT "联合索引,video用的较多"
);
