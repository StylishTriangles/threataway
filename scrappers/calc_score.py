import pymysql

def calc_score(rows):
  print(rows)
  malicious = int(rows[1])
  web_score_name = rows[3].lower()
  email_score_name = rows[4].lower()
  monthly_spam = int(rows[5])
  honeypot_score = float(rows[6])
  shodan_malware = int(rows[7])
  shodan_creds = int(rows[8])
 
  print("malicious: " + str(malicious))
  print("web_score_name: " + web_score_name)
  print("email_score_name: " + email_score_name)
  print("monthly_spam: " + str(monthly_spam))
  print("honeypot_score: " + str(honeypot_score))
  print("shodan_malware: " + str(shodan_malware))
  print("shodan_creds: " + str(shodan_creds))
 
  if shodan_malware == 1 or malicious == 1:
    return 0.0

  web_coeff=5.0
  if web_score_name == "poor":
    web_coeff = 1.0
  elif web_score_name == "neutral":
    web_coeff = 4.0
  
  email_coeff = 2.0
  if email_score_name == "poor":
    email_coeff=1.0
  elif email_score_name == "neutral":
    email_coeff=1.8

  res = (1.0 - honeypot_score/2.0)*(1.0 - shodan_creds/2.0)*(1.0 - (monthly_spam != 0)/5.0) * web_coeff * email_coeff

  return res
  

def update_score(con, _id):
  print("score_calc")
  print(_id)
  cur = con.cursor()
  query = "SELECT rating, malicious, malicious_type, web_score_name, email_score_name, monthly_spam_level, honeypot_score, shodan_malware, shodan_creds  FROM urls WHERE idUrl=%s"
  cur.execute(query, _id)
  rows = cur.fetchall()[0]
  print(rows)
  new_score = calc_score(rows)
  print(rows[0])
  print(new_score)

  cur = con.cursor()
  update_query='UPDATE urls SET rating=%s  WHERE idUrl=%s'
  changed = cur.execute(update_query, (new_score, _id))
  con.commit()
