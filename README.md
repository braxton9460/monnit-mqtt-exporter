# monnit-mqtt-exporter
Prometheus exporter from mqtt for monnit sensors

The monnit sensors from monnit.com are excellent, and their edge gateway exports to mqtt. However the message format that they use is unsupported by any exporter I could find, and I want the data into prometheus for monitoring and alerting.

This exporter knows how to ingest the monnit specific topic & message structures, while retaining all parts of the sensors data in the prometheus output.

This are likely many golang improvements that could be made to this project. Contributions welcome.
