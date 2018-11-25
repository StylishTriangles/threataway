from shodan import Shodan

import socket
import requests
import base64
import time

from calc_score import update_score

def query_honeypot_score(ipaddr, key):
  url = "https://api.shodan.io/labs/honeyscore/"+ipaddr+"?key="+key
  return requests.get(url).text

def query_shodan(list_id_domain, con, lazy_rating=1):
  with open('../secret-dir/shodan.key', 'r') as f:
    key=f.read().replace('\n', '')
  api = Shodan(key)

  for domain_id,domain in list_id_domain:
    print("domain: " + domain)
    ipaddr=socket.gethostbyname(domain)
    res={}
    for i in range(0, 5):
      try:
        honeypot_score = query_honeypot_score(ipaddr, key)
        break
      except:
        time.sleep(1)
    res["honeypot_score"] = honeypot_score
    
    query_base="https://www.shodan.io/search?query="
    query_malware='category:malware ip:' + ipaddr 
    for i in range(0, 5):
      try:
        result_malware = api.search(query_malware)
        break
      except:
        time.sleep(1)
    
    try:
      if len(result_malware["matches"]) == 0:
        res["shodan_malware"]=0
        res["shodan_malware_query"]=""
      else:
        res["shodan_malware"]=1
        res["shodan_malware_query"]=base64.b64encode((query_base+query_malware).encode('ascii'))
    except:
      pass

    query_default_creds='"default password" ip:' + ipaddr
    for i in range(0, 5):
      try:
        result_creds = api.search(query_default_creds)
        break
      except:
        time.sleep(1)
     
    try:
      if len(result_creds["matches"]) == 0:
        res["shodan_creds"]=0
        res["shodan_creds_query"]=""
      else:
        res["shodan_creds"]=1
        res["shodan_creds_query"]=base64.b64encode((query_base+query_default_creds).encode('ascii'))
    except:
      pass
    res = {k: v for k, v in res.items() if v is not None}
    update_query = 'UPDATE urls SET {} WHERE idUrl=%s'.format(', '.join('{}=%s'.format(k) for k in res))
    
    cur = con.cursor()
    changed = cur.execute(update_query, tuple(list(res.values())) + (domain_id,))
    con.commit()
    if lazy_rating == 0 or changed != 0:
      update_score(con, domain_id)
