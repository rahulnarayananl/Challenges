import java.util.*;

class GCObject {
    String name;
    List<GCObject> references = new ArrayList<>();
    boolean marked = false; // Used in mark phase
    int age = 0; // For promotion tracking

    public GCObject(String name) {
        this.name = name;
    }

    public void addReference(GCObject obj) {
        references.add(obj);
    }

    @Override
    public String toString() {
        return name;
    }
}
