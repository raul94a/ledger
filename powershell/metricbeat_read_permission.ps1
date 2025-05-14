icacls "..\elk\metricbeat.yml" /inheritance:r
icacls "..\elk\metricbeat.yml" /grant:r "$($env:USERNAME):(R,W)"


icacls "..\elk\filebeat.yml" /inheritance:r
icacls "..\elk\filebeat.yml" /grant:r "$($env:USERNAME):(R,W)"