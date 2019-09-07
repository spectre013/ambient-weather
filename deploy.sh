
kubectl delete -f weather-data/deployment.yaml
kubectl create -f weather-data/deployment.yaml

kubectl delete -f weather-server/deployment.yaml
kubectl create -f weather-server/deployment.yaml

kubectl delete -f weather-ui/deployment-only.yaml
kubectl create -f weather-ui/deployment-only.yaml