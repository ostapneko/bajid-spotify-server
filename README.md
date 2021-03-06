# bajid-spotify-server

A simple configurable Juke-Box to teach my toddler the alphabet listening to his favorite tunes!

![Screenshot](screenshot.jpg)

# Configuration

The application is built to run within GCP Cloud Run and is backed by Firestore.

It needs several parameters:
- `GCP_PROJECT_ID`, `SPOTIFY_CLIENT_ID` and `SPOTIFY_REDIRECT_URI` is passed at built time in the Cloud Run
- `SPOTIFY_CLIENT_SECRET` is fetched from Secrets Manager

# Development
## Running a local firestore emulator

Since the Firestore emulator is still rough around the edges, here are a few tips to help you operate it:

- Use GCloud CLI to install [the Beta Emulator for Firestore](https://cloud.google.com/sdk/gcloud/reference/beta/emulators/firestore)
- Once installed, you can run it locally using the following command:
```
# Development
gcloud beta emulators firestore start --log-http --project local --host-port localhost:8081
```

It is recommended to choose a port instead of having a new one generated each time.

- You will also need to set the two following environment variables in your project:
```
FIRESTORE_EMULATOR_HOST=localhost:8081
GCP_PROJECT_ID=local
```
If set correctly, the program will not try to connect to any cloud-running Firestore instance and will use the local version instead.
You can adjust other parameters like log verbosity levels by using Gcloud wide flags as specified in the doc
