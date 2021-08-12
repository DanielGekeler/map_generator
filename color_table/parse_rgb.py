import json

raw_file = "rgb_map.csv"


f = open(raw_file)
lines = f.readlines()

colors = {}

for line in lines:
    spl = line.split(";")
    a = spl[0]
    b = spl[1]

    a = a.split(" ")[0]
    
    b = b.strip()
    b = b.replace(", ", ":")

    colors[a] = b

json_object = json.dumps(colors, indent = 4)
print(json_object)