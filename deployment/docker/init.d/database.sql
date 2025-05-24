DROP DATABASE IF EXISTS lubricant;
CREATE DATABASE if not exists lubricant;
CREATE USER if not exists 'lubricant'@'%' IDENTIFIED BY '123456';
GRANT ALL PRIVILEGES ON lubricant.* TO 'lubricant'@'%';

DROP DATABASE IF EXISTS casdoor;
CREATE DATABASE casdoor;
CREATE USER if not exists 'casdoor'@'%' IDENTIFIED BY '123456';
GRANT ALL PRIVILEGES ON casdoor.* TO 'casdoor'@'%';
FLUSH PRIVILEGES;
