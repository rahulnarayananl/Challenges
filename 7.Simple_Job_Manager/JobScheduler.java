
import java.lang.reflect.Constructor;
import java.util.*;
import java.util.concurrent.*;

public class JobScheduler {
    private final Map<String, JobNode> jobMap = new HashMap<>();
    private final ScheduledExecutorService executorService = Executors.newScheduledThreadPool(5);

    public void addJob(String jobClassName, List<String> dependencies, long delayMillis, boolean isPeriodic, long periodMillis) {
        long executionTime = System.currentTimeMillis() + delayMillis;
        jobMap.put(jobClassName, new JobNode(jobClassName, dependencies, executionTime, isPeriodic, periodMillis));
    }

    // Resolve dependencies using Topological Sort
    private List<JobNode> resolveDependencies() throws Exception {
        Map<String, Integer> inDegree = new HashMap<>();
        Map<String, List<String>> adjList = new HashMap<>();
        
        for (String job : jobMap.keySet()) {
            inDegree.put(job, 0);
            adjList.put(job, new ArrayList<>());
        }

        // Build dependency graph
        for (JobNode jobNode : jobMap.values()) {
            for (String dep : jobNode.dependencies) {
                adjList.get(dep).add(jobNode.jobClassName);
                inDegree.put(jobNode.jobClassName, inDegree.get(jobNode.jobClassName) + 1);
            }
        }

        Queue<String> queue = new LinkedList<>();
        List<JobNode> sortedJobs = new ArrayList<>();
        
        for (Map.Entry<String, Integer> entry : inDegree.entrySet()) {
            if (entry.getValue() == 0) queue.add(entry.getKey());
        }

        while (!queue.isEmpty()) {
            String job = queue.poll();
            sortedJobs.add(jobMap.get(job));

            for (String neighbor : adjList.get(job)) {
                inDegree.put(neighbor, inDegree.get(neighbor) - 1);
                if (inDegree.get(neighbor) == 0) queue.add(neighbor);
            }
        }

        if (sortedJobs.size() != jobMap.size()) {
            throw new Exception("Cycle detected in job dependencies!");
        }

        return sortedJobs;
    }

    public void executeJobs() throws Exception {
        List<JobNode> sortedJobs = resolveDependencies();
        long currentTime = System.currentTimeMillis();

        for (JobNode jobNode : sortedJobs) {
            long delay = Math.max(jobNode.executionTime - currentTime, 0);
            Runnable jobTask = createJobTask(jobNode.jobClassName);

            if (jobNode.isPeriodic) {
                executorService.scheduleAtFixedRate(jobTask, delay, jobNode.period, TimeUnit.MILLISECONDS);
            } else {
                executorService.schedule(jobTask, delay, TimeUnit.MILLISECONDS);
            }
        }
    }

    private Runnable createJobTask(String jobClassName) throws Exception {
        Class<?> clazz = Class.forName(jobClassName);
        Constructor<?> constructor = clazz.getConstructor();
        Job jobInstance = (Job) constructor.newInstance();
        return jobInstance::execute;
    }

    public void shutdown() {
        executorService.shutdown();
    }
}