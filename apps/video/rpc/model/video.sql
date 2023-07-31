CREATE TABLE `video` (
                         `video_id` bigint(20) unsigned NOT NULL COMMENT '唯一视频ID，使用雪花算法生成',
                         `user_id` bigint(20) unsigned NOT NULL COMMENT '视频作者用户ID',
                         `play_url` varchar(256) NOT NULL COMMENT '视频播放地址, key',
                         `cover_url` varchar(256) NOT NULL COMMENT '视频封面地址, key',
                         `title` varchar(256) NOT NULL COMMENT '视频标题',
                         `create_time` timestamp NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
                         PRIMARY KEY (`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='TikTok视频表';