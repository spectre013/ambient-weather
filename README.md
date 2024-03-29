# Ambient Weather

This is a custom display for data from the Ambient Weather PWS api, and requires an Ambient Weather Station. view live @ https://weather.zoms.net

# Status
Currently the status of this repo is in a beta state and is not complete. Bug reports are welcome but many features are not complete or still changing at this time 

![GitHub Logo](image.png)

# Setup
Recommended set up is  using docker in a linux environment that is operational 24 hours a day. Though this set up could be applied to any system. 
These are basic instructions that require some knowledge of setting up websites. More detailed instructions for all systems coming soon.

### PostgreSQL Setup
	
    psql set up comming soon. 

###API KEYS
Api keys will need to be acquired for the following interfaces

- Virtual Crossing - https://www.visualcrossing.com/
- IPGeolocation - https://ipgeolocation.io/astronomy-api.html
- Ambient Weather - through your ambient console

### Docker

Install Docker (https://www.docker.com/) follow instructions for your platform

To build the docker containers in the main directory run  

```$bash
$ docker-compose build
```

To push to a repository 

 ```$bash
 $ docker-compose push
 ```

execute run.sh file in each folder and update the environment variables with your own values. Windows user rename the file to .bat and then run.


# This repository is maintained Brian (aka Spectre013)

Brian will be maintaining on behalf of Brian Underdown.
Please use the support available via the Github issues function. This project is completely separate from any of the Meteobridge-Weather34-Template repositories and any support for this version should be requested in this repository.  


#LICENSE
<!--
Copyright (c) 2019-2023 by Brian Paulson (https://weather.zoms.net) JS/SQL/HTML/CSS/SVG
-->
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Template”), to deal in the Template without restriction, including without limitation the rights to, can use, can not copy without prior permission, can modify for personal use, can use and publish for personal use ,can not distribute without prior permission, can not sublicense without prior permission, and can not sell copies of the Template, and subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Template.

THE TEMPLATE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE TEMPLATE OR THE USE OR OTHER DEALINGS IN THE TEMPLATE.


Attribution-NonCommercial 4.0 International based on a work at https://weather34.com/homeweatherstation
