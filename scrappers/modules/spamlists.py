from utils import *
from bs4 import BeautifulSoup
import requests
import socket

def check_spam_lists(list_id_domain, con):
  url="https://www.ultratools.com/tools/spamDBLookupResult"

  for domain_id,domain in list_id_domain:
    ipaddr=socket.gethostbyname(domain)
    r = requests.post(url, data={'domainName':ipaddr})
    if "resultstable" in r.text:
      print("yo")
    else:
      print("ehhh")
    soup = BeautifulSoup(r.text, 'html.parser')
    v = soup.find(class_='resultstable')
    for x in v.contents[3].contents:
      if not "<td>" in x:
        continue
      print(x)
      t = x.split('</td>')[0]
      
      print("********")
      print(t)
