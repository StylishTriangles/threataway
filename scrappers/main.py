import sys

from modules.ssltest import *
from modules.talos import *

def main():
  if len(sys.argv) != 3:
    print("main.py domain service")
    sys.exit(-1)
  domain = sys.argv[1]
  service = sys.argv[2]
  print(domain + " " + service)
  if service == "ssllabs":
    ssl = ssltest(domain)
    print("ssl: " + ssl)
  elif service == "talos":
    talos(domain)  
    #print("talos:" + talos_report)
  else:
    print("unknown service")
    sys.exit(-2)

     
if __name__ == "__main__":
  main()
