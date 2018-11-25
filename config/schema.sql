CREATE DATABASE  IF NOT EXISTS `gowebapp` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `gowebapp`;
-- MySQL dump 10.13  Distrib 8.0.13, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: gowebapp
-- ------------------------------------------------------
-- Server version	8.0.13

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
 SET NAMES utf8 ;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `confirm`
--

DROP TABLE IF EXISTS `confirm`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `confirm` (
  `idConfirm` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `idUser` int(10) unsigned NOT NULL,
  `key` varchar(64) NOT NULL,
  `expires` datetime NOT NULL DEFAULT '2137-12-31 23:59:59' COMMENT 'Defaults to never',
  `used` tinyint(1) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`idConfirm`),
  UNIQUE KEY `idConfirm_UNIQUE` (`idConfirm`),
  UNIQUE KEY `key_UNIQUE` (`key`),
  KEY `idUser_idx` (`idUser`),
  CONSTRAINT `idUser` FOREIGN KEY (`idUser`) REFERENCES `users` (`iduser`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `deployments`
--

DROP TABLE IF EXISTS `deployments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `deployments` (
  `deploymentID` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `url` text NOT NULL,
  `listID` int(10) unsigned NOT NULL,
  `templateID` int(10) unsigned NOT NULL,
  PRIMARY KEY (`deploymentID`),
  KEY `listID` (`listID`),
  KEY `templateID` (`templateID`),
  CONSTRAINT `deployments_ibfk_1` FOREIGN KEY (`listID`) REFERENCES `lists` (`idlist`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `deployments_ibfk_2` FOREIGN KEY (`templateID`) REFERENCES `templates` (`templateid`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `listlists`
--

DROP TABLE IF EXISTS `listlists`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `listlists` (
  `idURL` int(11) unsigned NOT NULL,
  `idList` int(11) unsigned NOT NULL,
  `dirty` int(11) NOT NULL DEFAULT '0',
  UNIQUE KEY `uniquePairs` (`idURL`,`idList`),
  KEY `listID_idx` (`idList`),
  KEY `urlID_idx` (`idURL`),
  CONSTRAINT `idList` FOREIGN KEY (`idList`) REFERENCES `lists` (`idlist`) ON DELETE CASCADE,
  CONSTRAINT `idURL` FOREIGN KEY (`idURL`) REFERENCES `urls` (`idurl`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `lists`
--

DROP TABLE IF EXISTS `lists`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `lists` (
  `idList` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `ownerID` int(10) unsigned NOT NULL,
  PRIMARY KEY (`idList`),
  UNIQUE KEY `idList_UNIQUE` (`idList`),
  UNIQUE KEY `name_UNIQUE` (`name`),
  KEY `ownerID_idx` (`ownerID`),
  KEY `name_idx` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `templates`
--

DROP TABLE IF EXISTS `templates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `templates` (
  `templateID` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(256) NOT NULL,
  `header` text NOT NULL,
  `footer` text NOT NULL,
  `urlTemplate` text NOT NULL,
  PRIMARY KEY (`templateID`),
  UNIQUE KEY `templateID_UNIQUE` (`templateID`),
  KEY `name_idx` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `urls`
--

DROP TABLE IF EXISTS `urls`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `urls` (
  `idUrl` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `domain` varchar(256) NOT NULL,
  `rating` float NOT NULL DEFAULT '0',
  `malicious` int(11) NOT NULL DEFAULT '0',
  `malicious_type` text NOT NULL,
  `hostname` text NOT NULL,
  `organization` text NOT NULL,
  `web_score_name` varchar(16) NOT NULL DEFAULT '',
  `email_score_name` varchar(16) NOT NULL DEFAULT '',
  `monthly_spam_level` int(11) NOT NULL DEFAULT '0',
  `honeypot_score` float NOT NULL DEFAULT '0',
  `shodan_malware` int(11) NOT NULL DEFAULT '0',
  `shodan_malware_query` text NOT NULL,
  `shodan_creds` int(11) NOT NULL DEFAULT '0',
  `shodan_creds_query` text NOT NULL,
  PRIMARY KEY (`idUrl`),
  UNIQUE KEY `url_UNIQUE` (`domain`),
  UNIQUE KEY `idUrl_UNIQUE` (`idUrl`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `users` (
  `idUser` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `login` varchar(32) NOT NULL,
  `email` varchar(256) NOT NULL,
  `password` varbinary(64) NOT NULL COMMENT 'Bcrypt hashed password',
  `firstName` varchar(45) NOT NULL DEFAULT '' COMMENT 'First name, defaults to empty string',
  `lastName` varchar(45) NOT NULL DEFAULT '' COMMENT 'Last name, defaults to empty string',
  `lastLogin` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `createDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `active` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `role` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '0 - readonly\n1 - r/w',
  PRIMARY KEY (`idUser`),
  UNIQUE KEY `idUser_UNIQUE` (`idUser`),
  UNIQUE KEY `login_UNIQUE` (`login`),
  UNIQUE KEY `email_UNIQUE` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-11-25 10:18:07
-- Insert admin account

LOCK TABLES users WRITE;

INSERT INTO `users` VALUES (1,'admin','admin@gmail.com',_binary '$2a$10$Oa9eDK/T9lUADqd/QYtZLu7/z.r3eKlpLNSVYMltYktEUZWpoyI/i','Łukasz','Ptak','0000-00-00 00:00:00','2018-01-06 15:11:53',1,1)