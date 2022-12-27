[![Actions Status](https://github.com/tommzn/hob-apigw-handler/actions/workflows/go.image.build.yml/badge.svg)](https://github.com/tommzn/hob-apigw-handler/actions)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/hob-apigw-handler)

# HomeOffice Button - API Gateway Request Handler
A handler to process time tracking requests from AWS API Gateway.  
This handler belongs to the [HomeOffice Button - Time Tracking](https://github.com/tommzn/hob-timetracker) Project.

## API Description
This handler can be used with an API which provides access to a capture and a generatereport endpoint.
### Example
'''yaml
swagger: "2.0"
info:
  title: "Time Tracking API"
basePath: "/v1"
schemes:
- "https"
paths:
  /capture:
    post:
      consumes:
      - "application/json"
      produces:
      - "application/json"
      parameters:
      - in: "body"
        name: "TimeTrackingRecord"
        required: true
        schema:
          $ref: "#/definitions/TimeTrackingRecord"
      responses:
        "200":
          description: "200 response"
          schema:
            $ref: "#/definitions/Empty"
  /generatereport:
    post:
      consumes:
      - "application/json"
      produces:
      - "application/json"
      parameters:
      - in: "body"
        name: "ReportGenerateRequest"
        required: true
        schema:
          $ref: "#/definitions/ReportGenerateRequest"
      responses:
        "200":
          description: "200 response"
          schema:
            $ref: "#/definitions/Empty"
      
definitions:
  Empty:
    type: "object"
    title: "Empty Schema"
  ReportGenerateRequest:
    type: "object"
    required:
    - "month"
    - "type"
    - "year"
    properties:
      type:
        type: "string"
        description: "Type of a report. Atm a monthly report is supported, only."
      year:
        type: "string"
        description: "Year, a report should be generated for."
      month:
        type: "string"
        description: "Month, a report should be generated for."
    title: "ReportGenerateRequest"
    description: "Request to start gennerating a time tracking report."
  TimeTrackingRecord:
    type: "object"
    required:
    - "clicktype"
    - "deviceid"
    properties:
      deviceid:
        type: "string"
        description: "Id of a device which has captured this time tracking event."
      clicktype:
        type: "string"
        description: "Type of click on a device."
      timestamp:
        type: "string"
        format: "date-time"
        description: "timestamp this time tracking event has been captured. (Optional)"
    title: "TimeTrackingRecord"
    description: "Single time tracking event."
'''

# Links
[HomeOffice Button - Time Tracking](https://github.com/tommzn/hob-timetracker)  
[AWS IoT 1-Click](https://aws.amazon.com/iot-1-click/?nc1=h_ls)  
