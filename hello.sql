-- phpMyAdmin SQL Dump
-- version 4.8.5
-- https://www.phpmyadmin.net/
--
-- 主机： localhost
-- 生成日期： 2023-02-28 13:46:06
-- 服务器版本： 5.7.39-log
-- PHP 版本： 7.3.31

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 数据库： `hello`
--

-- --------------------------------------------------------

--
-- 表的结构 `h_message`
--

CREATE TABLE `h_message` (
  `id` int(11) NOT NULL,
  `users_id` int(11) NOT NULL,
  `from_id` int(11) NOT NULL,
  `to_id` int(11) NOT NULL,
  `content` text,
  `is_read` tinyint(4) NOT NULL COMMENT '0 未读 1 已读',
  `create_at` int(11) NOT NULL,
  `update_at` int(11) NOT NULL DEFAULT '0',
  `delete_at` int(11) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天记录';

--
-- 转存表中的数据 `h_message`
--

INSERT INTO `h_message` (`id`, `users_id`, `from_id`, `to_id`, `content`, `is_read`, `create_at`, `update_at`, `delete_at`) VALUES
(1, 2, 1, 2, '%E4%BD%A0%E5%A5%BD%E5%95%8A', 0, 1677506658, 0, 0),
(2, 1, 1, 2, '%E4%BD%A0%E5%A5%BD%E5%95%8A', 1, 1677506658, 0, 0);

-- --------------------------------------------------------

--
-- 表的结构 `h_users`
--

CREATE TABLE `h_users` (
  `id` int(11) NOT NULL,
  `username` varchar(30) NOT NULL COMMENT '用户名',
  `password` varchar(32) NOT NULL COMMENT '密码',
  `nickname` varchar(20) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(120) NOT NULL DEFAULT 'http://fan.jx.cn:81/hello/drawable/default.png' COMMENT '头像',
  `cover` varchar(200) NOT NULL DEFAULT 'http://fan.jx.cn:81/hello/drawable/cover.png' COMMENT '封面图',
  `star` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
  `sex` tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别 0 保密 1 男 2 女',
  `birthday` int(11) NOT NULL COMMENT '出生日期',
  `identity` varchar(30) DEFAULT NULL COMMENT '身份',
  `motto` varchar(50) DEFAULT NULL COMMENT '格言',
  `address` varchar(120) DEFAULT NULL,
  `info` varchar(200) DEFAULT NULL COMMENT '个人简介',
  `is_ban` tinyint(4) NOT NULL DEFAULT '0' COMMENT '禁止登录 0 否 1 是',
  `last_login` int(11) DEFAULT NULL,
  `device` varchar(60) DEFAULT NULL,
  `create_at` int(11) NOT NULL DEFAULT '0' COMMENT '注册时间',
  `update_at` int(11) NOT NULL DEFAULT '0',
  `delete_at` int(11) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

--
-- 转存表中的数据 `h_users`
--

INSERT INTO `h_users` (`id`, `username`, `password`, `nickname`, `avatar`, `cover`, `star`, `sex`, `birthday`, `identity`, `motto`, `address`, `info`, `is_ban`, `last_login`, `device`, `create_at`, `update_at`, `delete_at`) VALUES
(1, '441479573@qq.com', '3c2233e28ece95ddf0d16748fe0c5870', 'Break', 'http://fan.jx.cn:81/hello/upload/f97783599b1b28c5d32dfea98ca6323c.jpg', 'http://fan.jx.cn:81/hello/upload/df66e8c22d2704e68647d3c93618fff1.png', 180, 1, 1647670198, 'Hello开发者', NULL, '九江', 'Hello开发者，95后程序员，江西人，现居杭州', 0, 1677510409, 'realme RMX2202', 1647670198, 0, 0),
(2, '474024153@qq.com', '3c2233e28ece95ddf0d16748fe0c5870', '杰瑞', 'http://fan.jx.cn:81/hello/images/buxsren.jpg', 'http://fan.jx.cn:81/hello/images/cover.png', 20, 0, 1647671110, NULL, NULL, '江西九江', NULL, 0, 1677501129, 'Netease cancro_x86_64', 1647671110, 0, 0),
(3, '2793469806@qq.com', '3c2233e28ece95ddf0d16748fe0c5870', '小冰', 'http://fan.jx.cn:81/hello/images/xiaobing.jpg', 'http://fan.jx.cn:81/hello/images/cover.png', 0, 2, 1647680794, NULL, NULL, '江西九江', NULL, 0, 1647680794, 'Android', 1647680794, 0, 0),
(4, 'buxsren@qq.com', '3c2233e28ece95ddf0d16748fe0c5870', 'ikun', 'http://fan.jx.cn:81/hello/drawable/default.png', 'http://fan.jx.cn:81/hello/drawable/cover.png', 0, 0, 1677518323, '学生', NULL, NULL, NULL, 0, 1677562904, 'realme RMX2202', 1677518323, 0, 0);

--
-- 转储表的索引
--

--
-- 表的索引 `h_message`
--
ALTER TABLE `h_message`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `h_users`
--
ALTER TABLE `h_users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `username` (`username`),
  ADD KEY `create_at` (`create_at`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `h_message`
--
ALTER TABLE `h_message`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- 使用表AUTO_INCREMENT `h_users`
--
ALTER TABLE `h_users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
