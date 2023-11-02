FROM gradle:6.9.4-jdk 
RUN gradle init --type java-application
COPY build.gradle app/
RUN gradle war
