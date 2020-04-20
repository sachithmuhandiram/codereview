# Code Reviever 

This is a hobby project developed by myself and [Vikum](https://www.linkedin.com/in/vikum-dheemantha-b2449a146/).

Design and architectural advices from [Anjana](https://www.linkedin.com/in/anjana-kasun-arayawatta-55114b48/)

This system will facilitate people to submit their codes to their peers and get their comments.

No need to have locally installed `golang` instance. If you have locally installed `golang` ,some libries may show as not found locally in your IDE. It wont be a problem. 

Project developed and tested in Debian environment

## Prerequest

* [Docker](https://www.docker.com/) and [Docker-Compose](https://docs.docker.com/compose/).

* MySQL 5.7.* ( 5.7.29 used for development)

Linux environment is prefered. Didnt test on Windows.

## How to run

`git clone https://github.com/sachithmuhandiram/codereview.git`

`cd codereview`

`sudo docker-compose up --build`


## Configuration

### Networking

In order to run services, add followings to your `/etc/hosts` file 
```
0.0.0.0         user                # runs users service
0.0.0.0         notification        # notification service
0.0.0.0         localhost           # gateway
```

### Databases

Database files are located in `resources/` folder. For this project, I have used `root` user with a simple password `7890`. 

To add your used and password, please edit `api.env` file.

```
 MYSQLDBGATEWAY=database_user:password@tcp(127.0.0.1:3306)/api_gateway
 MYSQLDBUSERS=database_user:password@tcp(127.0.0.1:3306)/codereview_users
```

If you run your MySQL instance other than localhost (127.0.0.1), please add that IP address to `api.env`.

* Import `sql` schema files to your MySQL database.

### Notifications

To send emails, you need to configure your email server, or gmail. I have used gmail for test.

* Create a file in `notification` folder `emailData.json` and add :

```
{
    "From" : "your gmail",
    "Parse" : "parse given by gmail",
}
```

## Important Notes

To change Docker time, you need to add you timezone to each service's Dockerfile.
Here we used `Asia/Colombo` timezone, add your time zone and continet for following configs.

```
# # setting container time zone
RUN \
    apk --update add curl bash nano tzdata && \
    cp /usr/share/zoneinfo/Asia/Colombo /etc/localtime && \
    echo "Asia/Colombo" > /etc/timezone && \
    apk del tzdata && \
    rm -r /var/cache/apk/* && \
    mkdir -p /usr/share/zoneinfo/Asia && \
    ln -s /etc/localtime /usr/share/zoneinfo/Asia/Colombo
```