# Aurora Test Authentication Server

## Overview

Path Check, Inc.'s prototype test authentication server provides the vital functionality of
generating and validating authorization codes that allow a mobile app user to upload their temporary exposure keys.

This test authentication server seeks to complement
[Google's exposure notifications server](https://github.com/google/exposure-notifications-server) by mirroring
its architecture. It stands as an independent microservice, yet it is designed to seamlessly integrate into the
exposure notifications server.

## Architecture

The Aurora Test Authentication Server is highly modular to ensure maximum compatibility with different systems.
As of now, it only supports Redis for caching authorization codes. However, the server's structure allows for
additional databases, which will be added in the future.
