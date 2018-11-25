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

def merge(list1, list2):
  return list(set(list1+list2))

def set_dirty(ids, con):
  for id_ in ids:
    update_query = 'UPDATE listlists SET dirty = 1 WHERE idURL = %s'
    cur = con.cursor()
    changed = cur.execute(update_query, (id_,))
    con.commit()

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
    domains = [(int(sys.argv[1]), sys.argv[2])]
  ids=[]
  if len(sys.argv) == 2:
    if service == "talos":
      ids=merge(ids, talos(domains, con, 1))
    elif service == "shodan":
      ids=merge(ids, query_shodan(domains, con, 1))
    elif service == "all":
      ids=merge(ids, query_shodan(domains, con, 1))
      ids=merge(ids, talos(domains, con, 1))
    elif service == "recalc":
      ids=merge(ids, query_shodan(domains, con, 0))
      ids=merge(ids, talos(domains, con, 0))
    else:
      print("unknown service")
      sys.exit(-2)
  else:
    ids=merge(ids,talos(domains, con, 1))
    ids=merge(ids,query_shodan(domains, con, 1))
  print("ids:")
  print(ids)
  set_dirty(ids, con)
  con.close()
     
if __name__ == "__main__":
  main()
