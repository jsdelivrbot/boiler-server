/*
 Navicat Premium Data Transfer

 Source Server         : MySQL @Aliyun
 Source Server Type    : MySQL
 Source Server Version : 50634
 Source Host           : rm-uf63a6pe498w6mi00o.mysql.rds.aliyuncs.com
 Source Database       : boiler_main

 Target Server Type    : MySQL
 Target Server Version : 50634
 File Encoding         : utf-8

 Date: 01/02/2018 12:41:30 PM
*/

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `address`
-- ----------------------------
DROP TABLE IF EXISTS `address`;
CREATE TABLE `address` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `location_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `address` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `zip_code` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `longitude` double NOT NULL DEFAULT '0',
  `latitude` double NOT NULL DEFAULT '0',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `address_created_date` (`created_date`),
  KEY `address_is_deleted` (`is_deleted`),
  KEY `address_name` (`name`),
  KEY `address_is_demo` (`is_demo`),
  KEY `address_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `application`
-- ----------------------------
DROP TABLE IF EXISTS `application`;
CREATE TABLE `application` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `platform` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `app` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `identity` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `domain` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `path` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `app_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `app_secret` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `api_token` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `aes_key` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `origin_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`uid`),
  KEY `application_name` (`name`),
  KEY `application_created_date` (`created_date`),
  KEY `application_is_demo` (`is_demo`),
  KEY `application_is_deleted` (`is_deleted`),
  KEY `application_platform` (`platform`),
  KEY `application_app` (`app`),
  KEY `application_identity` (`identity`),
  KEY `application_domain` (`domain`),
  KEY `application_app_id` (`app_id`),
  KEY `application_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler`
-- ----------------------------
DROP TABLE IF EXISTS `boiler`;
CREATE TABLE `boiler` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `fuel_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `template_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `factory_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `enterprise_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `installed_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `address_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `factory_number` varchar(30) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `register_code` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `register_org_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `certificate_number` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `device_code` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `model_code` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `evaporating_capacity` bigint(20) NOT NULL DEFAULT '0',
  `rated_capacity_level_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `contact_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `terminal_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `rated_capacity_range_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `terminal_code` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `terminal_set_id` int(11) NOT NULL DEFAULT '0',
  `form_id` bigint(20) NOT NULL DEFAULT '0',
  `medium_id` bigint(20) NOT NULL DEFAULT '0',
  `usage_id` bigint(20) NOT NULL DEFAULT '0',
  `inspect_date_next` datetime DEFAULT NULL,
  `inspect_inner_date_next` datetime DEFAULT NULL,
  `inspect_outer_date_next` datetime DEFAULT NULL,
  `inspect_valve_date_next` datetime DEFAULT NULL,
  `inspect_gauge_date_next` datetime DEFAULT NULL,
  PRIMARY KEY (`uid`),
  KEY `boiler_is_demo` (`is_demo`),
  KEY `boiler_created_date` (`created_date`),
  KEY `boiler_is_deleted` (`is_deleted`),
  KEY `boiler_name` (`name`),
  KEY `boiler_rated_capacity_range_id` (`rated_capacity_range_id`),
  KEY `boiler_contact_id` (`contact_id`),
  KEY `boiler_terminal_id` (`terminal_id`),
  KEY `boiler_updated_date` (`updated_date`),
  KEY `boiler_terminal_code` (`terminal_code`),
  KEY `boiler_terminal_set_id` (`terminal_set_id`),
  KEY `boiler_evaporating_capacity` (`evaporating_capacity`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_alarm`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_alarm`;
CREATE TABLE `boiler_alarm` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `start_date` datetime NOT NULL,
  `end_date` datetime DEFAULT NULL,
  `confirmed_date` datetime DEFAULT NULL,
  `confirmed_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `verified_date` datetime DEFAULT NULL,
  `verified_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `trigger_rule_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) NOT NULL DEFAULT '0',
  `state` int(11) NOT NULL DEFAULT '0',
  `priority` int(11) NOT NULL DEFAULT '0',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `need_send` tinyint(1) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `boiler_id` (`boiler_id`,`parameter_id`,`trigger_rule_id`,`start_date`),
  KEY `boiler_alarm_name` (`name`),
  KEY `boiler_alarm_created_date` (`created_date`),
  KEY `boiler_alarm_updated_date` (`updated_date`),
  KEY `boiler_alarm_is_demo` (`is_demo`),
  KEY `boiler_alarm_is_deleted` (`is_deleted`),
  KEY `boiler_alarm_boiler_id` (`boiler_id`),
  KEY `boiler_alarm_parameter_id` (`parameter_id`),
  KEY `boiler_alarm_start_date` (`start_date`),
  KEY `boiler_alarm_end_date` (`end_date`),
  KEY `boiler_alarm_confirmed_date` (`confirmed_date`),
  KEY `boiler_alarm_confirmed_by_id` (`confirmed_by_id`),
  KEY `boiler_alarm_verified_date` (`verified_date`),
  KEY `boiler_alarm_verified_by_id` (`verified_by_id`),
  KEY `boiler_alarm_trigger_rule_id` (`trigger_rule_id`),
  KEY `boiler_alarm_alarm_level` (`alarm_level`),
  KEY `boiler_alarm_status` (`state`),
  KEY `boiler_alarm_priority` (`priority`),
  KEY `boiler_alarm_boiler_id_parameter_id_trigger_rule_id` (`boiler_id`,`parameter_id`,`trigger_rule_id`),
  KEY `boiler_alarm_state` (`state`),
  KEY `boiler_alarm_need_send` (`need_send`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_alarm_feedback`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_alarm_feedback`;
CREATE TABLE `boiler_alarm_feedback` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `type` int(11) NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `state` int(11) NOT NULL DEFAULT '0',
  `content` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`uid`),
  KEY `boiler_alarm_feedback_name` (`name`),
  KEY `boiler_alarm_feedback_created_date` (`created_date`),
  KEY `boiler_alarm_feedback_is_demo` (`is_demo`),
  KEY `boiler_alarm_feedback_is_deleted` (`is_deleted`),
  KEY `boiler_alarm_feedback_type` (`type`),
  KEY `boiler_alarm_feedback_alarm_id` (`alarm_id`),
  KEY `boiler_alarm_feedback_updated_date` (`updated_date`),
  KEY `boiler_alarm_feedback_state` (`state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_alarm_history`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_alarm_history`;
CREATE TABLE `boiler_alarm_history` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `start_date` datetime NOT NULL,
  `end_date` datetime DEFAULT NULL,
  `confirmed_date` datetime DEFAULT NULL,
  `confirmed_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `verified_date` datetime DEFAULT NULL,
  `verified_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `trigger_rule_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) NOT NULL DEFAULT '0',
  `priority` int(11) NOT NULL DEFAULT '0',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `boiler_id` (`boiler_id`,`parameter_id`,`trigger_rule_id`,`start_date`),
  KEY `boiler_alarm_history_name` (`name`),
  KEY `boiler_alarm_history_created_date` (`created_date`),
  KEY `boiler_alarm_history_updated_date` (`updated_date`),
  KEY `boiler_alarm_history_is_demo` (`is_demo`),
  KEY `boiler_alarm_history_is_deleted` (`is_deleted`),
  KEY `boiler_alarm_history_boiler_id` (`boiler_id`),
  KEY `boiler_alarm_history_parameter_id` (`parameter_id`),
  KEY `boiler_alarm_history_start_date` (`start_date`),
  KEY `boiler_alarm_history_end_date` (`end_date`),
  KEY `boiler_alarm_history_confirmed_date` (`confirmed_date`),
  KEY `boiler_alarm_history_confirmed_by_id` (`confirmed_by_id`),
  KEY `boiler_alarm_history_verified_date` (`verified_date`),
  KEY `boiler_alarm_history_verified_by_id` (`verified_by_id`),
  KEY `boiler_alarm_history_trigger_rule_id` (`trigger_rule_id`),
  KEY `boiler_alarm_history_alarm_level` (`alarm_level`),
  KEY `boiler_alarm_history_priority` (`priority`),
  KEY `boiler_alarm_history_boiler_id_parameter_id_trigger_rule_id` (`boiler_id`,`parameter_id`,`trigger_rule_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_calculate_parameter`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_calculate_parameter`;
CREATE TABLE `boiler_calculate_parameter` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `coal_qnetvar` double NOT NULL DEFAULT '0',
  `coal_aar` double NOT NULL DEFAULT '0',
  `coal_mar` double NOT NULL DEFAULT '0',
  `coal_vdaf` double NOT NULL DEFAULT '0',
  `coal_clz` double NOT NULL DEFAULT '0',
  `coal_clm` double NOT NULL DEFAULT '0',
  `coal_cfh` double NOT NULL DEFAULT '0',
  `coal_ded` double NOT NULL DEFAULT '0',
  `coal_dsc` double NOT NULL DEFAULT '0',
  `coal_alz` double NOT NULL DEFAULT '0',
  `coal_alm` double NOT NULL DEFAULT '0',
  `coal_afh` double NOT NULL DEFAULT '0',
  `coal_q3` double NOT NULL DEFAULT '0',
  `coal_m` double NOT NULL DEFAULT '0',
  `coal_n` double NOT NULL DEFAULT '0',
  `coal_tlz` double NOT NULL DEFAULT '0',
  `coal_ct_lz` double NOT NULL DEFAULT '0',
  `gas_ded` double NOT NULL DEFAULT '0',
  `gas_dsc` double NOT NULL DEFAULT '0',
  `gas_apy` double NOT NULL DEFAULT '0',
  `gas_q3` double NOT NULL DEFAULT '0',
  `gas_m` double NOT NULL DEFAULT '0',
  `gas_n` double NOT NULL DEFAULT '0',
  `conf_param1` double NOT NULL DEFAULT '0',
  `conf_param2` double NOT NULL DEFAULT '0',
  `conf_param3` double NOT NULL DEFAULT '0',
  `conf_param4` double NOT NULL DEFAULT '0',
  `conf_param5` double NOT NULL DEFAULT '0',
  `conf_param6` double NOT NULL DEFAULT '0',
  `alarm_threshold1` double NOT NULL DEFAULT '0',
  `alarm_threshold2` double NOT NULL DEFAULT '0',
  `alarm_threshold3` double NOT NULL DEFAULT '0',
  `alarm_threshold4` double NOT NULL DEFAULT '0',
  `alarm_threshold5` double NOT NULL DEFAULT '0',
  `alarm_threshold6` double NOT NULL DEFAULT '0',
  `alarm_threshold7` double NOT NULL DEFAULT '0',
  `alarm_threshold8` double NOT NULL DEFAULT '0',
  `reserved1` double NOT NULL DEFAULT '0',
  `reserved2` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `reserved3` double NOT NULL DEFAULT '0',
  `reserved4` double NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `boiler_calculate_parameter_name` (`name`),
  KEY `boiler_calculate_parameter_created_date` (`created_date`),
  KEY `boiler_calculate_parameter_updated_date` (`updated_date`),
  KEY `boiler_calculate_parameter_is_demo` (`is_demo`),
  KEY `boiler_calculate_parameter_is_deleted` (`is_deleted`),
  KEY `boiler_calculate_parameter_boiler_id` (`boiler_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_calculate_result`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_calculate_result`;
CREATE TABLE `boiler_calculate_result` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `fuel_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `based_parameter_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `qnetvar` double NOT NULL DEFAULT '0',
  `aar` double NOT NULL DEFAULT '0',
  `mar` double NOT NULL DEFAULT '0',
  `vdaf` double NOT NULL DEFAULT '0',
  `clz` double NOT NULL DEFAULT '0',
  `clm` double NOT NULL DEFAULT '0',
  `cfh` double NOT NULL DEFAULT '0',
  `ded` double NOT NULL DEFAULT '0',
  `dsc` double NOT NULL DEFAULT '0',
  `alz` double NOT NULL DEFAULT '0',
  `alm` double NOT NULL DEFAULT '0',
  `afh` double NOT NULL DEFAULT '0',
  `tlz` double NOT NULL DEFAULT '0',
  `ct_lz` double NOT NULL DEFAULT '0',
  `m` double NOT NULL DEFAULT '0',
  `n` double NOT NULL DEFAULT '0',
  `q2` double NOT NULL DEFAULT '0',
  `q3` double NOT NULL DEFAULT '0',
  `q4` double NOT NULL DEFAULT '0',
  `q5` double NOT NULL DEFAULT '0',
  `q6` double NOT NULL DEFAULT '0',
  `excess_air` double NOT NULL DEFAULT '0',
  `heat` double NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `boiler_calculate_result_name` (`name`),
  KEY `boiler_calculate_result_created_date` (`created_date`),
  KEY `boiler_calculate_result_updated_date` (`updated_date`),
  KEY `boiler_calculate_result_is_demo` (`is_demo`),
  KEY `boiler_calculate_result_is_deleted` (`is_deleted`),
  KEY `boiler_calculate_result_boiler_id` (`boiler_id`),
  KEY `boiler_calculate_result_fuel_id` (`fuel_id`),
  KEY `boiler_calculate_result_based_parameter_id` (`based_parameter_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1859330 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_config`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_config`;
CREATE TABLE `boiler_config` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `is_generate_data` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `boiler_id` (`boiler_id`),
  KEY `boiler_config_name` (`name`),
  KEY `boiler_config_created_date` (`created_date`),
  KEY `boiler_config_updated_date` (`updated_date`),
  KEY `boiler_config_is_demo` (`is_demo`),
  KEY `boiler_config_is_deleted` (`is_deleted`),
  KEY `boiler_config_boiler_id` (`boiler_id`),
  KEY `boiler_config_is_generate_data` (`is_generate_data`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_fuel_record`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_fuel_record`;
CREATE TABLE `boiler_fuel_record` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `fuel_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `start_date` datetime NOT NULL,
  `end_date` datetime NOT NULL,
  `total_flow` double NOT NULL DEFAULT '0',
  `fuel_amount` double NOT NULL DEFAULT '0',
  `rate` double NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `boiler_fuel_record_name` (`name`),
  KEY `boiler_fuel_record_created_date` (`created_date`),
  KEY `boiler_fuel_record_updated_date` (`updated_date`),
  KEY `boiler_fuel_record_is_demo` (`is_demo`),
  KEY `boiler_fuel_record_is_deleted` (`is_deleted`),
  KEY `boiler_fuel_record_boiler_id` (`boiler_id`),
  KEY `boiler_fuel_record_fuel_id` (`fuel_id`),
  KEY `boiler_fuel_record_start_date` (`start_date`),
  KEY `boiler_fuel_record_end_date` (`end_date`),
  KEY `boiler_fuel_record_rate` (`rate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_m160`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_m160`;
CREATE TABLE `boiler_m160` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `Boiler_sn` int(11) NOT NULL,
  `Boiler_status_code` char(2) CHARACTER SET utf8 NOT NULL,
  `Boiler_term_id` char(6) CHARACTER SET utf8 NOT NULL,
  `TS` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20261 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_m163`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_m163`;
CREATE TABLE `boiler_m163` (
  `Boiler_term_id` char(6) COLLATE utf8_unicode_ci NOT NULL,
  `Boiler_boiler_id` char(2) COLLATE utf8_unicode_ci NOT NULL,
  `Term_sys_time` char(19) CHARACTER SET utf8 NOT NULL,
  `Boiler_data_fmt_ver` char(2) CHARACTER SET utf8 NOT NULL,
  `Boiler_sn` int(11) NOT NULL,
  `Boiler_status_code` char(2) CHARACTER SET utf8 NOT NULL,
  `Temper1_channel` int(11) NOT NULL,
  `Temper2_channel` int(11) NOT NULL,
  `Temper3_channel` int(11) NOT NULL,
  `Temper4_channel` int(11) NOT NULL,
  `Temper5_channel` int(11) NOT NULL,
  `Temper6_channel` int(11) NOT NULL,
  `Temper7_channel` int(11) NOT NULL,
  `Temper8_channel` int(11) NOT NULL,
  `Temper9_channel` int(11) NOT NULL,
  `Temper10_channel` int(11) NOT NULL,
  `Temper11_channel` int(11) NOT NULL,
  `Temper12_channel` int(11) NOT NULL,
  `Analog1_channel` int(11) NOT NULL,
  `Analog2_channel` int(11) NOT NULL,
  `Analog3_channel` int(11) NOT NULL,
  `Analog4_channel` int(11) NOT NULL,
  `Analog5_channel` int(11) NOT NULL,
  `Analog6_channel` int(11) NOT NULL,
  `Analog7_channel` int(11) NOT NULL,
  `Analog8_channel` int(11) NOT NULL,
  `Analog9_channel` int(11) NOT NULL,
  `Analog10_channel` int(11) NOT NULL,
  `Analog11_channel` int(11) NOT NULL,
  `Analog12_channel` int(11) NOT NULL,
  `Switch_in_1_16_channel` int(11) NOT NULL,
  `Switch_in_17_32_channel` int(11) NOT NULL,
  `Switch_out_1_16_channel` int(11) NOT NULL,
  `C1_calculate_parm` decimal(10,4) NOT NULL,
  `C2_calculate_parm` decimal(10,4) NOT NULL,
  `C3_calculate_parm` decimal(10,4) NOT NULL,
  `C4_calculate_parm` decimal(10,4) NOT NULL,
  `C5_calculate_parm` decimal(10,4) NOT NULL,
  `C6_calculate_parm` decimal(10,4) NOT NULL,
  `C7_calculate_parm` decimal(10,4) NOT NULL,
  `C8_calculate_parm` decimal(10,4) NOT NULL,
  `C9_calculate_parm` decimal(10,4) NOT NULL,
  `C10_calculate_parm` decimal(10,4) NOT NULL,
  `C11_calculate_parm` decimal(10,4) NOT NULL,
  `C12_calculate_parm` decimal(10,4) NOT NULL,
  `Boiler_term_err_code` char(2) CHARACTER SET utf8 NOT NULL,
  `Reserved1_filler` int(11) NOT NULL,
  `Reserved2_filler` int(11) NOT NULL,
  `Reserved3_filler` int(11) NOT NULL,
  `Reserved4_filler` int(11) NOT NULL,
  `TS` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `need_reload` tinyint(1) unsigned NOT NULL DEFAULT '1',
  `uid` char(36) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`Boiler_term_id`,`Boiler_boiler_id`,`Term_sys_time`,`Boiler_data_fmt_ver`,`Boiler_sn`),
  UNIQUE KEY `uid_unique` (`uid`) USING BTREE,
  KEY `reload_index` (`need_reload`) USING BTREE,
  KEY `terminal_index` (`Boiler_term_id`,`Boiler_boiler_id`) USING BTREE,
  KEY `time_index` (`TS`) USING BTREE,
  KEY `terminal_code_index` (`Boiler_term_id`) USING BTREE,
  KEY `terminal_set_index` (`Boiler_boiler_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_m176`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_m176`;
CREATE TABLE `boiler_m176` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `Boiler_sn` int(11) NOT NULL,
  `Boiler_status_code` char(2) CHARACTER SET utf8 NOT NULL,
  `Boiler_term_id` char(6) CHARACTER SET utf8 NOT NULL,
  `Boiler_resp_data` char(1) CHARACTER SET utf8 NOT NULL,
  `Boiler_pf_pwd` int(2) NOT NULL,
  `TS` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=46 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_maintenance`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_maintenance`;
CREATE TABLE `boiler_maintenance` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `inspect_date` datetime NOT NULL,
  `target` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `inspector` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `content` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `attachment` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_done` tinyint(1) NOT NULL DEFAULT '0',
  `burner` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `import_grate` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `water_softener` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `water_pump` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `boiler_body` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `energy_saver` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `air_pre_heater` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `dust_catcher` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `draught_fan` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`uid`),
  KEY `boiler_maintenance_created_date` (`created_date`),
  KEY `boiler_maintenance_is_deleted` (`is_deleted`),
  KEY `boiler_maintenance_name` (`name`),
  KEY `boiler_maintenance_is_demo` (`is_demo`),
  KEY `boiler_maintenance_boiler_id` (`boiler_id`),
  KEY `boiler_maintenance_inspect_date` (`inspect_date`),
  KEY `boiler_maintenance_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_medium`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_medium`;
CREATE TABLE `boiler_medium` (
  `id` bigint(20) NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `medium_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `boiler_medium_name` (`name`),
  KEY `boiler_medium_created_date` (`created_date`),
  KEY `boiler_medium_updated_date` (`updated_date`),
  KEY `boiler_medium_is_demo` (`is_demo`),
  KEY `boiler_medium_is_deleted` (`is_deleted`),
  KEY `boiler_medium_medium_id` (`medium_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_message_subscriber`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_message_subscriber`;
CREATE TABLE `boiler_message_subscriber` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `user_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  KEY `boiler_message_subscriber_boiler_id` (`boiler_id`),
  KEY `boiler_message_subscriber_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=149 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime`;
CREATE TABLE `boiler_runtime` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `value` bigint(20) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`,`parameter_id`,`boiler_id`,`created_date`),
  KEY `boiler_runtime_is_demo` (`is_demo`),
  KEY `boiler_runtime_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_created_date` (`created_date`) USING HASH,
  KEY `boiler_runtime_updated_date` (`updated_date`) USING HASH,
  KEY `boiler_runtime_boiler_id` (`boiler_id`) USING HASH,
  KEY `boiler_runtime_parameter_id` (`parameter_id`) USING HASH,
  KEY `boiler_runtime_alarm_id` (`alarm_id`) USING HASH,
  KEY `boiler_runtime_name` (`name`),
  KEY `boiler_runtime_uid` (`uid`) USING BTREE,
  KEY `boiler_runtime_main` (`boiler_id`,`parameter_id`,`created_date`) USING BTREE,
  KEY `boiler_runtime_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=284974829 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci
/*!50100 PARTITION BY RANGE (YEARWEEK(`created_date`))
SUBPARTITION BY HASH (`parameter_id`)
SUBPARTITIONS 6
(PARTITION p0 VALUES LESS THAN (201720) ENGINE = InnoDB,
 PARTITION p1 VALUES LESS THAN (201724) ENGINE = InnoDB,
 PARTITION p2 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_environment_temperature`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_environment_temperature`;
CREATE TABLE `boiler_runtime_cache_environment_temperature` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) NOT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `cache_env_temper_boiler_id` (`boiler_id`,`created_date`) USING BTREE,
  UNIQUE KEY `cache_env_temper_main` (`boiler_id`,`parameter_id`,`created_date`) USING BTREE,
  KEY `boiler_runtime_cache_environment_temperature_name` (`name`),
  KEY `boiler_runtime_cache_environment_temperature_created_date` (`created_date`),
  KEY `boiler_runtime_cache_environment_temperature_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_environment_temperature_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_environment_temperature_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_environment_temperature_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_environment_temperature_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_environment_temperature_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_environment_temperature_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_environment_temperature_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_environment_temperature_alarm_description` (`alarm_description`)
) ENGINE=InnoDB AUTO_INCREMENT=17991988 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_excess_air`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_excess_air`;
CREATE TABLE `boiler_runtime_cache_excess_air` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) NOT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `cache_exc_air_boiler_id` (`boiler_id`,`created_date`) USING BTREE,
  UNIQUE KEY `cache_exc_air_main` (`boiler_id`,`parameter_id`,`created_date`) USING BTREE,
  KEY `boiler_runtime_cache_excess_air_name` (`name`),
  KEY `boiler_runtime_cache_excess_air_created_date` (`created_date`),
  KEY `boiler_runtime_cache_excess_air_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_excess_air_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_excess_air_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_excess_air_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_excess_air_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_excess_air_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_excess_air_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_excess_air_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_excess_air_alarm_description` (`alarm_description`)
) ENGINE=InnoDB AUTO_INCREMENT=16466288 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_flow`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_flow`;
CREATE TABLE `boiler_runtime_cache_flow` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) NOT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `cache_flow_boiler_id` (`boiler_id`,`created_date`) USING BTREE,
  UNIQUE KEY `cache_flow_main` (`boiler_id`,`parameter_id`,`created_date`) USING BTREE,
  KEY `boiler_runtime_cache_flow_name` (`name`),
  KEY `boiler_runtime_cache_flow_created_date` (`created_date`),
  KEY `boiler_runtime_cache_flow_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_flow_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_flow_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_flow_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_flow_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_flow_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_flow_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_flow_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_flow_alarm_description` (`alarm_description`)
) ENGINE=InnoDB AUTO_INCREMENT=18361344 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_flow_daily`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_flow_daily`;
CREATE TABLE `boiler_runtime_cache_flow_daily` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) DEFAULT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `date` datetime NOT NULL,
  `hours` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `boiler_id` (`boiler_id`,`date`),
  UNIQUE KEY `boiler_id_2` (`boiler_id`,`parameter_id`,`date`),
  KEY `boiler_runtime_cache_flow_daily_name` (`name`),
  KEY `boiler_runtime_cache_flow_daily_created_date` (`created_date`),
  KEY `boiler_runtime_cache_flow_daily_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_flow_daily_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_flow_daily_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_flow_daily_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_flow_daily_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_flow_daily_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_flow_daily_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_flow_daily_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_flow_daily_alarm_description` (`alarm_description`),
  KEY `boiler_runtime_cache_flow_daily_date` (`date`)
) ENGINE=InnoDB AUTO_INCREMENT=78663578 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_heat`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_heat`;
CREATE TABLE `boiler_runtime_cache_heat` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) NOT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `cache_heat_boiler_id` (`boiler_id`,`created_date`) USING BTREE,
  UNIQUE KEY `cache_heat_main` (`boiler_id`,`parameter_id`,`created_date`) USING BTREE,
  KEY `boiler_runtime_cache_heat_name` (`name`),
  KEY `boiler_runtime_cache_heat_created_date` (`created_date`),
  KEY `boiler_runtime_cache_heat_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_heat_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_heat_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_heat_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_heat_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_heat_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_heat_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_heat_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_heat_alarm_description` (`alarm_description`)
) ENGINE=InnoDB AUTO_INCREMENT=16173506 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_heat_daily`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_heat_daily`;
CREATE TABLE `boiler_runtime_cache_heat_daily` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) DEFAULT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `date` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `boiler_id` (`boiler_id`,`date`),
  UNIQUE KEY `boiler_id_2` (`boiler_id`,`parameter_id`,`date`),
  KEY `boiler_runtime_cache_heat_daily_name` (`name`),
  KEY `boiler_runtime_cache_heat_daily_created_date` (`created_date`),
  KEY `boiler_runtime_cache_heat_daily_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_heat_daily_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_heat_daily_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_heat_daily_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_heat_daily_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_heat_daily_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_heat_daily_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_heat_daily_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_heat_daily_alarm_description` (`alarm_description`),
  KEY `boiler_runtime_cache_heat_daily_date` (`date`)
) ENGINE=InnoDB AUTO_INCREMENT=42092262 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_history`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_history`;
CREATE TABLE `boiler_runtime_cache_history` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `p1001` double NOT NULL DEFAULT '0',
  `p1002` double NOT NULL DEFAULT '0',
  `p1003` double NOT NULL DEFAULT '0',
  `p1004` double NOT NULL DEFAULT '0',
  `p1005` double NOT NULL DEFAULT '0',
  `p1006` double NOT NULL DEFAULT '0',
  `p1007` double NOT NULL DEFAULT '0',
  `p1008` double NOT NULL DEFAULT '0',
  `p1009` double NOT NULL DEFAULT '0',
  `p1010` double NOT NULL DEFAULT '0',
  `p1011` double NOT NULL DEFAULT '0',
  `p1012` double NOT NULL DEFAULT '0',
  `p1013` double NOT NULL DEFAULT '0',
  `p1014` double NOT NULL DEFAULT '0',
  `p1015` double NOT NULL DEFAULT '0',
  `p1016` double NOT NULL DEFAULT '0',
  `p1017` double NOT NULL DEFAULT '0',
  `p1018` double NOT NULL DEFAULT '0',
  `p1019` double NOT NULL DEFAULT '0',
  `p1020` double NOT NULL DEFAULT '0',
  `p1021` double NOT NULL DEFAULT '0',
  `p1022` double NOT NULL DEFAULT '0',
  `p1023` double NOT NULL DEFAULT '0',
  `p1024` double NOT NULL DEFAULT '0',
  `p1025` double NOT NULL DEFAULT '0',
  `p1026` double NOT NULL DEFAULT '0',
  `p1027` double NOT NULL DEFAULT '0',
  `p1028` double NOT NULL DEFAULT '0',
  `p1029` double NOT NULL DEFAULT '0',
  `p1030` double NOT NULL DEFAULT '0',
  `p1031` double NOT NULL DEFAULT '0',
  `p1032` double NOT NULL DEFAULT '0',
  `p1033` double NOT NULL DEFAULT '0',
  `p1034` double NOT NULL DEFAULT '0',
  `p1035` double NOT NULL DEFAULT '0',
  `p1036` double NOT NULL DEFAULT '0',
  `p1037` double NOT NULL DEFAULT '0',
  `p1038` double NOT NULL DEFAULT '0',
  `p1039` double NOT NULL DEFAULT '0',
  `p1040` double NOT NULL DEFAULT '0',
  `p1041` double NOT NULL DEFAULT '0',
  `p1042` double NOT NULL DEFAULT '0',
  `p1043` double NOT NULL DEFAULT '0',
  `p1044` double NOT NULL DEFAULT '0',
  `p1045` double NOT NULL DEFAULT '0',
  `p1046` double NOT NULL DEFAULT '0',
  `p1047` double NOT NULL DEFAULT '0',
  `p1048` double NOT NULL DEFAULT '0',
  `p1049` double NOT NULL DEFAULT '0',
  `p1050` double NOT NULL DEFAULT '0',
  `p1051` double NOT NULL DEFAULT '0',
  `p1052` double NOT NULL DEFAULT '0',
  `p1053` double NOT NULL DEFAULT '0',
  `p1054` double NOT NULL DEFAULT '0',
  `p1055` double NOT NULL DEFAULT '0',
  `p1056` double NOT NULL DEFAULT '0',
  `p1057` double NOT NULL DEFAULT '0',
  `p1058` double NOT NULL DEFAULT '0',
  `p1059` double NOT NULL DEFAULT '0',
  `p1060` double NOT NULL DEFAULT '0',
  `p1061` double NOT NULL DEFAULT '0',
  `p1062` double NOT NULL DEFAULT '0',
  `p1063` double NOT NULL DEFAULT '0',
  `p1064` double NOT NULL DEFAULT '0',
  `p1065` double NOT NULL DEFAULT '0',
  `p1066` double NOT NULL DEFAULT '0',
  `p1067` double NOT NULL DEFAULT '0',
  `p1068` double NOT NULL DEFAULT '0',
  `p1069` double NOT NULL DEFAULT '0',
  `p1070` double NOT NULL DEFAULT '0',
  `p1071` double NOT NULL DEFAULT '0',
  `p1072` double NOT NULL DEFAULT '0',
  `p1073` double NOT NULL DEFAULT '0',
  `p1101` double NOT NULL DEFAULT '0',
  `p1102` double NOT NULL DEFAULT '0',
  `p1103` double NOT NULL DEFAULT '0',
  `p1104` double NOT NULL DEFAULT '0',
  `p1105` double NOT NULL DEFAULT '0',
  `p1106` double NOT NULL DEFAULT '0',
  `p1107` double NOT NULL DEFAULT '0',
  `p1108` double NOT NULL DEFAULT '0',
  `p1109` double NOT NULL DEFAULT '0',
  `p1110` double NOT NULL DEFAULT '0',
  `p1201` double NOT NULL DEFAULT '0',
  `p1202` double NOT NULL DEFAULT '0',
  `p1203` double NOT NULL DEFAULT '0',
  `p1204` double NOT NULL DEFAULT '0',
  `p1205` double NOT NULL DEFAULT '0',
  `p1206` double NOT NULL DEFAULT '0',
  `p1207` double NOT NULL DEFAULT '0',
  `a1001` int(11) NOT NULL DEFAULT '0',
  `a1002` int(11) NOT NULL DEFAULT '0',
  `a1003` int(11) NOT NULL DEFAULT '0',
  `a1004` int(11) NOT NULL DEFAULT '0',
  `a1005` int(11) NOT NULL DEFAULT '0',
  `a1006` int(11) NOT NULL DEFAULT '0',
  `a1007` int(11) NOT NULL DEFAULT '0',
  `a1008` int(11) NOT NULL DEFAULT '0',
  `a1009` int(11) NOT NULL DEFAULT '0',
  `a1010` int(11) NOT NULL DEFAULT '0',
  `a1011` int(11) NOT NULL DEFAULT '0',
  `a1012` int(11) NOT NULL DEFAULT '0',
  `a1013` int(11) NOT NULL DEFAULT '0',
  `a1014` int(11) NOT NULL DEFAULT '0',
  `a1015` int(11) NOT NULL DEFAULT '0',
  `a1016` int(11) NOT NULL DEFAULT '0',
  `a1017` int(11) NOT NULL DEFAULT '0',
  `a1018` int(11) NOT NULL DEFAULT '0',
  `a1019` int(11) NOT NULL DEFAULT '0',
  `a1020` int(11) NOT NULL DEFAULT '0',
  `a1021` int(11) NOT NULL DEFAULT '0',
  `a1022` int(11) NOT NULL DEFAULT '0',
  `a1023` int(11) NOT NULL DEFAULT '0',
  `a1024` int(11) NOT NULL DEFAULT '0',
  `a1025` int(11) NOT NULL DEFAULT '0',
  `a1026` int(11) NOT NULL DEFAULT '0',
  `a1027` int(11) NOT NULL DEFAULT '0',
  `a1028` int(11) NOT NULL DEFAULT '0',
  `a1029` int(11) NOT NULL DEFAULT '0',
  `a1030` int(11) NOT NULL DEFAULT '0',
  `a1031` int(11) NOT NULL DEFAULT '0',
  `a1032` int(11) NOT NULL DEFAULT '0',
  `a1033` int(11) NOT NULL DEFAULT '0',
  `a1034` int(11) NOT NULL DEFAULT '0',
  `a1035` int(11) NOT NULL DEFAULT '0',
  `a1036` int(11) NOT NULL DEFAULT '0',
  `a1037` int(11) NOT NULL DEFAULT '0',
  `a1038` int(11) NOT NULL DEFAULT '0',
  `a1039` int(11) NOT NULL DEFAULT '0',
  `a1040` int(11) NOT NULL DEFAULT '0',
  `a1041` int(11) NOT NULL DEFAULT '0',
  `a1042` int(11) NOT NULL DEFAULT '0',
  `a1043` int(11) NOT NULL DEFAULT '0',
  `a1044` int(11) NOT NULL DEFAULT '0',
  `a1045` int(11) NOT NULL DEFAULT '0',
  `a1046` int(11) NOT NULL DEFAULT '0',
  `a1047` int(11) NOT NULL DEFAULT '0',
  `a1048` int(11) NOT NULL DEFAULT '0',
  `a1049` int(11) NOT NULL DEFAULT '0',
  `a1050` int(11) NOT NULL DEFAULT '0',
  `a1051` int(11) NOT NULL DEFAULT '0',
  `a1052` int(11) NOT NULL DEFAULT '0',
  `a1053` int(11) NOT NULL DEFAULT '0',
  `a1054` int(11) NOT NULL DEFAULT '0',
  `a1055` int(11) NOT NULL DEFAULT '0',
  `a1056` int(11) NOT NULL DEFAULT '0',
  `a1057` int(11) NOT NULL DEFAULT '0',
  `a1058` int(11) NOT NULL DEFAULT '0',
  `a1059` int(11) NOT NULL DEFAULT '0',
  `a1060` int(11) NOT NULL DEFAULT '0',
  `a1061` int(11) NOT NULL DEFAULT '0',
  `a1062` int(11) NOT NULL DEFAULT '0',
  `a1063` int(11) NOT NULL DEFAULT '0',
  `a1064` int(11) NOT NULL DEFAULT '0',
  `a1065` int(11) NOT NULL DEFAULT '0',
  `a1066` int(11) NOT NULL DEFAULT '0',
  `a1067` int(11) NOT NULL DEFAULT '0',
  `a1068` int(11) NOT NULL DEFAULT '0',
  `a1069` int(11) NOT NULL DEFAULT '0',
  `a1070` int(11) NOT NULL DEFAULT '0',
  `a1071` int(11) NOT NULL DEFAULT '0',
  `a1072` int(11) NOT NULL DEFAULT '0',
  `a1073` int(11) NOT NULL DEFAULT '0',
  `a1101` int(11) NOT NULL DEFAULT '0',
  `a1102` int(11) NOT NULL DEFAULT '0',
  `a1103` int(11) NOT NULL DEFAULT '0',
  `a1104` int(11) NOT NULL DEFAULT '0',
  `a1105` int(11) NOT NULL DEFAULT '0',
  `a1106` int(11) NOT NULL DEFAULT '0',
  `a1107` int(11) NOT NULL DEFAULT '0',
  `a1108` int(11) NOT NULL DEFAULT '0',
  `a1109` int(11) NOT NULL DEFAULT '0',
  `a1110` int(11) NOT NULL DEFAULT '0',
  `a1201` int(11) NOT NULL DEFAULT '0',
  `a1202` int(11) NOT NULL DEFAULT '0',
  `a1203` int(11) NOT NULL DEFAULT '0',
  `a1204` int(11) NOT NULL DEFAULT '0',
  `a1205` int(11) NOT NULL DEFAULT '0',
  `a1206` int(11) NOT NULL DEFAULT '0',
  `a1207` int(11) NOT NULL DEFAULT '0',
  `p1080` double NOT NULL DEFAULT '0',
  `a1080` int(11) NOT NULL DEFAULT '0',
  `p1090` double NOT NULL,
  `p1091` double NOT NULL,
  `p1092` double NOT NULL,
  `p1093` double NOT NULL,
  `p1094` double NOT NULL,
  `p1095` double NOT NULL,
  `p1096` double NOT NULL,
  `p1097` double NOT NULL,
  `p1098` double NOT NULL,
  `p1099` double NOT NULL,
  `a1090` int(11) NOT NULL,
  `a1091` int(11) NOT NULL,
  `a1092` int(11) NOT NULL,
  `a1093` int(11) NOT NULL,
  `a1094` int(11) NOT NULL,
  `a1095` int(11) NOT NULL,
  `a1096` int(11) NOT NULL,
  `a1097` int(11) NOT NULL,
  `a1098` int(11) NOT NULL,
  `a1099` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `boiler_runtime_cache_history_name` (`name`),
  KEY `boiler_runtime_cache_history_created_date` (`created_date`),
  KEY `boiler_runtime_cache_history_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_history_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_history_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_history_boiler_id` (`boiler_id`),
  KEY `history_main` (`boiler_id`,`updated_date`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=15266250 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_instant`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_instant`;
CREATE TABLE `boiler_runtime_cache_instant` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) NOT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_valid` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `cache_instant_main` (`boiler_id`,`parameter_id`) USING BTREE,
  KEY `boiler_runtime_cache_instant_name` (`name`),
  KEY `boiler_runtime_cache_instant_created_date` (`created_date`),
  KEY `boiler_runtime_cache_instant_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_instant_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_instant_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_instant_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_instant_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_instant_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_instant_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_instant_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_instant_boiler_id_parameter_id_updated_date` (`boiler_id`,`parameter_id`,`updated_date`),
  KEY `boiler_runtime_cache_instant_alarm_description` (`alarm_description`)
) ENGINE=InnoDB AUTO_INCREMENT=202896892 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_smoke_component`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_smoke_component`;
CREATE TABLE `boiler_runtime_cache_smoke_component` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) NOT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `cache_smoke_component_unique` (`boiler_id`,`parameter_id`,`created_date`) USING BTREE,
  KEY `boiler_runtime_cache_smoke_component_name` (`name`),
  KEY `boiler_runtime_cache_smoke_component_created_date` (`created_date`),
  KEY `boiler_runtime_cache_smoke_component_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_smoke_component_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_smoke_component_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_smoke_component_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_smoke_component_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_smoke_component_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_smoke_component_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_smoke_component_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_smoke_component_alarm_description` (`alarm_description`)
) ENGINE=InnoDB AUTO_INCREMENT=60938367 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_smoke_temperature`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_smoke_temperature`;
CREATE TABLE `boiler_runtime_cache_smoke_temperature` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) NOT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `cache_smoke_temper_main` (`boiler_id`,`parameter_id`,`created_date`) USING BTREE,
  KEY `boiler_runtime_cache_smoke_temperature_name` (`name`),
  KEY `boiler_runtime_cache_smoke_temperature_created_date` (`created_date`),
  KEY `boiler_runtime_cache_smoke_temperature_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_smoke_temperature_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_smoke_temperature_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_smoke_temperature_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_smoke_temperature_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_smoke_temperature_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_smoke_temperature_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_smoke_temperature_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_smoke_temperature_alarm_description` (`alarm_description`)
) ENGINE=InnoDB AUTO_INCREMENT=35882396 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_status`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_status`;
CREATE TABLE `boiler_runtime_cache_status` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) DEFAULT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `boiler_id` (`boiler_id`,`parameter_id`,`created_date`),
  KEY `boiler_runtime_cache_status_name` (`name`),
  KEY `boiler_runtime_cache_status_created_date` (`created_date`),
  KEY `boiler_runtime_cache_status_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_status_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_status_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_status_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_status_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_status_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_status_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_status_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_status_alarm_description` (`alarm_description`)
) ENGINE=InnoDB AUTO_INCREMENT=31870833 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_status_running`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_status_running`;
CREATE TABLE `boiler_runtime_cache_status_running` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) DEFAULT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `date` datetime NOT NULL,
  `duration` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `boiler_id` (`boiler_id`,`parameter_id`,`created_date`),
  KEY `boiler_runtime_cache_status_running_name` (`name`),
  KEY `boiler_runtime_cache_status_running_created_date` (`created_date`),
  KEY `boiler_runtime_cache_status_running_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_status_running_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_status_running_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_status_running_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_status_running_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_status_running_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_status_running_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_status_running_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_status_running_alarm_description` (`alarm_description`),
  KEY `boiler_runtime_cache_status_running_date` (`date`)
) ENGINE=InnoDB AUTO_INCREMENT=2093861 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_steam_pressure`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_steam_pressure`;
CREATE TABLE `boiler_runtime_cache_steam_pressure` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) NOT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `steam_press_boiler_id` (`boiler_id`,`created_date`) USING BTREE,
  UNIQUE KEY `steam_press_unique` (`boiler_id`,`parameter_id`,`created_date`) USING BTREE,
  KEY `boiler_runtime_cache_steam_pressure_name` (`name`),
  KEY `boiler_runtime_cache_steam_pressure_created_date` (`created_date`),
  KEY `boiler_runtime_cache_steam_pressure_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_steam_pressure_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_steam_pressure_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_steam_pressure_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_steam_pressure_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_steam_pressure_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_steam_pressure_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_steam_pressure_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_steam_pressure_alarm_description` (`alarm_description`)
) ENGINE=InnoDB AUTO_INCREMENT=18097169 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_steam_temperature`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_steam_temperature`;
CREATE TABLE `boiler_runtime_cache_steam_temperature` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) NOT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `steam_temper_boiler_id` (`boiler_id`,`created_date`) USING BTREE,
  UNIQUE KEY `steam_temper_unique` (`boiler_id`,`parameter_id`,`created_date`) USING BTREE,
  KEY `boiler_runtime_cache_steam_temperature_name` (`name`),
  KEY `boiler_runtime_cache_steam_temperature_created_date` (`created_date`),
  KEY `boiler_runtime_cache_steam_temperature_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_steam_temperature_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_steam_temperature_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_steam_temperature_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_steam_temperature_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_steam_temperature_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_steam_temperature_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_steam_temperature_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_steam_temperature_alarm_description` (`alarm_description`)
) ENGINE=InnoDB AUTO_INCREMENT=17939062 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_cache_water_temperature`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_cache_water_temperature`;
CREATE TABLE `boiler_runtime_cache_water_temperature` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `runtime_id` bigint(20) NOT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `parameter_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `alarm_level` int(11) DEFAULT NULL,
  `alarm_description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `water_temper_unique` (`boiler_id`,`parameter_id`,`created_date`) USING BTREE,
  KEY `boiler_runtime_cache_water_temperature_name` (`name`),
  KEY `boiler_runtime_cache_water_temperature_created_date` (`created_date`),
  KEY `boiler_runtime_cache_water_temperature_updated_date` (`updated_date`),
  KEY `boiler_runtime_cache_water_temperature_is_demo` (`is_demo`),
  KEY `boiler_runtime_cache_water_temperature_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_cache_water_temperature_runtime_id` (`runtime_id`),
  KEY `boiler_runtime_cache_water_temperature_boiler_id` (`boiler_id`),
  KEY `boiler_runtime_cache_water_temperature_parameter_id` (`parameter_id`),
  KEY `boiler_runtime_cache_water_temperature_alarm_id` (`alarm_id`),
  KEY `boiler_runtime_cache_water_temperature_alarm_level` (`alarm_level`),
  KEY `boiler_runtime_cache_water_temperature_alarm_description` (`alarm_description`)
) ENGINE=InnoDB AUTO_INCREMENT=35800754 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_runtime_copy`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_copy`;
CREATE TABLE `boiler_runtime_copy` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `parameter_id` bigint(20) NOT NULL,
  `alarm_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `value` bigint(20) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`,`parameter_id`,`boiler_id`,`created_date`),
  UNIQUE KEY `boiler_runtime_main` (`boiler_id`,`parameter_id`,`created_date`) USING HASH,
  KEY `boiler_runtime_is_demo` (`is_demo`),
  KEY `boiler_runtime_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_created_date` (`created_date`) USING HASH,
  KEY `boiler_runtime_updated_date` (`updated_date`) USING HASH,
  KEY `boiler_runtime_boiler_id` (`boiler_id`) USING HASH,
  KEY `boiler_runtime_parameter_id` (`parameter_id`) USING HASH,
  KEY `boiler_runtime_alarm_id` (`alarm_id`) USING HASH,
  KEY `boiler_runtime_name` (`name`),
  KEY `boiler_runtime_uid` (`uid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=183360995 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci
/*!50100 PARTITION BY RANGE (YEARWEEK(`created_date`))
SUBPARTITION BY HASH (`parameter_id`)
SUBPARTITIONS 6
(PARTITION p0 VALUES LESS THAN (201720) ENGINE = InnoDB,
 PARTITION p1 VALUES LESS THAN (201724) ENGINE = InnoDB,
 PARTITION p2 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */;

-- ----------------------------
--  Table structure for `boiler_runtime_history`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_runtime_history`;
CREATE TABLE `boiler_runtime_history` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `json_data` longtext COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `boiler_runtime_history_main` (`boiler_id`,`created_date`) USING BTREE,
  KEY `boiler_runtime_history_name` (`name`),
  KEY `boiler_runtime_history_created_date` (`created_date`),
  KEY `boiler_runtime_history_updated_date` (`updated_date`),
  KEY `boiler_runtime_history_is_demo` (`is_demo`),
  KEY `boiler_runtime_history_is_deleted` (`is_deleted`),
  KEY `boiler_runtime_history_boiler_id` (`boiler_id`)
) ENGINE=InnoDB AUTO_INCREMENT=20126608 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_template`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_template`;
CREATE TABLE `boiler_template` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `template_id` int(11) NOT NULL DEFAULT '0',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `boiler_template_template_id` (`template_id`),
  KEY `boiler_template_created_date` (`created_date`),
  KEY `boiler_template_is_deleted` (`is_deleted`),
  KEY `boiler_template_name` (`name`),
  KEY `boiler_template_is_demo` (`is_demo`),
  KEY `boiler_template_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_term_status`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_term_status`;
CREATE TABLE `boiler_term_status` (
  `Boiler_term_id` char(6) COLLATE utf8_unicode_ci NOT NULL,
  `Boiler_term_ip` char(15) COLLATE utf8_unicode_ci NOT NULL,
  `Boiler_term_pwd` smallint(6) NOT NULL,
  `Boiler_term_status` char(10) COLLATE utf8_unicode_ci NOT NULL,
  `TS` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`Boiler_term_id`),
  UNIQUE KEY `boiler_term_status_unique_idx1` (`Boiler_term_ip`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_terminal_combined`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_terminal_combined`;
CREATE TABLE `boiler_terminal_combined` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `terminal_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `terminal_code` bigint(20) NOT NULL DEFAULT '0',
  `terminal_set_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `boiler_id` (`boiler_id`,`terminal_id`),
  UNIQUE KEY `terminal_code` (`terminal_code`,`terminal_set_id`),
  KEY `boiler_terminal_combined_boiler_id` (`boiler_id`),
  KEY `boiler_terminal_combined_terminal_id` (`terminal_id`),
  KEY `boiler_terminal_combined_terminal_code` (`terminal_code`),
  KEY `boiler_terminal_combined_terminal_set_id` (`terminal_set_id`)
) ENGINE=InnoDB AUTO_INCREMENT=239955 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_type`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_type`;
CREATE TABLE `boiler_type` (
  `id` bigint(20) NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `type_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `boiler_type_name` (`name`),
  KEY `boiler_type_created_date` (`created_date`),
  KEY `boiler_type_updated_date` (`updated_date`),
  KEY `boiler_type_is_demo` (`is_demo`),
  KEY `boiler_type_is_deleted` (`is_deleted`),
  KEY `boiler_type_type_id` (`type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_type_form`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_type_form`;
CREATE TABLE `boiler_type_form` (
  `id` bigint(20) NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `type_id` bigint(20) NOT NULL,
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `form_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `boiler_type_form_name` (`name`),
  KEY `boiler_type_form_created_date` (`created_date`),
  KEY `boiler_type_form_updated_date` (`updated_date`),
  KEY `boiler_type_form_is_demo` (`is_demo`),
  KEY `boiler_type_form_is_deleted` (`is_deleted`),
  KEY `boiler_type_form_type_id` (`type_id`),
  KEY `boiler_type_form_form_id` (`form_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `boiler_usage`
-- ----------------------------
DROP TABLE IF EXISTS `boiler_usage`;
CREATE TABLE `boiler_usage` (
  `id` bigint(20) NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `usage_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `boiler_usage_name` (`name`),
  KEY `boiler_usage_created_date` (`created_date`),
  KEY `boiler_usage_updated_date` (`updated_date`),
  KEY `boiler_usage_is_demo` (`is_demo`),
  KEY `boiler_usage_is_deleted` (`is_deleted`),
  KEY `boiler_usage_usage_id` (`usage_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `contact`
-- ----------------------------
DROP TABLE IF EXISTS `contact`;
CREATE TABLE `contact` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `contact_id` int(11) NOT NULL DEFAULT '0',
  `phone_number` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `mobile_number` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `email` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `contact_contact_id` (`contact_id`),
  KEY `contact_created_date` (`created_date`),
  KEY `contact_is_deleted` (`is_deleted`),
  KEY `contact_name` (`name`),
  KEY `contact_is_demo` (`is_demo`),
  KEY `contact_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `dialogue`
-- ----------------------------
DROP TABLE IF EXISTS `dialogue`;
CREATE TABLE `dialogue` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `dialogue_id` bigint(20) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `dialogue_name` (`name`),
  KEY `dialogue_created_date` (`created_date`),
  KEY `dialogue_is_demo` (`is_demo`),
  KEY `dialogue_is_deleted` (`is_deleted`),
  KEY `dialogue_dialogue_id` (`dialogue_id`),
  KEY `dialogue_status` (`status`),
  KEY `dialogue_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `dialogue_comment`
-- ----------------------------
DROP TABLE IF EXISTS `dialogue_comment`;
CREATE TABLE `dialogue_comment` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `dialogue_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `from_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `to_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `content` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `attachment` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`uid`),
  KEY `dialogue_comment_name` (`name`),
  KEY `dialogue_comment_created_date` (`created_date`),
  KEY `dialogue_comment_is_demo` (`is_demo`),
  KEY `dialogue_comment_is_deleted` (`is_deleted`),
  KEY `dialogue_comment_dialogue_id` (`dialogue_id`),
  KEY `dialogue_comment_from_id` (`from_id`),
  KEY `dialogue_comment_to_id` (`to_id`),
  KEY `dialogue_comment_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `fuel`
-- ----------------------------
DROP TABLE IF EXISTS `fuel`;
CREATE TABLE `fuel` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `fuel_id` int(11) NOT NULL DEFAULT '0',
  `q_net_v_ar_min` bigint(20) NOT NULL DEFAULT '0',
  `q_net_v_ar_max` bigint(20) NOT NULL DEFAULT '0',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `type_id` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `fuel_fuel_id` (`fuel_id`),
  KEY `fuel_created_date` (`created_date`),
  KEY `fuel_is_deleted` (`is_deleted`),
  KEY `fuel_name` (`name`),
  KEY `fuel_is_demo` (`is_demo`),
  KEY `fuel_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `fuel_type`
-- ----------------------------
DROP TABLE IF EXISTS `fuel_type`;
CREATE TABLE `fuel_type` (
  `id` bigint(20) NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `type_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fuel_type_name` (`name`),
  KEY `fuel_type_created_date` (`created_date`),
  KEY `fuel_type_updated_date` (`updated_date`),
  KEY `fuel_type_is_demo` (`is_demo`),
  KEY `fuel_type_is_deleted` (`is_deleted`),
  KEY `fuel_type_type_id` (`type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `location`
-- ----------------------------
DROP TABLE IF EXISTS `location`;
CREATE TABLE `location` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `location_id` bigint(20) NOT NULL DEFAULT '0',
  `super_id` bigint(20) NOT NULL DEFAULT '0',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `location_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`uid`),
  KEY `location_location_id` (`location_id`),
  KEY `location_super_id` (`super_id`),
  KEY `location_created_date` (`created_date`),
  KEY `location_is_deleted` (`is_deleted`),
  KEY `location_name` (`name`),
  KEY `location_is_demo` (`is_demo`),
  KEY `location_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `message`
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `formatter_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `header` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `content_length` int(11) NOT NULL DEFAULT '0',
  `content` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `serial_number` int(11) NOT NULL DEFAULT '0',
  `status` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `identifier` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `response_code` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `server_date` datetime NOT NULL,
  `password` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `crc` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `tailer` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `body` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `message_created_date` (`created_date`),
  KEY `message_is_deleted` (`is_deleted`),
  KEY `message_name` (`name`),
  KEY `message_is_demo` (`is_demo`),
  KEY `message_updated_date` (`updated_date`),
  KEY `message_server_date` (`server_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `message_16bit`
-- ----------------------------
DROP TABLE IF EXISTS `message_16bit`;
CREATE TABLE `message_16bit` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_date` datetime NOT NULL,
  `terminal_code` int(11) NOT NULL DEFAULT '0',
  `terminal_set_id` int(11) NOT NULL DEFAULT '0',
  `server_date` varchar(19) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `version` int(11) NOT NULL DEFAULT '0',
  `serial_number` int(11) NOT NULL DEFAULT '0',
  `error_code` varchar(2) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `channel_temperature1` int(11) NOT NULL DEFAULT '0',
  `channel_temperature2` int(11) NOT NULL DEFAULT '0',
  `channel_temperature3` int(11) NOT NULL DEFAULT '0',
  `channel_temperature4` int(11) NOT NULL DEFAULT '0',
  `channel_temperature5` int(11) NOT NULL DEFAULT '0',
  `channel_temperature6` int(11) NOT NULL DEFAULT '0',
  `channel_temperature7` int(11) NOT NULL DEFAULT '0',
  `channel_temperature8` int(11) NOT NULL DEFAULT '0',
  `channel_temperature9` int(11) NOT NULL DEFAULT '0',
  `channel_temperature10` int(11) NOT NULL DEFAULT '0',
  `channel_temperature11` int(11) NOT NULL DEFAULT '0',
  `channel_temperature12` int(11) NOT NULL DEFAULT '0',
  `channel_analog1` int(11) NOT NULL DEFAULT '0',
  `channel_analog2` int(11) NOT NULL DEFAULT '0',
  `channel_analog3` int(11) NOT NULL DEFAULT '0',
  `channel_analog4` int(11) NOT NULL DEFAULT '0',
  `channel_analog5` int(11) NOT NULL DEFAULT '0',
  `channel_analog6` int(11) NOT NULL DEFAULT '0',
  `channel_analog7` int(11) NOT NULL DEFAULT '0',
  `channel_analog8` int(11) NOT NULL DEFAULT '0',
  `channel_analog9` int(11) NOT NULL DEFAULT '0',
  `channel_analog10` int(11) NOT NULL DEFAULT '0',
  `channel_analog11` int(11) NOT NULL DEFAULT '0',
  `channel_analog12` int(11) NOT NULL DEFAULT '0',
  `channel_switch_in1_to16` int(11) NOT NULL DEFAULT '0',
  `channel_switch_in17_to32` int(11) NOT NULL DEFAULT '0',
  `channel_switch_out1_to16` int(11) NOT NULL DEFAULT '0',
  `channel_calculate1` double NOT NULL DEFAULT '0',
  `channel_calculate2` double NOT NULL DEFAULT '0',
  `channel_calculate3` double NOT NULL DEFAULT '0',
  `channel_calculate4` double NOT NULL DEFAULT '0',
  `channel_calculate5` double NOT NULL DEFAULT '0',
  `channel_calculate6` double NOT NULL DEFAULT '0',
  `channel_calculate7` double NOT NULL DEFAULT '0',
  `channel_calculate8` double NOT NULL DEFAULT '0',
  `channel_calculate9` double NOT NULL DEFAULT '0',
  `channel_calculate10` double NOT NULL DEFAULT '0',
  `channel_calculate11` double NOT NULL DEFAULT '0',
  `channel_calculate12` double NOT NULL DEFAULT '0',
  `reserved1` int(11) NOT NULL DEFAULT '0',
  `reserved2` int(11) NOT NULL DEFAULT '0',
  `reserved3` int(11) NOT NULL DEFAULT '0',
  `reserved4` int(11) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `message_16bit_created_date` (`created_date`),
  KEY `message_16bit_terminal_code` (`terminal_code`),
  KEY `message_16bit_terminal_set_id` (`terminal_set_id`),
  KEY `message_16bit_version` (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `message_32bit`
-- ----------------------------
DROP TABLE IF EXISTS `message_32bit`;
CREATE TABLE `message_32bit` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_date` datetime NOT NULL,
  `terminal_code` int(11) NOT NULL DEFAULT '0',
  `terminal_set_id` int(11) NOT NULL DEFAULT '0',
  `server_date` varchar(19) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `version` int(11) NOT NULL DEFAULT '0',
  `serial_number` int(11) NOT NULL DEFAULT '0',
  `error_code` varchar(2) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `channel1` int(11) NOT NULL DEFAULT '0',
  `channel2` int(11) NOT NULL DEFAULT '0',
  `channel3` int(11) NOT NULL DEFAULT '0',
  `channel4` int(11) NOT NULL DEFAULT '0',
  `channel5` int(11) NOT NULL DEFAULT '0',
  `channel6` int(11) NOT NULL DEFAULT '0',
  `channel7` int(11) NOT NULL DEFAULT '0',
  `channel8` int(11) NOT NULL DEFAULT '0',
  `channel9` int(11) NOT NULL DEFAULT '0',
  `channel10` int(11) NOT NULL DEFAULT '0',
  `channel11` int(11) NOT NULL DEFAULT '0',
  `channel12` int(11) NOT NULL DEFAULT '0',
  `channel13` int(11) NOT NULL DEFAULT '0',
  `channel14` int(11) NOT NULL DEFAULT '0',
  `channel15` int(11) NOT NULL DEFAULT '0',
  `channel16` int(11) NOT NULL DEFAULT '0',
  `channel17` int(11) NOT NULL DEFAULT '0',
  `channel18` int(11) NOT NULL DEFAULT '0',
  `channel19` int(11) NOT NULL DEFAULT '0',
  `channel20` int(11) NOT NULL DEFAULT '0',
  `channel21` int(11) NOT NULL DEFAULT '0',
  `channel22` int(11) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `message_32bit_created_date` (`created_date`),
  KEY `message_32bit_terminal_code` (`terminal_code`),
  KEY `message_32bit_terminal_set_id` (`terminal_set_id`),
  KEY `message_32bit_version` (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `message_formatter`
-- ----------------------------
DROP TABLE IF EXISTS `message_formatter`;
CREATE TABLE `message_formatter` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `formatter_id` int(11) NOT NULL DEFAULT '0',
  `type_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `tag_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `length` int(11) NOT NULL DEFAULT '0',
  `sequence_number` int(11) NOT NULL DEFAULT '0',
  `start_point` int(11) NOT NULL DEFAULT '0',
  `default` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `message_formatter_formatter_id` (`formatter_id`),
  KEY `message_formatter_created_date` (`created_date`),
  KEY `message_formatter_is_deleted` (`is_deleted`),
  KEY `message_formatter_name` (`name`),
  KEY `message_formatter_is_demo` (`is_demo`),
  KEY `message_formatter_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `message_log`
-- ----------------------------
DROP TABLE IF EXISTS `message_log`;
CREATE TABLE `message_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `body` longtext COLLATE utf8_unicode_ci NOT NULL,
  `error` longtext COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  KEY `message_log_name` (`name`),
  KEY `message_log_created_date` (`created_date`),
  KEY `message_log_updated_date` (`updated_date`),
  KEY `message_log_is_demo` (`is_demo`),
  KEY `message_log_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `message_tag`
-- ----------------------------
DROP TABLE IF EXISTS `message_tag`;
CREATE TABLE `message_tag` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `tag_id` int(11) NOT NULL DEFAULT '0',
  `data_type` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `column` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `content` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `length` int(11) NOT NULL DEFAULT '0',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `message_tag_created_date` (`created_date`),
  KEY `message_tag_is_deleted` (`is_deleted`),
  KEY `message_tag_name` (`name`),
  KEY `message_tag_is_demo` (`is_demo`),
  KEY `message_tag_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `message_type`
-- ----------------------------
DROP TABLE IF EXISTS `message_type`;
CREATE TABLE `message_type` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `type_id` int(11) NOT NULL DEFAULT '0',
  `from` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `message_type_type_id` (`type_id`),
  KEY `message_type_created_date` (`created_date`),
  KEY `message_type_is_deleted` (`is_deleted`),
  KEY `message_type_name` (`name`),
  KEY `message_type_is_demo` (`is_demo`),
  KEY `message_type_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `organization`
-- ----------------------------
DROP TABLE IF EXISTS `organization`;
CREATE TABLE `organization` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `type_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `address_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `contact_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `brand_name` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `brand_image_url` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `show_brand` tinyint(1) NOT NULL DEFAULT '0',
  `is_supervisor` tinyint(1) NOT NULL DEFAULT '0',
  `super_organization_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`uid`),
  KEY `organization_created_date` (`created_date`),
  KEY `organization_is_deleted` (`is_deleted`),
  KEY `organization_name` (`name`),
  KEY `organization_is_demo` (`is_demo`),
  KEY `organization_type_id` (`type_id`),
  KEY `organization_updated_date` (`updated_date`),
  KEY `organization_super_organization_id` (`super_organization_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `organization_type`
-- ----------------------------
DROP TABLE IF EXISTS `organization_type`;
CREATE TABLE `organization_type` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `type_id` int(11) NOT NULL DEFAULT '0',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `organization_type_type_id` (`type_id`),
  KEY `organization_type_created_date` (`created_date`),
  KEY `organization_type_is_deleted` (`is_deleted`),
  KEY `organization_type_name` (`name`),
  KEY `organization_type_is_demo` (`is_demo`),
  KEY `organization_type_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `pf_atomic_data`
-- ----------------------------
DROP TABLE IF EXISTS `pf_atomic_data`;
CREATE TABLE `pf_atomic_data` (
  `id` int(11) NOT NULL,
  `name` char(30) NOT NULL,
  `data_desc` char(60) NOT NULL,
  `atom_len` int(11) NOT NULL,
  `atom_type` char(1) NOT NULL,
  `dis_len` int(11) NOT NULL,
  `dis_type` char(1) NOT NULL,
  `dec_len` int(11) NOT NULL DEFAULT '0',
  `mark` char(120) NOT NULL,
  `lst_upd_dt` date NOT NULL,
  `ts` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`,`name`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `pf_message_data`
-- ----------------------------
DROP TABLE IF EXISTS `pf_message_data`;
CREATE TABLE `pf_message_data` (
  `id` int(11) NOT NULL,
  `name` char(30) NOT NULL,
  `elem` char(30) NOT NULL,
  `elem_sseq` int(11) NOT NULL,
  `attr` char(30) NOT NULL DEFAULT '',
  `def_data` varchar(20) NOT NULL,
  `data_desc` char(60) NOT NULL,
  `mark` char(120) NOT NULL,
  `lst_upd_dt` date NOT NULL,
  `ts` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`,`name`,`elem`,`elem_sseq`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `pf_message_list`
-- ----------------------------
DROP TABLE IF EXISTS `pf_message_list`;
CREATE TABLE `pf_message_list` (
  `id` int(11) NOT NULL,
  `name` char(30) NOT NULL,
  `data_desc` char(60) NOT NULL,
  `data_count` int(11) NOT NULL,
  `data_fmt` char(200) NOT NULL,
  `func_code` char(4) NOT NULL,
  `mark` char(120) NOT NULL,
  `lst_upd_dt` date NOT NULL,
  `ts` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for `runtime_alarm_rule`
-- ----------------------------
DROP TABLE IF EXISTS `runtime_alarm_rule`;
CREATE TABLE `runtime_alarm_rule` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `parameter_id` bigint(20) NOT NULL,
  `boiler_form_id` bigint(20) DEFAULT NULL,
  `boiler_medium_id` bigint(20) DEFAULT NULL,
  `boiler_fuel_type_id` bigint(20) DEFAULT NULL,
  `boiler_capacity_min` int(11) NOT NULL DEFAULT '0',
  `boiler_capacity_max` int(11) NOT NULL DEFAULT '0',
  `normal` double NOT NULL DEFAULT '0',
  `warning` double NOT NULL DEFAULT '0',
  `danger` double NOT NULL DEFAULT '0',
  `priority` int(11) NOT NULL DEFAULT '0',
  `scope` int(11) NOT NULL DEFAULT '0',
  `need_send` tinyint(1) NOT NULL DEFAULT '0',
  `delay` bigint(20) NOT NULL DEFAULT '0',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`uid`),
  KEY `runtime_alarm_rule_name` (`name`),
  KEY `runtime_alarm_rule_created_date` (`created_date`),
  KEY `runtime_alarm_rule_updated_date` (`updated_date`),
  KEY `runtime_alarm_rule_is_demo` (`is_demo`),
  KEY `runtime_alarm_rule_is_deleted` (`is_deleted`),
  KEY `runtime_alarm_rule_parameter_id` (`parameter_id`),
  KEY `runtime_alarm_rule_boiler_form_id` (`boiler_form_id`),
  KEY `runtime_alarm_rule_boiler_medium_id` (`boiler_medium_id`),
  KEY `runtime_alarm_rule_boiler_fuel_type_id` (`boiler_fuel_type_id`),
  KEY `runtime_alarm_rule_boiler_capacity_min` (`boiler_capacity_min`),
  KEY `runtime_alarm_rule_boiler_capacity_max` (`boiler_capacity_max`),
  KEY `runtime_alarm_rule_priority` (`priority`),
  KEY `runtime_alarm_rule_scope` (`scope`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `runtime_parameter`
-- ----------------------------
DROP TABLE IF EXISTS `runtime_parameter`;
CREATE TABLE `runtime_parameter` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `param_id` int(11) NOT NULL DEFAULT '0',
  `category_id` bigint(20) NOT NULL,
  `medium_id` bigint(20) NOT NULL,
  `length` int(11) NOT NULL DEFAULT '2',
  `scale` double NOT NULL DEFAULT '1',
  `fix` int(11) NOT NULL DEFAULT '2',
  `unit` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_short` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `runtime_parameter_id` (`param_id`,`category_id`) USING BTREE,
  KEY `runtime_parameter_name` (`name`),
  KEY `runtime_parameter_created_date` (`created_date`),
  KEY `runtime_parameter_updated_date` (`updated_date`),
  KEY `runtime_parameter_is_demo` (`is_demo`),
  KEY `runtime_parameter_is_deleted` (`is_deleted`),
  KEY `runtime_parameter_param_id` (`param_id`),
  KEY `runtime_parameter_category_id` (`category_id`),
  KEY `runtime_parameter_medium_id` (`medium_id`)
) ENGINE=InnoDB AUTO_INCREMENT=13103 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `runtime_parameter_boiler_mediums`
-- ----------------------------
DROP TABLE IF EXISTS `runtime_parameter_boiler_mediums`;
CREATE TABLE `runtime_parameter_boiler_mediums` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `runtime_parameter_id` bigint(20) NOT NULL,
  `boiler_medium_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=599 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `runtime_parameter_category`
-- ----------------------------
DROP TABLE IF EXISTS `runtime_parameter_category`;
CREATE TABLE `runtime_parameter_category` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `runtime_parameter_category_name` (`name`),
  KEY `runtime_parameter_category_created_date` (`created_date`),
  KEY `runtime_parameter_category_updated_date` (`updated_date`),
  KEY `runtime_parameter_category_is_demo` (`is_demo`),
  KEY `runtime_parameter_category_is_deleted` (`is_deleted`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `runtime_parameter_channel_config`
-- ----------------------------
DROP TABLE IF EXISTS `runtime_parameter_channel_config`;
CREATE TABLE `runtime_parameter_channel_config` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `parameter_id` bigint(20) NOT NULL,
  `boiler_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `terminal_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `channel_type` int(11) NOT NULL DEFAULT '0',
  `channel_number` int(11) NOT NULL DEFAULT '0',
  `signed` tinyint(1) NOT NULL DEFAULT '1',
  `negative_threshold` int(11) NOT NULL DEFAULT '65000',
  `is_default` tinyint(1) NOT NULL DEFAULT '0',
  `scale` double NOT NULL DEFAULT '0',
  `length` int(11) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  `sequence_number` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `runtime_parameter_channel_config_name` (`name`),
  KEY `runtime_parameter_channel_config_created_date` (`created_date`),
  KEY `runtime_parameter_channel_config_updated_date` (`updated_date`),
  KEY `runtime_parameter_channel_config_is_demo` (`is_demo`),
  KEY `runtime_parameter_channel_config_is_deleted` (`is_deleted`),
  KEY `runtime_parameter_channel_config_parameter_id` (`parameter_id`),
  KEY `runtime_parameter_channel_config_boiler_id` (`boiler_id`),
  KEY `runtime_parameter_channel_config_terminal_id` (`terminal_id`),
  KEY `runtime_parameter_channel_config_channel_type` (`channel_type`),
  KEY `runtime_parameter_channel_config_channel_number` (`channel_number`),
  KEY `runtime_parameter_channel_config_is_default` (`is_default`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `runtime_parameter_channel_config_range`
-- ----------------------------
DROP TABLE IF EXISTS `runtime_parameter_channel_config_range`;
CREATE TABLE `runtime_parameter_channel_config_range` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `channel_config_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `min` bigint(20) NOT NULL DEFAULT '0',
  `max` bigint(20) NOT NULL DEFAULT '0',
  `value` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `channel_config_id` (`channel_config_id`,`min`,`max`),
  KEY `runtime_parameter_channel_config_range_name` (`name`),
  KEY `runtime_parameter_channel_config_range_created_date` (`created_date`),
  KEY `runtime_parameter_channel_config_range_updated_date` (`updated_date`),
  KEY `runtime_parameter_channel_config_range_is_demo` (`is_demo`),
  KEY `runtime_parameter_channel_config_range_is_deleted` (`is_deleted`),
  KEY `runtime_parameter_channel_config_range_channel_config_id` (`channel_config_id`),
  KEY `runtime_parameter_channel_config_range_min` (`min`),
  KEY `runtime_parameter_channel_config_range_max` (`max`),
  KEY `runtime_parameter_channel_config_range_value` (`value`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `runtime_parameter_medium`
-- ----------------------------
DROP TABLE IF EXISTS `runtime_parameter_medium`;
CREATE TABLE `runtime_parameter_medium` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `runtime_parameter_medium_name` (`name`),
  KEY `runtime_parameter_medium_created_date` (`created_date`),
  KEY `runtime_parameter_medium_updated_date` (`updated_date`),
  KEY `runtime_parameter_medium_is_demo` (`is_demo`),
  KEY `runtime_parameter_medium_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `terminal`
-- ----------------------------
DROP TABLE IF EXISTS `terminal`;
CREATE TABLE `terminal` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `terminal_code` bigint(20) NOT NULL,
  `ip_address` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `upload_flag` tinyint(1) NOT NULL DEFAULT '0',
  `upload_period` int(11) NOT NULL DEFAULT '0',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `organization_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `local_ip` varchar(60) COLLATE utf8_unicode_ci DEFAULT NULL,
  `remote_ip` varchar(60) COLLATE utf8_unicode_ci DEFAULT NULL,
  `remote_port` int(11) NOT NULL DEFAULT '0',
  `sim_number` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `installed_by` varchar(60) COLLATE utf8_unicode_ci DEFAULT NULL,
  `installed_date` datetime DEFAULT NULL,
  `is_online` tinyint(1) NOT NULL DEFAULT '0',
  `terminal_id` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `terminal_terminal_id` (`terminal_code`) USING BTREE,
  KEY `terminal_created_date` (`created_date`),
  KEY `terminal_is_deleted` (`is_deleted`),
  KEY `terminal_name` (`name`),
  KEY `terminal_is_demo` (`is_demo`),
  KEY `terminal_organization_id` (`organization_id`),
  KEY `terminal_updated_date` (`updated_date`),
  KEY `terminal_terminal_code` (`terminal_code`),
  KEY `terminal_is_online` (`is_online`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `user`
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `supervisor_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `role_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `organization_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `username` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `password` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `picture` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'avatar0.png',
  `status` int(11) NOT NULL DEFAULT '0',
  `address_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `phone_number` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `mobile_number` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `email` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `register_ip` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `gender` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `username_unique` (`username`) USING HASH,
  KEY `user_created_date` (`created_date`),
  KEY `user_is_deleted` (`is_deleted`),
  KEY `user_name` (`name`),
  KEY `user_is_demo` (`is_demo`),
  KEY `user_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `user_login`
-- ----------------------------
DROP TABLE IF EXISTS `user_login`;
CREATE TABLE `user_login` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `user_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `login_password` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_success` tinyint(1) NOT NULL DEFAULT '0',
  `login_ip` varchar(20) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `login_date` datetime NOT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `login_method` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_login` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `user_login_created_date` (`created_date`),
  KEY `user_login_is_deleted` (`is_deleted`),
  KEY `user_login_name` (`name`),
  KEY `user_login_is_demo` (`is_demo`),
  KEY `user_login_updated_date` (`updated_date`),
  KEY `user_login_is_login` (`is_login`),
  KEY `user_login_is_success` (`is_success`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `user_role`
-- ----------------------------
DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `role_id` int(11) NOT NULL DEFAULT '0',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `user_role_role_id` (`role_id`),
  KEY `user_role_created_date` (`created_date`),
  KEY `user_role_is_deleted` (`is_deleted`),
  KEY `user_role_name` (`name`),
  KEY `user_role_is_demo` (`is_demo`),
  KEY `user_role_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `user_session`
-- ----------------------------
DROP TABLE IF EXISTS `user_session`;
CREATE TABLE `user_session` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `user_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `login_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `application_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_actived` tinyint(1) NOT NULL DEFAULT '0',
  `token` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`uid`),
  KEY `user_session_name` (`name`),
  KEY `user_session_created_date` (`created_date`),
  KEY `user_session_is_demo` (`is_demo`),
  KEY `user_session_is_deleted` (`is_deleted`),
  KEY `user_session_user_id` (`user_id`),
  KEY `user_session_login_id` (`login_id`),
  KEY `user_session_application_id` (`application_id`),
  KEY `user_session_is_actived` (`is_actived`),
  KEY `user_session_token` (`token`),
  KEY `user_session_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  Table structure for `user_third`
-- ----------------------------
DROP TABLE IF EXISTS `user_third`;
CREATE TABLE `user_third` (
  `uid` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `name_en` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL,
  `created_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_date` datetime NOT NULL,
  `updated_by_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `user_id` varchar(36) COLLATE utf8_unicode_ci DEFAULT NULL,
  `application_id` varchar(36) COLLATE utf8_unicode_ci NOT NULL,
  `platform` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `app` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `identity` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `language` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `sex` int(11) NOT NULL DEFAULT '0',
  `province` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `city` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `country` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `head_image_url` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `open_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `union_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `access_token` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `refresh_token` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `session_key` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `expires_in` datetime DEFAULT NULL,
  PRIMARY KEY (`uid`),
  KEY `user_third_name` (`name`),
  KEY `user_third_created_date` (`created_date`),
  KEY `user_third_is_demo` (`is_demo`),
  KEY `user_third_is_deleted` (`is_deleted`),
  KEY `user_third_user_id` (`user_id`),
  KEY `user_third_application_id` (`application_id`),
  KEY `user_third_platform` (`platform`),
  KEY `user_third_app` (`app`),
  KEY `user_third_identity` (`identity`),
  KEY `user_third_open_id` (`open_id`),
  KEY `user_third_union_id` (`union_id`),
  KEY `user_third_access_token` (`access_token`),
  KEY `user_third_refresh_token` (`refresh_token`),
  KEY `user_third_expires_in` (`expires_in`),
  KEY `user_third_updated_date` (`updated_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
--  View structure for `boiler_runtime_flow_daily`
-- ----------------------------
DROP VIEW IF EXISTS `boiler_runtime_flow_daily`;
CREATE ALGORITHM=UNDEFINED DEFINER=`azureadmin`@`%` SQL SECURITY DEFINER VIEW `boiler_runtime_flow_daily` AS select `boiler_runtime`.`boiler_id` AS `boiler_id`,(avg(`boiler_runtime`.`value`) * hour(max(`boiler_runtime`.`created_date`))) AS `flow`,cast(`boiler_runtime`.`created_date` as date) AS `date` from `boiler_runtime` where (`boiler_runtime`.`parameter_id` = 1003) group by `boiler_runtime`.`boiler_id`,cast(`boiler_runtime`.`created_date` as date);

-- ----------------------------
--  View structure for `boiler_runtime_heat_avg_yearweek`
-- ----------------------------
DROP VIEW IF EXISTS `boiler_runtime_heat_avg_yearweek`;
CREATE ALGORITHM=UNDEFINED DEFINER=`azureadmin`@`%` SQL SECURITY DEFINER VIEW `boiler_runtime_heat_avg_yearweek` AS select `boiler`.`uid` AS `uid`,`boiler`.`name` AS `name`,`boiler`.`evaporating_capacity` AS `evaporte`,`fuel`.`type_id` AS `fuel_type_id`,avg(`heat`.`value`) AS `heats`,yearweek(`heat`.`created_date`,0) AS `week` from ((`boiler` join `fuel`) join `boiler_runtime_cache_heat` `heat`) where ((`boiler`.`fuel_id` = `fuel`.`uid`) and (`boiler`.`uid` = `heat`.`boiler_id`)) group by `heat`.`boiler_id`,yearweek(`heat`.`created_date`,0);

-- ----------------------------
--  View structure for `boiler_runtime_heat_month`
-- ----------------------------
DROP VIEW IF EXISTS `boiler_runtime_heat_month`;
CREATE ALGORITHM=UNDEFINED DEFINER=`azureadmin`@`%` SQL SECURITY DEFINER VIEW `boiler_runtime_heat_month` AS select `boiler_runtime`.`boiler_id` AS `boiler_id`,cast(`boiler_runtime`.`created_date` as date) AS `date`,avg(`boiler_runtime`.`value`) AS `heat` from `boiler_runtime` where (`boiler_runtime`.`parameter_id` = 1201) group by `boiler_runtime`.`boiler_id`,cast(`boiler_runtime`.`created_date` as date);

-- ----------------------------
--  View structure for `boiler_runtime_view`
-- ----------------------------
DROP VIEW IF EXISTS `boiler_runtime_view`;
CREATE ALGORITHM=UNDEFINED DEFINER=`azureadmin`@`%` SQL SECURITY DEFINER VIEW `boiler_runtime_view` AS select `runtime`.`id` AS `Id`,`runtime`.`boiler_id` AS `Boiler`,`runtime`.`parameter_id` AS `Parameter`,`runtime`.`alarm_id` AS `Alarm`,round((`param`.`scale` * `runtime`.`value`),`param`.`fix`) AS `Value`,`param`.`name` AS `ParameterName`,`param`.`unit` AS `Unit`,`alarm`.`alarm_level` AS `AlarmLevel`,`alarm`.`description` AS `AlarmDescription`,`runtime`.`created_date` AS `CreatedDate`,`runtime`.`created_by_id` AS `CreatedBy`,`runtime`.`updated_date` AS `UpdatedDate`,`runtime`.`updated_by_id` AS `UpdatedBy`,`runtime`.`is_deleted` AS `IsDeleted`,`runtime`.`is_demo` AS `IsDemo` from ((`boiler_runtime` `runtime` left join `runtime_parameter` `param` on((`runtime`.`parameter_id` = `param`.`id`))) left join `boiler_alarm` `alarm` on((`runtime`.`alarm_id` = `alarm`.`uid`))) where (`runtime`.`created_date` >= cast(now() as date)) order by `runtime`.`created_date`;

-- ----------------------------
--  Function structure for `check_boilerid_null`
-- ----------------------------
DROP FUNCTION IF EXISTS `check_boilerid_null`;
delimiter ;;
CREATE DEFINER=`azureadmin`@`%` FUNCTION `check_boilerid_null`(boilerId varchar(36)) RETURNS int(10)
BEGIN
    DECLARE ts INT(10);

    IF !ISNULL(boilerId) THEN
	SET ts = 1;
    ELSE
	SET ts = 10;
    END IF;

    RETURN ts;
END
 ;;
delimiter ;

-- ----------------------------
--  Event structure for `delete_cache_remain_today`
-- ----------------------------
DROP EVENT IF EXISTS `delete_cache_remain_today`;
delimiter ;;
CREATE DEFINER=`azureadmin`@`%` EVENT `delete_cache_remain_today` ON SCHEDULE EVERY 1 DAY STARTS '2017-05-01 04:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN
	DELETE 
	FROM	`boiler_runtime_cache_environment_temperature`
	WHERE 	`created_date` < DATE(NOW());

	DELETE 
	FROM	`boiler_runtime_cache_excess_air`
	WHERE 	`created_date` < DATE(NOW());

	DELETE 
	FROM	`boiler_runtime_cache_smoke_component`
	WHERE 	`created_date` < DATE(NOW());

	DELETE 
	FROM	`boiler_runtime_cache_smoke_temperature`
	WHERE 	`created_date` < DATE(NOW());

	DELETE 
	FROM	`boiler_runtime_cache_steam_temperature`
	WHERE 	`created_date` < DATE(NOW());

	DELETE 
	FROM	`boiler_runtime_cache_steam_pressure`
	WHERE 	`created_date` < DATE(NOW());

	DELETE 
	FROM	`boiler_runtime_cache_water_temperature`
	WHERE 	`created_date` < DATE(NOW());
END
 ;;
delimiter ;

-- ----------------------------
--  Event structure for `delete_demo_runtime_physical`
-- ----------------------------
DROP EVENT IF EXISTS `delete_demo_runtime_physical`;
delimiter ;;
CREATE DEFINER=`azureadmin`@`%` EVENT `delete_demo_runtime_physical` ON SCHEDULE EVERY 1 DAY STARTS '2017-05-01 03:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN
	DELETE 
	FROM	`boiler_runtime`
	WHERE 	`is_demo` = 1 AND `is_deleted` = 1;

	DELETE 
	FROM	`boiler_alarm`
	WHERE 	`is_demo` = 1 AND `is_deleted` = 1;

	DELETE 
	FROM	`boiler_alarm_history`
	WHERE 	`is_demo` = 1 AND `is_deleted` = 1;
END
 ;;
delimiter ;

-- ----------------------------
--  Event structure for `update_daily_flow_and_heat`
-- ----------------------------
DROP EVENT IF EXISTS `update_daily_flow_and_heat`;
delimiter ;;
CREATE DEFINER=`azureadmin`@`%` EVENT `update_daily_flow_and_heat` ON SCHEDULE EVERY 4 HOUR STARTS '2017-05-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN
	DELETE
	FROM	`boiler_runtime_cache_heat_daily`;

	INSERT IGNORE `boiler_runtime_cache_heat_daily`
		( `name`, `remark`, `is_demo`, `boiler_id`, 
		`parameter_id`, `parameter_name`, `unit`, 
		`value`, `date`) 
	SELECT 	`name`, `remark`, `is_demo`, `boiler_id`, 
		`parameter_id`, '每日平均热效率', `unit`, 
		AVG(`value`) AS `value`, DATE(`created_date`) AS `day`
	FROM  	`boiler_runtime_cache_heat`
	WHERE 	`value` > 0
	GROUP BY `day`, `boiler_id`;

	DELETE
	FROM `boiler_runtime_cache_flow_daily`;

	INSERT IGNORE `boiler_runtime_cache_flow_daily`
		( `name`, `remark`, `is_demo`, `boiler_id`, 
		`parameter_id`, `parameter_name`, `unit`, 
		`value`, `hours`, `date`) 
	SELECT 	`flow`.`name`, `flow`.`remark`, `flow`.`is_demo`, `flow`.`boiler_id`, 
		1003, '每日累计流量', 't', 
		ROUND(SUM(`flow`.`value`), 3) AS `total`, COUNT(DISTINCT `flow`.`hour`) AS `hours`, `flow`.`day` AS `date`
	FROM 	(SELECT DATE(`created_date`) AS `day`, HOUR(`created_date`) AS `hour`, AVG(`value`) AS `value`, `boiler_id`, `name`, `is_demo`, `remark`
		FROM `boiler_runtime_cache_flow`
		GROUP BY `day`, `hour`, `boiler_id`, `is_demo`) AS `flow`
	GROUP BY `flow`.`day`, `flow`.`boiler_id`, `flow`.`is_demo`;
END
 ;;
delimiter ;

-- ----------------------------
--  Event structure for `move_old_alarm_to_history`
-- ----------------------------
DROP EVENT IF EXISTS `move_old_alarm_to_history`;
delimiter ;;
CREATE DEFINER=`azureadmin`@`%` EVENT `move_old_alarm_to_history` ON SCHEDULE EVERY 1 HOUR STARTS '2017-05-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN
	INSERT IGNORE `boiler_alarm_history`
	(`uid`, `name`, `name_en`, `remark`,
	 `is_demo`, `is_deleted`,
	 `boiler_id`, `parameter_id`, `trigger_rule_id`,
	 `start_date`, `end_date`,
	 `confirmed_date`, `confirmed_by_id`, `verified_date`, `verified_by_id`,
	 `alarm_level`, `priority`, `description`)
	SELECT 	
	 `uid`, `name`, `name_en`, `remark`,
	 `is_demo`, `is_deleted`,
	 `boiler_id`, `parameter_id`, `trigger_rule_id`,
	 `start_date`, `end_date`,
	 `confirmed_date`, `confirmed_by_id`, `verified_date`, `verified_by_id`,
	 `alarm_level`, `priority`, `description`
	FROM	`boiler_alarm`
	WHERE	`end_date` < (NOW() - INTERVAL 4 HOUR);

	DELETE  `alarm`.*
	FROM	`boiler_alarm` AS `alarm`, `boiler_alarm_history` AS `history`
	WHERE	`alarm`.`uid` = `history`.`uid`;
END
 ;;
delimiter ;

-- ----------------------------
--  Event structure for `set_delete_sign_of_alarm_orphans`
-- ----------------------------
DROP EVENT IF EXISTS `set_delete_sign_of_alarm_orphans`;
delimiter ;;
CREATE DEFINER=`azureadmin`@`%` EVENT `set_delete_sign_of_alarm_orphans` ON SCHEDULE EVERY 6 HOUR STARTS '2017-06-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN
	UPDATE	`boiler_alarm` AS `alarm`
	SET	`is_deleted` = true
	WHERE	NOT EXISTS 
		(SELECT DISTINCT(`alarm`.`uid`) 
		FROM	`boiler_runtime` AS `runtime`
		WHERE	`runtime`.`alarm_id` = `alarm`.`uid`);

	UPDATE	`boiler_alarm_history` AS `alarm`
	SET	`is_deleted` = true
	WHERE	NOT EXISTS 
		(SELECT DISTINCT(`alarm`.`uid`) 
		FROM	`boiler_runtime` AS `runtime`
		WHERE	`runtime`.`alarm_id` = `alarm`.`uid`);
END
 ;;
delimiter ;

-- ----------------------------
--  Event structure for `delete_runtime_is_demo_over_month`
-- ----------------------------
DROP EVENT IF EXISTS `delete_runtime_is_demo_over_month`;
delimiter ;;
CREATE DEFINER=`azureadmin`@`%` EVENT `delete_runtime_is_demo_over_month` ON SCHEDULE EVERY 1 DAY STARTS '2017-11-02 00:00:00' ON COMPLETION PRESERVE ENABLE DO BEGIN
	DELETE FROM `boiler_runtime` WHERE `is_demo` = 1 AND `created_date` < CURDATE() - INTERVAL 1 MONTH LIMIT 1000000;
	DELETE FROM `boiler_runtime_cache_history_backup` WHERE `is_demo` = 1 AND `created_date` < CURDATE() - INTERVAL 1 MONTH LIMIT 500000;
	DELETE FROM `boiler_runtime_cache_history` WHERE `is_demo` = 1 AND `created_date` < CURDATE() - INTERVAL 1 MONTH LIMIT 500000;
END
 ;;
delimiter ;

-- ----------------------------
--  Event structure for `insert_m163_range_channel_sample`
-- ----------------------------
DROP EVENT IF EXISTS `insert_m163_range_channel_sample`;
delimiter ;;
CREATE DEFINER=`azureadmin`@`%` EVENT `insert_m163_range_channel_sample` ON SCHEDULE EVERY 15 SECOND STARTS '2017-12-29 23:29:00' ON COMPLETION PRESERVE ENABLE DO BEGIN
	INSERT INTO `boiler_main`.`boiler_m163` 
		(`Boiler_term_id`, `Boiler_boiler_id`, `Term_sys_time`, `Boiler_sn`, `Boiler_data_fmt_ver`, `Boiler_status_code`, `Boiler_term_err_code`, 
		`Temper1_channel`, `Temper2_channel`, `Temper3_channel`, `Temper4_channel`, `Temper5_channel`, `Temper6_channel`, `Temper7_channel`, `Temper8_channel`, `Temper9_channel`, `Temper10_channel`, `Temper11_channel`, `Temper12_channel`, 
		`Analog1_channel`, `Analog2_channel`, `Analog3_channel`, `Analog4_channel`, `Analog5_channel`, `Analog6_channel`, `Analog7_channel`, `Analog8_channel`, `Analog9_channel`, `Analog10_channel`, `Analog11_channel`, `Analog12_channel`, 
		`C1_calculate_parm`, `C2_calculate_parm`, `C3_calculate_parm`, `C4_calculate_parm`, `C5_calculate_parm`, `C6_calculate_parm`, `C7_calculate_parm`, `C8_calculate_parm`, `C9_calculate_parm`, `C10_calculate_parm`, `C11_calculate_parm`, `C12_calculate_parm`, 
		`Switch_in_1_16_channel`, `Switch_in_17_32_channel`, `Switch_out_1_16_channel`, 
		`Reserved1_filler`, `Reserved2_filler`, `Reserved3_filler`, `Reserved4_filler`, 
		`need_reload`, `TS`)
	VALUES	('700010', '01', NOW(), 668, 80, 'a3', 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		444, 333, 444, 555, 666, 777, 888, 999, 10101, 20202, 30303, 40404,
		711, 554, 29837,
		0, 0, 0, 0,
		TRUE, NOW());
END
 ;;
delimiter ;

-- ----------------------------
--  Triggers structure for table boiler
-- ----------------------------
DROP TRIGGER IF EXISTS `boiler_after_insert`;
delimiter ;;
CREATE TRIGGER `boiler_after_insert` AFTER INSERT ON `boiler` FOR EACH ROW BEGIN
	INSERT INTO `boiler_config`
	SET	`uid` = UUID(),
		`boiler_id` = NEW.`uid`,
		`name` = NEW.`name`,
		`name_en` = NEW.`name_en`,
		`created_by_id` = NEW.`created_by_id`,
		`updated_by_id` = NEW.`updated_by_id`,
		`is_demo` = `is_demo`
	ON DUPLICATE KEY 
	UPDATE	`name` = NEW.`name`,
		`name_en` = NEW.`name_en`,
		`updated_by_id` = NEW.`updated_by_id`,
		`is_demo` = `is_demo`;

	INSERT INTO `boiler_calculate_parameter`
	SET	`uid` = UUID(),
		`boiler_id` = NEW.`uid`,
		`name` = NEW.`name`,
		`name_en` = NEW.`name_en`,
		`created_by_id` = NEW.`created_by_id`,
		`updated_by_id` = NEW.`updated_by_id`,
		`is_demo` = `is_demo`
	ON DUPLICATE KEY 
	UPDATE	`name` = NEW.`name`,
		`name_en` = NEW.`name_en`,
		`updated_by_id` = NEW.`updated_by_id`,
		`is_demo` = `is_demo`;
END
 ;;
delimiter ;
DROP TRIGGER IF EXISTS `boiler_after_update`;
delimiter ;;
CREATE TRIGGER `boiler_after_update` AFTER UPDATE ON `boiler` FOR EACH ROW BEGIN
	IF NEW.`is_deleted` = 1 THEN
		UPDATE 	`boiler_runtime`
		SET	`is_deleted` = 1
		WHERE	`boiler_id` = NEW.`uid`;

		UPDATE 	`boiler_alarm`
		SET	`is_deleted` = 1
		WHERE	`boiler_id` = NEW.`uid`;

		UPDATE 	`boiler_alarm_history`
		SET	`is_deleted` = 1
		WHERE	`boiler_id` = NEW.`uid`;

		UPDATE 	`boiler_runtime_cache_instant`
		SET	`is_deleted` = 1
		WHERE	`boiler_id` = NEW.`uid`;

		UPDATE 	`boiler_runtime_cache_history`
		SET	`is_deleted` = 1
		WHERE	`boiler_id` = NEW.`uid`;

		UPDATE 	`boiler_maintenance`
		SET	`is_deleted` = 1
		WHERE	`boiler_id` = NEW.`uid`;

		UPDATE 	`boiler_calculate_parameter`
		SET	`is_deleted` = 1
		WHERE	`boiler_id` = NEW.`uid`;

		UPDATE 	`boiler_config`
		SET	`is_deleted` = 1
		WHERE	`boiler_id` = NEW.`uid`;
	END IF;
END
 ;;
delimiter ;
DROP TRIGGER IF EXISTS `boiler_after_delete`;
delimiter ;;
CREATE TRIGGER `boiler_after_delete` AFTER DELETE ON `boiler` FOR EACH ROW BEGIN
	DELETE FROM `boiler_runtime`
	WHERE	`boiler_id` = OLD.`uid`;

	DELETE FROM `boiler_alarm`
	WHERE	`boiler_id` = OLD.`uid`;

	DELETE FROM `boiler_alarm_history`
	WHERE	`boiler_id` = OLD.`uid`;

	DELETE FROM `boiler_runtime_cache_instant`
	WHERE	`boiler_id` = OLD.`uid`;

	DELETE FROM `boiler_runtime_cache_history`
	WHERE	`boiler_id` = OLD.`uid`;

	DELETE FROM `boiler_maintenance`
	WHERE	`boiler_id` = OLD.`uid`;

	DELETE FROM `boiler_calculate_parameter`
	WHERE	`boiler_id` = OLD.`uid`;

	DELETE FROM `boiler_config`
	WHERE	`boiler_id` = OLD.`uid`;
END
 ;;
delimiter ;

delimiter ;;
-- ----------------------------
--  Triggers structure for table boiler_alarm
-- ----------------------------
 ;;
delimiter ;
DROP TRIGGER IF EXISTS `alarm_after_updated`;
delimiter ;;
CREATE TRIGGER `alarm_after_updated` AFTER UPDATE ON `boiler_alarm` FOR EACH ROW BEGIN
	IF NEW.`is_deleted` = 1 THEN
		UPDATE	`boiler_alarm_feedback`
		SET	`is_deleted` = 1
		WHERE	`alarm_id` = NEW.`uid`;
	END IF;
END
 ;;
delimiter ;

delimiter ;;
-- ----------------------------
--  Triggers structure for table boiler_m163
-- ----------------------------
 ;;
delimiter ;
DROP TRIGGER IF EXISTS `uid`;
delimiter ;;
CREATE TRIGGER `uid` BEFORE INSERT ON `boiler_m163` FOR EACH ROW BEGIN
	SET NEW.`uid` = UUID();
END
 ;;
delimiter ;
DROP TRIGGER IF EXISTS `runtime`;
delimiter ;;
CREATE TRIGGER `runtime` AFTER INSERT ON `boiler_m163` FOR EACH ROW BEGIN
	SET 	@boilerId = NULL;
	SELECT 	`boiler`.`uid` INTO @boilerId
	FROM 	`boiler`, `terminal`
	WHERE 	`boiler`.`terminal_code` = CAST(NEW.`Boiler_term_id` AS SIGNED) AND
	    	`boiler`.`terminal_set_id` = CAST(NEW.`Boiler_boiler_id` AS SIGNED)
	LIMIT	1;

	IF !ISNULL(@boilerId) THEN
		IF NEW.`Boiler_term_id` = '030075' AND NEW.`Boiler_boiler_id` = '01' THEN
			INSERT IGNORE `boiler_runtime`
			SET `uid` = UUID(),
			`boiler_id` = @boilerId,
			`remark` = 'm163.Analog3_Sakura_030075_channel',
			`created_date` = NEW.`TS`,
			`parameter_id` = 1003,
			`value` = NEW.`Analog3_channel` * 0.6;
		END IF;

		SET	@fuelType = 0;
		SET	@q2 = 0;
		SET	@q3 = 0;
		SET	@q4 = 0;
		SET	@q5 = 0;
		SET	@q6 = 0;
		SET	@excAir = 0;
		SET	@heat = 0;
		SET	@remark = '';

		SET	@calcId = '';
		SET	@qnetvar = 0, @aar = 0, @mar = 0, @vdaf = 0, 
			@clz = 0, @clm = 0, @cfh = 0, @ded = 0, @dsc = 0, 
			@alz = 0, @alm = 0, @afh = 0, @tlz = 0, @ctlz = 0,
			@m = 0, @n = 0;

		SELECT 	`param`.`scale` INTO @scaleO2
		FROM	`runtime_parameter` AS `param`
		WHERE	`param`.`id` = 1016
		LIMIT	1;
		SELECT 	`param`.`scale` INTO @scaleSmoke
		FROM	`runtime_parameter` AS `param`
		WHERE	`param`.`id` = 1014
		LIMIT	1;
		SELECT 	`param`.`scale` INTO @scaleWind
		FROM	`runtime_parameter` AS `param`
		WHERE	`param`.`id` = 1021
		LIMIT	1;
		SELECT 	`param`.`scale` INTO @scaleExcAir
		FROM	`runtime_parameter` AS `param`
		WHERE	`param`.`id` = 1202
		LIMIT	1;
		SELECT 	`param`.`scale` INTO @scaleHeat
		FROM	`runtime_parameter` AS `param`
		WHERE	`param`.`id` = 1201
		LIMIT	1;
		

		SELECT  21 / (21 - NEW.`Analog2_channel` * @scaleO2) INTO @excAir;

		INSERT IGNORE `boiler_runtime`
		SET	`uid` = UUID(),
			`boiler_id` = @boilerId,
			`remark` = 'm163.calculate_excess_air',
			`created_date` = NEW.`TS`,
			`parameter_id` = 1202,
			`value` = @excAir / @scaleExcAir;
		
		SELECT	`fuel`.`type_id` INTO @fuelType
		FROM	`boiler`, `fuel`
		WHERE	`boiler`.`uid` = @boilerId AND `boiler`.`fuel_id` = `fuel`.`uid`
		LIMIT	1;

		IF @fuelType = 1 OR @fuelType = 4 THEN
			SELECT  `calc`.`uid`,
				`calc`.`coal_qnetvar`, `calc`.`coal_aar`, `calc`.`coal_mar`, `calc`.`coal_vdaf`,
				`calc`.`coal_clz`, `calc`.`coal_clm`, `calc`.`coal_cfh`, `calc`.`coal_ded`, `calc`.`coal_dsc`, 
				`calc`.`coal_alz`, `calc`.`coal_alm`, `calc`.`coal_afh`, `calc`.`coal_tlz`, `calc`.`coal_ct_lz`, 
				`calc`.`coal_m`, `calc`.`coal_n`, `calc`.`coal_q3`, `calc`.`conf_param1`
			INTO	@calcId,
				@qnetvar, @aar, @mar, @vdaf, 
				@clz, @clm, @cfh, @ded, @dsc, 
				@alz, @alm, @afh, @tlz, @ctlz,
				@m, @n, @q3, @q5
			FROM	`boiler_calculate_parameter` AS `calc`
			WHERE	`calc`.`boiler_id` = @boilerId
			LIMIT	1;
			
			IF @qnetvar > 0 THEN 
				SET	@c = (@alz * @clz) / (100.0 - @clz) + (@alm * @clm) / (100.0 - @clm) + (@afh * @cfh) / (100.0 - @cfh);
				SET	@q4 = (328.66 * @aar * @c) / @qnetvar;
				SET	@q2 = (@m + @n * @excAir) * ((NEW.`Temper5_channel` * @scaleSmoke - NEW.`Temper6_channel` * @scaleWind) / 100.0) * (1.0 - @q4 / 100.0);
				SET	@q6 = (100.0 - @afh * @aar * @ctlz) / (@qnetvar * (100.0 - @clz));
				SET	@heat = 100.0 - @q2 - @q3 - @q4 - @q5 - @q6; 
				SET	@remark = 'm163.calculate_coal_heat';

				INSERT IGNORE `boiler_calculate_result` ( 
					`name`, `name_en`, `remark`, 
					`boiler_id`, `fuel_id`, `based_parameter_id`, 
					`qnetvar`, `aar`, `mar`, `vdaf`, `clz`, `clm`, `cfh`, `ded`, `dsc`, `alz`, `alm`, `afh`, `tlz`, `ct_lz`,
					`m`, `n`, `q2`, `q3`, `q4`, `q5`, `q6`, 
					`excess_air`, `heat`) 
				SELECT 	`boiler`.`name`, `boiler`.`name_en`, 'auto',
					`boiler`.`uid`, `boiler`.`fuel_id`, @calcId,
					@qnetvar, @aar, @mar, @vdaf, @clz, @clm, @cfh, @ded, @dsc, @alz, @alm, @afh, @tlz, @ctlz,
					@m, @n, @q2, @q3, @q4, @q5, @q6, 
					@excAir, @heat
				FROM	`boiler`
				WHERE	`boiler`.`uid` = @boilerId;
			END IF;
		ELSE
			SELECT  `calc`.`uid`,
				`calc`.`coal_qnetvar`, `calc`.`gas_ded`,
				`calc`.`gas_m`, `calc`.`gas_n`, 
				`calc`.`gas_q3`, `calc`.`conf_param1`
			INTO	@calcId,
				@qnetvar, @ded, 
				@m, @n, @q3, @q5
			FROM	`boiler_calculate_parameter` AS `calc`
			WHERE	`calc`.`boiler_id` = @boilerId
			LIMIT	1;

			SET	@q2 = (@m + @n * @excAir) * ((NEW.`Temper5_channel` * @scaleSmoke - NEW.`Temper6_channel` * @scaleWind) / 100.0) * (1.0 - @q4 / 100.0);
			SET	@heat = 100.0 - @q2 - @q3 - @q4 - @q5 - @q6; 
			SET	@remark = 'm163.calculate_gas_heat';

			INSERT IGNORE `boiler_calculate_result` ( 
				`name`, `name_en`, `remark`, 
				`boiler_id`, `fuel_id`, `based_parameter_id`, 
				`qnetvar`, `ded`, 
				`m`, `n`, `q2`, `q3`, `q4`, `q5`, `q6`, 
				`excess_air`, `heat`) 
			SELECT 	`boiler`.`name`, `boiler`.`name_en`, 'auto',
				`boiler`.`uid`, `boiler`.`fuel_id`, @calcId,
				@qnetvar, @ded,
				@m, @n, @q2, @q3, @q4, @q5, @q6, 
				@excAir, @heat
			FROM	`boiler`
			WHERE	`boiler`.`uid` = @boilerId;
		END IF;

		IF @heat > 0.0 && @heat <= 100.0 THEN
			INSERT IGNORE `boiler_runtime`
			SET	`uid` = UUID(),
				`boiler_id` = @boilerId,
				`remark` = @remark,
				`created_date` = NEW.`TS`,
				`parameter_id` = 1201,
				`value` = @heat / @scaleHeat;
		END IF;
	END IF;
END
 ;;
delimiter ;

delimiter ;;
-- ----------------------------
--  Triggers structure for table boiler_runtime
-- ----------------------------
 ;;
delimiter ;
DROP TRIGGER IF EXISTS `alarm`;
delimiter ;;
CREATE TRIGGER `alarm` BEFORE INSERT ON `boiler_runtime` FOR EACH ROW BEGIN
	SELECT '' INTO @ruleId;
	SELECT '' INTO @existAlarmId;
	SELECT 0 INTO @priority;
	SELECT 0 INTO @alarmLevel;
	SELECT '' INTO @alarmDesc;
	SELECT 0 INTO @needSend;
	
	SELECT	`rule`.`uid`, `rule`.`priority`, `rule`.`need_send` INTO @ruleId, @priority, @needSend
	FROM	`boiler_runtime` AS `runtime`,
		(SELECT NEW.`value` * `param`.`scale` AS `value`
		 FROM 	`runtime_parameter` AS `param`
		 WHERE	`param`.`id` = NEW.`parameter_id`) 
		AS `rtm`,
		(SELECT	DISTINCT(`rule`.`uid`), `rule`.`parameter_id` AS `param_id`, 
			`rule`.`normal` AS `normal`, `rule`.`warning` AS `warning`, `rule`.`danger` AS `danger`,
			`rule`.`priority` AS `priority`, `rule`.`delay` AS `delay`, `rule`.`need_send` AS `need_send`
		 FROM	`runtime_alarm_rule` AS `rule`, 
			`runtime_parameter` AS `param`, 
			`boiler`, `boiler_type_form` AS `form`, `boiler_medium` AS `medium`,
			`fuel`, `fuel_type` AS `ftype`
		 WHERE   ((`rule`.`boiler_form_id` = `form`.`id` AND `form`.`id` = 0) 
			    OR `rule`.`boiler_form_id` = `boiler`.`form_id`)
			AND ((`rule`.`boiler_medium_id` = `medium`.`id` AND `medium`.`id` = 0) 
			    OR `rule`.`boiler_medium_id` = `boiler`.`medium_id`)
			AND ((`rule`.`boiler_fuel_type_id` = `ftype`.`id` AND `ftype`.`id` = 0) 
			    OR (`rule`.`boiler_fuel_type_id` = `fuel`.`type_id` AND `boiler`.`fuel_id` = `fuel`.`uid`))
			AND ((`rule`.`boiler_capacity_min` = 0 AND `rule`.`boiler_capacity_max` = 0) 
			    OR (`rule`.`boiler_capacity_min` <= `boiler`.`evaporating_capacity`)
			    AND (`rule`.`boiler_capacity_max` >= `boiler`.`evaporating_capacity` OR `rule`.`boiler_capacity_max` <= 0))
			AND `boiler`.`uid` = NEW.`boiler_id`
			AND `rule`.`parameter_id` = NEW.`parameter_id`
			AND `rule`.`is_deleted` = 0) 
		AS `rule`
	WHERE	((`rule`.`normal` < `rule`.`warning` AND `rtm`.`value` > `rule`.`warning`)
		OR (`rule`.`normal` > `rule`.`warning` AND `rtm`.`value` < `rule`.`warning`))
	LIMIT	1;

	SELECT	`boiler`.`name` INTO @boilerName
	FROM	`boiler`
	WHERE	`boiler`.`uid` = NEW.`boiler_id`
	LIMIT	1;
	
	IF @ruleId <> '' THEN
		SELECT	`alarm`.`uid` INTO @existAlarmId
		FROM	`boiler_alarm` AS `alarm`
		WHERE	`alarm`.`trigger_rule_id` = @ruleId
			AND `alarm`.`boiler_id` = NEW.`boiler_id`
			AND `alarm`.`end_date` > NOW() - INTERVAL 4 HOUR
		LIMIT	1;
		
		IF @existAlarmId <> '' THEN
			UPDATE	`boiler_alarm` AS `alarm`
			SET	`alarm`.`end_date` = NOW()
			WHERE	`alarm`.`uid` = @existAlarmId;

			SET	NEW.`alarm_id` = @existAlarmId;
		ELSE	
			SELECT UUID() INTO @alarmId;
			INSERT IGNORE `boiler_alarm`
			SET	`uid` = @alarmId,
				`name` = (SELECT `boiler`.`name` FROM `boiler` WHERE `boiler`.`uid` = NEW.`boiler_id`),
				`boiler_id` = NEW.`boiler_id`,
				`parameter_id` = NEW.`parameter_id`,
				`start_date` = NOW(),
				`end_date` = NOW(),
				`trigger_rule_id` = @ruleId,
				`alarm_level` = 1,
				`state` = 1,
				`priority` = @priority,
				`need_send` = @needSend,
				`description` = (SELECT CONCAT(`param`.`name`, IF(`rule`.`normal` < `rule`.`warning`, '过高', '过低'))
						 FROM	`runtime_parameter` AS `param`, `runtime_alarm_rule` AS `rule`
						 WHERE	`param`.`id` = `rule`.`parameter_id` AND `rule`.`uid` = @ruleId
						 LIMIT	1),
				`is_demo` = NEW.`is_demo`;

			SET	NEW.`alarm_id` = @alarmId;
		END IF;
	END IF;
END
 ;;
delimiter ;
DROP TRIGGER IF EXISTS `caches`;
delimiter ;;
CREATE TRIGGER `caches` AFTER INSERT ON `boiler_runtime` FOR EACH ROW BEGIN
	/* INSTANTS */
	INSERT IGNORE `boiler_runtime_cache_instant`(
		`runtime_id` ,
		`boiler_id` ,
		`parameter_id` ,
		`alarm_id` ,
		`created_date` ,
		`updated_date` ,
		`name` ,
		`value` ,
		`parameter_name` ,
		`unit` ,
		`alarm_level` ,
		`alarm_description`,
		`remark`
		)
	VALUES
		(
		NEW.`id` ,
		NEW.`boiler_id` ,
		NEW.`parameter_id` ,
		NEW.`alarm_id` ,
		NEW.`created_date` ,
		NEW.`created_date` ,
		(SELECT `boiler`.`name`
		 FROM	`boiler`
		 WHERE	`boiler`.`uid` = NEW.`boiler_id`),
		(SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`)
		 FROM	`runtime_parameter` AS `param`
		 WHERE	`param`.id = NEW.`parameter_id`) ,
		(SELECT `param`.`name` FROM `runtime_parameter` AS `param` WHERE `param`.id = NEW.`parameter_id`) ,
		(SELECT `param`.`unit` FROM `runtime_parameter` AS `param` WHERE `param`.id = NEW.`parameter_id`) ,
		(SELECT `alarm`.`priority` FROM `boiler_alarm` AS `alarm` WHERE `alarm`.uid = NEW.`alarm_id`) ,
		(SELECT `alarm`.`description` FROM `boiler_alarm` AS `alarm` WHERE `alarm`.uid = NEW.`alarm_id`),
		NEW.`remark`
	);

	UPDATE 	`boiler_runtime_cache_instant`
	SET 	`runtime_id` = NEW.`id` ,
		`alarm_id` = NEW.`alarm_id` ,
		`updated_date` = NEW.`created_date` ,
		`value`	= (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`) 
			FROM `runtime_parameter` AS `param` 
		   	WHERE `param`.id = NEW.`parameter_id`) ,
		`name` = (SELECT `param`.`name` FROM `runtime_parameter` AS `param` WHERE `param`.id = NEW.`parameter_id`) ,
		`unit` = (SELECT `param`.`unit` FROM `runtime_parameter` AS `param` WHERE `param`.id = NEW.`parameter_id`) ,
		`alarm_level` = IFNULL((SELECT `alarm`.`priority` FROM `boiler_alarm` AS `alarm` WHERE `alarm`.`uid` = NEW.alarm_id LIMIT 1), 0) ,
		`alarm_description` = IFNULL((SELECT `alarm`.`description` FROM `boiler_alarm` AS `alarm` WHERE `alarm`.`uid` = NEW.alarm_id LIMIT 1), ''),
		`remark` = NEW.`remark`
	WHERE	`boiler_id` = NEW.`boiler_id` 
	AND	`parameter_id` = NEW.`parameter_id` 
	AND	NEW.`created_date` > `updated_date`;

	UPDATE `boiler_runtime_cache_instant` SET `alarm_level` = 0 WHERE `alarm_level` IS NULL;

	UPDATE `boiler_runtime_cache_instant` SET `alarm_description` = '' WHERE `alarm_description` IS NULL;

	/* HISTORY */
	INSERT IGNORE `boiler_runtime_history`(
		`name` ,
		`created_date` ,
		`boiler_id`,
		`json_data`,
		`is_demo`, 
		`remark`) 
	SELECT * 
	FROM (SELECT `boiler`.`name`, NEW.`created_date`, `boiler`.`uid` , '[]', NEW.`is_demo`, 'triggered' FROM `boiler` WHERE `boiler`.`uid` = NEW.`boiler_id`) AS `tmp` 
	WHERE NOT EXISTS
		(SELECT `history`.* 
		 FROM 	`boiler_runtime_history` AS `history` , `boiler_runtime` AS `runtime` 
		 WHERE 	`history`.`boiler_id` = NEW.`boiler_id` AND `history`.`created_date` >(NOW() - INTERVAL 5 MINUTE) AND `history`.`created_date` <= NOW()) 
	LIMIT 1;
	
	/* CACHES */
	IF NEW.`parameter_id` = 1001 THEN
		SELECT 'boiler_runtime_cache_steam_temperature' INTO @tableName;
		INSERT INTO `boiler_runtime_cache_steam_temperature`
		SET 	`runtime_id` = NEW.`id`,
			`boiler_id` = NEW.`boiler_id`, 
			`parameter_id` = NEW.`parameter_id`,
			`alarm_id` = NEW.`alarm_id`,

			`created_date` = NEW.`created_date`,
			`created_by_id` = NEW.`created_by_id`,
			`updated_date` = NEW.`updated_date`,
			`updated_by_id` = NEW.`updated_by_id`,
			`is_deleted` = NEW.`is_deleted`,
			`is_demo` = NEW.`is_demo`,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name` 
				FROM	`boiler`
				WHERE	`boiler`.`uid` = NEW.`boiler_id`),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`) 
				FROM	`runtime_parameter` AS `param`
				WHERE	`param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name` 
				FROM	`runtime_parameter` AS `param`
				WHERE	`param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit` 
				FROM	`runtime_parameter` AS `param`
				WHERE	`param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level` 
				FROM	`boiler_alarm` AS `alarm`
				WHERE	`alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description` 
				FROM	`boiler_alarm` AS `alarm`
				WHERE	`alarm`.uid = NEW.alarm_id)
			;

		UPDATE 	`boiler_runtime_cache_steam_temperature`
		SET 	`alarm_level` = 0
		WHERE 	`alarm_level` IS NULL;

		UPDATE 	`boiler_runtime_cache_steam_temperature`
		SET 	`alarm_description` = ''
		WHERE 	`alarm_description` IS NULL;
	END IF;

	IF NEW.parameter_id = 1002 THEN
		INSERT IGNORE `boiler_runtime_cache_steam_pressure`
		SET runtime_id = NEW.id,
			boiler_id = NEW.boiler_id, 
			parameter_id = NEW.parameter_id,
			alarm_id = NEW.alarm_id,

			created_date = NEW.created_date,
			created_by_id = NEW.created_by_id,
			updated_date = NEW.updated_date,
			updated_by_id = NEW.updated_by_id,
			is_deleted = NEW.is_deleted,
			is_demo = NEW.is_demo,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name` 
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`) 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE 	`boiler_runtime_cache_steam_pressure`
		SET 	`alarm_level` = 0
		WHERE 	`alarm_level` IS NULL;

		UPDATE 	`boiler_runtime_cache_steam_pressure`
		SET 	`alarm_description` = ''
		WHERE 	`alarm_description` IS NULL;
	END IF;

	IF NEW.`parameter_id` = 1003 THEN
		INSERT IGNORE `boiler_runtime_cache_flow`
		SET 	`runtime_id` = NEW.`id`,
			`boiler_id` = NEW.`boiler_id`, 
			`parameter_id` = NEW.`parameter_id`,
			`alarm_id` = NEW.`alarm_id`,

			created_date = NEW.created_date,
			created_by_id = NEW.created_by_id,
			updated_date = NEW.updated_date,
			updated_by_id = NEW.updated_by_id,
			is_deleted = NEW.is_deleted,
			is_demo = NEW.is_demo,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name` 
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`) 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE `boiler_runtime_cache_flow`
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE `boiler_runtime_cache_flow`
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;
	END IF;


	IF NEW.parameter_id = 1005 OR NEW.parameter_id = 1006 THEN
		INSERT IGNORE `boiler_runtime_cache_water_temperature`
		SET 	`runtime_id` = NEW.id,
			`boiler_id` = NEW.boiler_id, 
			parameter_id = NEW.parameter_id,
			alarm_id = NEW.alarm_id,

			created_date = NEW.created_date,
			created_by_id = NEW.created_by_id,
			updated_date = NEW.updated_date,
			updated_by_id = NEW.updated_by_id,
			is_deleted = NEW.is_deleted,
			is_demo = NEW.is_demo,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name` 
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`) 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE boiler_runtime_cache_water_temperature
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE boiler_runtime_cache_water_temperature
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;
	END IF;

	IF NEW.parameter_id = 1014 OR NEW.parameter_id = 1015 THEN
		INSERT IGNORE boiler_runtime_cache_smoke_temperature
		SET runtime_id = NEW.id,
			boiler_id = NEW.boiler_id, 
			parameter_id = NEW.parameter_id,
			alarm_id = NEW.alarm_id,

			created_date = NEW.created_date,
			created_by_id = NEW.created_by_id,
			updated_date = NEW.updated_date,
			updated_by_id = NEW.updated_by_id,
			is_deleted = NEW.is_deleted,
			is_demo = NEW.is_demo,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name` 
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`) 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE boiler_runtime_cache_smoke_temperature
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE boiler_runtime_cache_smoke_temperature
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;
	END IF;

	IF NEW.parameter_id = 1016 OR NEW.parameter_id = 1017 OR NEW.parameter_id = 1018 OR NEW.parameter_id = 1019 THEN
		INSERT IGNORE boiler_runtime_cache_smoke_component
		SET runtime_id = NEW.id,
			boiler_id = NEW.boiler_id, 
			parameter_id = NEW.parameter_id,
			alarm_id = NEW.alarm_id,

			created_date = NEW.created_date,
			created_by_id = NEW.created_by_id,
			updated_date = NEW.updated_date,
			updated_by_id = NEW.updated_by_id,
			is_deleted = NEW.is_deleted,
			is_demo = NEW.is_demo,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name` 
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`) 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE boiler_runtime_cache_smoke_component
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE boiler_runtime_cache_smoke_component
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;
	END IF;

	IF NEW.`parameter_id` = 1021 OR NEW.`parameter_id` = 1022 THEN
		INSERT IGNORE `boiler_runtime_cache_environment_temperature`
		SET 	`runtime_id` = NEW.`id`,
			`boiler_id` = NEW.`boiler_id`, 
			`parameter_id` = NEW.`parameter_id`,
			`alarm_id` = NEW.`alarm_id`,

			`created_date` = NEW.`created_date`,
			`created_by_id` = NEW.`created_by_id`,
			`updated_date` = NEW.`updated_date`,
			`updated_by_id` = NEW.`updated_by_id`,
			`is_deleted` = NEW.`is_deleted`,
			`is_demo` = NEW.`is_demo`,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name` 
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`) 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.`parameter_id`),
			`parameter_name` = (SELECT `param`.`name` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.`parameter_id`),
			`unit` = (SELECT `param`.`unit` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.`parameter_id`),

			`alarm_level` = (SELECT `alarm`.`alarm_level` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.`uid` = NEW.`alarm_id`),
			`alarm_description` = (SELECT `alarm`.`description` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.`uid` = NEW.`alarm_id`)
			;

		UPDATE 	`boiler_runtime_cache_environment_temperature`
		SET 	`alarm_level` = 0
		WHERE 	`alarm_level` IS NULL;

		UPDATE 	`boiler_runtime_cache_environment_temperature`
		SET 	`alarm_description` = ''
		WHERE 	`alarm_description` IS NULL;
	END IF;

	IF NEW.`parameter_id` = 1201 THEN
		INSERT IGNORE `boiler_runtime_cache_heat`
		SET 	`runtime_id` = NEW.`id`,
			`boiler_id` = NEW.`boiler_id`, 
			`parameter_id` = NEW.`parameter_id`,
			`alarm_id` = NEW.`alarm_id`,

			`created_date` = NEW.`created_date`,
			`created_by_id` = NEW.`created_by_id`,
			`updated_date` = NEW.`updated_date`,
			`updated_by_id` = NEW.`updated_by_id`,
			`is_deleted` = NEW.`is_deleted`,
			`is_demo` = NEW.`is_demo`,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name` 
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`) 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE 	`boiler_runtime_cache_heat`
		SET 	`alarm_level` = 0
		WHERE 	`alarm_level` IS NULL;

		UPDATE 	`boiler_runtime_cache_heat`
		SET 	`alarm_description` = ''
		WHERE 	`alarm_description` IS NULL;
	END IF;

	IF NEW.parameter_id = 1202 THEN
		INSERT IGNORE boiler_runtime_cache_excess_air
		SET 	`runtime_id` = NEW.`id`,
			`boiler_id` = NEW.`boiler_id`, 
			`parameter_id` = NEW.`parameter_id`,
			`alarm_id` = NEW.`alarm_id`,

			`created_date` = NEW.`created_date`,
			`created_by_id` = NEW.`created_by_id`,
			`updated_date` = NEW.`updated_date`,
			`updated_by_id` = NEW.`updated_by_id`,
			`is_deleted` = NEW.`is_deleted`,
			`is_demo` = NEW.`is_demo`,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name` 
				FROM `boiler`
				WHERE `boiler`.uid = NEW.boiler_id),

			`value` = (SELECT ROUND(`param`.`scale` * NEW.`value`, `param`.`fix`) 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`parameter_name` = (SELECT `param`.`name` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE 	`boiler_runtime_cache_excess_air`
		SET 	`alarm_level` = 0
		WHERE 	`alarm_level` IS NULL;

		UPDATE 	`boiler_runtime_cache_excess_air`
		SET 	`alarm_description` = ''
		WHERE 	`alarm_description` IS NULL;
	END IF;

	IF NEW.`parameter_id` DIV 100 = 11 THEN
		INSERT IGNORE `boiler_runtime_cache_status`
		SET 	`runtime_id` = NEW.`id`,
			`boiler_id` = NEW.`boiler_id`, 
			`parameter_id` = NEW.`parameter_id`,
			`alarm_id` = NEW.`alarm_id`,

			`created_date` = NEW.`created_date`,
			`created_by_id` = NEW.`created_by_id`,
			`updated_date` = NEW.`updated_date`,
			`updated_by_id` = NEW.`updated_by_id`,
			`is_deleted` = NEW.`is_deleted`,
			`is_demo` = NEW.`is_demo`,

			`name_en` = NEW.`name_en`,
			`remark` = NEW.`remark`,
			`name` = (SELECT `boiler`.`name` 
				FROM `boiler`
				WHERE `boiler`.`uid` = NEW.`boiler_id`),

			`value` = IF(NEW.`value` > 0, 1, 0),
			`parameter_name` = (SELECT `param`.`name` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),
			`unit` = (SELECT `param`.`unit` 
				FROM `runtime_parameter` AS `param`
				WHERE `param`.id = NEW.parameter_id),

			`alarm_level` = (SELECT `alarm`.`alarm_level` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id),
			`alarm_description` = (SELECT `alarm`.`description` 
				FROM `boiler_alarm` AS `alarm`
				WHERE `alarm`.uid = NEW.alarm_id)
			;

		UPDATE `boiler_runtime_cache_status`
		SET `alarm_level` = 0
		WHERE `alarm_level` IS NULL;

		UPDATE `boiler_runtime_cache_status`
		SET `alarm_description` = ''
		WHERE `alarm_description` IS NULL;
	END IF;
END
 ;;
delimiter ;

delimiter ;;
-- ----------------------------
--  Triggers structure for table boiler_runtime_cache_status
-- ----------------------------
 ;;
delimiter ;
DROP TRIGGER IF EXISTS `running_duration`;
delimiter ;;
CREATE TRIGGER `running_duration` BEFORE INSERT ON `boiler_runtime_cache_status` FOR EACH ROW BEGIN
	IF NEW.`parameter_id` = 1107 THEN
		SELECT 	`status`.`id`, `status`.`value` INTO @statusId, @statusValue
		FROM	`boiler_runtime_cache_status_running` AS `status`
		WHERE	`status`.`boiler_id` = NEW.`boiler_id` AND `status`.`parameter_id` = 1107
		ORDER BY `status`.`created_date` DESC
		LIMIT	1;

		UPDATE 	`boiler_runtime_cache_status_running` AS `status`
		SET	`status`.`duration` = TIMESTAMPDIFF(MICROSECOND, `status`.`created_date`, NEW.`created_date`)
		WHERE	`status`.`id` = @statusId;

		IF @statusValue <> NEW.`value` THEN
			INSERT IGNORE `boiler_runtime_cache_status_running`
			SET	`name` = NEW.`name`,
				`name_en` = NEW.`name_en`,
				`remark` = NEW.`remark`,
				`created_date` = NEW.`created_date`,
				`is_demo` = NEW.`is_demo`,
				`runtime_id` = NEW.`runtime_id`,
				`boiler_id` = NEW.`boiler_id`,
				`parameter_id` = NEW.`parameter_id`,
				`parameter_name` = NEW.`parameter_name`,
				`value` = NEW.`value`,
				`alarm_level` = 0,
				`alarm_description` = '',
				`date` = DATE(NEW.`created_date`),
				`duration` = 0;
		END IF;
	END IF;
END
 ;;
delimiter ;

SET FOREIGN_KEY_CHECKS = 1;
