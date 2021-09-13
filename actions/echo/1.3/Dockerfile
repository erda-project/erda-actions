FROM alpine

WORKDIR /opt/action

RUN echo -e "#!/bin/sh" > /opt/action/run && \
    echo -e "for i in \$(seq 1 \${ACTION_COUNT-10}); do echo \${ACTION_WHAT-'Nothing I can echo :('}; done" >> /opt/action/run && \
    chmod +x /opt/action/run

COPY actions/echo/1.3/internal/.keep /tmp/.keep

CMD /opt/action/run