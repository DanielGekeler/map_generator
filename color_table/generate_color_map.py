import json

raw_file = "map_colors_1-17-1.csv"


f = open(raw_file)
lines = f.readlines()

colors = {}

for line in lines:
    spl = line.split(";")
    a = spl[0]
    b = spl[1]

    if "[*]" in a:
        a = a.replace("[*]", "")
    elif "Block{" in a:
        a = a.replace("Block{", "")
        a = a.split("}[")[0]
    
    b = int(b.strip())
    colors[a] = b

json_object = json.dumps(colors, indent = 4)
print(json_object)