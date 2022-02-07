#!/bin/bash
# 数据库信息
DB_USER="root"
DB_PASS="123456"
DB_HOST="localhost"
DB_NAME=(daxigua bpbxq caichengyu fangkuainiao gongjianchuanshu paopaolong)   #需要备份的数据库名称,注意中间用空格隔开.
 
# 其他设置
BIN_DIR="/usr/bin"   # mysqldump命令执行路径
BACK_DIR="/home/backup/mysql"    #备份目录，这里设为/home/backup/mysql
DATE=`date +%Y%m%d%H%M%S`       #显示备份时间，格式为20180808122556
# 备份所有指定数据库
for backdb in ${DB_NAME[@]}
#也可以写成for eachdb in ${DB_NAME[*]}
do
$BIN_DIR/mysqldump --opt -u$DB_USER -p$DB_PASS -h$DB_HOST -B ${backdb} | gzip > $BACK_DIR/db_${backdb}_$DATE.sql.gz
done
# 删除5天之前的备份文件，但保留日期为1号的文件（用于手动删除）
find $BACK_DIR/* -regextype "posix-extended" -not -regex ".*[0-9]{6}01[0-9]{6}\.sql\.gz$" -mtime +5 -exec rm {} \;