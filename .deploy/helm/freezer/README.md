# Updating the SOPS secrets

1. Install the latest [AWS-CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
2. Ask owners for access to AWS
3. Configure it [AWS-CLI Setup](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-quickstart.html)
    * Make sure you're using the default profile and NOT creating a new one (I.E. `aws configure --profile new-profile-name`)
4. Install [sops](https://github.com/mozilla/sops)
    * macOS only : `brew install sops`
5. cd `.deploy/helm/$SERVICE_NAME`
6. run sops to decrypt
    * staging
        * `sops -d -i staging/secret-values.yaml`
    * production
        * `sops -d -i production/secret-values.yaml`
7. replace whichever value you need with the plain text value
    * See info about the yaml structure & examples [here](https://github.com/ICE-Blockchain/helm-charts/tree/master/generic-service-chart)
8. run sops to encrypt
    * staging
        * `sops -e -i staging/secret-values.yaml`
    * production
        * `sops -e -i production/secret-values.yaml`
9. inspect the changes and if everything is ok, commit

# Testing the helm charts

1. Install [Helm](https://helm.sh)
2. Get access to [Organization Helm Chart Repository](https://github.com/ICE-Blockchain/helm-charts)
3. Generate [Github Personal Token](https://github.com/settings/tokens)
4. `helm repo add ICE-Blockchain https://helm-charts.CHANGE_ME.CHANGE_ME --username <your-login> --password <personal-token>`
5. `helm repo update`
6. `cd .deploy/helm/$SERVICE_NAME`
7. `helm pull ICE-Blockchain/generic-service-chart --version 0.1.0`
    * Use the proper version, from `Chart.yaml`
8. `sops -d -i staging/secret-values.yaml`
9. Setup [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
10. Ask for credentials for k8s staging clusters
11. copy the registry credentials from the real namespace
    1. ```kubectl get secret registry-credentials -o yaml -n $SERVICE_NAME```
    2. save it yo a yaml: I.E. `registry-credentials.yaml`
    3. remove everything from `metadata`, except `name` and `namespace`
    4. `kubectl apply -f registry-credentials.yaml`
12. ```helm upgrade --debug --atomic --wait --wait-for-jobs --cleanup-on-fail --install $SERVICE_NAME-tmp . -f common-values.yaml -f staging/common-values.yaml -f staging/secret-values.yaml -f staging/fra1-values.yaml -n $SERVICE_NAME-tmp --create-namespace --force```
13. After you're done (!!!this requires knowledge of k8s and kubectl), `helm delete $SERVICE_NAME-tmp -n $SERVICE_NAME-tmp`