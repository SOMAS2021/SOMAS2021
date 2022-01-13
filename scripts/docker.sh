set -euxo pipefail

docker build -t ui-test .
docker run -it -p 9000:9000 ui-test\
