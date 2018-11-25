from shodan import Shodan
import socket
import requests
import base64

def query_honeypot_score(ipaddr, key):
  url = "https://api.shodan.io/labs/honeyscore/"+ipaddr+"?key="+key
  return requests.get(url).text

def query_shodan(list_id_domain, con):
  with open('../secret-dir/shodan.key', 'r') as f:
    key=f.read().replace('\n', '')
  api = Shodan(key)

  for domain_id,domain in list_id_domain:
    ipaddr=socket.gethostbyname(domain)
    res={}
    honeypot_score = query_honeypot_score(ipaddr, key)
    res["honeypot_score"] = honeypot_score
    
    query_base="https://www.shodan.io/search?query="
    query_malware='category:malware ip:' + ipaddr 
    result_malware = api.search(query_malware)
    try:
      if len(result_malware["matches"]) == 0:
        res["shodan_malware"]=0
        res["shodan_malware_query"]=""
      else:
        res["shodan_malware"]=1
        res["shodan_malware_query"]=base64.b64encode((query_base+query_malware).encode('ascii'))
    except:
      pass
    #print(result_malware)
    

    query_default_creds='"default password" ip:' + ipaddr
    result_creds = api.search(query_default_creds)
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
    
    print(res)
    print(update_query)

    cur = con.cursor()
    cur.execute(update_query, tuple(list(res.values())) + (domain_id,))
    con.commit() 
