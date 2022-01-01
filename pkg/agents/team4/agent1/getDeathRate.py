import sys
import json

log_file_name = sys.argv[1]
number_of_agents = sys.argv[2]
agent_config_file_name = sys.argv[3]
best_agents_file_name = sys.argv[4]
current_iteration = sys.argv[5]

# reading contents of logfile
log_file = open(log_file_name, 'r')
lines = log_file.readlines()
log_file.close()

# convert log entries to json objects
logs = []
for i in range(len(lines)):
    logs.append(json.loads(lines[i]))

ourAgentDeaths = 0
totalAgentDeaths = 0
# create dict of agents with all their coressponding log entries
for i in logs:
    try:
        if i['agent_type'] == 'Team4':
            ourAgentDeaths+=1
    except:
        continue

totalAgentDeaths = logs[-1]['cumulativeDeaths']

print(ourAgentDeaths)
