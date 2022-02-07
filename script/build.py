#!/usr/bin/python
# -*- coding:utf-8 -*-
import os
import time
import shutil
import sys
import os.path
import json
import codecs
import re
import argparse

def copy_file(src, dst):
    shutil.copyfile(src, dst)
    
def copy_dir(src, dst):
    shutil.copytree(src, dst)

def read_file(src):
    f = open(src, 'r')
    s = f.read()
    f.close()
    return s

def write_file(dst, s):
    f = open(dst, 'w')
    s = f.write(s)
    f.close()

def mkdir(dir):
    os.mkdir(dir)

class Script:
    def init_env(self, env):
        self.env = env
        self.project_dir = os.path.dirname(os.path.dirname(os.path.abspath(sys.argv[0])))
        print(self.project_dir)
        localtime = time.localtime(time.time())
        self.build_dir = 'build/%s-%04d%02d%02d-%02d%02d%02d'%(env, localtime.tm_year, localtime.tm_mon, localtime.tm_mday, localtime.tm_hour, localtime.tm_min, localtime.tm_sec)
        self.runtime_dir = '%s/runtime'%self.build_dir
        self.script_dir = '%s/script'%self.build_dir
        self.bin_dir = '%s/bin'%self.build_dir
        self.data_dir = '%s/data'%self.build_dir
        self.tlogxml_dir = '%s/tlog.xml'%self.data_dir
        self.config_dir = '%s/config'%self.data_dir
        self.crt_dir = '%s/crt'%self.data_dir
        self.sql_dir = '%s/sql'%self.data_dir
        self.update_sql_path = '%s/update.sql'%self.sql_dir
        self.log_dir = '%s/log'%self.build_dir
        self.tlog_dir = '%s/tlog'%self.build_dir
        self.supervisord_dir = '%s/supervisord'%self.build_dir
        self.template_dir = "%s/build-templates/%s"%(self.project_dir, self.env)
        self.template_supervisord_dir = "%s/supervisord"%self.template_dir
        self.template_gamesvr_config_path = "%s/gamesvr/config.ini"%self.template_dir
        self.template_coord_config_path = "%s/coord/config.ini"%self.template_dir
        self.template_update_sql_path = '%s/update.sql'%self.template_dir        
        self.server_host = "193.112.209.40"
        if env == "minigame-prod":
            self.server_host = "129.204.136.126"

    def command_list(self, args):
        self.init_env(args.env)
        f = codecs.open('%s/deploy.conf'%self.template_dir)
        deploy_conf = f.read()
        f.close()
        deploy_conf = json.loads(deploy_conf)
        #print(deploy_conf)
        for host_name, app_conf_arr in deploy_conf.items(): 
            for app_conf in app_conf_arr:
                print('%s#%s'%(host_name, app_conf['name']))

    def command_pack(self, args):
        #print(args.env)
        #print(args.proc)
        pack_proc_dict = {}
        for x in args.proc:
            pack_proc_dict[x] = True
        print(pack_proc_dict, len(pack_proc_dict))
        #return
        self.init_env(args.env[0])
        os.chdir(self.project_dir)
        os.mkdir(self.build_dir)
        print(self.build_dir)
        os.mkdir(self.log_dir)
        os.mkdir(self.tlog_dir)
        #构建bin目录
        copy_dir('%s/bin'%self.project_dir, self.bin_dir)
        os.system('chmod +x %s/bin/*'%self.project_dir)
        #构建script目录
        os.mkdir(self.script_dir)
        copy_file('%s/script/deploy.py'%self.project_dir, '%s/deploy.py'%self.script_dir)
        copy_file('%s/script/export_rank.py'%self.project_dir, '%s/export_rank.py'%self.script_dir)
        os.system('chmod +x %s/export_rank.py'%self.script_dir)
        #构建supervisord目录
        self.make_supervisord(pack_proc_dict)
        #构建data目录
        os.mkdir(self.data_dir)
        copy_dir('%s/data/config'%self.project_dir, self.config_dir)
        copy_dir('%s/data/sql'%self.project_dir, self.sql_dir)
        copy_dir('%s/data/crt'%self.project_dir, self.crt_dir)
        copy_file('%s/data/tlog.xml'%self.project_dir, self.tlogxml_dir)
        copy_file(self.template_update_sql_path, self.update_sql_path)
        os.system('tar czf %s.tar.gz %s'%(self.build_dir, self.build_dir))
        os.system("scp %s.tar.gz root@%s:~"%(self.build_dir, self.server_host))

    def make_supervisord(self, deploy_proc_list): 
        #创建runtime目录
        os.mkdir(self.runtime_dir)
        #创建supervisord目录
        os.mkdir(self.supervisord_dir)
        #copy_file('%s/supervisord.conf'%template_supervisord_dir, '%s/supervisord.conf'%supervisord_dir)
        f = codecs.open('%s/deploy.conf'%self.template_dir)
        deploy_conf = f.read()
        f.close()
        deploy_conf = json.loads(deploy_conf)
        #print(deploy_conf)
        for host_name, app_conf_arr in deploy_conf.items(): 
            #重命名和拷贝supervisord.conf
            s = read_file('%s/supervisord.conf'%(self.template_supervisord_dir))
            s = s.replace('files = supervisord/*.ini', 'files = %s/*.ini'%host_name)
            write_file('%s/%s.conf'%(self.supervisord_dir, host_name), s)
            #创建host目录
            host_dir = '%s/%s'%(self.supervisord_dir, host_name)
            os.mkdir(host_dir)
            for app_conf in app_conf_arr:
                if len(deploy_proc_list) != 0 and not deploy_proc_list.has_key('%s#%s'%(host_name, app_conf['name'])):
                    continue
                print('make %s'%app_conf['name'])
                #修改和拷贝supervisord.ini
                s = read_file('%s/supervisord/%s.ini'%(self.template_supervisord_dir, app_conf['bin']))
                matches = re.findall(r'\{(\w+)\}', s)
                for m in matches:
                    s = s.replace('{'+m+'}', app_conf[m])
                write_file('%s/%s.ini'%(host_dir, app_conf['name']), s)
                #创建proc目录
                copy_dir('%s/runtime/%s'%(self.template_dir, app_conf['bin']), '%s/%s'%(self.runtime_dir, app_conf['name']))
                #修改和拷贝config.ini
                s = read_file('%s/runtime/%s/config.ini'%(self.template_dir, app_conf['bin']))
                matches = re.findall(r'\{(\w+)\}', s)
                for m in matches:
                    s = s.replace('{'+m+'}', app_conf[m])
                write_file('%s/%s/config.ini'%(self.runtime_dir, app_conf['name']), s)
                #拷贝bin
                copy_file('%s/%s'%(self.bin_dir, app_conf['bin']), '%s/%s/%s'%(self.runtime_dir, app_conf['name'], app_conf['name']))
                os.system('chmod +x %s/%s/%s'%(self.runtime_dir, app_conf['name'], app_conf['name']))
        #os.mkdir('%s/supervisord'%supervisord_dir)
        #copy_file('%s/supervisord/tlogsync.ini'%template_supervisord_dir, '%s/supervisord/tlogsync.ini'%supervisord_dir)

    def command_upload(self, args):
        pass

    def command_hotfix(self, args):
        pass

script = Script()
parser = argparse.ArgumentParser(description='Build program')
#parser.add_argument('--env', '-e', dest='env', required=True)
subparsers = parser.add_subparsers(title='subcommands', description='description', help='sub-command help')
parser_list = subparsers.add_parser('list', help='list all proc')
parser_list.add_argument('env')
parser_list.set_defaults(func=script.command_list)

parser_pack = subparsers.add_parser('pack', help='pack some proc')
parser_pack.add_argument('env', nargs=1, help='minigame-test2')
parser_pack.add_argument('proc', nargs="*", help='gamesvr_1 coord')
parser_pack.set_defaults(func=script.command_pack)

parser_pack = subparsers.add_parser('upload', help='upload pack to server')
parser_pack.add_argument('file', nargs=1, help='xxxx')
parser_pack.set_defaults(func=script.command_upload)

parser_pack = subparsers.add_parser('hotfix', help='hotfix')
parser_pack.add_argument('file', nargs=1, help='xxxx')
parser_pack.set_defaults(func=script.command_hotfix)

args = parser.parse_args()
args.func(args)

#sys.exit(0)



