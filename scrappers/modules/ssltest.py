from utils import *
from bs4 import BeautifulSoup

def fetchSSLLabs(domain):
  src = get("https://www.ssllabs.com/ssltest/analyze.html?d="+domain)
  if not "rating_g" in src:
    return False
  soup = BeautifulSoup(src, 'html.parser')
  grade = soup.find(class_='rating_g ').contents[0].strip()
  return grade 

def ssltest(domain):
  grade = fetchSSLLabs(domain)
  print("grade=|"+grade+"|")
  return grade
