#!/usr/bin/python
# -*- coding:utf-8 -*-
import os
import time
import shutil
import sys
import os.path
import hashlib

minigame_dir = "minigame"
build_dir = os.path.dirname(os.path.dirname(os.path.abspath(sys.argv[0])))

if not os.path.exists(build_dir):
    print(build_dir, 'not exists')
    sys.exit(0)

if not os.path.exists('%s/runtime'%build_dir):
    print(build_dir, 'not a build dir')
    sys.exit(0)


logfile = open('deploy.log', 'a+')
def log(str):
    print(str)
    logfile.write(str+'\n')

def copy_file(src, dst):
    shutil.copy(src, dst)
    
def copy_dir(src, dst):
    if os.path.exists(dst):
        shutil.rmtree(dst)
    shutil.copytree(src, dst)

def mkdir(dir):
    if os.path.exists(dir):
        return
    os.mkdir(dir)

def rmdir(dir):
    shutil.rmtree(dir)

def getFileHash(filePath):
    f = open(filePath, 'r')
    md5_hash = hashlib.md5()
    while True:
        content = f.read(1024)
        if not content:
            break
        md5_hash.update(content)
    f.close()
    return md5_hash.hexdigest()

def compareFile(srcFile, dstFile):
    srcHash = getFileHash(srcFile)
    dstHash = getFileHash(dstFile)
    return srcHash == dstHash

def updateOrAddFile(filepath):
    srcFile = filepath
    dstFile = minigame_dir+fullpath.replace(build_dir, '')
    if not os.path.exists(dstFile):
        if os.path.isdir(srcFile):
            log('add dir:\t%s'%dstFile)
            mkdir(dstFile)
        else:
            log('add file:\t%s'%dstFile)
            copy_file(srcFile, dstFile)
    else:
        if not os.path.isdir(srcFile):
            if not compareFile(srcFile, dstFile):
                log('update file:\t%s'%dstFile)
                copy_file(srcFile, dstFile)

    #print(srcFile)
    #print(dstFile)

log('\ndate:\t%s'%time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(time.time())))
log('build_dir:\t%s'%build_dir)

total = 0
for dirpath, dirnames, filenames in os.walk(build_dir):
    for file in dirnames:
        total = total + 1
    for file in filenames:
        total = total + 1

progress = 0
for dirpath, dirnames, filenames in os.walk(build_dir):
    for file in dirnames:
        fullpath = os.path.join(dirpath, file)
        updateOrAddFile(fullpath)
        progress = progress + 1
    for file in filenames:
        fullpath = os.path.join(dirpath, file)
        updateOrAddFile(fullpath)
        progress = progress + 1

#mkdir(minigame_dir)
#mkdir("%s/log"%minigame_dir)
#mkdir("%s/tlog"%minigame_dir)
#mkdir("%s/tlogbak"%minigame_dir)
#mkdir("%s/tlogsync"%minigame_dir)   
#copy_dir('%s/bin'%build_dir, '%s/bin'%minigame_dir)
#copy_dir('%s/data'%build_dir, '%s/data'%minigame_dir)
#copy_dir('%s/supervisord'%build_dir, '%s/supervisord'%minigame_dir)
#copy_dir('%s/script'%build_dir, '%s/script'%minigame_dir)
#copy_dir('%s/runtime'%build_dir, '%s/runtime'%minigame_dir)
