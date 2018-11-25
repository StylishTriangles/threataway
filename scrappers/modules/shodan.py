from shodan import Shodan

def query_shodan(domains, con):
  api = ''
  ipinfo = api.scan('8.8.8.8')
  print(ipinfo)
  #for banner in api.search_cursor('http.title:"hacked by"'):
  #  print(banner)
  ics_services = api.count('tag:ics')
  print('Industrial Control Systems: {}'.format(ics_services['total']))

