FROM gradle:8.3.0-jdk 
RUN gradle init --type java-application
COPY build.gradle app/
RUN gradle war
