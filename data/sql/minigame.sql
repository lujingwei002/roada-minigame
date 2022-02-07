CREATE DATABASE IF NOT EXISTS minigame DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
USE minigame;

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

DROP TABLE IF EXISTS `benpaobaxiaoqie`;
CREATE TABLE `benpaobaxiaoqie` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `skinarr` varchar(64) NOT NULL DEFAULT '0' COMMENT '皮肤列表',
    `skin` int(11) NOT NULL DEFAULT 0 COMMENT '皮肤',
    `lastsigntime` int(11) NOT NULL DEFAULT 0 COMMENT '签到时间',
    `signtimes` int(11) NOT NULL DEFAULT 0 COMMENT '签到次数',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='奔跑吧小企鹅';

DROP TABLE IF EXISTS `caichengyu`;
CREATE TABLE `caichengyu` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `rolelevel` int(11) NOT NULL DEFAULT 0 COMMENT '角色等级',
    `buildlevel` int(11) NOT NULL DEFAULT 0 COMMENT '建筑等级',
    `level` int(11) NOT NULL DEFAULT 1 COMMENT '当前关卡',
    `leveltype` int(11) NOT NULL DEFAULT 0 COMMENT '关卡状态',
    `leveltip` int(11) NOT NULL DEFAULT 0 COMMENT '提示次数',
    `hp` bigint(11) NOT NULL DEFAULT 0 COMMENT '体力',
    `hpdate` int(11) NOT NULL DEFAULT 0 COMMENT '上次体力恢复时间',
    `gethpcount` int(11) NOT NULL DEFAULT 1 COMMENT '领取体力次数',
    `gethpday` varchar(8) NOT NULL DEFAULT 0 COMMENT '领取体力日期',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='猜成语';

DROP TABLE IF EXISTS `fangkuainiao`;
CREATE TABLE `fangkuainiao` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `birdarr` varchar(64) NOT NULL DEFAULT '0' COMMENT '皮肤列表',
    `birdid` int(11) NOT NULL DEFAULT 0 COMMENT '皮肤',
    `signtime` int(11) NOT NULL DEFAULT 0 COMMENT '签到时间',
    `signday` int(11) NOT NULL DEFAULT 0 COMMENT '签到次数',
    `getgoldcount` int(11) NOT NULL DEFAULT 0 COMMENT '领取金币次数',
    `getgoldtime` int(11) NOT NULL DEFAULT 0 COMMENT '领取金币时间',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='方块鸟';

DROP TABLE IF EXISTS `gongjianchuanshu`;
CREATE TABLE `gongjianchuanshu` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `level` int(11) NOT NULL DEFAULT 0 COMMENT '当前关卡',
    `skinarr` varchar(64) NOT NULL DEFAULT '[1,0,0,0,0,0,0,0,0,0]' COMMENT '皮肤列表',
    `shoparr` varchar(64) NOT NULL DEFAULT '[2,3,4,5,6,7,8,9,10]' COMMENT '商品列表',
    `skinid` int(11) NOT NULL DEFAULT 1 COMMENT '皮肤',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='弓箭传说';

DROP TABLE IF EXISTS `paopaolong`;
CREATE TABLE `paopaolong` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `level` int(11) NOT NULL DEFAULT 1 COMMENT '当前关卡',
    `itemarr` varchar(64) NOT NULL DEFAULT '[0,0,0,0,5,0,0]' COMMENT '道具',
    `hp` bigint(11) NOT NULL DEFAULT 5 COMMENT '体力',
    `freedraw_time` int(11) NOT NULL DEFAULT 0 COMMENT '免费抽奖时间',
    `new_pack_redeemed` int(11) NOT NULL DEFAULT 0 COMMENT '新手礼包领取状态',
    `shop_free_diamond_time` int(11) NOT NULL DEFAULT 0 COMMENT '商城免费领取金币状态',
    `shop_free_diamond_time2` int(11) NOT NULL DEFAULT 0 COMMENT '商城免费领取金币状态2',
    `last_sign_time` int(11) NOT NULL DEFAULT 0 COMMENT '上次签到时间',
    `signed_time` int(11) NOT NULL DEFAULT -1 COMMENT '签到天数',
    `first_sign_time` int(11) NOT NULL DEFAULT 0 COMMENT '签到开启时间',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='泡泡龙';

DROP TABLE IF EXISTS `paopaolong_level`;
CREATE TABLE `paopaolong_level` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `level` int(11) NOT NULL DEFAULT 0 COMMENT '当前关卡',
    `sec` int(11) NOT NULL DEFAULT 0 COMMENT '时间',
    `lose` int(11) NOT NULL DEFAULT 0 COMMENT '失败次数',
    `score` bigint(11) NOT NULL DEFAULT 0 COMMENT '分数',
    `star` bigint(11) NOT NULL DEFAULT 0 COMMENT '星星数',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`, `level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='泡泡龙关卡';

DROP TABLE IF EXISTS `tanchishedazuozhan`;
CREATE TABLE `tanchishedazuozhan` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `skinarr` varchar(64) NOT NULL DEFAULT '1' COMMENT '皮肤列表',
    `skinid` int(11) NOT NULL DEFAULT 1 COMMENT '皮肤',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='贪吃蛇';

DROP TABLE IF EXISTS `tiantianpaoku`;
CREATE TABLE `tiantianpaoku` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `skinarr` varchar(64) NOT NULL DEFAULT '1' COMMENT '皮肤列表',
    `skinid` int(11) NOT NULL DEFAULT 1 COMMENT '皮肤',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='天天跑酷';


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