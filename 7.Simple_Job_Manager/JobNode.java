import java.util.*;

public class JobNode {
    String jobClassName;
    List<String> dependencies;
    Long executionTime; 
    boolean isPeriodic;
    Long period;

    public JobNode(String jobClassName, List<String> dependencies, Long executionTime, boolean isPeriodic, Long period) {
        this.jobClassName = jobClassName;
        this.dependencies = dependencies != null ? dependencies : new ArrayList<>();
        this.executionTime = executionTime;
        this.isPeriodic = isPeriodic;
        this.period = period;
    }
}