
import java.util.Arrays;
import java.util.Collections;

public class Main {
    public static void main(String[] args) {
        try {
            JobScheduler scheduler = new JobScheduler();

            scheduler.addJob("JobA", Collections.emptyList(), 2000, false, 0); // Runs in 2s
            scheduler.addJob("JobB", Arrays.asList("JobA"), 4000, false, 0); // Runs after A
            scheduler.addJob("JobC", Arrays.asList("JobA", "JobB"), 6000, false, 0); // Runs after A and B

            scheduler.addJob("JobA", Collections.emptyList(), 1000, true, 5000);

            scheduler.executeJobs();

            Thread.sleep(20000);
            scheduler.shutdown();

        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}

class JobA implements Job {
    @Override
    public void execute() {
        System.out.println("Executing Job A");
    }
}

class JobB implements Job {
    @Override
    public void execute() {
        System.out.println("Executing Job B");
    }
}

class JobC implements Job {
    @Override
    public void execute() {
        System.out.println("Executing Job C");
    }
}

/*
Good question! While periodic jobs generally won’t interfere with each other, here are some edge cases where conflicts could happen:

🔴 Potential Issues & Conflicts

1️⃣ Circular Dependencies (Job Deadlock)

🚨 Problem: If two jobs depend on each other directly or indirectly, the scheduler will detect a cycle and throw an exception.

Example:
	•	JobA depends on JobB
	•	JobB depends on JobC
	•	JobC depends on JobA (💥 Circular Dependency)

scheduler.addJob("JobA", Arrays.asList("JobC"), 1000, true, 5000);
scheduler.addJob("JobB", Arrays.asList("JobA"), 2000, true, 5000);
scheduler.addJob("JobC", Arrays.asList("JobB"), 3000, true, 5000);

💥 Result: Cycle detected → Scheduler throws an exception.

2️⃣ High Workload Leading to Job Overlaps

🚨 Problem: If a periodic job takes longer to execute than its scheduled period, overlapping executions might occur, leading to:
	•	Increased CPU & memory usage
	•	Thread starvation

Example:

scheduler.addJob("HeavyJob", Collections.emptyList(), 1000, true, 2000); // Runs every 2s

But HeavyJob takes 5s to execute, meaning it starts piling up, leading to:

Executing HeavyJob... (takes 5s)
[2s later] HeavyJob starts again before the previous one finishes! 💥

🔹 Solution: Use scheduleWithFixedDelay() instead of scheduleAtFixedRate().

3️⃣ Dependency Job Hasn’t Finished Before Next Cycle

🚨 Problem: If JobB depends on JobA, but JobA is periodic and hasn’t finished before JobB’s turn, JobB might get stuck or run with outdated data.

Example:

scheduler.addJob("JobA", Collections.emptyList(), 1000, true, 3000); // Runs every 3s
scheduler.addJob("JobB", Arrays.asList("JobA"), 2000, true, 5000);   // Runs every 5s, depends on JobA

💥 Issue:
	•	JobA starts at t=1s, takes 4s to finish.
	•	JobB is scheduled at t=2s but JobA is still running.
	•	JobB either waits indefinitely or runs with incomplete JobA results.

🔹 Solution: Make sure dependent jobs check if the parent job has completed before running.

4️⃣ Thread Pool Exhaustion

🚨 Problem: The default thread pool (ScheduledExecutorService) has limited threads. If too many long-running jobs are scheduled, new jobs will be delayed or never run.

🔹 Solution: Use a larger thread pool:

ScheduledExecutorService executorService = Executors.newScheduledThreadPool(10); // Increase threads

✅ How to Prevent These Issues?

1️⃣ Check for Circular Dependencies before scheduling.
2️⃣ Use scheduleWithFixedDelay() if jobs take long to execute.
3️⃣ Increase the thread pool if many jobs run concurrently.
4️⃣ Ensure dependent jobs wait for completion before executing.
5️⃣ Monitor execution times and adjust scheduling dynamically.

Final Thought

For independent periodic jobs, you’re good! 🚀
For dependent periodic jobs, be cautious about execution time & dependencies!
*/