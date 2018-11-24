from utils import *
from bs4 import BeautifulSoup
import requests


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

def talos(domain):
  ipaddr='31.13.81.9'
  #ipaddr=domain
  #ipaddr='31.185.104.19'
  blacklist = talos_query_blacklist(ipaddr)
  print(blacklist)
  print("********")
  details = talos_query_details(ipaddr)
  print(details)
  print("********")
  wscore = talos_query_wscore(ipaddr)
  print(wscore)
  
