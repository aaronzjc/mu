# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Database: mu
# Generation Time: 2019-09-22 15:34:09 +0000
# ************************************************************

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE DATABASE IF NOT EXISTS `mu` default CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

# Dump of table node
# ------------------------------------------------------------

CREATE TABLE `node` IF NOT EXISTS (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(64) NOT NULL DEFAULT '' COMMENT '节点名字',
    `addr` varchar(128) NOT NULL DEFAULT '' COMMENT '节点IP',
    `type` tinyint(4) NOT NULL COMMENT '1:海外;2:大陆',
    `enable` tinyint(4) NOT NULL COMMENT '是否开启',
    `ping` tinyint(4) NOT NULL COMMENT 'pong',
    `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `ip_idx` (`addr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table site
# ------------------------------------------------------------

CREATE TABLE `site` IF NOT EXISTS  (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '站点名字',
  `root` varchar(256) NOT NULL DEFAULT '' COMMENT '站点根地址',
  `key` varchar(64) NOT NULL DEFAULT '' COMMENT '站点KEY',
  `desc` varchar(256) NOT NULL DEFAULT '' COMMENT '站点描述',
  `type` tinyint(4) NOT NULL COMMENT '1:html;2:json',
  `tags` text NOT NULL COMMENT '标签',
  `cron` varchar(64) NOT NULL DEFAULT '' COMMENT '频率',
  `enable` tinyint(4) NOT NULL COMMENT '是否开启',
  `node_option` tinyint(4) NOT NULL COMMENT '1:节点类型;2:节点IP',
  `node_type` tinyint(4) NOT NULL COMMENT '爬取的服务器类型',
  `node_hosts` text NOT NULL COMMENT '爬取的服务器IP列表',
  PRIMARY KEY (`id`),
  UNIQUE KEY `key_idx` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

# Dump of table user
# ------------------------------------------------------------

CREATE TABLE `user` IF NOT EXISTS  (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(64) NOT NULL DEFAULT '',
    `nickname` varchar(64) NOT NULL DEFAULT '',
    `avatar` varchar(256) NOT NULL DEFAULT '',
    `auth_type` tinyint(4) NOT NULL,
    `auth_time` datetime NOT NULL,
    `token` varchar(64) NOT NULL DEFAULT '',
    `expire_at` int(11) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

# Dump of table favor
# ------------------------------------------------------------

CREATE TABLE `favor` IF NOT EXISTS  (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `site` varchar(64) NOT NULL,
  `key` varchar(64) NOT NULL DEFAULT '',
  `origin_url` text NOT NULL,
  `title` text NOT NULL,
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `sk_idx` (`user_id`,`site`,`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
