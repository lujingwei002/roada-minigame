CREATE DATABASE IF NOT EXISTS yangzhunongchang DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
USE yangzhunongchang;

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

DROP TABLE IF EXISTS `yangzhunongchang`;
CREATE TABLE `yangzhunongchang` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `farm_lv` int(11) NOT NULL DEFAULT 0 COMMENT '猪场等级',
    `farm_lv_name` varchar(64) NOT NULL DEFAULT '' COMMENT '猪场等级称号',
    `farm_lv_exp` int(11) NOT NULL DEFAULT 0 COMMENT '猪场等级累计经验',
    `farm_lv_exp_cur` int(11) NOT NULL DEFAULT 0 COMMENT '猪场等级当前等级累计经验',
    `award_num` int(11) NOT NULL DEFAULT 0 COMMENT '抽奖次数',
    `award_time` int(11) NOT NULL DEFAULT 0 COMMENT '抽奖时间',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='养猪农场';

DROP TABLE IF EXISTS `yangzhunongchang_item`;
CREATE TABLE `yangzhunongchang_item` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `id` int(11) NOT NULL DEFAULT 0 COMMENT '道具id',
    `num` int(11) NOT NULL DEFAULT 0 COMMENT '道具数量',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='养猪农场道具表';

DROP TABLE IF EXISTS `yangzhunongchang_pig`;
CREATE TABLE `yangzhunongchang_pig` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `id` varchar(32) NOT NULL DEFAULT 0 COMMENT '猪id',
    `data` varchar(256) NOT NULL DEFAULT 0 COMMENT '猪数据',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='养猪农场猪栏表';


DROP TABLE IF EXISTS `yangzhunongchang_breedpig`;
CREATE TABLE `yangzhunongchang_breedpig` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `id` varchar(32) NOT NULL DEFAULT 0 COMMENT '猪id',
    `data` varchar(256) NOT NULL DEFAULT 0 COMMENT '猪数据',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='养猪农场育种表';


DROP TABLE IF EXISTS `yangzhunongchang_usu`;
CREATE TABLE `yangzhunongchang_usu` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `id` int(11) NOT NULL DEFAULT 0 COMMENT '猪id',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='养猪农场图鉴表';

DROP TABLE IF EXISTS `yangzhunongchang_food`;
CREATE TABLE `yangzhunongchang_food` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `id` varchar(32) NOT NULL DEFAULT 0 COMMENT '食物id',
    `data` varchar(256) NOT NULL DEFAULT 0 COMMENT '猪数据',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='养猪农场食物表';

DROP TABLE IF EXISTS `yangzhunongchang_task`;
CREATE TABLE `yangzhunongchang_task` (
    `userid` bigint(11) AUTO_INCREMENT COMMENT '用户id',
    `id` varchar(32) NOT NULL DEFAULT 0 COMMENT '食物id',
    `index` int(11) NOT NULL DEFAULT 0 COMMENT '下标',
    `count` varchar(256) NOT NULL DEFAULT 0 COMMENT '进度',
    `createtime` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updatetime` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='养猪农场任务表';