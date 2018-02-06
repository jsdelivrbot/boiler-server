docker run -it --rm mysql mysql -h172.17.0.2 -uroot -P3306 -p

docker run -it --rm mysql mysql -h183.195.133.16 -uroot -P3306 -p

docker run -itd --name boilermysql -p 12006:3306 -p 12007:3306 -p 12008:3306 azuretech/boiler-mysql /bin/bash

service mysql stop

mysqld_safe --skip-grant-tables &

mysql -u root

UPDATE mysql.user
    SET authentication_string = PASSWORD('hold2017'), password_expired = 'N'
    WHERE User = 'root';
FLUSH PRIVILEGES;
