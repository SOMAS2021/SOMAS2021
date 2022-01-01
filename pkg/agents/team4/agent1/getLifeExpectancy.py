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

ticks_per_floor = 10
ticks_per_day = int(number_of_agents) * ticks_per_floor

# create dict of agents with all their coressponding log entries
agents = {}
for i in logs:
    try:
        if i['agent_id'] not in agents and i['msg'] == "team4EvoAgent reporting status:":
            agents[i['agent_id']] = []
        if i['msg'] == "team4EvoAgent reporting status:":
            agents[i['agent_id']].append(i)
    except:
        continue

# get agentid we want to create a population with
agent = agents[list(agents.keys())[0]]

# get values for the agent
values = (agent[0]["FoodToEat"], agent[0]["DaysToWait"])

# get avg number of days lived by the different instantiated agents (all have the same coeffs)
counts = []
for _, value in agents.items():
    counts.append(len(value)//ticks_per_day)

avg_days_lived = sum(counts)/len(counts)
print(avg_days_lived)

# read best_agent file
best_agent_json_file = open(best_agents_file_name)
best_agent_json = json.load(best_agent_json_file)
best_agent_json_file.close()

try:
    # read agent at current_iteration+1
    next_agent = best_agent_json[int(current_iteration) + 1]
    # print("Changing agent config to: {0}".format(next_agent))

    # pass agent to agent_config file to create population for next run
    agent_config_file = open(agent_config_file_name, 'w')
    agent_config_file.write(json.dumps(next_agent, indent=4))
    agent_config_file.close()
except:
    pass

#