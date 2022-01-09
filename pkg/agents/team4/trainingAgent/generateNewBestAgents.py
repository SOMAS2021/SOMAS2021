from datetime import datetime
import sys
import random
import json
import math
import numpy as np

# input files
best_agents_file_name = sys.argv[1]
agent_life_expectancies_file_name = sys.argv[2]
agent_death_rate_file_name = sys.argv[3]
agent_our_life_expectancies_file_name = sys.argv[4]
agent_other_life_expectancies_file_name = sys.argv[5]

# input flags and numerical parameters 
number_of_health_levels = int(sys.argv[6])
selfish_flag = sys.argv[7]
selfless_flag = sys.argv[8]
current_iteration = int(sys.argv[9])

def generate_mutated_agent(agent, attribute, population_size): # We define a mutated angent from a select base agent in the previous iteration.
    
    noise_multiplier_for_attribute = {
        "FoodToEat": math.sqrt(population_size)*5/math.sqrt(current_iteration), # Value will tend to 10
        "WaitProbability": math.sqrt(population_size)*10/math.sqrt(current_iteration)
    }

    mu, sigma = 0, 1  # mean and standard deviation

    for i in attribute:
        for j in range(number_of_health_levels):
            rand = int(np.random.normal(mu, noise_multiplier_for_attribute[i]*sigma))
            agent[i][j] += rand
            
            if ( agent[i][j] > 60):
                agent[i][j] = 60
            if (agent[i][j] < 0):
                agent[i][j] = 0        
    return agent

def generate_mixed_agent(agent_list, number_of_agents_to_use, attribute): # We define a mixed agent from parameters of the previous best agents

    picked_best_agents = agent_list[:number_of_agents_to_use]

    bands = []
    for i in range(number_of_agents_to_use): # determines probability of each agent within the number of agents to use
        bands.append(((1/number_of_agents_to_use) * (1/(i+1))))
        
    sumOf = sum(bands)
    for i in range(len(bands)): # sets bands in range totalling 1
        bands[i] *= (1/sumOf)

    for i in range(1,len(bands)): # makes bands cumulative
        bands[i] += bands[i-1]

    agent_mix = {}
    for i in attribute:
        values = []
        for j in range(number_of_health_levels):
            temp_val = 0.0
            rand = random.random()
            for k in range(number_of_agents_to_use):
                if rand < bands[k]:
                    temp_val = picked_best_agents[k][i][j]
                    break
            values.append(temp_val)

        agent_mix[i] = values
    return agent_mix


def generate_random_agent(): # We define an entirely new random agent
    random_mutation_agent = {
    "FoodToEat": [],
    "WaitProbability": [],
    }

    random_mutation_agent["FoodToEat"] = [random.randint(0, 60)
                                        for _ in range(int(number_of_health_levels))]
    random_mutation_agent["WaitProbability"] = [random.randint(0, 100)
                                        for _ in range(int(number_of_health_levels))]
                                        
    return random_mutation_agent

# read bestAgent.config
best_agents_file = open(best_agents_file_name)
best_agents_list = json.load(best_agents_file)
best_agents_file.close()

# read the array of average time survived by each agent
agent_life_expectancies_file = open(agent_life_expectancies_file_name)
performance_list_life = json.load(agent_life_expectancies_file)
agent_life_expectancies_file.close()

# read the array of average time survived by each agent of our type
agent_our_life_expectancies_file = open(agent_our_life_expectancies_file_name)
performance_list_our_life = json.load(agent_our_life_expectancies_file)
agent_our_life_expectancies_file.close()

agent_other_life_expectancies_file = open(agent_other_life_expectancies_file_name)
performance_list_other_life = json.load(agent_other_life_expectancies_file)
agent_other_life_expectancies_file.close()

# read the array of average death by each agent
agent_death_rate_file = open(agent_death_rate_file_name)
performance_list_death = json.load(agent_death_rate_file)
agent_death_rate_file.close()

# sort array but keep a track of index to obtain best agent
performance_list_indices = range(len(performance_list_life))
    
if selfless_flag == "True":
    print("This is selfless") # TODO: this to change
    zipped_performance = zip(performance_list_other_life, performance_list_our_life,  performance_list_life,
                             performance_list_death, performance_list_indices)
    sorted_zipped_performance = sorted(zipped_performance, reverse=True)
    tuples = zip(*sorted_zipped_performance)
    sorted_performance_list_other_life, sorted_performance_list_our_life, sorted_performance_list_life, sorted_performance_list_death, sorted_performance_list_indices = [
        list(tuple) for tuple in tuples]
elif selfish_flag=="True":
    print("This is selfish")
    zipped_performance = zip(performance_list_our_life, performance_list_other_life, performance_list_life,
                             performance_list_death, performance_list_indices)
    sorted_zipped_performance = sorted(zipped_performance, reverse=True)
    tuples = zip(*sorted_zipped_performance)
    sorted_performance_list_our_life, sorted_performance_list_other_life, sorted_performance_list_life, sorted_performance_list_death, sorted_performance_list_indices = [
        list(tuple) for tuple in tuples]
else:
    print("This is neutral")
    zipped_performance = zip(performance_list_life, performance_list_other_life, performance_list_our_life,
                             performance_list_death, performance_list_indices)
    sorted_zipped_performance = sorted(zipped_performance, reverse=True)
    tuples = zip(*sorted_zipped_performance)
    sorted_performance_list_life, sorted_performance_list_other_life, sorted_performance_list_our_life, sorted_performance_list_death, sorted_performance_list_indices = [
        list(tuple) for tuple in tuples]
print("The current global life expectancies are: {0}\nThe current our life expectancies are: {1} \nThe current other life expectancies are: {2}".format(
    sorted_performance_list_life, sorted_performance_list_our_life, sorted_performance_list_other_life))

# get all best agents
agent_list = []
for i in range(len(sorted_performance_list_indices)):
    agent_list.append(best_agents_list[sorted_performance_list_indices[i]])

number_of_agents_to_use = int((2/5)*len(agent_list))
attribute_arr = ["FoodToEat", "WaitProbability"]

################################################################################################
agents_of_time = agent_list
for i in range(len(agents_of_time)):
    agents_of_time[i]['life'] = performance_list_life[sorted_performance_list_indices[i]]
    agents_of_time[i]['death'] = performance_list_death[sorted_performance_list_indices[i]]
    agents_of_time[i]['our_life'] = performance_list_our_life[sorted_performance_list_indices[i]]
    agents_of_time[i]['other_life'] = performance_list_other_life[sorted_performance_list_indices[i]]

timestamp = datetime.now(tz=None)
file_name_out = "pkg/agents/team4/trainingAgent/"+"{0}/{1}_{2}_{3}.json".format(
    "configs/storedConfigs", timestamp, performance_list_life[0], performance_list_death[0])
file_name_out = file_name_out.replace(" ", "_")
old_agents_file = open(file_name_out, 'w')
old_agents_file.write(json.dumps(agents_of_time, indent=4))
old_agents_file.close()
################################################################################################

new_best_agents = []

for i in range(len(agent_list)):
    number_of_best_agent = int((2/10) * len(agent_list))
    number_of_mixed_agents = int((4/10) * len(agent_list))
    number_of_mutated_agents = int((2/10) * len(agent_list))
    number_of_random_agents = len(agent_list) - number_of_best_agent - number_of_mixed_agents - number_of_mutated_agents

    for j in range(number_of_best_agent):
        new_best_agents.append(agent_list[j])

    for j in range(number_of_mixed_agents):
        new_best_agents.append(generate_mixed_agent(agent_list, number_of_agents_to_use, attribute_arr))

    for j in range(number_of_mutated_agents):
        new_best_agents.append(generate_mutated_agent(agent_list[j%3], attribute_arr, len(agent_list)))

    for j in range(number_of_random_agents):
        new_best_agents.append(generate_random_agent())

for best_agent in new_best_agents:
    best_agent["FoodToEat"][0] = -1
    best_agent["WaitProbability"][0] = -1

best_agents_file = open(best_agents_file_name, 'w')
best_agents_file.write(json.dumps(new_best_agents, indent=4))
best_agents_file.close()
