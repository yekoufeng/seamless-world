ALTER TABLE `playerrounddata` ADD COLUMN `UserName` varchar(50) DEFAULT NULL COMMENT '玩家名字';
ALTER TABLE `playerrounddata` ADD COLUMN `SkyBox` int(10) unsigned DEFAULT NULL COMMENT '本局天空盒';
ALTER TABLE `playerdaydata` ADD COLUMN `UserName` varchar(50) DEFAULT NULL COMMENT '玩家名字';

ALTER TABLE playerdaydata ADD INDEX ix_UID(UID);
ALTER TABLE playerdaydata ADD INDEX ix_UserName(UserName);
ALTER TABLE playerdaydata ADD INDEX ix_Model(Model);
ALTER TABLE playerdaydata ADD INDEX ix_StartTime(StartTime);

ALTER TABLE playersearchnum ADD INDEX ix_UserName(UserName);

ALTER TABLE playerrounddata ADD INDEX ix_UID(UID);
ALTER TABLE playerrounddata ADD INDEX ix_UserName(UserName);
ALTER TABLE playerrounddata ADD INDEX ix_StartTime(StartTime);

ALTER TABLE playersearchnum MODIFY COLUMN UID BIGINT(30) UNSIGNED;

ALTER TABLE playerdaydata MODIFY COLUMN DayID BIGINT(30) UNSIGNED;
ALTER TABLE playerdaydata MODIFY COLUMN UID BIGINT(30) UNSIGNED;

ALTER TABLE playerrounddata MODIFY COLUMN GameID BIGINT(30) UNSIGNED;
ALTER TABLE playerrounddata MODIFY COLUMN UID BIGINT(30) UNSIGNED;

ALTER TABLE `playerrounddata` ADD COLUMN `PlatID` INT(10) unsigned DEFAULT NULL COMMENT '平台ID';
ALTER TABLE `playerrounddata` ADD COLUMN `LoginChannel` INT(10) unsigned DEFAULT NULL COMMENT '登录渠道';