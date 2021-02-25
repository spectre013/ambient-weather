CREATE DATABASE IF NOT EXISTS ambient;
USE ambient;

CREATE TABLE IF NOT EXISTS `ambient.records` (
        `id` int NOT NULL AUTO_INCREMENT,
        `mac` varchar(255) DEFAULT NULL,
        `date` datetime DEFAULT NULL,
        `baromabsin` double DEFAULT NULL,
        `baromrelin` double DEFAULT NULL,
        `battout` int DEFAULT NULL,
        `Batt1` int DEFAULT NULL,
        `Batt2` int DEFAULT NULL,
        `Batt3` int DEFAULT NULL,
        `Batt4` int DEFAULT NULL,
        `Batt5` int DEFAULT NULL,
        `Batt6` int DEFAULT NULL,
        `Batt7` int DEFAULT NULL,
        `Batt8` int DEFAULT NULL,
        `Batt9` int DEFAULT NULL,
        `Batt10` int DEFAULT NULL,
        `co2` double DEFAULT NULL,
        `dailyrainin` double DEFAULT NULL,
        `dewpoint` double DEFAULT NULL,
        `eventrainin` double DEFAULT NULL,
        `feelslike` double DEFAULT NULL,
        `hourlyrainin` double DEFAULT NULL,
        `hourlyrain` double DEFAULT NULL,
        `humidity` int DEFAULT NULL,
        `humidity1` int DEFAULT NULL,
        `humidity2` int DEFAULT NULL,
        `humidity3` int DEFAULT NULL,
        `humidity4` int DEFAULT NULL,
        `humidity5` int DEFAULT NULL,
        `humidity6` int DEFAULT NULL,
        `humidity7` int DEFAULT NULL,
        `humidity8` int DEFAULT NULL,
        `humidity9` int DEFAULT NULL,
        `humidity10` int DEFAULT NULL,
        `humidityin` int DEFAULT NULL,
        `lastrain` datetime DEFAULT NULL,
        `maxdailygust` double DEFAULT NULL,
        `relay1` int DEFAULT NULL,
        `relay2` int DEFAULT NULL,
        `relay3` int DEFAULT NULL,
        `relay4` int DEFAULT NULL,
        `relay5` int DEFAULT NULL,
        `relay6` int DEFAULT NULL,
        `relay7` int DEFAULT NULL,
        `relay8` int DEFAULT NULL,
        `relay9` int DEFAULT NULL,
        `relay10` int DEFAULT NULL,
        `monthlyrainin` double DEFAULT NULL,
        `solarradiation` double DEFAULT NULL,
        `tempf` double DEFAULT NULL,
        `temp1f` double DEFAULT NULL,
        `temp2f` double DEFAULT NULL,
        `temp3f` double DEFAULT NULL,
        `temp4f` double DEFAULT NULL,
        `temp5f` double DEFAULT NULL,
        `temp6f` double DEFAULT NULL,
        `temp7f` double DEFAULT NULL,
        `temp8f` double DEFAULT NULL,
        `temp9f` double DEFAULT NULL,
        `temp10f` double DEFAULT NULL,
        `tempinf` double DEFAULT NULL,
        `totalrainin` double DEFAULT NULL,
        `uv` double DEFAULT NULL,
        `weeklyrainin` double DEFAULT NULL,
        `winddir` int DEFAULT NULL,
        `windgustmph` double DEFAULT NULL,
        `windgustdir` int DEFAULT NULL,
        `windspeedmph` double DEFAULT NULL,
        `yearlyrainin` double DEFAULT NULL,
        `battlightning` int DEFAULT NULL,
        `lightningday` int DEFAULT '0',
        `lightninghour` int DEFAULT '0',
        `lightningtime` datetime DEFAULT NULL,
        `lightningdistance` double DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `date` (`date`),
        KEY `DateTempfIDX` (`date`,`tempf`),
        KEY `tempfIDX` (`tempf`,`date`),
        KEY `tempinfIDX` (`tempinf`,`date`),
        KEY `temp1fIDX` (`temp1f`,`date`),
        KEY `temp2fIDX` (`temp2f`,`date`),
        KEY `baromrelinIDX` (`baromrelin`,`date`),
        KEY `uvIDX` (`uv`,`date`),
        KEY `humidityIDX` (`humidity`,`date`),
        KEY `windspeedmphIDX` (`windspeedmph`,`date`),
        KEY `windgustmphIDX` (`windgustmph`,`date`),
        KEY `dewpointIDX` (`dewpoint`,`date`),
        KEY `humidityinIDX` (`humidityin`,`date`),
        KEY `humidity1IDX` (`humidity1`,`date`),
        KEY `humidity2IDX` (`humidity2`,`date`),
        KEY `dailyraininIDX` (`dailyrainin`,`date`),
        KEY `lightningdayIDX` (`lightningday`,`date`)
)


CREATE TABLE IF NOT EXISTS `ambient.stat` (
    `id` varchar(100) NOT NULL DEFAULT '',
    `date` datetime DEFAULT NULL,
    `value` double DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;