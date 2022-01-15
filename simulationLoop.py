import csv
from encodings import utf_8
import os

NUM_ITERS = 10

DEATH_FILE = "csvFiles/deaths.csv"
SOCIAL_MOTIVES_FILE = "csvFiles/socialMotives.csv"
SOCIAL_MOTIVES_CHANGE_FILE = "csvFiles/socialMotivesChange.csv"
UTILITY_FILE = "csvFiles/utility.csv"

TOTAL_DEATH_FILE = "csvFiles/avgDeaths.csv"
TOTAL_SOCIAL_MOTIVES_FILE = "csvFiles/avgSocialMotives.csv"
TOTAL_SOCIAL_MOTIVES_CHANGE_FILE = "csvFiles/avgSocialMotivesChange.csv"
TOTAL_UTILITY_FILE = "csvFiles/avgUtility.csv"

# totalDeathFile = open('csvFiles/totalDeaths.csv', mode='w+')
totalSocialMotivesFile = open(TOTAL_SOCIAL_MOTIVES_FILE, mode='w+')
totalSocialMotivesChangeFile = open(TOTAL_SOCIAL_MOTIVES_CHANGE_FILE, mode='w+')
totalUtilityFile = open(TOTAL_UTILITY_FILE, mode='w+')

deathList = []
smList = []
smChangeList = []
utilityList = []

def createTotalFile(mode):
    global deathList, smList, smChangeList, utilityList
    fileName = ""
    totalFileName = ""
    ls = []

    if mode == "death":
        fileName = DEATH_FILE
        totalFileName = TOTAL_DEATH_FILE
        ls = deathList
    elif mode == "SM":
        fileName = SOCIAL_MOTIVES_FILE
        totalFileName = TOTAL_SOCIAL_MOTIVES_FILE
        ls = smList
    elif mode == "SMChange":
        fileName = SOCIAL_MOTIVES_CHANGE_FILE
        totalFileName = TOTAL_SOCIAL_MOTIVES_CHANGE_FILE
        ls = smChangeList
    elif mode == "utility":
        fileName = UTILITY_FILE
        totalFileName = TOTAL_UTILITY_FILE
        ls = utilityList
    else:
        return
    
    totalFile = open(totalFileName, mode='w+')
    writer = csv.writer(totalFile)
    with open(fileName, "r") as f:
        reader = csv.reader(f)
        for row in reader:
            ls.append(row)
            writer.writerow(row)
    totalFile.close()

    if mode == "death":
        deathList = ls
    elif mode == "SM":
        smList = ls
    elif mode == "SMChange":
        smChangeList = ls
    elif mode == "utility":
        utilityList = ls
    else:
        return

def updateTotalFile(mode):
    global deathList, smList, smChangeList, utilityList
    fileName = ""
    totalFileName = ""
    ls = []

    headers = {
        "death": ["New deaths", "Cumulative deaths"],
        "SM": [
            "Altruist",
            "Collectivist",
            "Selfish",
            "Narcissist"
        ],
        "SMChange": [
            "A2A",
            "A2C",
            "A2S",
            "A2N",
            "C2A",
            "C2C",
            "C2S",
            "C2N",
            "S2A",
            "S2C",
            "S2S",
            "S2N",
            "N2A",
            "N2C",
            "N2S",
            "N2N"
        ],
        "utility": ["Average Utility"]
    }

    if mode == "death":
        fileName = DEATH_FILE
        totalFileName = TOTAL_DEATH_FILE
        ls = deathList
    elif mode == "SM":
        fileName = SOCIAL_MOTIVES_FILE
        totalFileName = TOTAL_SOCIAL_MOTIVES_FILE
        ls = smList
    elif mode == "SMChange":
        fileName = SOCIAL_MOTIVES_CHANGE_FILE
        totalFileName = TOTAL_SOCIAL_MOTIVES_CHANGE_FILE
        ls = smChangeList
    elif mode == "utility":
        fileName = UTILITY_FILE
        totalFileName = TOTAL_UTILITY_FILE
        ls = utilityList
    else:
        return

    totalFile = open(totalFileName, mode='w+')
    writer = csv.writer(totalFile)
    with open(fileName, "r") as f:
        reader = csv.reader(f)
        for j, row in enumerate(reader):
            if j == 0:
                writer.writerow(headers[mode])
                continue
            writeRow = []
            for k, elem in enumerate(row):
                modeRow = ls[j]
                writeRow.append(float(elem) + float(modeRow[k]))
            writer.writerow(writeRow)
            ls[j] = writeRow
    totalFile.close()

    if mode == "death":
        deathList = ls
    elif mode == "SM":
        smList = ls
    elif mode == "SMChange":
        smChangeList = ls
    elif mode == "utility":
        utilityList = ls
    else:
        return

def divideByNumIters(mode):
    totalFileName = ""
    ls = []

    if mode == "death":
        totalFileName = TOTAL_DEATH_FILE
        ls = deathList
    elif mode == "SM":
        totalFileName = TOTAL_SOCIAL_MOTIVES_FILE
        ls = smList
    elif mode == "SMChange":
        totalFileName = TOTAL_SOCIAL_MOTIVES_CHANGE_FILE
        ls = smChangeList
    elif mode == "utility":
        totalFileName = TOTAL_UTILITY_FILE
        ls = utilityList
    else:
        return

    with open(totalFileName, 'w+') as totalFile:
        writer = csv.writer(totalFile)
        for j, row in enumerate(ls):
            if j == 0:
                continue
            for k, elem in enumerate(row):
                row[k] = float(row[k]) / NUM_ITERS
            writer.writerow(row)
            ls[j] = row


for i in range(NUM_ITERS):

    os.system("go run .")

    if i == 0:
        createTotalFile("death")
        createTotalFile("SM")
        createTotalFile("SMChange")
        createTotalFile("utility")
        continue

    updateTotalFile("death")
    updateTotalFile("SM")
    updateTotalFile("SMChange")
    updateTotalFile("utility")

divideByNumIters("death")
divideByNumIters("SM")
divideByNumIters("SMChange")
divideByNumIters("utility")