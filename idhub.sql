DROP DATABASE IF EXISTS `idhub`;
CREATE DATABASE `idhub`
  DEFAULT CHARACTER SET utf8
  DEFAULT COLLATE utf8_general_ci;
USE idhub;
SET NAMES utf8;

DROP TABLE IF EXISTS `idhub`.`users`;
CREATE TABLE `idhub`.`users` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) CHARACTER SET utf8 NOT NULL,
  `idnumber` varchar(20) DEFAULT NULL,
  `phone` varchar(16) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  `country` varchar(20) DEFAULT NULL,
  `region` varchar(20) DEFAULT NULL,
  `locality` varchar(20) DEFAULT NULL,
  `street_address` varchar(100) DEFAULT NULL,
  `postal_code` varchar(20) DEFAULT NULL,
  `privatekey` varchar(120) DEFAULT NULL,
  `publickey` varchar(150) NOT NULL,
  `address` varchar(70) NOT NULL,
  `proxy` varchar(50) NOT NULL,
  `controller` varchar(50) DEFAULT NULL,
  `recovery` varchar(50) DEFAULT NULL,
  `ipfs` varchar(60) DEFAULT NULL,
  `description` varchar(300) DEFAULT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `idhub`.`tokens`;
CREATE TABLE `idhub`.`tokens` (
  `token` varchar(40) NOT NULL,
  `valid` tinyint(1) UNSIGNED DEFAULT NULL,
  `proxy` varchar(50) NOT NULL,
  `scope` varchar(100) DEFAULT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `idhub`.`claims`;
CREATE TABLE `idhub`.`claims` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `proxy` varchar(50) NOT NULL,
  `type` varchar(30) DEFAULT NULL,
  `status` varchar(15) NOT NULL,
  `claim` text DEFAULT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `idhub`.`attestations`;
CREATE TABLE `idhub`.`attestations` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `claimid` int(10) NOT NULL,
  `attestant` varchar(30) DEFAULT NULL,
  `attestation` text NOT NULL,
  `status` varchar(15) NOT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
