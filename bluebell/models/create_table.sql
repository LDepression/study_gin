
create database bluebell,

create table `user`(
    `id` bigint(20) not null  auto_increment,
    `user_id` bigint(20) not null ,
    `username`  varchar(64) collate utf8mb4_general_ci not null,
    `password` varchar(64) collate utf8mb4_general_ci not null ,
    `email` varchar(64) collate utf8mb4_general_ci,
    `gender` tinyint(4) not null default '0',
    `create_time` timestamp null default current_timestamp,
    `update_time` timestamp null default current_timestamp on update current_timestamp,
    primary key (`id`),
    unique key `idx_username` (`username`) using btree ,
    unique key `idx_user_id`(`user_id`) using btree
) engine =InnoDB default charset=utf8mb4 collate utf8mb4_general_ci

create table `community`(
    `id` int(11) not null auto_increment,
    `community_id` int(10) unsigned not null ,
    `community_name` varchar(128) collate utf8mb4_general_ci not null ,
    `introduction` varchar(256) collate utf8mb4_general_ci not null ,
    `create_time` timestamp not null default current_timestamp,
    `update_time` timestamp not null default current_timestamp on update current_timestamp,
    primary key (`id`),
     unique key `idx_community_id` (`community_id`),
     unique key `idx_community_name` (`community_name`)
)engine =InnoDB default charset=utf8mb4 collate =utf8mb4_general_ci;
insert into `community` values('1','1','Go','Golang','2016-11-01 08:10:10','2016-11-01 08:10:10');
insert into `community` values('2','2','leetcode','刷穿leetcode','2020-01-01 08:10:10','2022-11-01 08:10:10');
insert into `community` values('3','3','CF','codeforce','2018-11-01 08:10:10','2020-11-01 08:10:10');
select * form community;


create table `post`(
    `id` bigint(20) not null auto_increment,
    `post_id` bigint(20) not null comment '帖子id',
    `title` varchar(128) collate utf8mb4_general_ci not null comment '标题',
    `content` varchar(9000) collate utf8mb4_general_ci not null comment '内容',
    `author_id` bigint(20) not null comment '作者用户的id',
    `community_id` bigint(20) not null comment '所属社区',
    `status` tinyint(4) not null default '1' comment '帖子的状态',
    `create_time` timestamp not null default current_timestamp comment '创建时间',
    `update_time` timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (`id`),
    unique key `idx_post_id` (`post_id`),
    key `idx_author_id` (`author_id`),
    key `idx_community_id` (`community_id`)
)engine =InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;