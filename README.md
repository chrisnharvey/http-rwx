# HTTP-RWX

A HTTP server that reads query string variables, passes them to a template, writes the
output to disk, then executes a command.

## WTF?

I needed a way to update a config file on one of my servers based on the data passed
from a webhook. This tool allows me to do exactly that, then restart the service to
reload the config file on the server.