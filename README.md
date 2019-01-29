# synectiks-monitoring
This is a fork from grafana monitoring platform. The synectiks monitoring platform stack primariliy have extended the 
following tools & utilities  and composed of some of our own written components:
1) The collection stack uses prometheus / some prometheus exporters ( written by us) / Telegraph collector / some plugins for telegraph (written by us),
and some scrap jobs that run in jenkin container cluster.
2) The Data storage stack primarily consist of influxdb/ OpenTSDB for realtime data store , Elasticsaearch cluster for some meta data and structured data store 

3) The UI layers dashboard & alrets are mainly written as grafana plugins , primarily written by using synectiks react libarary , and we have used many cases to create 
dashboard based on graphql queries.
