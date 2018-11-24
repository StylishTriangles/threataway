# this script acquires passwords of length at least 8
# sorts and saves them in common.txt

fi = open("100k_common.txt", "r")
fo = open("common.txt", "w")

s = []

for line in fi:
    line = line.strip()
    if len(line) >= 8:
        s.append(line)

s.sort()

for line in s:
    fo.write(line + '\n')
