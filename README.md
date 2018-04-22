# gdriver-go

Used for uploading files from a directory and its sub-directories to Google Drive.


## Setup
A `etc/client_secret.json` file will need to be created. Follow these instructions:

1. Use [this wizard](https://console.developers.google.com/start/api?id=drive) to create or select a project in the Google Developers Console and automatically turn on the API. Click Continue, then Go to credentials.
2. On the Add credentials to your project page, click the Cancel button.
3. At the top of the page, select the OAuth consent screen tab. Select an Email address, enter a Product name if not already set, and click the Save button.
4. Select the Credentials tab, click the Create credentials button and select OAuth client ID.
5. Select the application type Other, enter your desired name, and click the Create button.
6. Click OK to dismiss the resulting dialog.
7. Click the (Download JSON) button to the right of the client ID.
8. Move this file to your `etc/` directory and rename it to `client_secret.json`.

On first run, follow the instructions in the program to allow access to your Google Drive account.
