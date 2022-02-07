use bpbxq;
alter table gatestat modify column `host` varchar(32) NOT NULL DEFAULT '' COMMENT 'Host';

use caichengyu;
alter table gatestat modify column `host` varchar(32) NOT NULL DEFAULT '' COMMENT 'Host';

use daxigua;
alter table gatestat modify column `host` varchar(32) NOT NULL DEFAULT '' COMMENT 'Host';

use minigame;
alter table gatestat modify column `host` varchar(32) NOT NULL DEFAULT '' COMMENT 'Host';
