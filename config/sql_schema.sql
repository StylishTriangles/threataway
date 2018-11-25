CREATE DATABASE  IF NOT EXISTS `gowebapp` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `gowebapp`;
-- MySQL dump 10.13  Distrib 5.7.17, for Win64 (x86_64)
--
-- Host: localhost    Database: gowebapp
-- ------------------------------------------------------
-- Server version	5.7.20-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `balance`
--

DROP TABLE IF EXISTS `balance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `balance` (
  `idBalance` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `USD` decimal(24,12) NOT NULL DEFAULT '0.000000000000',
  `BTC` decimal(24,12) NOT NULL DEFAULT '0.000000000000',
  `ETH` decimal(24,12) NOT NULL DEFAULT '0.000000000000',
  PRIMARY KEY (`idBalance`),
  UNIQUE KEY `idUser_UNIQUE` (`idBalance`),
  KEY `idUser_idx` (`idBalance`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `balance`
--

LOCK TABLES `balance` WRITE;
/*!40000 ALTER TABLE `balance` DISABLE KEYS */;
INSERT INTO `balance` VALUES (1,10.000000000000,10.001000000000,10.000000000000),(3,0.000000000000,0.000000000000,0.000000000000);
/*!40000 ALTER TABLE `balance` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `confirm`
--

DROP TABLE IF EXISTS `confirm`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
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
  CONSTRAINT `idUser` FOREIGN KEY (`idUser`) REFERENCES `users` (`idUser`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `confirm`
--

LOCK TABLES `confirm` WRITE;
/*!40000 ALTER TABLE `confirm` DISABLE KEYS */;
INSERT INTO `confirm` VALUES (1,1,'22','2137-12-31 23:59:59',1);
/*!40000 ALTER TABLE `confirm` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
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
  `idBalance` int(10) unsigned NOT NULL,
  PRIMARY KEY (`idUser`),
  UNIQUE KEY `idUser_UNIQUE` (`idUser`),
  UNIQUE KEY `login_UNIQUE` (`login`),
  UNIQUE KEY `email_UNIQUE` (`email`),
  KEY `idBalance_idx` (`idBalance`),
  CONSTRAINT `idBalance` FOREIGN KEY (`idBalance`) REFERENCES `balance` (`idBalance`) ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'admin','admin@gmail.com','$2a$10$Oa9eDK/T9lUADqd/QYtZLu7/z.r3eKlpLNSVYMltYktEUZWpoyI/i','Łukasz','Ptak','0000-00-00 00:00:00','2018-01-06 15:11:53',1,1),(3,'domingo1337','dgrasz69@gmail.com','$2a$12$M//1XeGaBH.xlbGS5la7m.L5nYd6T6JcQAA9t/9.RbYl54zPhSkLG','','','0000-00-00 00:00:00','2018-01-10 21:10:40',0,3);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-01-14 23:49:12
