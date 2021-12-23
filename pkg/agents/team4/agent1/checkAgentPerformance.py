import sys
import json

logfile = sys.argv[1]
number_of_agents = sys.argv[2]
config_file_name = sys.argv[3]

print("loading logfile: {0}".format(logfile))

file = open(logfile, 'r')
lines = file.readlines()

logs = []
for i in range(len(lines)):
    logs.append(json.loads(lines[i]))

ticks_per_floor = 10
ticks_per_day = int(number_of_agents) * ticks_per_floor

agents = {}

for i in logs:
    try:
        if i['agent_id'] not in agents:
            agents[i['agent_id']] = []
            agents[i['agent_id']].append(i)
        else:
            agents[i['agent_id']].append(i)
    except:
        continue

counts = []
coefficients = []
for key, value in agents.items():
    counts.append(((key, value[8]["currentFloorScore"], value[8]
                  ["currentHpScore"]), len(value)//ticks_per_day))

sorted_counts = sorted(counts, key=lambda x: x[1], reverse=True)
print(sorted_counts[:3])

json_var = {
    "Floor": [
        [
            sorted_counts[0][0][1][0],
            sorted_counts[1][0][1][0],
            sorted_counts[2][0][1][0],
        ],
        [
            sorted_counts[0][0][1][1],
            sorted_counts[1][0][1][1],
            sorted_counts[2][0][1][1],
        ],
        [
            sorted_counts[0][0][1][2],
            sorted_counts[1][0][1][2],
            sorted_counts[2][0][1][2],
        ],
    ],
    "Hp": [
        [
            sorted_counts[0][0][2][0],
            sorted_counts[1][0][2][0],
            sorted_counts[2][0][2][0],
        ],
        [
            sorted_counts[0][0][2][1],
            sorted_counts[1][0][2][1],
            sorted_counts[2][0][2][1],
        ],
        [
            sorted_counts[0][0][2][2],
            sorted_counts[1][0][2][2],
            sorted_counts[2][0][2][2],
        ],
    ],
}


config_file = open(config_file_name, 'w')
config_file.write(json.dumps(json_var, indent=4))
config_file.close()
