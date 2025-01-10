package dummy;

import java.io.*;
import java.net.*;
import java.util.concurrent.*;

public class Server {
    private static final int QUEUE_SIZE = 10000;
    private static final int CONNECTION_TIMEOUT = 30000; // 30 seconds
    
    public void init() {
        int processors = Runtime.getRuntime().availableProcessors();
        ThreadPoolExecutor threadPoolExecutor = new ThreadPoolExecutor(
            processors * 2,                  
            processors * 4,                  
            60L, TimeUnit.SECONDS,
            new LinkedBlockingQueue<>(QUEUE_SIZE),
            new ThreadPoolExecutor.CallerRunsPolicy()
        );

        try (ServerSocket serverSocket = new ServerSocket(8080)) {
            System.out.println("Server started on port 8080");
            while(true) {
                Socket socket = serverSocket.accept();
                socket.setSoTimeout(CONNECTION_TIMEOUT);
                threadPoolExecutor.submit(() -> streamSSE(socket));
            }   
        } catch (Exception e) {
            System.err.println("Server error: " + e.getMessage());
        } finally {
            shutdownThreadPool(threadPoolExecutor);
        }
    }

    private void shutdownThreadPool(ThreadPoolExecutor executor) {
        executor.shutdown();
        try {
            if (!executor.awaitTermination(60, TimeUnit.SECONDS)) {
                executor.shutdownNow();
            }
        } catch (InterruptedException e) {
            executor.shutdownNow();
        }
    }

    public void streamSSE(Socket socket){
        String messages[] = new String[]{"this", "is", "data"};
        try (OutputStream outputStream = socket.getOutputStream()){
            PrintWriter writer = new PrintWriter(outputStream, true);
            writer.println("HTTP/1.1 200 OK");
            writer.println("Content-Type: text/event-stream");
            for(String msg : messages){
                writer.println();
                writer.println(String.format("data: %s\n\n", msg));
                writer.flush();
                Thread.sleep(1000);
            }
        } catch(Exception e){

        }
    }

}