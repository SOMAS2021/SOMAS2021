<<<<<<< HEAD
=======
rm -r configs
mkdir configs
python3 configGen.py

>>>>>>> adding config genrator scripts and a way to run all of them
CONFIG_DIR="configs/*.json"
for f in $CONFIG_DIR
do
    echo "Running $f"
    go run main.go --configpath=$f
    rm $f
done