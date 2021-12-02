# Host: 127.0.0.1  (Version: 5.7.26)
# Date: 2021-12-02 22:20:19
# Generator: MySQL-Front 5.3  (Build 4.234)

/*!40101 SET NAMES utf8 */;

#
# Structure for table "auth_session"
#

DROP TABLE IF EXISTS `auth_session`;
CREATE TABLE `auth_session` (
  `time_hash` varchar(64) COLLATE utf8_unicode_ci NOT NULL,
  `last_visit` varchar(30) COLLATE utf8_unicode_ci NOT NULL,
  `username` varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`time_hash`),
  UNIQUE KEY `uix_auth_session_time_hash` (`time_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

#
# Data for table "auth_session"
#


#
# Structure for table "notification"
#

DROP TABLE IF EXISTS `notification`;
CREATE TABLE `notification` (
  `number` varchar(8) COLLATE utf8_unicode_ci NOT NULL,
  `user_id` varchar(8) COLLATE utf8_unicode_ci NOT NULL,
  `title` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  `content` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `status` varchar(4) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`number`),
  UNIQUE KEY `uix_notification_number` (`number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

#
# Data for table "notification"
#

INSERT INTO `notification` VALUES ('1','SA-0246','测试1','测试测试1','未读'),('2','SA-0246','测试2','测试测试2','未读');

#
# Structure for table "parts_overview"
#

DROP TABLE IF EXISTS `parts_overview`;
CREATE TABLE `parts_overview` (
  `parts_number` varchar(8) COLLATE utf8_unicode_ci NOT NULL,
  `parts_name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  `unit` varchar(6) COLLATE utf8_unicode_ci NOT NULL,
  `parts_cost` double(8,2) NOT NULL,
  PRIMARY KEY (`parts_number`),
  UNIQUE KEY `uix_parts_overview_parts_number` (`parts_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

#
# Data for table "parts_overview"
#


#
# Structure for table "repairman"
#

DROP TABLE IF EXISTS `repairman`;
CREATE TABLE `repairman` (
  `number` varchar(8) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  `password` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `type` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
  `current_work_hour` int(4) NOT NULL,
  `status` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`number`),
  UNIQUE KEY `uix_repairman_number` (`number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

#
# Data for table "repairman"
#

INSERT INTO `repairman` VALUES ('RE-0246','奚嘉骏','$2a$10$6NBVE248Ig9/TFz9QW9vquTsTH.HLlK2y8M.z4D/FaYA8wsnAnePq','机修工',0,'正常'),('RE-0248','伍慕庭','$2a$10$pfPLL4MczpKHpmaALWSbCO90wBz47v/Ce6ssN6YUAMVPNMWUsgHVe','电工',0,'正常'),('RE-0249','夏逸凡','$2a$10$XtptbBnCexoEUmAzLV1TCeqrJTepnKJZsiMgtpYNd76a9/ti6ti.m','钣金工',0,'正常'),('RE-0251','张霖锋','$2a$10$U159VfMPR9IzZpxlCvdLH.DbUL5Rr.ufYfReUDpWrWlV3wU57Nibi','喷漆工',0,'正常');

#
# Structure for table "salesman"
#

DROP TABLE IF EXISTS `salesman`;
CREATE TABLE `salesman` (
  `number` varchar(8) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  `password` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`number`),
  UNIQUE KEY `uix_salesman_number` (`number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

#
# Data for table "salesman"
#

INSERT INTO `salesman` VALUES ('SA-0246','奚嘉骏','$2a$10$twmzVNJP1m09vTIysSB8K.939/zQ666zv.w96ox2nbr70hqys6NqW'),('SA-0248','伍慕庭','$2a$10$ZUnlWs99.r2rT4vXekXLru4hRRYNVcBKZ/EOs/o.8NoCRh1muE9jq'),('SA-0249','夏逸凡','$2a$10$RNyBPg7160koaffiF.DdIO90b5Fzf1zIcWN8qLvl.PWjsjwjatJrK'),('SA-0251','张霖锋','$2a$10$dsp2i6FCvYzlFvU8pD8R7O/nj/ok5I/R.qSfKlBrrA.rbeQomQOvu');

#
# Structure for table "type_overview"
#

DROP TABLE IF EXISTS `type_overview`;
CREATE TABLE `type_overview` (
  `project_number` varchar(8) COLLATE utf8_unicode_ci NOT NULL,
  `project_name` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `type` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`project_number`,`type`),
  UNIQUE KEY `uix_type_overview_project_number` (`project_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

#
# Data for table "type_overview"
#


#
# Structure for table "user"
#

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `number` varchar(8) COLLATE utf8_unicode_ci NOT NULL,
  `password` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `property` varchar(4) COLLATE utf8_unicode_ci DEFAULT NULL,
  `discount_rate` int(2) DEFAULT NULL,
  `contact_person` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `contact_tel` varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`number`),
  UNIQUE KEY `uix_user_number` (`number`),
  UNIQUE KEY `uix_user_contact_tel` (`contact_tel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

#
# Data for table "user"
#


#
# Structure for table "vehicle"
#

DROP TABLE IF EXISTS `vehicle`;
CREATE TABLE `vehicle` (
  `number` varchar(17) COLLATE utf8_unicode_ci NOT NULL,
  `license_number` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
  `user_id` varchar(8) COLLATE utf8_unicode_ci NOT NULL,
  `color` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
  `model` varchar(40) COLLATE utf8_unicode_ci NOT NULL,
  `type` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
  `time` varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`number`,`user_id`),
  KEY `vehicle_user_id_user_number_foreign` (`user_id`),
  CONSTRAINT `vehicle_user_id_user_number_foreign` FOREIGN KEY (`user_id`) REFERENCES `user` (`number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

#
# Data for table "vehicle"
#


#
# Structure for table "attorney"
#

DROP TABLE IF EXISTS `attorney`;
CREATE TABLE `attorney` (
  `number` varchar(11) COLLATE utf8_unicode_ci NOT NULL,
  `user_id` varchar(8) COLLATE utf8_unicode_ci NOT NULL,
  `vehicle_number` varchar(17) COLLATE utf8_unicode_ci NOT NULL,
  `repair_type` varchar(4) COLLATE utf8_unicode_ci DEFAULT NULL,
  `classification` varchar(4) COLLATE utf8_unicode_ci DEFAULT NULL,
  `pay_method` varchar(4) COLLATE utf8_unicode_ci DEFAULT NULL,
  `start_time` varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  `salesman_id` varchar(8) COLLATE utf8_unicode_ci DEFAULT NULL,
  `predict_finish_time` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `actual_finish_time` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `rough_problem` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `specific_problem` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `progress` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
  `total_cost` double(6,2) NOT NULL,
  `start_petrol` double(5,2) NOT NULL,
  `start_mile` double(8,2) NOT NULL,
  `end_petrol` double(5,2) DEFAULT NULL,
  `end_mile` double(8,2) DEFAULT NULL,
  `out_range` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`number`),
  UNIQUE KEY `uix_attorney_number` (`number`),
  KEY `attorney_user_id_user_number_foreign` (`user_id`),
  KEY `attorney_vehicle_number_vehicle_number_foreign` (`vehicle_number`),
  KEY `attorney_salesman_id_salesman_number_foreign` (`salesman_id`),
  CONSTRAINT `attorney_salesman_id_salesman_number_foreign` FOREIGN KEY (`salesman_id`) REFERENCES `salesman` (`number`),
  CONSTRAINT `attorney_user_id_user_number_foreign` FOREIGN KEY (`user_id`) REFERENCES `user` (`number`),
  CONSTRAINT `attorney_vehicle_number_vehicle_number_foreign` FOREIGN KEY (`vehicle_number`) REFERENCES `vehicle` (`number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

#
# Data for table "attorney"
#


#
# Structure for table "arrangement"
#

DROP TABLE IF EXISTS `arrangement`;
CREATE TABLE `arrangement` (
  `order_number` varchar(11) COLLATE utf8_unicode_ci NOT NULL,
  `project_number` varchar(8) COLLATE utf8_unicode_ci NOT NULL,
  `predict_time` int(3) NOT NULL,
  `actual_t_ime` int(3) NOT NULL,
  `repairman_number` varchar(8) COLLATE utf8_unicode_ci NOT NULL,
  `parts_number` varchar(8) COLLATE utf8_unicode_ci NOT NULL,
  `parts_count` int(2) NOT NULL,
  `progress` varchar(6) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`order_number`,`project_number`,`repairman_number`),
  KEY `arrangement_repairman_number_repairman_number_foreign` (`repairman_number`),
  KEY `arrangement_project_number_type_overview_project_number_foreign` (`project_number`),
  KEY `arrangement_parts_number_parts_overview_parts_number_foreign` (`parts_number`),
  CONSTRAINT `arrangement_order_number_attorney_number_foreign` FOREIGN KEY (`order_number`) REFERENCES `attorney` (`number`),
  CONSTRAINT `arrangement_parts_number_parts_overview_parts_number_foreign` FOREIGN KEY (`parts_number`) REFERENCES `parts_overview` (`parts_number`),
  CONSTRAINT `arrangement_project_number_type_overview_project_number_foreign` FOREIGN KEY (`project_number`) REFERENCES `type_overview` (`project_number`),
  CONSTRAINT `arrangement_repairman_number_repairman_number_foreign` FOREIGN KEY (`repairman_number`) REFERENCES `repairman` (`number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

#
# Data for table "arrangement"
#

