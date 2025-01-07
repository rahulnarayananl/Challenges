import java.io.*;
import java.net.*;

public class SimpleServer {

    public static void main(String[] args) {
        int port = 8080; // Server will listen on this port

        try (ServerSocket serverSocket = new ServerSocket(port)) {
            System.out.println("Server is listening on port " + port);

            while (true) {
                Socket socket = serverSocket.accept();
                System.out.println("New client connected");

                // Handle the client's request in a separate thread
                new Thread(() -> handleClient(socket)).start();
                // TODO: move to ThreadPool
            }
        } catch (IOException ex) {
            System.out.println("Server exception: " + ex.getMessage());
            ex.printStackTrace();
        }
    }

    private static void handleClient(Socket socket) {
        try (InputStream input = socket.getInputStream();
             BufferedReader reader = new BufferedReader(new InputStreamReader(input));
             OutputStream output = socket.getOutputStream()) {

            // Parse the HTTP request
            Request request = parseRequest(reader);

            // Prepare the HTTP response
            Response response = new Response();
            handleRequest(request, response);

            // Send the response
            sendResponse(response, output);

        } catch (IOException ex) {
            System.out.println("Client handling exception: " + ex.getMessage());
            ex.printStackTrace();
        } finally {
            try {
                socket.close();
            } catch (IOException ex) {
                ex.printStackTrace();
            }
        }
    }

    private static Request parseRequest(BufferedReader reader) throws IOException {
        // Read the first line of the request (e.g., "GET / HTTP/1.1")
        String line = reader.readLine();
        if (line == null || line.isEmpty()) return null;

        String[] parts = line.split(" ");
        String method = parts[0];
        String path = parts[1];
        String version = parts[2];

        // Read headers (optional in this example)
        String header;
        StringBuilder headers = new StringBuilder();
        while ((header = reader.readLine()) != null && !header.isEmpty()) {
            headers.append(header).append("\n");
        }

        return new Request(method, path, version, headers.toString());
    }

    private static void handleRequest(Request request, Response response) {
        if (request == null) {
            response.setStatusCode(400);
            response.setContent("Bad Request");
            return;
        }

        // Handle different paths
        if (request.getPath().equals("/")) {
            response.setStatusCode(200);
            response.setContent("Welcome to the Simple HTTP Server!");
        } else if (request.getPath().equals("/hello")) {
            response.setStatusCode(200);
            response.setContent("Hello, World!");
        } else {
            response.setStatusCode(404);
            response.setContent("Page not found");
        }
    }

    private static void sendResponse(Response response, OutputStream output) throws IOException {
        PrintWriter writer = new PrintWriter(output, true);
        writer.println("HTTP/1.1 " + response.getStatusCode() + " " + response.getStatusMessage());
        writer.println("Content-Type: text/plain");
        writer.println("Content-Length: " + response.getContent().length());
        writer.println();
        writer.println(response.getContent());
        writer.flush();
    }

    // Request object
    static class Request {
        private final String method;
        private final String path;
        private final String version;
        private final String headers;

        public Request(String method, String path, String version, String headers) {
            this.method = method;
            this.path = path;
            this.version = version;
            this.headers = headers;
        }

        public String getMethod() {
            return method;
        }

        public String getPath() {
            return path;
        }

        public String getVersion() {
            return version;
        }

        public String getHeaders() {
            return headers;
        }
    }

    // Response object
    static class Response {
        private int statusCode;
        private String content;

        public int getStatusCode() {
            return statusCode;
        }

        public void setStatusCode(int statusCode) {
            this.statusCode = statusCode;
        }

        public String getStatusMessage() {
            String responseCode = "Internal Server Error";
            switch (statusCode) {
                case 200 : 
                    responseCode = "OK";
                    break;
                case 400: 
                    responseCode = "Bad Request";
                    break;
                case 404: 
                    responseCode = "Not Found";
                    break;
                default:
                    responseCode = "Internal Server Error";
                    break;
            }
            return responseCode;
        }

        public String getContent() {
            return content;
        }

        public void setContent(String content) {
            this.content = content;
        }
    }
}