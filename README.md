# LibreDocker
A simple go server that converts documents via LibreOffice.

## Getting Started
1) Build the docker and run it.
2) Make a POST call to the container `url:8081/conver-file/{originalFileType}/{newFileType}`. The body should be a file in binary format.
3) Enjoy your converted file.

## Developing 
Open a PR with new features, to develop its best to run `docker compose watch`. 
