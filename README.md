# bajid-spotify-server

# Configuration

The application is built to run within GCP Cloud Run and is backed by Firestore.

It needs several parameters:
- SPOTIFY_CLIENT_ID is passed at built time in the Cloud Run
- SPOTIFY_CLIENT_SECRET is fetched from Secrets Manager

## Running a local firestore emulator

- Use Gcloud CLI to install [the Beta Emulator for Firestore](https://cloud.google.com/sdk/gcloud/reference/beta/emulators/firestore)
- Once installed, you can run it locally using the following command:
```
gcloud beta emulators firestore start --log-http --project local --host-port localhost:8081
```

It is recommended to set the port to 8081 instead of having a new one generated each time.

- You will also need to set the two following environment variables in your project:
```
FIRESTORE_EMULATOR_HOST=localhost:8081
GCP_PROJECT_ID=local
```
If set correctly, the program will not try to connect to any cloud-running Firestore instance and will use the local version instead.
You can adjust other parameters like log verbosity levels by using Gcloud wide flags as specified in the doc