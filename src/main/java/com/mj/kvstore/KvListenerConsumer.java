package com.mj.kvstore;

import com.mj.distributed.tcp.nio.NioListenerConsumer;

import java.nio.ByteBuffer;
import java.nio.channels.SocketChannel;

public class KvListenerConsumer implements NioListenerConsumer {

    public void addedConnection(SocketChannel s) {

    }

    public void droppedConnection(SocketChannel s) {

    }

    public void consumeMessage(SocketChannel s, int numBytes, ByteBuffer b) {

    }
}
