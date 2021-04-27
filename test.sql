CREATE TABLE `channel_device_relation` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '通道和设备关联id',
  `channel_id` bigint(20) NOT NULL COMMENT '通道id',
  `device_rtsp_id` bigint(20) NOT NULL COMMENT '设备的rtsp_url_id',
  `type` tinyint(2) DEFAULT NULL COMMENT '设备的类型 0: 全局摄像头, 1: 车牌识别摄像头',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `time_version` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '数据的版本',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `agent` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '座席的id',
  `account` varchar(255) DEFAULT NULL COMMENT '座席名称',
  `pwd` varchar(255) NOT NULL COMMENT '座席账户的密码',
  `nick_name` varchar(100) DEFAULT NULL COMMENT '昵称',
  `employ_id` bigint(20) DEFAULT '0' COMMENT '座席对应的员工id',
  `sip_account` varchar(20) DEFAULT NULL COMMENT '座席登陆到sip的账号',
  `sip_pwd` varchar(50) DEFAULT NULL COMMENT '座席登陆到sip的密码',
  `sip_protocol` varchar(20) DEFAULT 'UDP' COMMENT 'sip的协议',
  `sip_server` varchar(50) DEFAULT '0.0.0.0' COMMENT 'sip服务器的ip或域名',
  `sip_port` int DEFAULT '5060' COMMENT 'sip的端口',
  `type` tinyint(2) DEFAULT NULL COMMENT '显示类型 0 园区座席 1 全网座席',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `time_version` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '数据的版本',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create TABLE 'channel' (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '通道id',
  `parking_id` bigint(20) DEFAULT '0' COMMENT '通道所属的停车场id',
  `name` varchar(255) DEFAULT NULL COMMENT '通道名称',
  `direction` tinyint(1) DEFAULT NULL COMMENT '通道的方向 0: 进口, 1: 出口',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `time_version` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '数据的版本',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `device_rtsp` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '设备id',
  `rtsp_url` varchar(255) DEFAULT NULL COMMENT '设备的rtsp_url',
  `name` varchar(255) DEFAULT NULL COMMENT '设备名称',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `time_version` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '数据的版本',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


ALTER TABLE `callcenter`.`device_rtsp` 
ADD UNIQUE INDEX `unique_device_rtsp`(`rtsp_url`) COMMENT 'rtsp_url不重复';