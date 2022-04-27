#!/bin/bash

# FIXME: Make these OS env vars w/ sane defaults
# OR: Parse them out of kubectl contexts
SERVER_URL="https://192.168.49.2:8443"
CACERT="/home/mat/.minikube/ca.crt"
CERT_PATH="/home/mat/.minikube/profiles/minikube/client.crt"
KEY_PATH="/home/mat/.minikube/profiles/minikube/client.key"

while IFS= read -r relurl
do
    url="${SERVER_URL}${relurl}"
    outName=$(echo "$relurl" | cut -d / -f 3,4 | tr / -)
    outPath="./third_party/minikube/v1.23.3/apis/${outName}.json"  # FIXME: support other k8s flavors/versions
    curl -sS --cacert ${CACERT} --cert ${CERT_PATH} --key ${KEY_PATH} ${url} -o ${outPath}
done < "./get_json_files.txt" # FIXME: Add top level jsons (above /apis/)
