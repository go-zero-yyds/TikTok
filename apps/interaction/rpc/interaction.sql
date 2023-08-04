create table Favorite(
    userId bigint,
    videoId bigint,
    primary key(userId , videoId)
);
create table Comment(
    commentId bigint primary key ,
    userId bigint not null ,
    videoId bigint not null ,
    createDate DATE DEFAULT CURDATE(), # mm-dd
    content text not null,
    INDEX idx_user_video (userId, videoId)
);
