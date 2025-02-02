import java.util.*;


class SimpleJavaGC {
    private static final int EDEN_SIZE = 10;
    private static final int TENURED_SIZE = 10;
    
    private List<GCObject> edenSpace = new ArrayList<>();
    private List<GCObject> tenuredSpace = new ArrayList<>();
    private Set<GCObject> rootSet = new HashSet<>();

    public GCObject createObject(String name, GCObject... references) {
        GCObject obj = new GCObject(name);
        Collections.addAll(obj.references, references);
        
        edenSpace.add(obj);
        return obj;
    }

    public void addRoot(GCObject obj) {
        rootSet.add(obj);
    }

    public void garbageCollect() {
        System.out.println("\n--- Starting Garbage Collection ---");
        markPhase();
        sweepPhase();
        promoteSurvivors();
    }

    private void markPhase() {
        System.out.println("Mark Phase: Traversing from roots.");
        for (GCObject root : rootSet) {
            traverseAndMark(root);
        }
    }

    private void traverseAndMark(GCObject obj) {
        if (obj == null || obj.marked) return;
        obj.marked = true;
        System.out.println("Marking: " + obj.name);
        for (GCObject ref : obj.references) {
            traverseAndMark(ref);
        }
    }

    private void sweepPhase() {
        System.out.println("Sweep Phase: Reclaiming unmarked objects.");
        edenSpace.removeIf(obj -> !obj.marked);
        tenuredSpace.removeIf(obj -> !obj.marked);
        resetMarks();
    }

    private void promoteSurvivors() {
        Iterator<GCObject> it = edenSpace.iterator();
        while (it.hasNext()) {
            GCObject obj = it.next();
            obj.age++;
            if (obj.age > 1 && tenuredSpace.size() < TENURED_SIZE) { // Age threshold for promotion
                tenuredSpace.add(obj);
                it.remove();
                System.out.println("Promoting to Tenured: " + obj.name);
            }
        }
    }

    private void resetMarks() {
        edenSpace.forEach(obj -> obj.marked = false);
        tenuredSpace.forEach(obj -> obj.marked = false);
    }
}
