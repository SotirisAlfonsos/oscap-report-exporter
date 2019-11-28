[![Build Status](https://travis-ci.org/SotirisAlfonsos/oscap-report-exporter.svg)](https://travis-ci.org/SotirisAlfonsos/oscap-report-exporter)
[![Report card](https://goreportcard.com/badge/github.com/SotirisAlfonsos/oscap-report-exporter)](https://goreportcard.com/report/github.com/SotirisAlfonsos/oscap-report-exporter)

# Oscap Report Exporter
The purpose of this project is to provide a fully configurable scheduler for oscap scans. 

## Functionality
- Schedule the date and time for the scan
- Configure the location of the global vulnerability report  
- Send the reports to the defined outputs

## Outputs
- Send the reports to a webhook. Request body format:
```
body{"results":< xml formated report >}
```
- Send the <i>.html</i> report to recipients via e-mail

## Installation
- Download the latest release containing the binary
```
wget https://github.com/SotirisAlfonsos/oscap-report-exporter/releases/download/<version>/oscap-exporter-<version>.<arch>..tar.gz
```
- Extract the binary to a folder
- Create configuration file
```
#Date on which the scan will be performed (Mon/Tue/Wed/Thu/Fri/Sat/Sun/Daily). Defaults to 'Sun'
scan_date: "Sun"

#Time on which the scan will be performed (HH:MM). Defaults to '23:00'
scan_time: "23:00"

#The working folder is the location where the global vulnerability reports file be downloaded,
#and the location when the results of the scan will be stored. Defaults to '/tmp/downloads/'
working_folder: "/tmp/downloads/"

#Webhook where the report will be sent. We send the report.xml file content by default.
#The webhook should be able to handle xml input, via a POST request. Defaults to 'http://localhost:8080' 
webhook: "http://localhost:8080"

#whether or not to clean the downloaded and created files after the scan. Defaults to true
clean_files: false

#Location of the global vulnerability report. Defaults to the official Red Hat location.
#You can also configure your own location. Basic auth credentials will be used if provided for the download.
vulnerability_report:
  global_vulnerability_report_https_location: "https://www.redhat.com/security/data/metrics/ds/com.redhat.rhsa-all.ds.xml"
  username: ""
  password: ""

#The settings for sending the html report to an email address. If left empty it is ignored.
email_config:
  smtp_smarthost: ""
  from: ""
  to: ""
  password: ""

```
- Create service file <i>oscap-exporter.service</i>. Service file content example
```
[Unit]
Description=Oscap report exporter
Wants=network-online.target
After=network-online.target

[Service]
User=root
Group=root
Type=simple
Restart=on-failure
ExecStart=</location/of/scheduler> --config.file=</location/of/oscap-config.yml>

[Install]
WantedBy=multi-user.target
```

## Usage
Run service
```
systemctl enable oscap-exporter
systemctl start oscap-exporter
```

## License
[MIT License](https://choosealicense.com/licenses/mit/)
