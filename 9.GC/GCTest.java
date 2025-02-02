

public class GCTest {
    public static void main(String[] args) {
        SimpleJavaGC gc = new SimpleJavaGC();

        // Structured object creation
        GCObject objA = gc.createObject("A");
        GCObject objB = gc.createObject("B", objA);
        GCObject objC = gc.createObject("C", objB);
        GCObject objD = gc.createObject("D", objC);
        GCObject objE = gc.createObject("E", objD); // Deeply nested reference

        // Adding some root objects
        gc.addRoot(objA);
        gc.addRoot(objC);

        // Trigger GC
        System.out.println("\nTriggering Garbage Collection...");
        gc.garbageCollect();
    }
}
