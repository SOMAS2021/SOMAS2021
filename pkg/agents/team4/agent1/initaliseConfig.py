import json
import sys
import random

configFileName = sys.argv[1]
configFile = open(configFileName, 'r+')

jsonVar = {
    "Floor": [[0.0, 0.0, 0.0], [0.0, 0.0, 0.0], [0.0, 0.0, 0.0]],
    "Hp": [[0.0, 0.0, 0.0], [0.0, 0.0, 0.0], [0.0, 0.0, 0.0]],
}


for k, v in jsonVar.items():
    random_coefs = []
    for i in range(len(jsonVar[k])):
        random_coefs.append([round(random.random(), 4)
                            for i in range(len(jsonVar[k][i]))])
    jsonVar[k] = random_coefs

jsonVarStr = json.dumps(jsonVar, indent=4)
configFile.write(jsonVarStr)
configFile.close()
