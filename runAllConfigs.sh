<<<<<<< HEAD
=======
rm -r configs
mkdir configs
python3 configGen.py
>>>>>>> made come configs and added some scripts
CONFIG_DIR="configs/*.json"
for f in $CONFIG_DIR
do
    echo "Running $f"
    go run main.go --configpath=$f
<<<<<<< HEAD
    rm $f
=======
>>>>>>> made come configs and added some scripts
done