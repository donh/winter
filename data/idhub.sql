DROP DATABASE IF EXISTS `idhub`;
CREATE DATABASE `idhub`
  DEFAULT CHARACTER SET utf8
  DEFAULT COLLATE utf8_general_ci;
USE idhub;
SET NAMES utf8;

DROP TABLE IF EXISTS `idhub`.`users`;
CREATE TABLE `idhub`.`users` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  -- `name` varchar(128) CHARACTER SET utf8 NOT NULL,
  `firstName` varchar(20) CHARACTER SET utf8 NOT NULL,
  `lastName` varchar(20) CHARACTER SET utf8 NOT NULL,
  -- `idnumber` varchar(20) DEFAULT NULL,
  `phone` varchar(16) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  -- `wallet` varchar(255) DEFAULT NULL,
  `country` varchar(20) DEFAULT NULL,
  `region` varchar(20) DEFAULT NULL,
  -- `locality` varchar(20) DEFAULT NULL,
  `city` varchar(20) DEFAULT NULL,
  -- `street_address` varchar(100) DEFAULT NULL,
  `street` varchar(100) DEFAULT NULL,
  -- `postal_code` varchar(20) DEFAULT NULL,
  `postalCode` varchar(20) DEFAULT NULL,
  -- `privatekey` varchar(120) DEFAULT NULL,
  `publickey` varchar(150) NOT NULL,
  `wallet` varchar(70) NOT NULL,
  -- `address` varchar(70) NOT NULL,
  `proxy` varchar(50) NOT NULL,
  -- `controller` varchar(50) DEFAULT NULL,
  -- `recovery` varchar(50) DEFAULT NULL,
  -- `ipfs` varchar(60) DEFAULT NULL,
  -- `description` varchar(300) DEFAULT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `idhub`.`payments`;
CREATE TABLE `idhub`.`payments` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `payid` varchar(30) NOT NULL,
  `companyCode` int(5) DEFAULT NULL,
  -- `payer` varchar(10) DEFAULT NULL,
  `paymentStatus` varchar(20) DEFAULT NULL,
  `currency` varchar(5) NOT NULL,
  `total` DECIMAL(13,2) NOT NULL,
  `subtotal` DECIMAL(13,2) NOT NULL,
  -- `handling_fee` DECIMAL(13,2) DEFAULT NULL,
  -- `insurance` DECIMAL(13,2) DEFAULT NULL,
  `shipping` DECIMAL(13,2) DEFAULT NULL,
  -- `shipping_discount` DECIMAL(13,2) DEFAULT NULL,
  `commission` DECIMAL(13,2) DEFAULT NULL,
  -- `tax` DECIMAL(13,2) DEFAULT NULL,
  `exchangeRate` DECIMAL(6,6) DEFAULT NULL,
  `eth` DECIMAL(10,6) DEFAULT NULL,
  `firstName` varchar(20) CHARACTER SET utf8 NOT NULL,
  `lastName` varchar(20) CHARACTER SET utf8 NOT NULL,
  `phone` varchar(16) NOT NULL,
  `wallet` varchar(70) NOT NULL,
  -- `return_url` varchar(150) DEFAULT NULL,
  -- `cancel_url` varchar(150) DEFAULT NULL,
  `note` varchar(200) DEFAULT NULL,
  `deadline` datetime NOT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `idhub`.`transactions`;
CREATE TABLE `idhub`.`transactions` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `payid` varchar(30) NOT NULL,
  `amount` DECIMAL(13,2) DEFAULT NULL,
  `custom` varchar(30) DEFAULT NULL,
  `description` varchar(50) DEFAULT NULL,
  `invoice_number` varchar(15) DEFAULT NULL,
  `soft_descriptor` varchar(15) DEFAULT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `idhub`.`amounts`;
CREATE TABLE `idhub`.`amounts` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `txid` int(10) NOT NULL,
  `currency` varchar(5) NOT NULL,
  `total` DECIMAL(13,2) NOT NULL,
  `subtotal` DECIMAL(13,2) NOT NULL,
  `handling_fee` DECIMAL(13,2) DEFAULT NULL,
  `insurance` DECIMAL(13,2) DEFAULT NULL,
  `shipping` DECIMAL(13,2) DEFAULT NULL,
  `shipping_discount` DECIMAL(13,2) DEFAULT NULL,
  `tax` DECIMAL(13,2) DEFAULT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `idhub`.`items`;
CREATE TABLE `idhub`.`items` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `txid` int(10) NOT NULL,
  `currency` varchar(5) NOT NULL,
  `description` varchar(100) DEFAULT NULL,
  `name` varchar(50) NOT NULL,
  `price` DECIMAL(13,2) NOT NULL,
  `quantity` int(8) NOT NULL,
  `sku` varchar(30) DEFAULT NULL,
  `tax` DECIMAL(13,2) DEFAULT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `idhub`.`shipping`;
CREATE TABLE `idhub`.`shipping` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `txid` int(10) NOT NULL,
  `city` varchar(20) NOT NULL,
  `country_code` varchar(3) NOT NULL,
  `line1` varchar(50) DEFAULT NULL,
  `line2` varchar(50) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `postal_code` varchar(8) DEFAULT NULL,
  `recipient_name` varchar(30) NOT NULL,
  `state` varchar(20) DEFAULT NULL,
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
