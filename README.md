# serverless-apis-with-gcp
Build a fully serverless APIs with Google Cloud Platform

Steps
* Open Cloud Shell
* Setup environment variables
```
export GCP_PROJECT=[your-project-id]
export REGION=europe-west1
```
* Configure context for your project
```
gcloud config set project $GCP_PROJECT
gcloud config set compute/region $REGION
```
* Configure default Cloud Run mode (Fully managed mode)
```
gcloud config set run/platform managed
```
* Configure default Cloud Run and Functons deploy default region
```
gcloud config set run/region $REGION
gcloud config set functions/region $REGION
```
* Enable Apis
```
gcloud services enable run.googleapis.com
gcloud services enable cloudfunctions.googleapis.com
gcloud services enable endpoints.googleapis.com

gcloud services enable firebaserules.googleapis.com
gcloud services enable firestore.googleapis.com

gcloud services enable iap.googleapis.com
```
* Pin the frequently used products in the codelab
    - Cloud Run
    - Cloud Functions
    - App Engine
    - Endpoints
    - Datastore
    - Storage
    - Container Registry
    
# 1. Init Datastore (NoSql database)
* As App Engine has been activated, this has automatiquely created : 
    - Datastore instance
    - Two buckets ($GCP_PROJECT.appspot.com and staging.$GCP_PROJECT.appspot.com) --> we will be using the first one
```
export GCP_BUCKET=$GCP_PROJECT.appspot.com
```
* Copy an existing datastore export into your project bucket
```
gsutil cp -r gs://serverless-apis-with-gcp-2020/export  gs://$GCP_BUCKET
```

Import sample data into your datastore instance
```
gcloud datastore import gs://$GCP_BUCKET/export/export.overall_export_metadata
```

# 2. Deploy ESP
* Create the ESP Service Account
```
gcloud iam service-accounts create esp-sa --display-name='ESP Service Account'
```
* Add the required IAM roles the ESP Service Account
```
gcloud projects add-iam-policy-binding $GCP_PROJECT --member="serviceAccount:esp-sa@$GCP_PROJECT.iam.gserviceaccount.com" --role="roles/viewer"
gcloud projects add-iam-policy-binding $GCP_PROJECT --member="serviceAccount:esp-sa@$GCP_PROJECT.iam.gserviceaccount.com" --role="roles/servicemanagement.serviceController"
```
* Deploy the Extensible Service Proxy Container to Cloud Run with the previously created service account as identity
```
gcloud run deploy endpoints-runtime-serverless \
--image=gcr.io/endpoints-release/endpoints-runtime-serverless:1.44 \
--service-account="esp-sa@$GCP_PROJECT.iam.gserviceaccount.com" \
--allow-unauthenticated
```
Note the generated url looks like this:
https://endpoints-runtime-serverless-[random]q-ew.a.run.app

# 3. Deploy Cloud Run API
* Clone Github Repository
```
git clone https://github.com/Sfeir/serverless-apis-with-gcp.git
```
* Explore the application file app.go and the associated docker file
```
cd ~/serverless-apis-with-gcp/run
cat app.go
cat Dockerfile
```
* Build the image and push it to the container registry
```
docker build -t gcr.io/$GCP_PROJECT/app-inventory .
docker push gcr.io/$GCP_PROJECT/app-inventory
```
* Check that the image is stored to GCR
TOOLS > Container Registry > Images 
* (Optional) Enable Vulnerability Scanning
TOOLS > Container Registry > Settings 
Click on the blue button 'Enable Vulnerability scanning'
You should re-push the image to run the scanning
```
docker push gcr.io/$GCP_PROJECT/app-inventory
```
* Deploy the previously built image to Cloud Run 
```
gcloud run deploy cloud-run-api \
--image=gcr.io/$GCP_PROJECT/app-inventory \
--allow-unauthenticated \
--set-env-vars PROJECT_ID=$GCP_PROJECT
```
Note the Cloud Run generated url looks like:
```
https://cloud-run-api-[random]-ew.a.run.app
https://cloud-run-api-[random]-ew.a.run.app?year=2018
https://cloud-run-api-[random]-ew.a.run.app?year=2019
https://cloud-run-api-[random]-ew.a.run.app?year=2020
```
* Remove public access to the Cloud Run API
```
gcloud run services remove-iam-policy-binding cloud-run-api \
--member="allUsers" \
--role="roles/run.invoker"
```
* Setup IAM Policy to allow only requests from ESP
```
gcloud run services add-iam-policy-binding cloud-run-api \
--member="serviceAccount:esp-sa@$GCP_PROJECT.iam.gserviceaccount.com" \
--role="roles/run.invoker"
```

# 4. Deploy Cloud Functions API
* Deploy the cloud function code using the NodeJs 8 runtime 
```
cd ~/serverless-apis-with-gcp/functions
gcloud functions deploy cloud-functions-api \
--runtime=nodejs8 \
--trigger-http \
--entry-point=appInventory \
--allow-unauthenticated
```
Note the Cloud Function urls that you can try to verify that everything is OK:
```
https://europe-west1-[GCP_PROJECT].cloudfunctions.net/cloud-functions-api
https://europe-west1-[GCP_PROJECT].cloudfunctions.net/cloud-functions-api?year=2018
https://europe-west1-[GCP_PROJECT].cloudfunctions.net/cloud-functions-api?year=2019
https://europe-west1-[GCP_PROJECT].cloudfunctions.net/cloud-functions-api?year=2020
```
* Remove public access to the Cloud Functions API
```
gcloud functions remove-iam-policy-binding cloud-functions-api \
--member="allUsers" \
--role="roles/cloudfunctions.invoker"
```
* Setup IAM to allow only requests from ESP
```
gcloud functions add-iam-policy-binding cloud-functions-api \
--member="serviceAccount:esp-sa@$GCP_PROJECT.iam.gserviceaccount.com" \
--role="roles/cloudfunctions.invoker"
```

# 5. Deploy App Engine API
* Deploy the App Engine code using the Python 3.7 runtime 
```
cd ~/serverless-apis-with-gcp/appengine
gcloud app deploy
```
Note the App Engine url:
```
https://[GCP_PROJECT].appspot.com/
https://[GCP_PROJECT].appspot.com/?year=2018
https://[GCP_PROJECT].appspot.com/?year=2019
https://[GCP_PROJECT].appspot.com/?year=2020
```

* Setup Identity-Aware Proxy (IAP) to allow only requests from ESP
```
Burger Menu > Security > Identity-Aware Proxy
Configure Consent Screen
> Choose external
> Choose public
> Application name --> Application Inventory
> Optionaly you can upload a logo for the app
> Save your changes

Burger Menu > Security > Identity-Aware Proxy
Turn-on IAP by activating the radio button in the iap column
Try to access again and this time you will be denied!

Select 'App Engine App'
Click 'Add Member' button on the right
- New members : Your Email Address
- Select a role : Cloud IAP > IAP-secured Web App User

Try to access again (May take a few seconds to take effect) and this time you will be authorized!

Now grant the same access to the ESP service account and remove your account.
```

Note the App Engine url:
```
https://[GCP_PROJECT].appspot.com/
https://[GCP_PROJECT].appspot.com/?year=2018
https://[GCP_PROJECT].appspot.com/?year=2019
https://[GCP_PROJECT].appspot.com/?year=2020
```

# 6. Deploy the API Specification to Cloud Endpoints
* Update the specification file
```
cd ~/serverless-apis-with-gcp/swagger
nano app-inventory-api.yaml
Replace the following placeholders:
- ESP_URL : the cloud run service 'endpoints-runtime-serverless' url (be careful without https://)
- CLOUD_FUNCTION_URL : the cloud function 'cloud-functions-api' url
- CLOUD_RUN_URL : the cloud run service 'cloud-run-api' url
- APP_ENGINE_URL : The App Engine application url
- Client_ID : 
  - Burger Menu > APIs & Services > credentials
  - Under the 'OAuth 2.0 Client IDs' copy the 'client id' of 'IAP-App-Engine-app'

(to save the file in nano use CTRL O, to quit CTRL X)
```

* Deploy the specification to Cloud Endpoints
```
gcloud endpoints services deploy app-inventory-api.yaml
```
If everything is okay you should get an output similar to this :
```
Service Configuration [2020-02-03r0] uploaded for service [endpoints-runtime-serverless-tpkdhd4z7q-uc.a.run.app]
```
Save the service name 'endpoints-runtime-serverless-[random]-ew.a.run.app' that we will be using in the next step.
```
export ENDPOINT_SERVICE_NAME=endpoints-runtime-serverless-[random]-uc.a.run.app
```

* Update the ESP with the Endpoint service name
```
gcloud run services update endpoints-runtime-serverless \
   --set-env-vars ENDPOINTS_SERVICE_NAME=$ENDPOINT_SERVICE_NAME \
   --project $GCP_PROJECT
```
* Test access to the different versions of API using the ESP runtime URL
```
https://endpoints-runtime-serverless-[random]-ew.a.run.app/run
https://endpoints-runtime-serverless-[random]-ew.a.run.app/functions
https://endpoints-runtime-serverless-[random]-ew.a.run.app/appengine
```

# 7. Restricting API access with API keys
* Configures basic authentication with an API key by adding this section to the end of the spec file
```
security:
  - api_key: []
  
securityDefinitions:
  api_key:
    type: "apiKey"
    name: "key"
    in: "query"
```
* Deploy the updated specification to Cloud Endpoints
```
gcloud endpoints services deploy app-inventory-api.yaml
```
* Try again to access the API and you should get this error
```
{
 "code": 16,
 "message": "Method doesn't allow unregistered callers (callers without established identity). Please use API Key or other form of API consumer identity to call this API.",
 "details": [
  {
   "@type": "type.googleapis.com/google.rpc.DebugInfo",
   "stackEntries": [],
   "detail": "service_control"
  }
 ]
}
```
* Create an API key for your project
```
Burger Menu > APIs & Services > credentials
> + CREATE CREDENTIALS
> Select 'API Key'
> Copy The generated 'API Key' in order to use it in the url
https://endpoints-runtime-serverless-[random]-ew.a.run.app/run?key=[API Key]
> Close the popup
```

# 8. Configuring 
* Define an api metric 
```
x-google-management:
  metrics:
    # Define a metric for read requests.
    - name: "read-requests"
      displayName: "Read requests"
      valueType: INT64
      metricKind: DELTA
```  
* Define the metric limit 
```
  quota:
    limits:
      # Define the limit or the read-requests metric.
      - name: "read-limit"
        metric: "read-requests"
        unit: "1/min/{project}"
        values:
          STANDARD: 100
``` 

* The final 'x-google-management' section should be :
```
x-google-management:
  metrics:
    # Define a metric for read requests.
    - name: "read-requests"
      displayName: "Read requests"
      valueType: INT64
      metricKind: DELTA
  quota:
    limits:
      # Define the limit on the read-requests metric.
      - name: "read-limit"
        metric: "read-requests"
        unit: "1/min/{project}"
        values:
          STANDARD: 100
```

* Apply the quota by adding 'x-google-quota' to the '/appengine' section in the same level with x-google-backend
```
x-google-quota:
  metricCosts:
    read-requests: 20
```
* The final '/appengine' section should be :
```
"/appengine":
  get:
    summary: "Get apps deployed to App Engine"
    operationId: "app-inventory-appengine"
    x-google-backend:
      address: "https://codelab-project-XX-20200214.appspot.com"
      jwt_audience: 671771450352-79a9ihrtgggm309n6n0g5ga64utctlq2.apps.googleusercontent.com
    x-google-quota:
      metricCosts:
        read-requests: 20
```
* Deploy the updated specification to Cloud Endpoints
```
gcloud endpoints services deploy app-inventory-api.yaml
```
* Verify that the quota is working by calling many times until getting this error (After five tries) :
```
{
 "code": 8,
 "message": "Quota exceeded for quota metric 'get_requests' and limit 'get-limit' of service 'endpoints-runtime-serverless-tpkdhd4z7q-uc.a.run.app' for consumer 'project_number:671771450352'.",
 "details": [
  {
   "@type": "type.googleapis.com/google.rpc.DebugInfo",
   "stackEntries": [],
   "detail": "internal"
  }
 ]
}
```
# 9. Generate a Developer portal
- Burger Menu > Endpoints > Developer Portal
- Click on the blue button "Create portal" to generate a portal for the APIs in your project (this will take 5 to 10 minutes)
- When it's ready please visit https://endpointsportal.codelab-project-XX-20200214.cloud.goog/
