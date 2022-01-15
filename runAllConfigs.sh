CONFIG_DIR="configs/*.json"
for f in $CONFIG_DIR
do
    echo "Running $f"
    go run main.go --configpath=$f
    rm $f
done