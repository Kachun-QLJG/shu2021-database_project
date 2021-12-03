# Host: 127.0.0.1  (Version: 5.7.26)
# Date: 2021-12-03 23:49:08
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

INSERT INTO `auth_session` VALUES ('825689dcd46dc225f5a06e92a40f900fd2f8090ee7ba8a2273cfc91319ab0102','2021-12-03 08:03:32','RE-0246'),('bd3189ef7f9b3a724073b081fcdc16eba5f4a10a603aa779c846036b4e889480','2021-12-03 23:48:24','SA-0246');

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

INSERT INTO `notification` VALUES ('1','SA-0246','测试1','测试测试1','已读'),('2','SA-0246','测试2','测试测试2','已读');

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
  `project_spelling` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`project_number`,`type`),
  UNIQUE KEY `uix_type_overview_project_number` (`project_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

#
# Data for table "type_overview"
#

INSERT INTO `type_overview` VALUES ('00000001','发动机大修','维修员','fdjdx'),('00000002','发动机中修','维修员','fdjzx'),('00000003','更换缸盖垫','维修员','ghggd'),('00000004','更换机油','维修员','ghjy'),('00000005','更换机油格（油箱外）','维修员','ghjygyxw'),('00000006','清洗喷油嘴（支）','维修员','qxpyzz'),('00000007','更换喷油嘴（支）','维修员','ghpyzz'),('00000008','更换发动机水泵','维修员','ghfdjsb'),('00000009','更换水箱及冷却液','维修员','ghsxjlqy'),('00000010','更换起动机','维修员','ghqdj'),('00000011','更换发电机','维修员','ghfdj'),('00000012','拆装汽油箱','维修员','czqyx'),('00000013','更换汽油箱','维修员','ghqyx'),('00000014','更换火花塞','维修员','ghhhs'),('00000015','更换汽油泵','维修员','ghqyb'),('00000016','更换空气格','维修员','ghkqg'),('00000017','更换汽油格（油箱内）','维修员','ghqygyxn'),('00000018','更换发动机皮带/条','维修员','ghfdjpdt'),('00000019','更换发动机分电器','维修员','ghfdjfdq'),('00000020','检修发动机','维修员','jxfdj'),('00000021','前/后封箱漏油','维修员','qhfxly'),('00000022','检修发动机电喷系统','维修员','jxfdjdpxt'),('00000023','更换正时皮带','维修员','ghzspd'),('00000024','更换机油泵','维修员','ghjyb'),('00000025','更换机油泵（V6）','维修员','ghjybv'),('00000026','清洗节气门','维修员','qxjqm'),('00000027','更换水箱','维修员','ghsx'),('00000028','更换水管','维修员','ghsg'),('00000029','清洗油电路','维修员','qxydl'),('00000030','解除机头故障灯','维修员','jcjtgzd'),('00000031','更换半轴/边','维修员','ghbzb'),('00000032','更换自动变速箱油','维修员','ghzdbsxy'),('00000033','自动变速箱大修','维修员','zdbsxdx'),('00000034','手动变速箱大修','维修员','sdbsxdx'),('00000035','更换前避震器/缓冲胶（支）','维修员','ghqbzqhcjz'),('00000036','更换后避震器/缓冲胶（支）','维修员','ghhbzqhcjz'),('00000037','拆装轮胎','维修员','czlt'),('00000038','平衡轮胎','维修员','phlt'),('00000039','更换轮胎','维修员','ghlt'),('00000040','四轮定位','维修员','sldw'),('00000041','更换前/后刹车片','维修员','ghqhscp'),('00000042','更换刹车碟','维修员','ghscd'),('00000043','维修手刹系统','维修员','wxssxt'),('00000044','更换刹车总泵','维修员','ghsczb'),('00000045','拆仪表盘','维修员','cybp'),('00000046','更换制动分泵','维修员','ghzdfb'),('00000047','更换平衡杆球头','维修员','ghphgqt'),('00000048','更换刹车油及排空','维修员','ghscyjpk'),('00000049','更换离合器片（拆装波箱）','维修员','ghlhqpczbx'),('00000050','更换下球头','维修员','ghxqt'),('00000051','更换方向内/外球头','维修员','ghfxnwqt'),('00000052','检修四轮刹车系统/保养','维修员','jxslscxtby'),('00000053','更换前上/下摆臂','维修员','ghqsxbb'),('00000054','更换前/后轮轴承','维修员','ghqhlzc'),('00000055','更换转动轴十字节','维修员','ghzdzszj'),('00000056','大修空调系统','维修员','dxkdxt'),('00000057','维修雨刮连杆','维修员','wxyglg'),('00000058','更换膨胀阀','维修员','ghpzf'),('00000059','拆蒸发器','维修员','czfq'),('00000060','更换空调冷凝器','维修员','ghkdlnq'),('00000061','更换空调干燥瓶','维修员','ghkdgzp'),('00000062','更换电子扇','维修员','ghdzs'),('00000063','更换暖水阀','维修员','ghnsf'),('00000064','更换鼓风机','维修员','ghgfj'),('00000065','更换喷水器','维修员','ghpsq'),('00000066','更换雨刮片','维修员','ghygp'),('00000067','更换电池','维修员','ghdc'),('00000068','更换汽车电喇叭','维修员','ghqcdlb'),('00000069','更换车门升降器','维修员','ghcmsjq'),('00000070','更换转向组合开关','维修员','ghzxzhkg'),('00000071','维修前/后雨刮马达线路','维修员','wxqhygmdxl'),('00000072','维修ABS系统警报灯线路','维修员','wxabsxtjbdxl'),('00000073','维修SRS气囊系统警报灯线路','维修员','wxsrsqnxtjbdxl'),('00000074','更换中央门锁系统','维修员','ghzymsxt'),('00000075','更换暖水箱','维修员','ghnsx'),('00000076','更换前/后挡风玻璃','维修员','ghqhdfbl'),('00000077','更换发动机盖或后箱盖','维修员','ghfdjghhxg'),('00000078','更换前叶子板','维修员','ghqyzb'),('00000079','更换车门玻璃','维修员','ghcmbl'),('00000080','调校/更换车门锁','维修员','dxghcms'),('00000081','更换（拆装）前/后保险杠','维修员','ghczqhbxg'),('00000082','更换全车锁','维修员','ghqcs'),('00000083','拆装门饰板','维修员','czmsb'),('00000084','更换车门拉手','维修员','ghcmls'),('00000085','更换车门铰链','维修员','ghcmjl'),('00000086','更换顶棚天花','维修员','ghdpth'),('00000087','拆装全车座椅','维修员','czqczy'),('00000088','更换后视镜','维修员','ghhsj'),('00000089','拆门饰板','维修员','cmsb'),('00000090','更换前大灯总成','维修员','ghqddzc'),('00000091','更换后尾灯总成','维修员','ghhwdzc');

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

