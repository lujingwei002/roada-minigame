CREATE DATABASE IF NOT EXISTS huanlemaomibei DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
USE huanlemaomibei;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `openid` varchar(64) NOT NULL  DEFAULT '' COMMENT 'openid',
    `nickname` varchar(128) NOT NULL  DEFAULT '' COMMENT '昵称',
    `avatar` varchar(128) NOT NULL  DEFAULT '' COMMENT '头像',
    `client_ip` varchar(32) NOT NULL  DEFAULT '' COMMENT '客户端地址',
    `logintime` int(11) NOT NULL  DEFAULT 0 COMMENT '登陆时间',
    `logouttime` int(11) NOT NULL  DEFAULT 0 COMMENT '登出时间',
    `roundstarttime` int(11) NOT NULL DEFAULT 0 COMMENT '游戏开始时间',
    `roundendtime` int(11) NOT NULL DEFAULT 0 COMMENT '游戏结束始时间',
    `lastroundstarttime` int(11) NOT NULL DEFAULT 0 COMMENT '游戏开始时间',
    `lastroundendtime` int(11) NOT NULL DEFAULT 0 COMMENT '游戏结束始时间',
    `score` bigint(11) NOT NULL DEFAULT 0 COMMENT '分数',
    `level` bigint(11) NOT NULL DEFAULT 0 COMMENT '关卡',
    `diamond` bigint(11) NOT NULL DEFAULT 0 COMMENT '钻石',
    `coin` bigint(11) NOT NULL DEFAULT 0 COMMENT '金币',
    `medal` bigint(11) NOT NULL DEFAULT 0 COMMENT '金牌',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    UNIQUE KEY (`openid`),
    KEY (`score`),
    KEY (`medal`),
    KEY (`level`),
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='用户表';

DROP TABLE IF EXISTS `scorerank`;
CREATE TABLE `scorerank` (
    `userid` bigint(11) NOT NULL DEFAULT 0 COMMENT '用户id',
    `score` bigint(11) NOT NULL DEFAULT 0 COMMENT '分数',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='积分总榜';
ALTER TABLE scorerank ADD INDEX i_score (`score`, `updatetime`) ;

DROP TABLE IF EXISTS `scoredayrank`;
CREATE TABLE `scoredayrank` (
    `rankid` bigint(11) NOT NULL DEFAULT 0 COMMENT '日期',
    `userid` bigint(11) NOT NULL DEFAULT 0 COMMENT '用户id',
    `score` bigint(11) NOT NULL DEFAULT 0 COMMENT '分数',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`rankid`, `userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='积分日榜';
ALTER TABLE scoredayrank ADD INDEX i_rankid_score (`rankid`, `score`, `updatetime`) ;

DROP TABLE IF EXISTS `medalrank`;
CREATE TABLE `medalrank` (
    `userid` bigint(11) NOT NULL DEFAULT 0 COMMENT '用户id',
    `score` bigint(11) NOT NULL DEFAULT 0 COMMENT '分数',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='金牌总榜';
ALTER TABLE medalrank ADD INDEX i_score (`score`, `updatetime`) ;

DROP TABLE IF EXISTS `medaldayrank`;
CREATE TABLE `medaldayrank` (
    `rankid` bigint(11) NOT NULL DEFAULT 0  COMMENT '日期',
    `userid` bigint(11) NOT NULL DEFAULT 0 COMMENT '用户id',
    `score` bigint(11) NOT NULL DEFAULT 0 COMMENT '分数',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`rankid`, `userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='金牌日榜';
ALTER TABLE medaldayrank ADD INDEX i_rankid_score (`rankid`, `score`, `updatetime`) ;

DROP TABLE IF EXISTS `levelrank`;
CREATE TABLE `levelrank` (
    `userid` bigint(11) NOT NULL DEFAULT 0 COMMENT '用户id',
    `score` bigint(11) NOT NULL DEFAULT 0 COMMENT '分数',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='等级总榜';
ALTER TABLE levelrank ADD INDEX i_score (`score`, `updatetime`) ;

DROP TABLE IF EXISTS `leveldayrank`;
CREATE TABLE `leveldayrank` (
    `rankid` bigint(11) NOT NULL DEFAULT 0 COMMENT '日期',
    `userid` bigint(11) NOT NULL DEFAULT 0 COMMENT '用户id',
    `score` bigint(11) NOT NULL DEFAULT 0 COMMENT '分数',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`rankid`, `userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='等级日榜';
ALTER TABLE leveldayrank ADD INDEX i_rankid_score (`rankid`, `score`, `updatetime`) ;

DROP TABLE IF EXISTS `gatestat`;
CREATE TABLE `gatestat` (
    `nodename` varchar(32) NOT NULL DEFAULT '' COMMENT '节点名',
    `nodefullname` varchar(64) NOT NULL DEFAULT '' COMMENT '节点全名',
    `gameid` int(11) NOT NULL DEFAULT '0' COMMENT '游戏id',
    `host` varchar(32) NOT NULL DEFAULT '' COMMENT 'Host',
    `port` int(11) NOT NULL DEFAULT 0 COMMENT 'Port',
    `onlinenum` int(11) NOT NULL DEFAULT 0 COMMENT '在线人数',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`nodename`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='服务器状态表';

DROP TABLE IF EXISTS `huanlemaomibei`;
CREATE TABLE `huanlemaomibei` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `freedraw_time` int(11) NOT NULL DEFAULT 0 COMMENT '免费抽奖时间',
    `inkarr` varchar(128) NOT NULL DEFAULT '0' COMMENT '墨水皮肤列表',
    `cuparr` varchar(128) NOT NULL DEFAULT '0' COMMENT '酒杯皮肤列表',
    `inkid` int(11) NOT NULL DEFAULT 0 COMMENT '墨水皮肤',
    `cupid` int(11) NOT NULL DEFAULT 0 COMMENT '酒杯皮肤',
    `signday` int(11) NOT NULL DEFAULT 0 COMMENT '签到天数',
    `lastsigntime` int(11) NOT NULL DEFAULT 0 COMMENT '上次签到时间',
    `signchecked` int(11) NOT NULL DEFAULT 0 COMMENT '是否检查签到跨天',
    `fly_time` int(11) NOT NULL DEFAULT 0 COMMENT '幸运蝴蝶抽奖时间',
    `offline_time` int(11) NOT NULL DEFAULT 0 COMMENT '离线收益时间',
    `hp` int(11) NOT NULL DEFAULT 0 COMMENT '体力',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='猫咪杯';

DROP TABLE IF EXISTS `huanlemaomibei_level`;
CREATE TABLE `huanlemaomibei_level` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `section` int(11) NOT NULL DEFAULT 0 COMMENT '章节',
    `level` int(11) NOT NULL DEFAULT 0 COMMENT '关卡',
    `unlocked` int(11) NOT NULL DEFAULT 0 COMMENT '是否解锁',
    `star` int(11) NOT NULL DEFAULT 0 COMMENT '星星数',
    `coin` int(11) NOT NULL DEFAULT 0 COMMENT '金币领取',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`, `section`, `level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='猫咪杯关卡';