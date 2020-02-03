# serverless-apis-with-gcp
Build a fully serverless APIs with Google Cloud Platform

Steps
* Open Cloud Shell
* Setup environment variables
```
export GCP_PROJECT=[your-project-id]
export PORT=9090
```
* Configure context for your project
```
gcloud config set project $GCP_PROJECT
```
* Configure default Cloud Run mode (Fully managed mode)
```
gcloud config set project $GCP_PROJECT
gcloud config set run/platform managed
```
* Enable Cloud Run Api
```
gcloud services enable run.googleapis.com
gcloud services enable cloudfunctions.googleapis.com
gcloud services enable endpoints.googleapis.com

For datas
gcloud services enable firebaserules.googleapis.com
gcloud services enable firestore.googleapis.com
```
* Pin the frequently used products in the codelab

* Clone Github Repository
```
git clone https://github.com/Sfeir/serverless-apis-with-gcp.git
```
# 0. Init database
Microservice names generator
https://project-names.herokuapp.com/names

# 1. Deploy ESP
* Create the ESP Service Account
```
gcloud iam service-accounts create esp-sa --display-name='ESP Service Account'
```
* Deploy the Extensible Service Proxy Container to Cloud Run with the previously created service account as identity
```
gcloud run deploy endpoints-runtime-serverless \
--image=gcr.io/endpoints-release/endpoints-runtime-serverless:1.44 \
--service-account="esp-sa@$GCP_PROJECT.iam.gserviceaccount.com"
```
Note the generated url looks like this:
https://endpoints-runtime-serverless-tpkdhd4z7q-uc.a.run.app

# 2. Deploy Cloud Run API
* Build the image and push it to the container registry
```
cd run
docker build -t gcr.io/$GCP_PROJECT/hello-devfest .
docker push gcr.io/$GCP_PROJECT/hello-devfest
```
* Deploy the previously built image to Cloud Run 
```
gcloud run deploy cloud-run-api \
--image=gcr.io/$GCP_PROJECT/hello-devfest \
--no-allow-unauthenticated
```
* Setup IAM to allow only requests from ESP
```
gcloud run services add-iam-policy-binding cloud-run-api \
--member="serviceAccount:esp-sa@$GCP_PROJECT.iam.gserviceaccount.com" \
--role="roles/run.invoker"
```
Note the Cloud Run generated url looks like:
https://cloud-run-api-tpkdhd4z7q-uc.a.run.app
https://cloud-run-api-tpkdhd4z7q-uc.a.run.app?year=2018
https://cloud-run-api-tpkdhd4z7q-uc.a.run.app?year=2019
https://cloud-run-api-tpkdhd4z7q-uc.a.run.app?year=2020

# 2. Deploy Cloud Functions API
* Deploy the cloud function code using the NodeJs 8 runtime 
```
cd ~/serverless-apis-with-gcp/functions
gcloud functions deploy cloud-functions-api --runtime=nodejs8 --trigger-http --entry-point=appInventory
```
* Setup IAM to allow only requests from ESP
```
gcloud functions add-iam-policy-binding cloud-functions-api \
--member="serviceAccount:esp-sa@$GCP_PROJECT.iam.gserviceaccount.com" \
--role="roles/cloudfunctions.invoker"

Note the Cloud Function url:

https://us-central1-[GCP_PROJECT].cloudfunctions.net/cloud-functions-api
https://us-central1-[GCP_PROJECT].cloudfunctions.net/cloud-functions-api?year=2018
https://us-central1-[GCP_PROJECT].cloudfunctions.net/cloud-functions-api?year=2019
https://us-central1-[GCP_PROJECT].cloudfunctions.net/cloud-functions-api?year=2020
```
