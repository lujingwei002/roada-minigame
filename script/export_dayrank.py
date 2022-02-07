#!/usr/bin/python
# -*- coding:utf-8 -*-
import os
import time
import shutil
import sys
import os.path
import codecs
import MySQLdb
import csv

mysql_host = sys.argv[1]
mysql_port = sys.argv[2]
mysql_user = sys.argv[3]
mysql_pwd = sys.argv[4]
mysql_db = sys.argv[5]
rank_table = sys.argv[6]
save_dir = sys.argv[7]

reload(sys) 
sys.setdefaultencoding('utf-8') 
now = time.localtime(time.time())
rankid = now.tm_year*10000+now.tm_mon*100+now.tm_mday
#rankid = 20210712 
db=MySQLdb.connect(host=mysql_host,user=mysql_user,passwd=mysql_pwd,db=mysql_db,charset="utf8") 
cursor = db.cursor()
f = codecs.open("%s/%s_%s_%d.csv"%(save_dir, mysql_db, rank_table, rankid), "w+")
writer = csv.writer(f)
writer.writerow(['openid', '昵称', '分数', '创建时间', '更新时间'])
print(rankid, rank_table)
n = cursor.execute('''SELECT a.rankid, b.openid, b.nickname, a.score, a.createtime, a.updatetime 
    FROM %s AS a
    LEFT JOIN user AS b ON a.userid=b.userid
    WHERE a.rankid=%d
    ORDER BY score desc, updatetime DESC limit 100'''%(rank_table, rankid))   

for row in cursor.fetchall():    
    writer.writerow([row[1], row[2], row[3], 
        time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(row[4])), 
        time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(row[5]))])
cursor.close()
db.close()
f.close()