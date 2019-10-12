FROM ubuntu:18.04

RUN apt-get update && \
    apt-get install -y stress-ng && \
    apt-get clean

WORKDIR /stress

COPY runstress runstress
COPY jobfile jobs/jobfile

CMD ["./runstress"]
