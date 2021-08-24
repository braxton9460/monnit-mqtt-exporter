# monnit-mqtt-exporter
Prometheus exporter from mqtt for monnit sensors

The monnit sensors from monnit.com are excellent, and their edge gateway exports to mqtt. However the message format that they use is unsupported by any exporter I could find, and I want the data into prometheus for monitoring and alerting.

This exporter knows how to ingest the monnit specific topic & message structures, while retaining all parts of the sensors data in the prometheus output.

There are likely many golang improvements that could be made to this project. Contributions welcome.

# Layout
How this code is structured

### internal/sensors
- Prometheus metrics to be utilized by the ingestor packages
- Lookup tables for the types of sensors that exist
  - Source: https://www.monnit.com/support/knowledgebase/on-premises/understanding-sensor-application-ids/

A key component is different types of sensors can report the same types of metrics.. i.e temperature 

Since metrics shared across sensor types should be of the same type (i.e. temp, humidity), the definitions of those can be shared. But there will have to be a "lookup table" that defines what sensor types have what metrics that the ingestor can reference by sensor type id, as that is available in the mqtt payload.
