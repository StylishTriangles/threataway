from utils import *
from bs4 import BeautifulSoup
import requests
import socket

def talos_query(endpoint, domain, query_params):
  search_string = domain
  base_url = 'https://talosintelligence.com/'
  referer = '%s/reputation_center/lookup?search=%s'%(base_url,
  search_string)
  user_agent = 'Mozilla/5.0 (Windows NT 6.1; WOW64; rv:64.0) Gecko/20100101 Firefox/64.0'
  default_headers={
    'Referer':referer,
    'User-Agent':user_agent
  }
 
  full_url=base_url+endpoint
  result = requests.get(full_url,
    headers={
      'Referer':referer,
      'User-Agent':user_agent
    },
    params = query_params
   ).json()
  return result

def talos_query_blacklist(query_string):
  return talos_query('sb_api/blacklist_lookup', query_string,
    {'query_type':'ipaddr', 'query_entry':query_string})

def talos_query_details(query_string):
  return talos_query('sb_api/query_lookup', query_string,
    {'query':'/api/v2/details/ip/',
    'query_entry':query_string})
 
def talos_query_wscore(query_string):
  return talos_query('sb_api/remote_lookup', query_string,
    {'hostname':'SDS', 'query_string':'/score/wbrs/json?url=%s' % query_string} 
  )

def talos(list_id_domain, con):
  for domain_id,domain in list_id_domain:
    ipaddr=socket.gethostbyname(domain)
    blacklist = talos_query_blacklist(ipaddr)
    details = talos_query_details(ipaddr)
    wscore = talos_query_wscore(ipaddr)
    res={}
    try:
      if blacklist['entry']['expiration'] == 'NEVER':
        res["malicious"]=0
        res["malicious_type"]=""
      else:
        res["malicious"]=1
        res["malicious_type"] = '|'.join(blacklist['entry']['classifications'])
    except:
      pass
    
    details_list=['hostname', 'monthly_spam_level', 'organization', 'web_score_name', 'email_score_name']
    for key in details_list:
      try:
        res[key] = details[key] 
      except:
        pass
    res = {k: v for k, v in res.items() if v is not None}
    update_query = 'UPDATE urls SET {} WHERE idUrl=%s'.format(', '.join('{}=%s'.format(k) for k in res))
    cur = con.cursor()
    cur.execute(update_query, tuple(list(res.values())) + (domain_id,))
    con.commit()
