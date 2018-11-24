import sys
import json
import pymysql

from modules.ssltest import *
from modules.talos import *

def open_db():
  with open('../config/dbconfig.json') as cfg:
    data = json.load(cfg)
  host=data["Hostname"]
  username=data["Username"]
  pw=data["Password"]
  db=data["Name"]
  try:
    con = pymysql.connect(host, username, pw, db)
  except:
    sys.exit(-3)
  return con

def get_domains_from_db(con):
  with con:
    cur = con.cursor()
    cur.execute("SELECT idUrl, domain FROM urls")
    rows = list(cur.fetchall())
    return rows

def main():
  if len(sys.argv) != 2:
    print("main.py service")
    sys.exit(-1)
  con = open_db()
  service = sys.argv[1]
  print("service: " + service)
  domains = get_domains_from_db(con)
  if service == "ssllabs":
    ssl = ssltest(domains)
    print("ssl: " + ssl)
  elif service == "talos":
    talos(domains, con)  
    #print(talos_report)
  else:
    print("unknown service")
    sys.exit(-2)
  con.close()
     
if __name__ == "__main__":
  main()
