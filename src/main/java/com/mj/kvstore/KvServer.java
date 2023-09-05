package com.mj.kvstore;


import com.mj.distributed.tcp.nio.NioListener;
import com.mj.distributed.tcp.nio.NioListenerConsumer;
import com.mj.kvstore.servlet.KVServlet;
import jakarta.servlet.Servlet;
import org.apache.catalina.Context;
import org.apache.catalina.startup.Tomcat;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.nio.ByteBuffer;
import java.nio.channels.SocketChannel;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentSkipListMap;

public class KvServer {

    private NioListener listener;
    private String host = "localhost";
    private int port = 7100;

    Map<String, String> map = new ConcurrentHashMap<>();

    static Logger LOGGER = LoggerFactory.getLogger(KvServer.class);

    public static void main(String[] args) throws Exception {

        LOGGER.info("Welcome to KvStore");

        KvServer kv = new KvServer();
        kv.start();

    }

    public void start() throws Exception{

        listener = new NioListener(host, port, new KvListenerConsumer());

        Tomcat tomcat = new Tomcat();

        String portString = System.getenv("KVSTOREPORT");

        int port = portString == null ? 8090 : Integer.getInteger(portString);

        tomcat.setPort(8090);
        tomcat.getConnector();

        Context ctx = tomcat.addContext("/", new File(".").getAbsolutePath());

        Tomcat.addServlet(ctx, "KVServlet", (Servlet) new KVServlet()) ;

        ctx.addServletMappingDecoded("/*", "KVServlet");

        tomcat.start();
        tomcat.getServer().await();

    }
}
