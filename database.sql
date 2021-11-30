# Host: 127.0.0.1  (Version: 5.7.26)
# Date: 2021-11-30 19:54:53
# Generator: MySQL-Front 5.3  (Build 4.234)

/*!40101 SET NAMES utf8 */;

#
# Structure for table "arrangement"
#

DROP TABLE IF EXISTS `arrangement`;
CREATE TABLE `arrangement` (
  `order_number` varchar(11) NOT NULL,
  `project_number` varchar(8) NOT NULL,
  `repairman_number` varchar(8) NOT NULL,
  `parts_number` varchar(8) NOT NULL,
  `parts_count` int(2) NOT NULL,
  `progress` varchar(6) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

#
# Data for table "arrangement"
#

/*!40000 ALTER TABLE `arrangement` DISABLE KEYS */;
/*!40000 ALTER TABLE `arrangement` ENABLE KEYS */;

#
# Structure for table "attorney"
#

DROP TABLE IF EXISTS `attorney`;
CREATE TABLE `attorney` (
  `number` varchar(11) NOT NULL,
  `user_id` varchar(8) NOT NULL,
  `vehicle_number` varchar(17) NOT NULL,
  `repair_type` varchar(4) DEFAULT NULL,
  `classification` varchar(4) DEFAULT NULL,
  `pay_method` varchar(4) DEFAULT NULL,
  `start_time` varchar(20) NOT NULL,
  `salesman_id` varchar(8) DEFAULT NULL,
  `predict_finish_time` varchar(20) DEFAULT NULL,
  `actual_finish_time` varchar(20) DEFAULT NULL,
  `rough_problem` varchar(255) NOT NULL,
  `specific_problem` varchar(255) NOT NULL,
  `progress` varchar(10) NOT NULL,
  `total_cost` double(6,2) NOT NULL,
  PRIMARY KEY (`number`),
  UNIQUE KEY `uix_attorney_number` (`number`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

#
# Data for table "attorney"
#

/*!40000 ALTER TABLE `attorney` DISABLE KEYS */;
/*!40000 ALTER TABLE `attorney` ENABLE KEYS */;

#
# Structure for table "auth_session"
#

DROP TABLE IF EXISTS `auth_session`;
CREATE TABLE `auth_session` (
  `time_hash` varchar(64) NOT NULL,
  `last_visit` varchar(30) NOT NULL,
  `username` varchar(20) NOT NULL,
  PRIMARY KEY (`time_hash`),
  UNIQUE KEY `uix_auth_session_time_hash` (`time_hash`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

#
# Data for table "auth_session"
#

/*!40000 ALTER TABLE `auth_session` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth_session` ENABLE KEYS */;

#
# Structure for table "parts_overview"
#

DROP TABLE IF EXISTS `parts_overview`;
CREATE TABLE `parts_overview` (
  `parts_number` varchar(8) NOT NULL,
  `parts_name` varchar(50) NOT NULL,
  `parts_cost` double(8,2) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

#
# Data for table "parts_overview"
#

/*!40000 ALTER TABLE `parts_overview` DISABLE KEYS */;
/*!40000 ALTER TABLE `parts_overview` ENABLE KEYS */;

#
# Structure for table "repairman"
#

DROP TABLE IF EXISTS `repairman`;
CREATE TABLE `repairman` (
  `number` varchar(8) NOT NULL,
  `name` varchar(20) NOT NULL,
  `password` varchar(255) NOT NULL,
  `type` varchar(10) NOT NULL,
  `current_work_hour` int(4) NOT NULL,
  PRIMARY KEY (`number`),
  UNIQUE KEY `uix_repairman_number` (`number`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

#
# Data for table "repairman"
#

/*!40000 ALTER TABLE `repairman` DISABLE KEYS */;
INSERT INTO `repairman` VALUES ('RE-0246','奚嘉骏','$2a$10$6NBVE248Ig9/TFz9QW9vquTsTH.HLlK2y8M.z4D/FaYA8wsnAnePq','机修工',0),('RE-0248','伍慕庭','$2a$10$pfPLL4MczpKHpmaALWSbCO90wBz47v/Ce6ssN6YUAMVPNMWUsgHVe','电工',0),('RE-0249','夏逸凡','$2a$10$XtptbBnCexoEUmAzLV1TCeqrJTepnKJZsiMgtpYNd76a9/ti6ti.m','钣金工',0),('RE-0251','张霖锋','$2a$10$U159VfMPR9IzZpxlCvdLH.DbUL5Rr.ufYfReUDpWrWlV3wU57Nibi','喷漆工',0);
/*!40000 ALTER TABLE `repairman` ENABLE KEYS */;

#
# Structure for table "salesman"
#

DROP TABLE IF EXISTS `salesman`;
CREATE TABLE `salesman` (
  `number` varchar(8) NOT NULL,
  `name` varchar(20) NOT NULL,
  `password` varchar(255) NOT NULL,
  PRIMARY KEY (`number`),
  UNIQUE KEY `uix_salesman_number` (`number`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

#
# Data for table "salesman"
#

/*!40000 ALTER TABLE `salesman` DISABLE KEYS */;
INSERT INTO `salesman` VALUES ('SA-0246','奚嘉骏','$2a$10$twmzVNJP1m09vTIysSB8K.939/zQ666zv.w96ox2nbr70hqys6NqW'),('SA-0248','伍慕庭','$2a$10$ZUnlWs99.r2rT4vXekXLru4hRRYNVcBKZ/EOs/o.8NoCRh1muE9jq'),('SA-0249','夏逸凡','$2a$10$RNyBPg7160koaffiF.DdIO90b5Fzf1zIcWN8qLvl.PWjsjwjatJrK'),('SA-0251','张霖锋','$2a$10$dsp2i6FCvYzlFvU8pD8R7O/nj/ok5I/R.qSfKlBrrA.rbeQomQOvu');
/*!40000 ALTER TABLE `salesman` ENABLE KEYS */;

#
# Structure for table "type_overview"
#

DROP TABLE IF EXISTS `type_overview`;
CREATE TABLE `type_overview` (
  `project_number` varchar(8) NOT NULL,
  `project_name` varchar(100) NOT NULL,
  `work_hour` int(3) NOT NULL,
  `type` varchar(10) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

#
# Data for table "type_overview"
#

/*!40000 ALTER TABLE `type_overview` DISABLE KEYS */;
/*!40000 ALTER TABLE `type_overview` ENABLE KEYS */;

#
# Structure for table "user"
#

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `number` varchar(8) NOT NULL,
  `password` varchar(255) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `property` varchar(4) DEFAULT NULL,
  `discount_rate` int(2) DEFAULT NULL,
  `contact_person` varchar(10) DEFAULT NULL,
  `contact_tel` varchar(20) NOT NULL,
  PRIMARY KEY (`number`),
  UNIQUE KEY `uix_user_number` (`number`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

#
# Data for table "user"
#

/*!40000 ALTER TABLE `user` DISABLE KEYS */;
/*!40000 ALTER TABLE `user` ENABLE KEYS */;

#
# Structure for table "vehicle"
#

DROP TABLE IF EXISTS `vehicle`;
CREATE TABLE `vehicle` (
  `number` varchar(17) NOT NULL,
  `license_number` varchar(10) NOT NULL,
  `user_id` varchar(8) NOT NULL,
  `color` varchar(10) NOT NULL,
  `model` varchar(40) NOT NULL,
  `type` varchar(10) NOT NULL,
  PRIMARY KEY (`number`),
  UNIQUE KEY `uix_vehicle_number` (`number`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

#
# Data for table "vehicle"
#

/*!40000 ALTER TABLE `vehicle` DISABLE KEYS */;
/*!40000 ALTER TABLE `vehicle` ENABLE KEYS */;
