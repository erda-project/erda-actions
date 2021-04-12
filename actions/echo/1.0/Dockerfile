FROM alpine

WORKDIR /opt/action

RUN echo -e "#!/bin/sh" > /opt/action/run && \
    echo -e "echo \${ACTION_WHAT-'Nothing I can echo :('}" >> /opt/action/run && \
    chmod +x /opt/action/run

COPY actions/echo/1.0/internal/.keep /tmp/.keep

CMD /opt/action/run