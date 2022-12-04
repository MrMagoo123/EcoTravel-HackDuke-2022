env GOARCH="amd64" GOOS="linux" go build -o build/debian_build

echo "Starting Service"
gcloud compute ssh source-aggregator --project=duke-hackathon-2022 -- 'sudo systemctl daemon-reload;sudo systemctl stop aggregator;pkill debian_build'

cd ..

echo "Uploading Executable"
gcloud compute scp --recurse ./source-aggregator/build source-aggregator:~ --project=duke-hackathon-2022

echo "Starting Service"
gcloud compute ssh source-aggregator --project=duke-hackathon-2022 -- 'sudo systemctl daemon-reload;sudo mv /home/administrator/build/aggregator.service /etc/systemd/system/aggregator.service;sudo systemctl start aggregator'
