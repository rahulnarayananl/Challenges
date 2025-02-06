import java.util.*;
import java.util.function.Predicate;
import java.util.function.Function;
import java.util.stream.Collectors;

interface GenericRepository<T, ID> {
    void create(T entity, ID id);
    void update(ID id, T entity);
    void delete(ID id);
    Optional<T> getById(ID id);
    List<T> getAll();
    List<T> filter(Predicate<T> condition);
    List<T> getAllSorted(Comparator<T> comparator); 
}

class InMemoryRepository<T, ID> implements GenericRepository<T, ID> {
    private final Map<ID, T> store = new HashMap<>();

    @Override
    public void create(T entity, ID id) {
        store.put(id, entity);
    }

    @Override
    public void update(ID id, T entity) {
        store.put(id, entity);
    }

    @Override
    public void delete(ID id) {
        store.remove(id);
    }

    @Override
    public Optional<T> getById(ID id) {
        return Optional.ofNullable(store.get(id));
    }

    @Override
    public List<T> getAll() {
        return new ArrayList<>(store.values());
    }

    @Override
    public List<T> filter(Predicate<T> condition) {
        return store.values().stream()
                .filter(condition)
                .collect(Collectors.toList());
    }

    @Override
    public List<T> getAllSorted(Comparator<T> comparator) {
        return store.values().stream()
                .sorted(comparator)
                .collect(Collectors.toList());
    }
}

class User {
    private final long id;
    private final String name;
    private final int age;

    public User(long id, String name, int age) {
        this.id = id;
        this.name = name;
        this.age = age;
    }

    public long getId() {
        return id;
    }

    public int getAge() {
        return age;
    }

    @Override
    public String toString() {
        return "User{id=" + id + ", name='" + name + "', age=" + age + "}";
    }
}

class Post {
    private final long postId;
    private final long userId;
    private final String content;

    public Post(long postId, long userId, String content) {
        this.postId = postId;
        this.userId = userId;
        this.content = content;
    }

    public long getUserId() {
        return userId;
    }

    public long getPostId() {
        return postId;
    }

    @Override
    public String toString() {
        return "Post{id=" + postId + ", content='" + content + "'}";
    }

    
}

class JoinUtils {
    public static <L, R, K> Map<L, List<R>> innerJoin(
            List<L> leftList, 
            Function<L, K> leftKeyExtractor, 
            List<R> rightList, 
            Function<R, K> rightKeyExtractor) {

        // Group rightList by join key
        Map<K, List<R>> rightGrouped = rightList.stream()
                .collect(Collectors.groupingBy(rightKeyExtractor));

        // Perform INNER JOIN
        return leftList.stream()
                .filter(left -> rightGrouped.containsKey(leftKeyExtractor.apply(left)))
                .collect(Collectors.toMap(
                        left -> left,
                        left -> rightGrouped.get(leftKeyExtractor.apply(left))
                ));
    }

    // Generic LEFT JOIN method
    public static <L, R, K> Map<L, List<R>> leftJoin(
            List<L> leftList, 
            Function<L, K> leftKeyExtractor, 
            List<R> rightList, 
            Function<R, K> rightKeyExtractor) {

        // Group rightList by join key
        Map<K, List<R>> rightGrouped = rightList.stream()
                .collect(Collectors.groupingBy(rightKeyExtractor));

        // Perform LEFT JOIN
        return leftList.stream()
                .collect(Collectors.toMap(
                        left -> left,
                        left -> rightGrouped.getOrDefault(leftKeyExtractor.apply(left), new ArrayList<>())
                ));
    }
}

public class GenericCrudExample {
    public static void main(String[] args) {

        GenericRepository<User, Long> userRepo = new InMemoryRepository<>();
        GenericRepository<Post, Long> postRepo = new InMemoryRepository<>();

        userRepo.create(new User(1L, "Alice", 25), 1L);
        userRepo.create(new User(2L, "Bob", 30), 2L);
        userRepo.create(new User(3L, "Charlie", 17), 3L); // No posts

        postRepo.create(new Post(101L, 1L, "Alice's first post"), 101L);
        postRepo.create(new Post(102L, 2L, "Bob's first post"), 102L);
        postRepo.create(new Post(103L, 1L, "Alice's second post"), 103L);

        // WHERE Clause: Filter users younger than 18
        List<User> youngUsers = userRepo.filter(user -> user.getAge() < 18);
        System.out.println("Users younger than 18: " + youngUsers);

        // INNER JOIN Users and Posts
        Map<User, List<Post>> userPosts = JoinUtils.innerJoin(
                userRepo.getAll(),
                User::getId,
                postRepo.getAll(),
                Post::getUserId
        );

        System.out.println("\nINNER JOIN (Users with Posts):");
        userPosts.forEach((user, posts) -> System.out.println(user + " -> " + posts));

        // LEFT JOIN Users and Posts
        Map<User, List<Post>> leftJoinResult = JoinUtils.leftJoin(
                userRepo.getAll(),
                User::getId,
                postRepo.getAll(),
                Post::getUserId
        );

        System.out.println("\nLEFT JOIN (All Users, Even Without Posts):");
        leftJoinResult.forEach((user, posts) -> System.out.println(user + " -> " + posts));

        System.out.println("\nUsers sorted by Age (Ascending):");
        userRepo.getAllSorted(Comparator.comparing(User::getAge))
                .forEach(System.out::println);

        System.out.println("\nPosts sorted by Post ID (Descending):");
        postRepo.getAllSorted(Comparator.comparing(Post::getPostId).reversed())
                        .forEach(System.out::println);
    }
}