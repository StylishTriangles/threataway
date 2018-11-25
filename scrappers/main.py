import sys
import json
import pymysql

#from modules.ssltest import *
from modules.talos import *
#from modules.spamlists import *
from modules.shodan import *

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
  if len(sys.argv) < 2 or len(sys.argv) > 3:
    print("main.py service")
    sys.exit(-1)
  con = open_db()
  service = sys.argv[1]
  print("service: " + service)
  if len(sys.argv) == 2:
    domains = get_domains_from_db(con)
  elif len(sys.argv) == 3:
    domains = sys.argv[2]
    _id = int(sys.argv[1])
  print(domains)
  if len(sys.argv) == 2:
    if service == "talos":
      talos(domains, con)  
    elif service == "shodan":
      query_shodan(domains, con)
    else:
      print("unknown service")
      sys.exit(-2)
  else:
    talos([(_id, domains)], con)
    #query_shodan(domains, con)
  con.close()
     
if __name__ == "__main__":
  main()
