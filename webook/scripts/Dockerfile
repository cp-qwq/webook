FROM williamyeh/wrk

COPY /wrk/login.lua /post.lua

ENTRYPOINT ["wrk"]


#测试指令
#docker run --network host --rm wrk-login -t4 -c20 -d10s -s /post.lua http://172.20.10.8:8080/users/login